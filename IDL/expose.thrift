namespace go api
include "struct.thrift"
struct Expose{
        1: string type
        2: string name
        3: map<string,string> content
        4: string source_name
}

struct ListExposesResponse {
        1: string message
        2: i32 code
        3: list<Expose> exposes
}

struct PullExposeStreamRequest {
        1: Expose expose
}

struct PullExposeStreamResponse {
        1: string message
        2: i32 code
        3: struct.SourceMsg source_msg
}