namespace go api
include "IDL/source.thrift"
include "IDL/pipeline.thrift"
include "IDL/expose.thrift"

service SDAS {
    source.AddSourceResponse AddSource(1: source.AddSourceRequest req)
    source.RemoveSourceResponse RemoveSource(1: source.RemoveSourceRequest req)
    source.ListSourcesResponse ListSources()

    pipeline.SetPipelineResponse SetPipeline(1: pipeline.SetPipelineRequest req)
    pipeline.QueryPipelineResponse QueryPipeline()

    expose.AddExposeResponse AddExpose(1: expose.AddExposeRequest req)
    expose.RemoveExposeResponse RemoveExpose(1: expose.RemoveExposeRequest req)
    expose.ListExposesResponse ListExposes()
}
