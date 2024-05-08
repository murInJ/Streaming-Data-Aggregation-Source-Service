namespace go api
include "IDL/source.thrift"

include "IDL/expose.thrift"

service SDAS {
    source.AddSourceResponse AddSource(1: source.AddSourceRequest req)
    source.RemoveSourceResponse RemoveSource(1: source.RemoveSourceRequest req)
    source.ListSourcesResponse ListSources()

    expose.ListExposesResponse ListExposes()

    expose.PullExposeStreamResponse PullExposeStream(1: expose.PullExposeStreamRequest req) (streaming.mode="bidirectional")
    expose.AddExposeResponse AddExpose(1: expose.AddExposeRequest req)
    expose.RemoveExposeResponse RemoveExpose(1: expose.RemoveExposeRequest req)
}
