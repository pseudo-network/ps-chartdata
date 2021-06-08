package model

type Config struct {
	SupportedResolutions   []string `json:"supported_resolutions"`
	SupportsGroupRequest   bool     `json:"supports_group_request"`
	SupportsMarks          bool     `json:"supports_marks"`
	SupportsSearch         bool     `json:"supports_search"`
	SupportsTimescaleMarks bool     `json:"supports_timescale_marks"`
	SupportsTime           bool     `json:"supports_time"`
}
