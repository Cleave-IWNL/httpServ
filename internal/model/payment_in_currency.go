package model

type PaymentInCurrency struct {
	ID               string  `json:"id"`
	OriginalAmount   int     `json:"original_amount"`
	OriginalCurrency string  `json:"original_currency"`
	TargetCurrency   string  `json:"target_currency"`
	ConvertedAmount  float64 `json:"converted_amount"`
	Rate             float64 `json:"rate"`
	RateFetchedAt    int64   `json:"rate_fetched_at"`
}
