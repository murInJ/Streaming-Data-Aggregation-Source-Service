namespace go api

struct Pipeline{
        1: string type
        2: string name
        3: string content
        4: map<string, string> entry_sources
        5: map<string, string> exit_sources
}

struct AddPipelineRequest {
        1: Pipeline Pipeline
}

struct AddPipelineResponse {
        1: string message
        2: i32 code
}


struct RemovePipelineRequest{
        1: string name
}

struct RemovePipelineResponse {
        1: string message
        2: i32 code
}

struct ListPipelinesResponse {
        1: string message
        2: i32 code
        3: list<Pipeline> pipelines
}
