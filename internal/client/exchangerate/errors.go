package exchangerate

import "errors"

var (
	ErrUnsupportedCurrency = errors.New("exchangerate: unsupported currency")
	ErrInvalidKey          = errors.New("exchangerate: invalid api key")
	ErrUpstream            = errors.New("exchangerate: upstream error")
)

func mapAPIError(code string) error {
	switch code {
	case "unsupported-code":
		return ErrUnsupportedCurrency
	case "invalid-key", "inactive-account":
		return ErrInvalidKey
	default:
		return ErrUpstream
	}
}
