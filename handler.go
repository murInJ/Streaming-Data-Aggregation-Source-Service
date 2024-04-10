package main

import (
	handler "SDAS/handlers"
	api "SDAS/kitex_gen/api"
	"context"
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
