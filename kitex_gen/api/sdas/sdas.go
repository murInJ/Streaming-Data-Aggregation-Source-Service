// Code generated by Kitex v0.9.1. DO NOT EDIT.

package sdas

import (
	api "SDAS/kitex_gen/api"
	"context"
	"errors"
	"fmt"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"AddSource": kitex.NewMethodInfo(
		addSourceHandler,
		newSDASAddSourceArgs,
		newSDASAddSourceResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"RemoveSource": kitex.NewMethodInfo(
		removeSourceHandler,
		newSDASRemoveSourceArgs,
		newSDASRemoveSourceResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"ListSources": kitex.NewMethodInfo(
		listSourcesHandler,
		newSDASListSourcesArgs,
		newSDASListSourcesResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"AddPipeline": kitex.NewMethodInfo(
		addPipelineHandler,
		newSDASAddPipelineArgs,
		newSDASAddPipelineResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"RemovePipeline": kitex.NewMethodInfo(
		removePipelineHandler,
		newSDASRemovePipelineArgs,
		newSDASRemovePipelineResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"ListPipeline": kitex.NewMethodInfo(
		listPipelineHandler,
		newSDASListPipelineArgs,
		newSDASListPipelineResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"AddExpose": kitex.NewMethodInfo(
		addExposeHandler,
		newSDASAddExposeArgs,
		newSDASAddExposeResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"RemoveExpose": kitex.NewMethodInfo(
		removeExposeHandler,
		newSDASRemoveExposeArgs,
		newSDASRemoveExposeResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"ListExposes": kitex.NewMethodInfo(
		listExposesHandler,
		newSDASListExposesArgs,
		newSDASListExposesResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"PullExposeStream": kitex.NewMethodInfo(
		pullExposeStreamHandler,
		newSDASPullExposeStreamArgs,
		newSDASPullExposeStreamResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingBidirectional),
	),
}

var (
	sDASServiceInfo                = NewServiceInfo()
	sDASServiceInfoForClient       = NewServiceInfoForClient()
	sDASServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return sDASServiceInfo
}

// for client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return sDASServiceInfoForStreamClient
}

// for stream client
func serviceInfoForClient() *kitex.ServiceInfo {
	return sDASServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(true, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "SDAS"
	handlerType := (*api.SDAS)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "api",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.9.1",
		Extra:           extra,
	}
	return svcInfo
}

func addSourceHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.SDASAddSourceArgs)
	realResult := result.(*api.SDASAddSourceResult)
	success, err := handler.(api.SDAS).AddSource(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSDASAddSourceArgs() interface{} {
	return api.NewSDASAddSourceArgs()
}

func newSDASAddSourceResult() interface{} {
	return api.NewSDASAddSourceResult()
}

func removeSourceHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.SDASRemoveSourceArgs)
	realResult := result.(*api.SDASRemoveSourceResult)
	success, err := handler.(api.SDAS).RemoveSource(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSDASRemoveSourceArgs() interface{} {
	return api.NewSDASRemoveSourceArgs()
}

func newSDASRemoveSourceResult() interface{} {
	return api.NewSDASRemoveSourceResult()
}

func listSourcesHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	_ = arg.(*api.SDASListSourcesArgs)
	realResult := result.(*api.SDASListSourcesResult)
	success, err := handler.(api.SDAS).ListSources(ctx)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSDASListSourcesArgs() interface{} {
	return api.NewSDASListSourcesArgs()
}

func newSDASListSourcesResult() interface{} {
	return api.NewSDASListSourcesResult()
}

func addPipelineHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.SDASAddPipelineArgs)
	realResult := result.(*api.SDASAddPipelineResult)
	success, err := handler.(api.SDAS).AddPipeline(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSDASAddPipelineArgs() interface{} {
	return api.NewSDASAddPipelineArgs()
}

func newSDASAddPipelineResult() interface{} {
	return api.NewSDASAddPipelineResult()
}

func removePipelineHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.SDASRemovePipelineArgs)
	realResult := result.(*api.SDASRemovePipelineResult)
	success, err := handler.(api.SDAS).RemovePipeline(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSDASRemovePipelineArgs() interface{} {
	return api.NewSDASRemovePipelineArgs()
}

func newSDASRemovePipelineResult() interface{} {
	return api.NewSDASRemovePipelineResult()
}

func listPipelineHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	_ = arg.(*api.SDASListPipelineArgs)
	realResult := result.(*api.SDASListPipelineResult)
	success, err := handler.(api.SDAS).ListPipeline(ctx)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSDASListPipelineArgs() interface{} {
	return api.NewSDASListPipelineArgs()
}

func newSDASListPipelineResult() interface{} {
	return api.NewSDASListPipelineResult()
}

func addExposeHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.SDASAddExposeArgs)
	realResult := result.(*api.SDASAddExposeResult)
	success, err := handler.(api.SDAS).AddExpose(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSDASAddExposeArgs() interface{} {
	return api.NewSDASAddExposeArgs()
}

func newSDASAddExposeResult() interface{} {
	return api.NewSDASAddExposeResult()
}

func removeExposeHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.SDASRemoveExposeArgs)
	realResult := result.(*api.SDASRemoveExposeResult)
	success, err := handler.(api.SDAS).RemoveExpose(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSDASRemoveExposeArgs() interface{} {
	return api.NewSDASRemoveExposeArgs()
}

func newSDASRemoveExposeResult() interface{} {
	return api.NewSDASRemoveExposeResult()
}

func listExposesHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	_ = arg.(*api.SDASListExposesArgs)
	realResult := result.(*api.SDASListExposesResult)
	success, err := handler.(api.SDAS).ListExposes(ctx)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSDASListExposesArgs() interface{} {
	return api.NewSDASListExposesArgs()
}

func newSDASListExposesResult() interface{} {
	return api.NewSDASListExposesResult()
}

func pullExposeStreamHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	st, ok := arg.(*streaming.Args)
	if !ok {
		return errors.New("SDAS.PullExposeStream is a thrift streaming method, please call with Kitex StreamClient")
	}
	stream := &sDASPullExposeStreamServer{st.Stream}
	return handler.(api.SDAS).PullExposeStream(stream)
}

type sDASPullExposeStreamClient struct {
	streaming.Stream
}

func (x *sDASPullExposeStreamClient) DoFinish(err error) {
	if finisher, ok := x.Stream.(streaming.WithDoFinish); ok {
		finisher.DoFinish(err)
	} else {
		panic(fmt.Sprintf("streaming.WithDoFinish is not implemented by %T", x.Stream))
	}
}
func (x *sDASPullExposeStreamClient) Send(m *api.PullExposeStreamRequest) error {
	return x.Stream.SendMsg(m)
}
func (x *sDASPullExposeStreamClient) Recv() (*api.PullExposeStreamResponse, error) {
	m := new(api.PullExposeStreamResponse)
	return m, x.Stream.RecvMsg(m)
}

type sDASPullExposeStreamServer struct {
	streaming.Stream
}

func (x *sDASPullExposeStreamServer) Send(m *api.PullExposeStreamResponse) error {
	return x.Stream.SendMsg(m)
}

func (x *sDASPullExposeStreamServer) Recv() (*api.PullExposeStreamRequest, error) {
	m := new(api.PullExposeStreamRequest)
	return m, x.Stream.RecvMsg(m)
}

func newSDASPullExposeStreamArgs() interface{} {
	return api.NewSDASPullExposeStreamArgs()
}

func newSDASPullExposeStreamResult() interface{} {
	return api.NewSDASPullExposeStreamResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) AddSource(ctx context.Context, req *api.AddSourceRequest) (r *api.AddSourceResponse, err error) {
	var _args api.SDASAddSourceArgs
	_args.Req = req
	var _result api.SDASAddSourceResult
	if err = p.c.Call(ctx, "AddSource", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) RemoveSource(ctx context.Context, req *api.RemoveSourceRequest) (r *api.RemoveSourceResponse, err error) {
	var _args api.SDASRemoveSourceArgs
	_args.Req = req
	var _result api.SDASRemoveSourceResult
	if err = p.c.Call(ctx, "RemoveSource", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ListSources(ctx context.Context) (r *api.ListSourcesResponse, err error) {
	var _args api.SDASListSourcesArgs
	var _result api.SDASListSourcesResult
	if err = p.c.Call(ctx, "ListSources", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) AddPipeline(ctx context.Context, req *api.AddPipelineRequest) (r *api.AddPipelineResponse, err error) {
	var _args api.SDASAddPipelineArgs
	_args.Req = req
	var _result api.SDASAddPipelineResult
	if err = p.c.Call(ctx, "AddPipeline", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) RemovePipeline(ctx context.Context, req *api.RemovePipelineRequest) (r *api.RemovePipelineResponse, err error) {
	var _args api.SDASRemovePipelineArgs
	_args.Req = req
	var _result api.SDASRemovePipelineResult
	if err = p.c.Call(ctx, "RemovePipeline", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ListPipeline(ctx context.Context) (r *api.ListPipelinesResponse, err error) {
	var _args api.SDASListPipelineArgs
	var _result api.SDASListPipelineResult
	if err = p.c.Call(ctx, "ListPipeline", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) AddExpose(ctx context.Context, req *api.AddExposeRequest) (r *api.AddExposeResponse, err error) {
	var _args api.SDASAddExposeArgs
	_args.Req = req
	var _result api.SDASAddExposeResult
	if err = p.c.Call(ctx, "AddExpose", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) RemoveExpose(ctx context.Context, req *api.RemoveExposeRequest) (r *api.RemoveExposeResponse, err error) {
	var _args api.SDASRemoveExposeArgs
	_args.Req = req
	var _result api.SDASRemoveExposeResult
	if err = p.c.Call(ctx, "RemoveExpose", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ListExposes(ctx context.Context) (r *api.ListExposesResponse, err error) {
	var _args api.SDASListExposesArgs
	var _result api.SDASListExposesResult
	if err = p.c.Call(ctx, "ListExposes", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) PullExposeStream(ctx context.Context) (SDAS_PullExposeStreamClient, error) {
	streamClient, ok := p.c.(client.Streaming)
	if !ok {
		return nil, fmt.Errorf("client not support streaming")
	}
	res := new(streaming.Result)
	err := streamClient.Stream(ctx, "PullExposeStream", nil, res)
	if err != nil {
		return nil, err
	}
	stream := &sDASPullExposeStreamClient{res.Stream}
	return stream, nil
}
