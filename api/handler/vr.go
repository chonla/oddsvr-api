package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chonla/oddsvr-api/jwt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

func (h *Handler) Vr(c echo.Context) error {
	id := c.Param("id")
	if h.vr.Exists(id) {
		vr, e := h.vr.FromLink(id)
		if e != nil {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.JSON(http.StatusOK, vr)
	}
	return c.NoContent(http.StatusNotFound)
}

func (h *Handler) Vrs(c echo.Context) error {
	vrs, e := h.vr.UnexpiredRuns()
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}
	return c.JSON(http.StatusOK, vrs)
}

func (h *Handler) CreateVr(c echo.Context) error {
	user := c.Get("user").(*jwtgo.Token)
	claims := user.Claims.(*jwt.Claims)
	uid := claims.ID

	vrContext, e := h.vr.FromContext(c)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}

	vrContext.ID = bson.NewObjectId()
	vrContext.Link = h.vr.CreateSafeVrLink()
	vrContext.CreatedBy = uid
	vrContext.CreatedDateTime = time.Now().Format(time.RFC3339)

	e = h.vr.SaveVr(vrContext)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}
	c.Response().Header().Add("Location", fmt.Sprintf("/vr/%s", vrContext.Link))
	c.Response().Header().Add("X-New-Vr-ID", vrContext.ID.String())

	return c.JSON(http.StatusCreated, vrContext)
}

func (h *Handler) JoinedVrs(c echo.Context) error {
	user := c.Get("user").(*jwtgo.Token)
	claims := user.Claims.(*jwt.Claims)
	uid := claims.ID

	vrs, e := h.vr.Joined(uid)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}
	return c.JSON(http.StatusOK, vrs)
}

func (h *Handler) JoinVr(c echo.Context) error {
	user := c.Get("user").(*jwtgo.Token)
	claims := user.Claims.(*jwt.Claims)
	uid := claims.ID

	vre := h.vr.NewEngagement()
	eng, e := vre.FromContext(c)
	if e != nil {
		return c.JSON(http.StatusInternalServerError, e)
	}

	id := c.Param("id")
	if h.vr.Exists(id) {

		me, e := h.vr.Profile(uid)
		if e != nil {
			return c.JSON(http.StatusInternalServerError, e)
		}

		eng.AthleteID = uid
		eng.AthleteName = strings.TrimSpace(fmt.Sprintf("%s %s", me.FirstName, me.LastName))

		e = h.vr.Join(id, eng)
		if e != nil {
			return c.JSON(http.StatusInternalServerError, e)
		}
		return c.NoContent(http.StatusCreated)
	}
	return c.NoContent(http.StatusNotFound)
}
