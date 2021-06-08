package model

type Symbol struct {
	Name                 string   `json:"name"`
	ExchangeTraded       string   `json:"exchange-traded"`
	ExchangeListed       string   `json:"exchange-listed"`
	Timezone             string   `json:"timezone"`
	PriceScale           int      `json:"pricescale"`
	Minmov               int      `json:"minmov"`
	Minmov2              int      `json:"minmov2"`
	PointValue           int      `json:"pointvalue"`
	Session              string   `json:"session"`
	HasIntraday          bool     `json:"has_intraday"`
	IntradayMultipliers  []string `json:"intraday_multipliers"`
	HasNoVolume          bool     `json:"has_no_volume"`
	Description          string   `json:"description"`
	Type                 string   `json:"type"`
	SupportedResolutions []string `json:"supported_resolutions"`
}
