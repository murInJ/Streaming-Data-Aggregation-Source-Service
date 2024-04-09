package test

import (
	api "SDAS/kitex_gen/api"
	sdas "SDAS/kitex_gen/api/sdas"
	"context"
	"testing"

	"github.com/cloudwego/kitex/client"
)

func NewClient(t *testing.T) sdas.Client {
	sdasClient, err := sdas.NewClient("SDAS", client.WithHostPorts("0.0.0.0:8088"))
	if err != nil {
		t.Fatal(err)
	}
	return sdasClient
}

func TestRtspSource(t *testing.T) {
	client := NewClient(t)
	resp2, err := client.ListSources(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp2)
	if resp2.Code != 0 {
		t.Fatal(resp2.Message)
	}
	if resp2.Sources == nil {
		t.Fatal("no sources")
	}
	l := len(resp2.Sources)

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
	resp2, err = client.ListSources(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp2)
	if resp2.Code != 0 {
		t.Fatal(resp2.Message)
	}
	if resp2.Sources == nil {
		t.Fatal("no sources")
	}
	if len(resp2.Sources) != l+1 {
		t.Fatal("sources length not equal 1")
	}
	if resp2.Sources[0].Name != "test" {
		t.Fatal("sources name not equal rtsp_source")
	}
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
	resp2, err = client.ListSources(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp2)
	if len(resp2.Sources) != l {
		t.Fatal("sources not empty")
	}
}
