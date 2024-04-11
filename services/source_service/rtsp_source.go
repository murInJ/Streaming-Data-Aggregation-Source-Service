package services

import (
	config "SDAS/config"
	decoder "SDAS/utils/decoder"
	"encoding/json"
	"errors"
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
	Img image.RGBA
	NTP int64
}

func (m MessageRtsp) GetNTP() int64 {
	return m.NTP
}

func (m MessageRtsp) GetImage() (image.Image, bool) {
	return &m.Img, true
}

type SourceEntityRtsp[T SOURCE_MESSAGE] struct {
	ControlChannel *chan int
	OutputChannel  *chan T
	Status         int
	Source         *config.SOURCE_RTSP
	Name           string
	Type           string
	Decoder        any
}

func NewSourceEntityRtsp(name string, source *config.SOURCE_RTSP) (*SourceEntityRtsp[MessageRtsp], error) {
	control_channel := make(chan int)
	output_channel := make(chan MessageRtsp, 1024)
	if source.Format == "h264" {
		frameDec := &decoder.H264Decoder{}
		err := frameDec.Initialize()
		if err != nil {
			return nil, err
		}
		frameDec.Initialize()
		entity := &SourceEntityRtsp[MessageRtsp]{
			Name:           name,
			Type:           "rtsp",
			ControlChannel: &control_channel,
			OutputChannel:  &output_channel,
			Status:         CLOSE,
			Source:         source,
			Decoder:        frameDec,
		}
		return entity, nil
	}
	return nil, errors.New("unknown source type")
}

func (e *SourceEntityRtsp[T]) GetSourceString() (string, error) {
	b, err := json.Marshal(e.Source)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *SourceEntityRtsp[T]) GetName() string {
	return e.Name
}

func (e *SourceEntityRtsp[T]) GetType() string {
	return e.Type
}

func (e *SourceEntityRtsp[T]) GetOutChannel() *chan T {
	return e.OutputChannel
}

func (e *SourceEntityRtsp[T]) send_to_out_channel(msg T) {
	*e.OutputChannel <- msg
}

func (e *SourceEntityRtsp[T]) Start() {
	go e.goroutine_rtsp_source()
}

func (e *SourceEntityRtsp[T]) Stop() {
	*e.ControlChannel <- CLOSE
	close(*e.OutputChannel)
	close(*e.ControlChannel)
}

func (e *SourceEntityRtsp[T]) goroutine_rtsp_source() {
	c, err := e.startup_rstp()
	if err != nil {
		e.Status = ERR
		klog.Errorf("source[rtsp]: %s open failed. %v", e.Name, err)
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
				switch e.Source.Format {
				case "h264":
					e.Decoder.(*decoder.H264Decoder).Close()
				}
				e.Status = CLOSE
				klog.Infof("source[rtsp]: %s closed.\n", e.Name)
				return
			}
		}
	}
}

func (e *SourceEntityRtsp[T]) startup_rstp() (*gortsplib.Client, error) {

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

func (e *SourceEntityRtsp[T]) handler_h264(c *gortsplib.Client, desc *description.Session) error {
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
		e.Decoder.(*decoder.H264Decoder).Decode(forma.SPS)
	}
	if forma.PPS != nil {
		e.Decoder.(*decoder.H264Decoder).Decode(forma.PPS)
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
		ntp, ntpAvailable := c.PacketNTP(medi, pkt)

		for _, nalu := range au {
			// convert NALUs into RGBA frames
			img, err := e.Decoder.(*decoder.H264Decoder).Decode(nalu)
			if err != nil {
				panic(err)
			}

			// wait for a frame
			if img == nil {
				runtime.Gosched()
				continue
			}
			rgba := img.(*image.RGBA)
			var msg SOURCE_MESSAGE
			if ntpAvailable {
				msg = MessageRtsp{Img: *rgba, NTP: ntp.UnixNano()}
			} else {
				msg = MessageRtsp{Img: *rgba, NTP: time.Now().UnixNano()}
			}

			e.send_to_out_channel(msg.(T))

		}
	})

	return nil
}
