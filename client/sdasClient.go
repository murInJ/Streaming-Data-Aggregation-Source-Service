package client

import (
	"SDAS/kitex_gen/api"
	"SDAS/kitex_gen/api/sdas"
	"context"
	"errors"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/streamclient"
	"strconv"
)

type SDASClient struct {
	c      sdas.Client
	sc     sdas.StreamClient
	stream sdas.SDAS_PullExposeStreamClient
	isc    bool
	issc   bool
}

func NewSDASClient(url string, isc bool, issc bool) (*SDASClient, error) {
	sdasClient := &SDASClient{
		isc:  isc,
		issc: issc,
	}
	if isc {
		c, err := sdas.NewClient("SDAS", client.WithHostPorts(url))
		if err != nil {
			return nil, err
		}
		sdasClient.c = c
	}
	if issc {
		sc := sdas.MustNewStreamClient("SDAS", streamclient.WithHostPorts(url))
		sdasClient.sc = sc
		stream, err := sc.PullExposeStream(context.Background())
		if err != nil {
			return nil, err
		}
		sdasClient.stream = stream
	}
	return sdasClient, nil
}

func (c *SDASClient) ListSources() ([]*api.Source, error) {
	resp, err := c.c.ListSources(context.Background())
	if err != nil {
		return nil, err
	}
	if resp.Code == -1 {
		return nil, errors.New(resp.Message)
	}
	return resp.Sources, nil
}

func (c *SDASClient) AddSource(Type, Name string, Expose bool, Content map[string]string) error {
	req := &api.AddSourceRequest{
		Source: &api.Source{
			Type:    Type,
			Name:    Name,
			Content: Content,
			Expose:  Expose,
		},
	}
	resp, err := c.c.AddSource(context.Background(), req)
	if err != nil {
		return err
	}
	if resp.Code == -1 {
		return errors.New(resp.Message)
	}
	return nil
}

func (c *SDASClient) RemoveSource(Name string) error {
	req := &api.RemoveSourceRequest{
		Name: Name,
	}
	resp, err := c.c.RemoveSource(context.Background(), req)
	if err != nil {
		return err
	}
	if resp.Code == -1 {
		return errors.New(resp.Message)
	}
	return nil
}

func (c *SDASClient) ListExposes() ([]*api.Expose, error) {
	resp, err := c.c.ListExposes(context.Background())
	if err != nil {
		return nil, err
	}
	if resp.Code == -1 {
		return nil, errors.New(resp.Message)
	}
	return resp.Exposes, nil
}

func (c *SDASClient) SendPullExposeStream(Name, SourceName string, Op int) error {
	req := &api.PullExposeStreamRequest{
		Expose: &api.Expose{
			Name:       Name,
			SourceName: SourceName,
			Type:       "pull",
			Content: map[string]string{
				"op": strconv.Itoa(Op),
			},
		},
	}
	if err := c.stream.Send(req); err != nil {
		return err
	}
	return nil
}

func (c *SDASClient) RecvPullExposeStream() (*api.SourceMsg, error) {
	resp, err := c.stream.Recv()
	if err != nil {
		return nil, err
	} else if resp.Code == -1 {
		return nil, errors.New(resp.Message)
	}
	return resp.SourceMsg, nil
}

func (c *SDASClient) AddExpose(Type, Name, SourceName string, Content map[string]string) error {
	req := &api.AddExposeRequest{
		Expose: &api.Expose{
			Type:       Type,
			Name:       Name,
			SourceName: SourceName,
			Content:    Content,
		},
	}
	resp, err := c.c.AddExpose(context.Background(), req)
	if err != nil {
		return err
	}
	if resp.Code == -1 {
		return errors.New(resp.Message)
	}
	return nil
}

func (c *SDASClient) RemoveExpose(Name string) error {
	req := &api.RemoveExposeRequest{
		Name: Name,
	}
	resp, err := c.c.RemoveExpose(context.Background(), req)
	if err != nil {
		return err
	}
	if resp.Code == -1 {
		return errors.New(resp.Message)
	}
	return nil
}

func (c *SDASClient) Close() error {
	if c.issc {
		err := c.stream.Close()
		if err != nil {
			return err
		}
	}
	c.issc = false
	c.isc = false
	return nil
}
