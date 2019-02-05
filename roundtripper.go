package roundtripper

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type roundtripper struct {
	sync.Mutex
	token                            string
	cognitoUsername, cognitoPassword string
	cognitoID, cognitoUserPool       string
	transport                        http.RoundTripper
}

// New returns RoundTripper implementation with JWT handling using AWS cognito
func New(rt http.RoundTripper, opts ...OptionRoundtripperFunc) (http.RoundTripper, error) {
	if rt == nil {
		rt = http.DefaultTransport
	}

	t := &roundtripper{
		transport: rt,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if err := opt(t); err != nil {
			return nil, err
		}
	}

	if err := t.GenerateToken(); err != nil {
		return nil, err
	}

	return t, nil
}

// generateToken gets a new JWT from Cognito and updates the transport with it
func (t *roundtripper) GenerateToken() error {
	t.Lock()
	defer t.Unlock()
	svc := cognitoidentityprovider.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})

	params := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String("ADMIN_NO_SRP_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(t.cognitoUsername),
			"PASSWORD": aws.String(t.cognitoPassword),
		},
		ClientId:   aws.String(t.cognitoID),
		UserPoolId: aws.String(t.cognitoUserPool),
	}

	resp, err := svc.AdminInitiateAuth(params)
	if err != nil {
		return fmt.Errorf("error init'ing auth with cognito: %s", err.Error())
	}

	t.token = *resp.AuthenticationResult.AccessToken

	return nil
}

// CurrentToken returns the current token from the roundtripper implementation
func (t *roundtripper) CurrentToken() string {
	t.Lock()
	defer t.Unlock()
	return t.token
}

// RoundTrip is an implementation of RoundTripper interface which sets the JWT
// authorization header in the request client
func (t *roundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	addAuthToken := func() {
		if token := t.CurrentToken(); token != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		}
	}

	addAuthToken()
	resp, err := t.transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	// if the response is a 401, the token has likely expired and we need
	// to regenerate it without erroring out
	if resp.StatusCode == http.StatusUnauthorized {
		if err := t.GenerateToken(); err != nil {
			return resp, err
		}

		addAuthToken()
		resp, err = t.transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
