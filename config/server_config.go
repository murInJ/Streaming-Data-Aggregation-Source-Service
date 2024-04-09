package config

type SERVER struct {
	Node string `json:"node"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type REGISTRY struct {
	Able        bool   `json:"able"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	User        string `json:"user"`
	Passwd      string `json:"passwd"`
	ServiceName string `json:"service_name"`
}

type MANAGE_SERVER struct {
	Url       string `json:"url"`
	SecretKey string `json:"secret_key"`
}
