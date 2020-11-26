package constErr

import (
	"errors"
)

//For Account
var (
	AccountExists      = errors.New("Account exists.")
	NotFoundAcc        = errors.New("Not found Account.")
	AccIsNil           = errors.New("Account is nil.")
	AccountBaseIsEmpty = errors.New("Base of accounts is empty.")
)

// For JSON
var (
	ErrorMarshal   = errors.New("Error JSON marshaling.")
	ErrorUnmarshal = errors.New("Error JSON unmarshaling.")
)

//For Ads
var (
	NotFoundAd    = errors.New("No Found ad.")
	AdBaseIsEmpty = errors.New("Ad base is empty.")
	AdIsNil       = errors.New("Ad is nil.")
)

// Other
var (
	ErrorConvertToInteger  = errors.New("Error convert to integer.")
	FailedToGenerateAToken = errors.New("Failed to generate a token.")
	InvalidID              = errors.New("Invalid ID.")
	EmptyFields            = errors.New("Empty fields.")
	TokenNotContain        = errors.New("Token not contain.")
	NoValidToken           = errors.New("No valid token.")
	YouRat                 = errors.New("...Ля ты крыса! (Comedy club).")
)
