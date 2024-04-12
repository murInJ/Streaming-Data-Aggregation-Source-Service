namespace go api

struct Source{
        1: string type
        2: string name
        3: map<string,string> content
        4: bool expose
}

struct AddSourceRequest {
        1: Source source
}

struct AddSourceResponse {
        1: string message
        2: i32 code
}

struct RemoveSourceRequest{
        1: string name
}

struct RemoveSourceResponse {
        1: string message
        2: i32 code
}

struct ListSourcesResponse {
        1: string message
        2: i32 code
        3: list<Source> sources
}

