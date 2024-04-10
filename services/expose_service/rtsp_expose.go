package services

import (
	config "SDAS/config"
	"SDAS/services/expose_service/encoder"
	source "SDAS/services/source_service"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"sync"
	"time"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/description"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pion/rtp"
)

var (
	CLOSE = 0
	ERR   = 1
	OPEN  = 2
)

type ExposeEntityRtsp struct {
	ControlChannel *chan int
	InputChannel   *chan source.MessageRtsp
	Status         int
	Expose         *config.EXPOSE_RTSP
	Name           string
	Type           string
	Server         *RtspServer
	Stream         *gortsplib.ServerStream
	Desc           *description.Session
	SourceName     string
	Encoder        any
	Height int,
	Width int,
}

func NewExposeEntityRtsp(name string, expose *config.EXPOSE_RTSP, sourceName string) (*ExposeEntityRtsp, error) {
	control_channel := make(chan int)

	i, ok := source.Sources.Load(sourceName)
	if !ok {
		err := fmt.Errorf("source not found")
		klog.Error(err)
		return nil, err
	}
	source_entity := i.(*source.SourceEntityRtsp)
	switch expose.Format {
	case "h264":
		h264encoder := &encoder.H264Encoder{}
		// err := encoder.Initialize(expose.Width, expose.Height)
		// if err != nil {
		// 	return nil, err
		// }
		entity := &ExposeEntityRtsp{
			Name:           name,
			Type:           "rtsp",
			ControlChannel: &control_channel,
			InputChannel:   source_entity.OutputChannel,
			Status:         CLOSE,
			Expose:         expose,
			SourceName:     sourceName,
			Encoder:        h264encoder,
			Height: 0,
			Width: 0,
		}
		return entity, nil
	}
	return nil, errors.New("expose Format type not support")

}

func (e *ExposeEntityRtsp) GetExposeString() (string, error) {
	b, err := json.Marshal(e.Expose)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (e *ExposeEntityRtsp) GetName() string {
	return e.Name
}

func (e *ExposeEntityRtsp) GetType() string {
	return e.Type
}

func (e *ExposeEntityRtsp) GetSourceName() string {
	return e.SourceName
}

func (e *ExposeEntityRtsp) Start() {
	go e.goroutine_rtsp_expose()
}

func (e *ExposeEntityRtsp) Stop() {
	*e.ControlChannel <- CLOSE
}

func (e *ExposeEntityRtsp) goroutine_rtsp_expose() {

	c, stream, err := e.startup_rstp()
	e.Stream = stream

	if err != nil {
		e.Status = ERR
		klog.Errorf("expose[rtsp]: %s open failed.", e.Name)
		return
	}
	e.Status = OPEN
	klog.Infof("expose[rtsp]: %s opened.\n", e.Name)
	for {
		select {
		case command := <-*e.ControlChannel:
			switch command {
			case CLOSE:
				e.Server.setStreamUnready()
				c.Close()
				switch e.Expose.Format{
					case "h264":
						e.Encoder.(*encoder.H264Encoder).Close()
				}
				e.Status = CLOSE
				klog.Infof("expose[rtsp]: %s closed.\n", e.Name)
				return
			}
		case msgRtsp := <-*e.InputChannel:
			switch e.Expose.Format {
			case "h264":
				err := e.handler_h264(*msgRtsp.Img, msgRtsp.NTP)
				if err != nil {
					klog.Errorf("expose[rtsp]: %s handler_h264 error: %v.\n", e.Name, err)
				}
			}
		}
	}
}

func (e *ExposeEntityRtsp) startup_rstp() (*gortsplib.Client, *gortsplib.ServerStream, error) {
	e.Server = NewRtspServer(e.Name, e.Expose)
	c := gortsplib.Client{}

	// parse URL
	u, err := base.ParseURL(e.Expose.RTSPAddress)
	if err != nil {
		klog.Error(err)
		return nil, nil, err
	}

	// connect to the server
	err = c.Start(u.Scheme, u.Host)
	if err != nil {
		klog.Error(err)
		return nil, nil, err
	}

	desc, _, err := c.Describe(u)
	if err != nil {
		klog.Error(err)
		return nil, nil, err
	}
	e.Desc = desc
	stream := e.Server.setStreamReady(desc)

	return &c, stream, nil
}

func (e *ExposeEntityRtsp) handler_h264(img image.Image, ntp int64) error {
	ntp_time := time.Unix(ntp/int64(time.Second), ntp%int64(time.Second))
	
	if e.Height == 0 && e.Width == 0 {
		e.Width = img.Bounds().Dx()
		e.Height = img.Bounds().Dy()
		e.Encoder.(*encoder.H264Encoder).Initialize(img.Bounds().Dx(), img.Bounds().Dy())
	}
	
	encoded_img, err := e.Encoder.(*encoder.H264Encoder).Encode(img)
	if err != nil {
		return err
	}
	packet := &rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			PayloadType:    96,
			SequenceNumber: 0,
			Timestamp:      uint32(ntp),
			SSRC:           123456,
		},
		Payload: encoded_img,
	}
	err = e.Stream.WritePacketRTPWithNTP(e.Desc.Medias[0], packet, ntp_time)
	if err != nil {
		return err
	}
	return nil
}

