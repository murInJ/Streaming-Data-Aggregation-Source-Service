package config

type EXPOSE struct {
	Type       string `json:"type"`
	Content    string `json:"content"`
	Name       string `json:"name"`
	SourceName string `json:"source_name"`
}

/*
**
EXPOSE TYPE DEFINITION
**
*/

type EXPOSE_RTSP struct {
	RTSPAddress       string `json:"rtsp_address"`
	UDPRTPAddress     string `json:"udp_rtp_address"`
	UDPRTCPAddress    string `json:"udp_rtcp_address"`
	MulticastIPRange  string `json:"multicast_ip_range"`
	MulticastRTPPort  int    `json:"multicast_rtp_port"`
	MulticastRTCPPort int    `json:"multicast_rtcp_port"`
	Format            string `json:"format"`
}

type EXPOSE_HTTP struct {
	//TODO: going to be defined
}

type EXPOSE_RPC struct {
	//TODO: going to be defined
}

type EXPOSE_PUSH struct {
	//TODO: going to be defined
}
