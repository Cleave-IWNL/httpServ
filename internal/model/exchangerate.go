package model

type Rate struct {
	From      string
	To        string
	Value     float64
	FetchedAt int64
}

type ExchangeRatePair struct {
	Result             string  `json:"result"`
	ErrorType          string  `json:"error-type,omitempty"`
	BaseCode           string  `json:"base_code"`
	TargetCode         string  `json:"target_code"`
	ConversionRate     float64 `json:"conversion_rate"`
	TimeLastUpdateUnix int64   `json:"time_last_update_unix"`
}
