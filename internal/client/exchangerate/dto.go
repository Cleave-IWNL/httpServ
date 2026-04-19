package exchangerate

type pairResponse struct {
	Result             string  `json:"result"`
	ErrorType          string  `json:"error-type,omitempty"`
	BaseCode           string  `json:"base_code"`
	TargetCode         string  `json:"target_code"`
	ConversionRate     float64 `json:"conversion_rate"`
	TimeLastUpdateUnix int64   `json:"time_last_update_unix"`
}

type Rate struct {
	From      string
	To        string
	Value     float64
	FetchedAt int64
}

type PaymentInCurrency struct {
	ID               string  `json:"id"`
	OriginalAmount   int     `json:"original_amount"`
	OriginalCurrency string  `json:"original_currency"`
	TargetCurrency   string  `json:"target_currency"`
	ConvertedAmount  float64 `json:"converted_amount"`
	Rate             float64 `json:"rate"`
	RateFetchedAt    int64   `json:"rate_fetched_at"`
}
