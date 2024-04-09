namespace go api

struct Pipeline{
        1: string type
        2: string name
        3: string content
}

struct SetPipelineRequest {
        1: Pipeline Pipeline
}

struct SetPipelineResponse {
        1: string message
        2: i32 code
}

struct QueryPipelineResponse {
        1: string message
        2: i32 code
        3: Pipeline Pipeline
}
