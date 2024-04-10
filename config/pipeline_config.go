package config

type PIPELINE struct {
	Type         string            `json:"type"`
	Content      string            `json:"content"`
	Name         string            `json:"name"`
	EntrySources map[string]string `json:"entry_sources"`
	ExitSources  map[string]string `json:"exit_sources"`
}

/*
**
PIPELINE TYPE DEFINITION
**
*/

type PIPELINE_PLUGIN struct {
	Name string `json:"name"`
}

type PIPELINE_INFERENCE_SERVER struct {
	//TODO: going to be defined
}

type PIPELINE_API struct {
	//TODO: going to be defined
}
