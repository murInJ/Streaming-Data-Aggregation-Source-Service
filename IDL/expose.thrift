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

struct AddExposeRequest{
    1: Expose expose
}

struct AddExposeResponse{
    1: string message
    2: i32 code
}

struct RemoveExposeRequest{
        1: string name
}

struct RemoveExposeResponse {
        1: string message
        2: i32 code
}

struct PullExposeStreamRequest {
        1: Expose expose
}

struct PullExposeStreamResponse {
        1: string message
        2: i32 code
        3: struct.SourceMsg source_msg
}