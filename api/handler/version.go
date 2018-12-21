package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func (h *Handler) Version(c echo.Context) error {
	version := map[string]string{
		"version": h.conf.AppVersion,
	}
	return c.JSON(http.StatusOK, version)
}
