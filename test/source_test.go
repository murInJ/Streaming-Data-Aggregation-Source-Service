package test

import (
	cli "SDAS/client"
	"encoding/json"
	"io"
	"testing"
)

func TestRtspSource(t *testing.T) {
	c, err := cli.NewSDASClient("0.0.0.0:8088", true, false)
	if err != nil {
		t.Fatal(err)
	}

	sources, err := c.ListSources()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sources)
	l := len(sources)

	err = c.AddSource("rtsp", "source_rtsp_test", true, map[string]string{
		"url":    "rtsp://admin:a12345678@192.168.0.238",
		"format": "h264",
	})
	if err != nil {
		t.Fatal(err)
	}

	sources, err = c.ListSources()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sources)
	if len(sources) != l+1 {
		t.Fatal("sources length not equal 1")
	}

	err = c.RemoveSource("source_rtsp_test")
	if err != nil {
		t.Fatal(err)
	}

	sources, err = c.ListSources()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sources)
	if len(sources) != l {
		t.Fatal("sources not empty")
	}
}

func TestPluginSource(t *testing.T) {
	c, err := cli.NewSDASClient("0.0.0.0:8088", true, true)
	if err != nil {
		t.Fatal(err)
	}

	sources, err := c.ListSources()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sources)

	err = c.AddSource("rtsp", "source_rtsp_test1", true, map[string]string{
		"url":    "rtsp://admin:a12345678@192.168.0.238",
		"format": "h264",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = c.AddSource("rtsp", "source_rtsp_test2", true, map[string]string{
		"url":    "rtsp://admin:a12345678@192.168.0.241",
		"format": "h264",
	})
	if err != nil {
		t.Fatal(err)
	}

	iem, err := json.Marshal(map[string]string{"source_rtsp_test1": "Entry1", "source_rtsp_test2": "Entry2"})
	if err != nil {
		t.Fatal(err)
	}
	eom, err := json.Marshal(map[string]string{"Exist": "source_VideoWall_test"})
	if err != nil {
		t.Fatal(err)
	}
	oem, err := json.Marshal(map[string]string{"source_VideoWall_test": "true"})
	if err != nil {
		t.Fatal(err)
	}

	err = c.AddSource("plugin", "VideoWall", true, map[string]string{
		"inSource_entry_map":   string(iem),
		"exist_outSource_map":  string(eom),
		"outSource_expose_map": string(oem),
		"args":                 ``,
	})
	if err != nil {
		t.Fatal(err)
	}

	sources, err = c.ListSources()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sources)

	//open
	err = c.SendPullExposeStream("expose_pull_test", "source_VideoWall_test", 0)
	if err != nil {
		t.Fatal(err)
	}

	//play
	err = c.SendPullExposeStream("expose_pull_test", "source_VideoWall_test", 1)
	if err != nil {
		t.Fatal(err)
	}

	//window := gocv.NewWindow("expose_pull_test")
	//window.ResizeWindow(512, 512)

	for i := 0; i < 100; i++ {
		msg, err := c.RecvPullExposeStream()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		t.Log(msg.Ntp, msg.DataType)

		//
		//var img image.RGBA
		//err = msgpack.Unmarshal(msg.Data, &img)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//
		//mat, err := gocv.ImageToMatRGB(&img)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//window.IMShow(mat)
		//
		//window.WaitKey(1)
		//err = mat.Close()
		//if err != nil {
		//	t.Fatal(err)
		//}
		//

	}

	//close
	err = c.SendPullExposeStream("expose_pull_test", "source_VideoWall_test", 3)
	if err != nil {
		t.Fatal(err)
	}

	err = c.RemoveSource("source_rtsp_test1")
	if err != nil {
		t.Fatal(err)
	}
	err = c.RemoveSource("source_rtsp_test2")
	if err != nil {
		t.Fatal(err)
	}

	sources, err = c.ListSources()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Close()
	if err != nil {
		t.Fatal(err)
	}
}
