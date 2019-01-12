package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chonla/oddsvr-api/httpclient"
	"github.com/labstack/echo"
)

func (h *Handler) Gateway(c echo.Context) error {
	code := c.QueryParam("code")

	token, e := h.strava.ExchangeToken(code)
	if e != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprint(e))
	}

	e = h.vr.SaveToken(token)
	if e != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprint(e))
	}

	jwtToken, e := h.jwt.Generate(token)
	if e != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprint(e))
	}

	jwtAge, _ := time.ParseDuration("168h")
	jwtCookie := httpclient.NewCookie("token", jwtToken, jwtAge, "/")
	c.SetCookie(jwtCookie)

	idAge, _ := time.ParseDuration("168h")
	idCookie := httpclient.NewCookie("me", fmt.Sprintf("%d", token.ID), idAge, "/")
	c.SetCookie(idCookie)

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/vr", h.conf.FrontBaseURL))
}
