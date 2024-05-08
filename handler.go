package main

import (
	config "SDAS/config"
	api "SDAS/kitex_gen/api"
	expose "SDAS/services/expose_service"
	source "SDAS/services/source_service"
	"context"
	"io"
	"strconv"
	"sync"
)

// SDASImpl implements the last service interface defined in the IDL.
type SDASImpl struct{}

// AddSource implements the SDASImpl interface.
func (s *SDASImpl) AddSource(ctx context.Context, req *api.AddSourceRequest) (resp *api.AddSourceResponse, err error) {
	err = source.AddSource(req.Source.Name, req.Source.Type, req.Source.Expose, req.Source.Content)
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
	source.RemoveSource(req.Name)
	resp = &api.RemoveSourceResponse{
		Code:    0,
		Message: "success",
	}
	return resp, nil
}

// ListSources implements the SDASImpl interface.
func (s *SDASImpl) ListSources(ctx context.Context) (resp *api.ListSourcesResponse, err error) {
	sources := source.ListSources()
	resp = &api.ListSourcesResponse{
		Code:    0,
		Message: "success",
		Sources: sources,
	}
	return resp, nil
}

// ListExposes implements the SDASImpl interface.
func (s *SDASImpl) ListExposes(ctx context.Context) (resp *api.ListExposesResponse, err error) {
	exposes := expose.ListExposes()
	resp = &api.ListExposesResponse{
		Code:    0,
		Message: "success",
		Exposes: exposes,
	}
	return resp, nil
}

func (s *SDASImpl) PullExposeStream(stream api.SDAS_PullExposeStreamServer) (err error) {
	// no need to call `stream.Close()` manually
	var entity *expose.ExposeEntityPull
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer func() {
			wg.Done()
		}()

		flag := false
		for {
			msg, recvErr := stream.Recv()
			// make sure you receive and io.EOF or other non-nil error
			// otherwise RPCFinish event will not be recorded
			if recvErr == io.EOF {
				if flag {
					entity.Stop()
				}
				return
			} else if recvErr != nil {
				if flag {
					entity.Stop()
					expose.Exposes.Delete(msg.Expose.Name)
				}
				err = recvErr
				return
			}
			op, err := strconv.Atoi(msg.Expose.Content["op"])
			if err != nil {
				resp := &api.PullExposeStreamResponse{
					Code:      -1,
					Message:   err.Error(),
					SourceMsg: nil,
				}
				if sendErr := stream.Send(resp); sendErr != nil {
					err = sendErr
					return
				}
			}
			if op == config.OPEN {
				e, err := expose.BuildExposeEntity(msg.Expose.Name, msg.Expose.Type, msg.Expose.SourceName, msg.Expose.Content)
				entity = e.(*expose.ExposeEntityPull)
				entity.Wg = &wg
				entity.Stream = stream
				err = entity.Start()
				if err != nil {
					resp := &api.PullExposeStreamResponse{
						Code:      -1,
						Message:   err.Error(),
						SourceMsg: nil,
					}
					if sendErr := stream.Send(resp); sendErr != nil {
						err = sendErr
						return
					}
				}
				expose.Exposes.Store(msg.Expose.Name, entity)
				flag = true
			} else if op == config.CLOSE {
				entity.Stop()
				expose.Exposes.Delete(msg.Expose.Name)
				flag = false
				return
			} else if flag && op == config.PLAY {
				*entity.ControlChannel <- config.PLAY
			} else if flag && op == config.PAUSE {
				*entity.ControlChannel <- config.PAUSE
			} else if flag && (entity == nil || entity.Status == config.CLOSE || entity.Status == config.ERR) {
				return
			}
		}
	}()
	wg.Wait()
	return
}

// AddExpose implements the SDASImpl interface.
func (s *SDASImpl) AddExpose(ctx context.Context, req *api.AddExposeRequest) (resp *api.AddExposeResponse, err error) {
	err = expose.AddExpose(req.Expose.Name, req.Expose.Type, req.Expose.SourceName, req.Expose.Content)
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
	expose.RemoveExpose(req.Name)
	resp = &api.RemoveExposeResponse{
		Code:    0,
		Message: "success",
	}
	return resp, nil
}
