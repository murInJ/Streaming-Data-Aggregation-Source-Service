package services

import (
	"SDAS/config"
	"SDAS/kitex_gen/api"
	decoder "SDAS/utils/decoder/DecoderRTSP"
	"errors"
	"fmt"
	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/description"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/bluenviron/gortsplib/v4/pkg/format/rtph264"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pion/rtp"
	"github.com/vmihailenco/msgpack/v5"
	"image"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

/**content
	Url
	Format
**/

type SourceEntityRtsp struct {
	ControlChannel *chan int
	OutputChannel  *chan *api.SourceMsg
	Status         int
	Content        map[string]string
	Name           string
	Type           string
	Decoder        decoder.DecoderRTSP
	Expose         bool
	once           sync.Once
	requested      atomic.Bool
}

func NewSourceEntityRtsp(name string, expose bool, content map[string]string) (*SourceEntityRtsp, error) {

	controlChannel := make(chan int)
	outputChannel := make(chan *api.SourceMsg, 1024)
	if content["format"] == "h264" {
		frameDec := decoder.H264Decoder{}
		err := frameDec.Initialize()
		if err != nil {
			return nil, err
		}
		entity := &SourceEntityRtsp{
			Name:           name,
			Type:           "rtsp",
			ControlChannel: &controlChannel,
			OutputChannel:  &outputChannel,
			Status:         config.CLOSE,
			Content:        content,
			Decoder:        &frameDec,
			Expose:         expose,
		}
		entity.requested.Store(false)
		return entity, nil
	} else {
		return nil, errors.New("unknown source type")
	}

}

func (e *SourceEntityRtsp) GetConfig() *api.Source {
	return &api.Source{
		Type:    e.Type,
		Name:    e.Name,
		Expose:  e.Expose,
		Content: e.Content,
	}
}

func (e *SourceEntityRtsp) RequestOutChannel() (*chan *api.SourceMsg, error) {
	if e.requested.CompareAndSwap(false, true) {
		return e.OutputChannel, nil
	} else {
		return nil, errors.New("request out channel already in use")
	}
}

func (e *SourceEntityRtsp) ReleaseOutChannel() {
	e.requested.Store(false)
}

func (e *SourceEntityRtsp) Start() error {
	if e.Status != config.CLOSE && e.Status != config.ERR {
		return errors.New("source already started")
	}
	go e.goroutineRtspSource()
	for {
		switch e.Status {
		case config.OPEN:
			return nil
		case config.ERR:
			err := errors.New("rtsp source start error")
			return err
		default:
			runtime.Gosched()
		}
	}
}

func (e *SourceEntityRtsp) Stop() {
	e.once.Do(func() {
		*e.ControlChannel <- config.CLOSE
		close(*e.OutputChannel)
		close(*e.ControlChannel)
	})
}

func (e *SourceEntityRtsp) GetName() string {
	return e.Name
}

func (e *SourceEntityRtsp) goroutineRtspSource() {
	c, err := e.startupRtsp()
	if err != nil {
		e.Status = config.ERR
		klog.Errorf("source[rtsp]: %s open failed. %v", e.Name, err)
		return
	}
	e.Status = config.OPEN
	klog.Infof("source[rtsp]: %s opened.\n", e.Name)
	for {
		select {
		case command := <-*e.ControlChannel:
			switch command {
			case config.CLOSE:
				c.Close()
				e.Decoder.Close()
				e.Status = config.CLOSE
				klog.Infof("source[rtsp]: %s closed.\n", e.Name)
				return
			}
		}
	}
}

func (e *SourceEntityRtsp) startupRtsp() (*gortsplib.Client, error) {

	c := gortsplib.Client{}

	// parse URL
	u, err := base.ParseURL(e.Content["url"])
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	// connect to the server
	err = c.Start(u.Scheme, u.Host)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	desc, _, err := c.Describe(u)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	switch e.Content["format"] {
	case "h264":
		err := e.handlerH264(&c, desc)
		if err != nil {
			klog.Error(err)
			return nil, err
		}
	default:
		err := fmt.Errorf("unsupported format")
		klog.Error(err)
		return nil, err
	}
	// start playing
	_, err = c.Play(nil)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	return &c, nil
}

func (e *SourceEntityRtsp) handlerH264(c *gortsplib.Client, desc *description.Session) error {

	var forma *format.H264
	medi := desc.FindFormat(&forma)
	if medi == nil {
		err := fmt.Errorf("media not found")
		klog.Error(err)
		return err
	}

	rtpDec, err := forma.CreateDecoder()
	if err != nil {
		klog.Error(err)
		return err
	}

	if forma.SPS != nil {
		e.Decoder.Decode(forma.SPS)
	}
	if forma.PPS != nil {
		e.Decoder.Decode(forma.PPS)
	}

	// setup a single media
	_, err = c.Setup(desc.BaseURL, medi, 0, 0)
	if err != nil {
		klog.Error(err)
		return err
	}

	c.OnPacketRTPAny(func(medi *description.Media, forma format.Format, pkt *rtp.Packet) {
		// extract access units from RTP packets
		au, err := rtpDec.Decode(pkt)
		if err != nil {
			if !errors.Is(err, rtph264.ErrNonStartingPacketAndNoPrevious) && !errors.Is(err, rtph264.ErrMorePacketsNeeded) {
				// log.Printf("ERR: %v", err)
				runtime.Gosched()
			}
			return
		}
		ntp, ntpAvailable := c.PacketNTP(medi, pkt)

		for _, nalu := range au {
			e.handlerNalu(nalu, ntp, ntpAvailable)
		}
	})

	return nil
}

func (e *SourceEntityRtsp) handlerNalu(nalu []byte, ntp time.Time, ntpAvailable bool) {
	defer func() {
		err := recover()
		if err != nil {
			klog.Error(err)
		}
		return
	}()

	// convert NALUs into RGBA frames
	img, err := e.Decoder.Decode(nalu)
	if err != nil {
		klog.Error(err)
		return
	}
	if img == nil {
		runtime.Gosched()
		return
	}
	rgba := img.(*image.RGBA)
	b, err := msgpack.Marshal(rgba)
	if err != nil {
		klog.Error(err)
		return
	}
	msg := &api.SourceMsg{
		Data:     b,
		DataType: "image.RGBA",
	}
	if ntpAvailable {
		msg.Ntp = ntp.UnixNano()
	} else {
		msg.Ntp = time.Now().UnixNano()
	}

	*e.OutputChannel <- msg

}
