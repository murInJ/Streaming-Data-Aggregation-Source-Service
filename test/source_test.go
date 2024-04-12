package test

import (
	cli "SDAS/client"
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

	err = c.AddSource("rtsp", "source_rtsp_test", "rtsp://admin:a12345678@192.168.0.238", "h264", true)
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
