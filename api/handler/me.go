package handler

import (
	"net/http"

	"github.com/chonla/oddsvr-api/jwt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func (h *Handler) Me(c echo.Context) error {
	user := c.Get("user").(*jwtgo.Token)
	claims := user.Claims.(*jwt.Claims)
	token := claims.StravaToken

	me, e := h.strava.Athlete(token)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}

	stats, e := h.vr.AthleteSummary(me.ID)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}

	me.Stats = stats

	return c.JSON(http.StatusOK, me)
}
