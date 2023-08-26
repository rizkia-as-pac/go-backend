package curr

// currency yang ingin kita support
const (
	USD = "USD"
	RUB = "RUB"
	JPY = "JPY"
)

// IsSupportedCurrency mengembalikan true jika currency disupport
func IsSupportedCurrency(currency string) bool {
	switch currency {
	// jika currency nilainya salah satu dari yang dibawah ini maka kembalikan true
	case USD, RUB, JPY:
		return true
	}
	return false
}
