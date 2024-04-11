namespace go api

struct Expose{
        1: string type
        2: string name
        3: string content
        4: string source_name
}

struct AddExposeRequest {
        1: Expose expose
}

struct AddExposeResponse {
        1: string message
        2: i32 code
}

struct RemoveExposeRequest {
        1: string name
}

struct RemoveExposeResponse {
        1: string message
        2: i32 code
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
        3: string data
        4: string type
}