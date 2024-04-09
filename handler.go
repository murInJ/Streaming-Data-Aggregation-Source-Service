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

// SetPipeline implements the SDASImpl interface.
func (s *SDASImpl) SetPipeline(ctx context.Context, req *api.SetPipelineRequest) (resp *api.SetPipelineResponse, err error) {
	// TODO: Your code here...
	return
}

// QueryPipeline implements the SDASImpl interface.
func (s *SDASImpl) QueryPipeline(ctx context.Context) (resp *api.QueryPipelineResponse, err error) {
	// TODO: Your code here...
	return
}

// AddExpose implements the SDASImpl interface.
func (s *SDASImpl) AddExpose(ctx context.Context, req *api.AddExposeRequest) (resp *api.AddExposeResponse, err error) {
	// TODO: Your code here...
	return
}

// RemoveExpose implements the SDASImpl interface.
func (s *SDASImpl) RemoveExpose(ctx context.Context, req *api.RemoveExposeRequest) (resp *api.RemoveExposeResponse, err error) {
	// TODO: Your code here...
	return
}

// ListExposes implements the SDASImpl interface.
func (s *SDASImpl) ListExposes(ctx context.Context) (resp *api.ListExposesResponse, err error) {
	// TODO: Your code here...
	return
}
