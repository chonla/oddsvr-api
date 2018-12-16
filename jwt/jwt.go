package jwt

import (
	"time"

	"github.com/chonla/oddsvr-api/run"
	jwt "github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID          uint32 `json:"id"`
	StravaToken string `json:"strava_token"`
	jwt.StandardClaims
}

type JWT struct {
	secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		secret: secret,
	}
}

func (j *JWT) Generate(token *run.Token) (string, error) {
	claims := &Claims{
		ID:          token.ID,
		StravaToken: token.AccessToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString([]byte(j.secret))
}
