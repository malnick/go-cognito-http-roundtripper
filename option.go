package roundtripper

import (
	"errors"
)

var (
	// ErrInvalidCredentials is the error returned by NewRoundTripper if user used empty string for a credentials.
	ErrInvalidCredentials = errors.New("uid, secret and loginEndpoit cannot be empty")

	// ErrInvalidUserAgent is the error returned by NewRoundTripper if user used empty string for a user agent.
	ErrInvalidUserAgent = errors.New("userAgent cannot be empty")

	// ErrInvalidExpireDuration is the error returned by NewRoundTripper if the token expire duration is negative or
	// zero value.
	ErrInvalidExpireDuration = errors.New("token expire duration must be positive non zero value")
)

// OptionRoundtripperFunc type sets options configuration
type OptionRoundtripperFunc func(*roundtripper) error

func errorOnEmpty(arg string) error {
	if len(arg) == 0 {
		return errors.New("Must pass non-empty string to this option")
	}
	return nil
}

// OptionCognitoUsername sets the cognito username
func OptionCognitoUsername(username string) OptionRoundtripperFunc {
	return func(o *roundtripper) error {
		if err := errorOnEmpty(username); err != nil {
			return err
		}
		o.cognitoUsername = username

		return nil
	}
}

// OptionCognitoPassword sets the cognito password
func OptionCognitoPassword(password string) OptionRoundtripperFunc {
	return func(o *roundtripper) error {
		if err := errorOnEmpty(password); err != nil {
			return err
		}
		o.cognitoPassword = password

		return nil
	}
}

// OptionCognitoClientID sets the cognito password
func OptionCognitoClientID(id string) OptionRoundtripperFunc {
	return func(o *roundtripper) error {
		if err := errorOnEmpty(id); err != nil {
			return err
		}
		o.cognitoID = id

		return nil
	}
}

// OptionCognitoUserPool sets the cognito password
func OptionCognitoUserPool(id string) OptionRoundtripperFunc {
	return func(o *roundtripper) error {
		if err := errorOnEmpty(id); err != nil {
			return err
		}
		o.cognitoUserPool = id

		return nil
	}
}
