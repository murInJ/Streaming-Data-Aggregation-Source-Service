package config

type PIPELINE struct {
	Type       string            `json:"type"`
	Pipeline   any               `json:"pipeline"`
	Name       string            `json:"name"`
	SourceName map[string]string `json:"source_name"`
}

/*
**
PIPELINE TYPE DEFINITION
**
*/

type PIPELINE_FUNCTIONAL struct {
	Name string `json:"name"`
}

type PIPELINE_INFERENCE_SERVER struct {
	//TODO: going to be defined
}

type PIPELINE_API struct {
	//TODO: going to be defined
}
