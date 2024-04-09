package config

type SOURCE struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

/*
**
SOURECE TYPE DEFINITION
**
*/
type SOURCE_RTSP struct {
	Url    string `json:"url"`
	Format string `json:"format"`
}

type SOURCE_SDAS struct {
	//TODO: going to be defined
}

type SOURCE_HTTP_CALL struct {
	//TODO: going to be defined
}

type SOURCE_RPC_CALL struct {
	//TODO: going to be defined
}

type SOURCE_LISTEN_SERVER struct {
	//TODO: going to be defined
}
