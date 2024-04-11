package test

import (
	api "SDAS/kitex_gen/api"
	sdas "SDAS/kitex_gen/api/sdas"
	"context"
	"encoding/gob"
	"image"
	"io"
	"testing"
	"time"

	"github.com/cloudwego/kitex/client/streamclient"
	"gocv.io/x/gocv"
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

func NewStreamClient() sdas.StreamClient {
	streamClient := sdas.MustNewStreamClient(
		"test",
		streamclient.WithHostPorts("127.0.0.1:8088"),
	)
	return streamClient
}

func TestPullExpose(t *testing.T) {
	//new client
	client := NewClient(t)
	//add source
	resp, err := client.AddSource(context.Background(), &api.AddSourceRequest{
		Source: &api.Source{
			Type:    "rtsp",
			Name:    "test",
			Content: `{"url":"rtsp://admin:a12345678@192.168.0.238","format":"h264"}`,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
	time.Sleep(time.Second)
	resp2_2, err := client.ListExposes(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp2_2)
	// l := len(resp2_2.Exposes)

	//stream client
	streamClient := NewStreamClient()
	stream, err := streamClient.PullExposeStream(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	//open
	req := &api.PullExposeStreamRequest{
		Expose: &api.Expose{
			Name:       "out",
			SourceName: "test",
			Type:       "pull",
			Content:    `{"op":0,"msg_type":"rtspMsg"}`,
		},
	}
	if err = stream.Send(req); err != nil {
		t.Errorf("echo.send: failed, err = " + err.Error())
	}

	//play
	req = &api.PullExposeStreamRequest{
		Expose: &api.Expose{
			Name:       "out",
			SourceName: "test",
			Type:       "pull",
			Content:    `{"op":1,"msg_type":"rtspMsg"}`,
		},
	}
	if err = stream.Send(req); err != nil {
		t.Errorf("echo.send: failed, err = " + err.Error())
	}
	time.Sleep(time.Second)

	gob.Register(image.RGBA{})

	window := gocv.NewWindow("test")

	window.ResizeWindow(512, 512)

	for i := 0; i < 1; i++ {
		resp, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			t.Errorf("echo.recv: failed, err = " + err.Error())
			continue
		} else if resp.Code != 0 {
			t.Errorf("echo.recv: failed, err = " + resp.Message)
			continue
		}

		data := resp.Data
		// buf := new(bytes.Buffer)
		// buf.WriteString(data)

		// dec := gob.NewDecoder(buf)

		// var msg image.RGBA
		// err = dec.Decode(&msg)
		// if err != nil {
		// 	t.Errorf("Error decoding:%v", err)
		// }
		t.Log(data)
		// 	mat, err := gocv.ImageToMatRGB(img)
		// 	defer mat.Close()

		// 	window.IMShow(mat)
		// 	// 等待一段时间或检测按键事件，以保持窗口打开并实时更新
		// 	window.WaitKey(1)
	}

	//close
	// req = &api.PullExposeStreamRequest{
	// 	Expose: &api.Expose{
	// 		Name:       "out",
	// 		SourceName: "test",
	// 		Type:       "pull",
	// 		Content:    `{"op":3,"msg_type":"rtspMsg"}`,
	// 	},
	// }
	// if err = stream.Send(req); err != nil {
	// 	t.Errorf("echo.send: failed, err = " + err.Error())
	// }
	// time.Sleep(time.Second)

	resp2_2, err = client.ListExposes(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp2_2)
	// if len(resp2_2.Exposes) != l {
	// 	t.Fatal("test len failed")
	// }

	resp3, err := client.RemoveSource(context.Background(), &api.RemoveSourceRequest{
		Name: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp3)
	if resp3.Code != 0 {
		t.Fatal(resp3.Message)
	}
}

// func TestRtspExpose(t *testing.T) {
// 	client := NewClient(t)
// 	resp, err := client.AddSource(context.Background(), &api.AddSourceRequest{
// 		Source: &api.Source{
// 			Type:    "rtsp",
// 			Name:    "test",
// 			Content: `{"url":"rtsp://admin:a12345678@192.168.0.238","format":"h264"}`,
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(resp)
// 	/**
// 	test expose
// 	**/
// 	resp2_2, err := client.ListExposes(context.Background())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log(resp2_2)
// 	l := len(resp2_2.Exposes)
// 	resp2_1, err := client.AddExpose(context.Background(), &api.AddExposeRequest{
// 		Expose: &api.Expose{
// 			Type: "rtsp",
// 			Name: "test",
// 			Content: `{"rtsp_address":":8554","udp_rtp_address":":10000","udp_rtcp_address":":10001","multicast_ip_range":"224.1.0.0/16","MulticastRTPPort":  "8002",
// 			"MulticastRTCPPort": "8003","format": "h264"}`,
// 			SourceName: "test",
// 		},
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(resp2_1)

// 	resp2_2, err = client.ListExposes(context.Background())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log(resp2_2)
// 	if resp2_2.Exposes[0].Name != "test" {
// 		t.Fatal("test name failed")
// 	}
// 	if len(resp2_2.Exposes) != l+1 {
// 		t.Fatal("test len failed")
// 	}

// 	resp2_3, err := client.RemoveExpose(context.Background(), &api.RemoveExposeRequest{
// 		Name: "test",
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(resp2_3)

// 	resp2_2, err = client.ListExposes(context.Background())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log(resp2_2)
// 	if len(resp2_2.Exposes) != l {
// 		t.Fatal("test len failed")
// 	}

// 	resp3, err := client.RemoveSource(context.Background(), &api.RemoveSourceRequest{
// 		Name: "test",
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(resp3)
// 	if resp3.Code != 0 {
// 		t.Fatal(resp3.Message)
// 	}
// }
