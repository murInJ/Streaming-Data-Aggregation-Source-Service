package main

import (
	config "SDAS/config"
	handler "SDAS/handlers"
	api "SDAS/kitex_gen/api"
	expose "SDAS/services/expose_service"
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
)

// SDASImpl implements the last service interface defined in the IDL.
type SDASImpl struct{}

// AddSource implements the SDASImpl interface.
func (s *SDASImpl) AddSource(ctx context.Context, req *api.AddSourceRequest) (resp *api.AddSourceResponse, err error) {
	err = handler.AddSourceHandler(req.Source.Name, req.Source.Type, req.Source.Content)
	if err != nil {
		resp = &api.AddSourceResponse{
			Code:    -1,
			Message: err.Error(),
		}
		return resp, err
	}

	resp = &api.AddSourceResponse{
		Code:    0,
		Message: "success",
	}

	return resp, nil
}

// RemoveSource implements the SDASImpl interface.
func (s *SDASImpl) RemoveSource(ctx context.Context, req *api.RemoveSourceRequest) (resp *api.RemoveSourceResponse, err error) {
	handler.RemoveSourceHandler(req.Name)
	resp = &api.RemoveSourceResponse{
		Code:    0,
		Message: "success",
	}

	return resp, nil
}

// ListSources implements the SDASImpl interface.
func (s *SDASImpl) ListSources(ctx context.Context) (resp *api.ListSourcesResponse, err error) {
	sources := handler.ListSourceHandler()
	apiSources := make([]*api.Source, 0)
	for _, source := range sources {
		apiSources = append(apiSources, &api.Source{
			Name:    source.Name,
			Type:    source.Type,
			Content: source.Content,
		})
	}
	resp = &api.ListSourcesResponse{
		Code:    0,
		Message: "success",
		Sources: apiSources,
	}
	return resp, nil
}

// AddExpose implements the SDASImpl interface.
func (s *SDASImpl) AddExpose(ctx context.Context, req *api.AddExposeRequest) (resp *api.AddExposeResponse, err error) {
	err = handler.AddExposeHandler(req.Expose.Name, req.Expose.Type, req.Expose.Content, req.Expose.SourceName)
	if err != nil {
		resp = &api.AddExposeResponse{
			Code:    -1,
			Message: err.Error(),
		}
		return resp, err
	}

	resp = &api.AddExposeResponse{
		Code:    0,
		Message: "success",
	}

	return resp, nil
}

// RemoveExpose implements the SDASImpl interface.
func (s *SDASImpl) RemoveExpose(ctx context.Context, req *api.RemoveExposeRequest) (resp *api.RemoveExposeResponse, err error) {
	handler.RemoveExposeHandler(req.Name)
	resp = &api.RemoveExposeResponse{
		Code:    0,
		Message: "success",
	}

	return resp, nil
}

// ListExposes implements the SDASImpl interface.
func (s *SDASImpl) ListExposes(ctx context.Context) (resp *api.ListExposesResponse, err error) {
	exposes := handler.ListExposeHandler()
	apiExposes := make([]*api.Expose, 0)
	for _, expose := range exposes {
		apiExposes = append(apiExposes, &api.Expose{
			Name:       expose.Name,
			Type:       expose.Type,
			Content:    expose.Content,
			SourceName: expose.SourceName,
		})
	}
	resp = &api.ListExposesResponse{
		Code:    0,
		Message: "success",
		Exposes: apiExposes,
	}
	return resp, nil
}

// AddPipeline implements the SDASImpl interface.
func (s *SDASImpl) AddPipeline(ctx context.Context, req *api.AddPipelineRequest) (resp *api.AddPipelineResponse, err error) {
	// TODO: Your code here...
	return
}

// RemovePipeline implements the SDASImpl interface.
func (s *SDASImpl) RemovePipeline(ctx context.Context, req *api.RemovePipelineRequest) (resp *api.RemovePipelineResponse, err error) {
	// TODO: Your code here...
	return
}

