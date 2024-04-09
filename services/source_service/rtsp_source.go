package services

import (
	config "SDAS/config"
	decoder "SDAS/services/source_service/decoder"
	"encoding/json"
	"fmt"
	"image"
	"runtime"
	"time"

	"github.com/aler9/gortsplib/pkg/rtpcodecs/rtph264"
	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/description"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pion/rtp"
)

var (
	CLOSE = 0
	ERR   = 1
	OPEN  = 2
)

type MessageRtsp struct {
	Img *image.Image
	NTP int64
}

type EntityRtsp struct {
	ControlChannel *chan int
	OutputChannel  *chan MessageRtsp
	Status         int
	Source         *config.SOURCE_RTSP
	Name           string
	Type           string
}

func NewEntityRtsp(name string, source *config.SOURCE_RTSP) *EntityRtsp {
	control_channel := make(chan int)
	output_channel := make(chan MessageRtsp, 1024)

	entity := &EntityRtsp{
		Name:           name,
		Type:           "rtsp",
		ControlChannel: &control_channel,
		OutputChannel:  &output_channel,
		Status:         CLOSE,
		Source:         source,
	}
	return entity
}

func (e *EntityRtsp) GetSourceString() (string, error) {
	b, err := json.Marshal(e.Source)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *EntityRtsp) GetName() string {
	return e.Name
}

func (e *EntityRtsp) GetType() string {
	return e.Type
}

func (e *EntityRtsp) Start() {
	go e.goroutine_rtsp_source()
}

func (e *EntityRtsp) Stop() {
	*e.ControlChannel <- CLOSE
}

func (e *EntityRtsp) goroutine_rtsp_source() {
	c, err := e.startup_rstp()
	if err != nil {
		e.Status = ERR
		klog.Errorf("source[rtsp]: %s open failed.\n", e.Name)
		return
	}
	e.Status = OPEN
	klog.Infof("source[rtsp]: %s opened.\n", e.Name)
	for {
		select {
		case command := <-*e.ControlChannel:
			switch command {
			case CLOSE:
				c.Close()
				e.Status = CLOSE
				klog.Infof("source[rtsp]: %s closed.\n", e.Name)
				return
			}
		}
	}
}

func (e *EntityRtsp) startup_rstp() (*gortsplib.Client, error) {

	c := gortsplib.Client{}

	// parse URL
	u, err := base.ParseURL(e.Source.Url)
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
	defer c.Close()

	desc, _, err := c.Describe(u)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	switch e.Source.Format {
	case "h264":
		e.handler_h264(&c, desc)
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

func (e *EntityRtsp) handler_h264(c *gortsplib.Client, desc *description.Session) error {
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

	frameDec := &decoder.H264Decoder{}
	err = frameDec.Initialize()
	if err != nil {
		klog.Error(err)
		return err
	}
	defer frameDec.Close()

	if forma.SPS != nil {
		frameDec.Decode(forma.SPS)
	}
	if forma.PPS != nil {
		frameDec.Decode(forma.PPS)
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
			if err != rtph264.ErrNonStartingPacketAndNoPrevious && err != rtph264.ErrMorePacketsNeeded {
				// log.Printf("ERR: %v", err)
				runtime.Gosched()
			}
			return
		}

		for _, nalu := range au {
			// convert NALUs into RGBA frames
			img, err := frameDec.Decode(nalu)
			if err != nil {
				panic(err)
			}

			// wait for a frame
			if img == nil {
				runtime.Gosched()
				continue
			}

			*e.OutputChannel <- MessageRtsp{Img: &img, NTP: time.Now().UnixNano()}

		}
	})

	return nil
}