/**
rtsp server
**/

type RtspServer struct {
	s      *gortsplib.Server
	mutex  sync.Mutex
	stream *gortsplib.ServerStream
	name   string
}

func NewRtspServer(name string, config *config.EXPOSE_RTSP) *RtspServer {
	s := &RtspServer{
		name: name,
	}
	// configure the server
	s.s = &gortsplib.Server{
		Handler:           s,
		RTSPAddress:       config.RTSPAddress,
		UDPRTPAddress:     config.UDPRTPAddress,
		UDPRTCPAddress:    config.UDPRTCPAddress,
		MulticastIPRange:  config.MulticastIPRange,
		MulticastRTPPort:  config.MulticastRTPPort,
		MulticastRTCPPort: config.MulticastRTCPPort,
	}
	return s
}

// called when a connection is opened.
func (s *RtspServer) OnConnOpen(ctx *gortsplib.ServerHandlerOnConnOpenCtx) {
	klog.Infof("expose[rtsp]: %s conn opened", s.name)
}

// called when a connection is closed.
func (s *RtspServer) OnConnClose(ctx *gortsplib.ServerHandlerOnConnCloseCtx) {
	klog.Infof("expose[rtsp]: %s conn closed (%v)", ctx.Error)
}

// called when a session is opened.
func (s *RtspServer) OnSessionOpen(ctx *gortsplib.ServerHandlerOnSessionOpenCtx) {
	klog.Infof("expose[rtsp]: %s session opened", s.name)
}

// called when a session is closed.
func (s *RtspServer) OnSessionClose(ctx *gortsplib.ServerHandlerOnSessionCloseCtx) {
	klog.Infof("expose[rtsp]: %s session closed", s.name)
}

// called when receiving a DESCRIBE request.
func (s *RtspServer) OnDescribe(ctx *gortsplib.ServerHandlerOnDescribeCtx) (*base.Response, *gortsplib.ServerStream, error) {
	klog.Infof("expose[rtsp]: %s describe request", s.name)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// stream is not available yet
	if s.stream == nil {
		return &base.Response{
			StatusCode: base.StatusNotFound,
		}, nil, nil
	}

	return &base.Response{
		StatusCode: base.StatusOK,
	}, s.stream, nil
}

// called when receiving a SETUP request.
func (s *RtspServer) OnSetup(ctx *gortsplib.ServerHandlerOnSetupCtx) (*base.Response, *gortsplib.ServerStream, error) {
	klog.Infof("expose[rtsp]: %s setup request", s.name)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// stream is not available yet
	if s.stream == nil {
		return &base.Response{
			StatusCode: base.StatusNotFound,
		}, nil, nil
	}

	return &base.Response{
		StatusCode: base.StatusOK,
	}, s.stream, nil
}

// called when receiving a PLAY request.
func (s *RtspServer) OnPlay(ctx *gortsplib.ServerHandlerOnPlayCtx) (*base.Response, error) {
	klog.Infof("expose[rtsp]: %s play request", s.name)

	return &base.Response{
		StatusCode: base.StatusOK,
	}, nil
}

func (s *RtspServer) setStreamReady(desc *description.Session) *gortsplib.ServerStream {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.stream = gortsplib.NewServerStream(s.s, desc)
	return s.stream
}

func (s *RtspServer) setStreamUnready() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.stream.Close()
	s.stream = nil
}