// ListPipeline implements the SDASImpl interface.
func (s *SDASImpl) ListPipeline(ctx context.Context) (resp *api.ListPipelinesResponse, err error) {
	// TODO: Your code here...
	return
}

func (s *SDASImpl) PullExposeStream(stream api.SDAS_PullExposeStreamServer) (err error) {
	// no need to call `stream.Close()` manually
	var name string
	var msgType string
	wg := sync.WaitGroup{}
	wg.Add(2)
	// Recv
	go func() {
		defer func() {
			// if p := recover(); p != nil {
			// 	err = fmt.Errorf("panic: %v", p)
			// }
			wg.Done()
		}()
		for {

			msg, recvErr := stream.Recv()
			// make sure you receive and io.EOF or other non-nil error
			// otherwise RPCFinish event will not be recorded
			if recvErr == io.EOF {
				handler.RemoveExposeHandler(name)
				return
			} else if recvErr != nil {
				err = recvErr
				return
			}
			var pull_expose config.EXPOSE_PULL
			json.Unmarshal([]byte(msg.Expose.Content), &pull_expose)
			if pull_expose.Op == expose.OPEN {
				err := handler.AddExposeHandler(msg.Expose.Name, msg.Expose.Type, pull_expose.MsgType, msg.Expose.SourceName)
				if err != nil {
					resp := &api.PullExposeStreamResponse{
						Code:    -1,
						Message: err.Error(),
						Data:    "",
						Type:    pull_expose.MsgType,
					}
					if sendErr := stream.Send(resp); sendErr != nil {
						err = sendErr
						return
					}
				}
				name = msg.Expose.Name
				msgType = pull_expose.MsgType
			} else if pull_expose.Op == expose.CLOSE {

				if name != msg.Expose.Name {
					resp := &api.PullExposeStreamResponse{
						Code:    -1,
						Message: "name not match",
						Data:    "",
						Type:    pull_expose.MsgType,
					}
					if sendErr := stream.Send(resp); sendErr != nil {
						err = sendErr
						return
					}
				}
				handler.RemoveExposeHandler(msg.Expose.Name)
			} else if pull_expose.Op == expose.PLAY {
				err := handler.PlayExposePull(name)
				if err != nil {
					resp := &api.PullExposeStreamResponse{
						Code:    -1,
						Message: err.Error(),
						Data:    "",
						Type:    pull_expose.MsgType,
					}
					if sendErr := stream.Send(resp); sendErr != nil {
						err = sendErr
						return
					}
				}
			} else if pull_expose.Op == expose.PAUSE {
				klog.Debug("recv pause")
				err := handler.PauseExposePull(name)
				if err != nil {
					resp := &api.PullExposeStreamResponse{
						Code:    -1,
						Message: err.Error(),
						Data:    "",
						Type:    pull_expose.MsgType,
					}
					if sendErr := stream.Send(resp); sendErr != nil {
						err = sendErr
						return
					}
				}
			}
		}
	}()
	// Send
	go func() {
		defer func() {
			// if p := recover(); p != nil {
			// 	err = fmt.Errorf("panic: %v", p)
			// }
			wg.Done()
		}()
		for {
			if data, err := handler.GetExposePullData(name); err != nil {
				resp := &api.PullExposeStreamResponse{
					Code:    -1,
					Message: err.Error(),
					Data:    "",
					Type:    msgType,
				}
				if sendErr := stream.Send(resp); sendErr != nil {
					err = sendErr
					return
				}
			} else {
				resp := &api.PullExposeStreamResponse{
					Code:    0,
					Message: "data",
					Data:    data,
					Type:    msgType,
				}
				if sendErr := stream.Send(resp); sendErr != nil {
					err = sendErr
					return
				}
			}
		}
	}()
	wg.Wait()
	return
}
