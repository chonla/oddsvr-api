package strava

import (
	"fmt"

	"github.com/chonla/oddsvr-api/httpclient"
	"github.com/chonla/oddsvr-api/run"
)

const (
	APIOAUTH = "https://www.strava.com/oauth/token"
	APIBASE  = "https://www.strava.com/api/v3"
)

type Strava struct {
	clientID     string
	clientSecret string
}

type TokenExchange struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

func NewStrava(clientID, clientSecret string) *Strava {
	return &Strava{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (s *Strava) ExchangeToken(code string) (*run.Token, error) {
	tokenEx := TokenExchange{
		ClientID:     s.clientID,
		ClientSecret: s.clientSecret,
		Code:         code,
	}
	token := &run.Token{}

	c := httpclient.NewClient()
	e := c.Post(APIOAUTH, tokenEx, token)
	if e != nil {
		return nil, e
	}

	return token, nil
}

func (s *Strava) Athlete(token string) (*run.Athlete, error) {
	c := httpclient.NewClientWithToken(token)
	me := &run.Athlete{}

	e := c.Get(fmt.Sprintf("%s/athlete", APIBASE), me)
	if e != nil {
		return nil, e
	}

	return me, nil
}
