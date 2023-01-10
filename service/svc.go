package service

import (
	"context"
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauth1tw "github.com/dghubble/oauth1/twitter"
)

const WYRE_ACCOUNT_ID = "AC_NLVGLA2DCCD"

var networks = map[string]string{
	"80001": "matic",
	"137":   "matic",
}

type TwitterService struct {
	consumerKey    string
	consumerSecret string
}

func NewTwitterService(consumerKey, consumerSecret string) *TwitterService {
	return &TwitterService{
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,
	}
}

func NewTwitterServiceWithENV() *TwitterService {
	return NewTwitterService(os.Getenv("twconsumerKey"), os.Getenv("twconsumerSecret"))
}

func (svc *TwitterService) Ping(ctx context.Context) error {
	return nil
}

func (svc *TwitterService) AuthURL(ctx context.Context, request *AuthURLRequest) (string, error) {
	config := oauth1.Config{
		ConsumerKey:    svc.consumerKey,
		ConsumerSecret: svc.consumerSecret,
		CallbackURL:    request.OauthCallback,
		Endpoint:       oauth1tw.AuthorizeEndpoint,
	}

	requestToken, _, err := config.RequestToken()

	if err != nil {
		return "", err
	}

	authorizationURL, err := config.AuthorizationURL(requestToken)

	if err != nil {
		return "", err
	}

	return authorizationURL.String(), nil
}

func (svc *TwitterService) Auth(ctx context.Context, request *AuthRequest) (*AuthResponse, error) {
	config := oauth1.Config{
		ConsumerKey:    svc.consumerKey,
		ConsumerSecret: svc.consumerSecret,
		CallbackURL:    "",
		Endpoint:       oauth1tw.AuthorizeEndpoint,
	}

	// Twitter ignores the oauth_signature on the access token request. The user
	// to which the request (temporary) token corresponds is already known on the
	// server. The request for a request token earlier was validated signed by
	// the consumer. Consumer applications can avoid keeping request token state
	// between authorization granting and callback handling.
	accessToken, accessSecret, err := config.AccessToken(
		request.Token,
		"secret does not matter",
		request.Verifier,
	)
	if err != nil {
		return nil, err
	}

	token := oauth1.NewToken(accessToken, accessSecret)

	// httpClient will automatically authorize http.Request's
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	twu, _, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{
		IncludeEmail: twitter.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		MessageHash: "",
		Signature:   "",
		AuthId:      fmt.Sprintf("twitter-%s", twu.IDStr),
		Email:       twu.Email,
	}, nil
}
