package handler

import (
	"fmt"
	"net/http"
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
