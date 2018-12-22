package api

import (
	"fmt"

	"github.com/chonla/oddsvr-api/api/handler"
	"github.com/chonla/oddsvr-api/database"
	"github.com/chonla/oddsvr-api/jwt"
	"github.com/chonla/oddsvr-api/logger"
	"github.com/chonla/oddsvr-api/run"
	"github.com/chonla/oddsvr-api/strava"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Conf struct {
	AppVersion         string
	DBConnection       string
	ServiceAddress     string
	FrontBaseURL       string
	JWTSecret          string
	StravaClientID     string
	StravaClientSecret string
}

type API struct {
	conf *Conf
}

func NewAPI(conf *Conf) *API {
	return &API{
		conf: conf,
	}
}

func (a *API) Start() {
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true

	server.Use(middleware.CORS())

	db, e := database.NewDatabase(a.conf.DBConnection, "vr")
	if e != nil {
		logger.Error(fmt.Errorf("Unable to connect database: %v", e).Error())
		return
	}
	s := strava.NewStrava(a.conf.StravaClientID, a.conf.StravaClientSecret)
	vr := run.NewVirtualRun(db)
	j := jwt.NewJWT(a.conf.JWTSecret)
	conf := &handler.Conf{
		AppVersion:   a.conf.AppVersion,
		FrontBaseURL: a.conf.FrontBaseURL,
	}

	h := handler.NewHandler(s, vr, j, conf)

	// Public endpoints
	r := server.Group("/api")
	r.GET("/gateway", h.Gateway)
	r.GET("/vr/:id", h.Vr)
	r.GET("/version", h.Version)

	// Private endpoints
	jwtConfig := middleware.JWTConfig{
		Claims:     &jwt.Claims{},
		SigningKey: []byte(a.conf.JWTSecret),
	}
	r.Use(middleware.JWTWithConfig(jwtConfig))
	r.GET("/me", h.Me)
	r.GET("/me/vr", h.JoinedVrs)
	r.GET("/vr", h.Vrs)
	r.POST("/vr", h.CreateVr)
	r.POST("/join/:id", h.JoinVr)
	r.POST("/leave/:id", h.LeaveVr)
	r.PATCH("/vr/:id", h.UpdateVr)
	r.DELETE("/vr/:id", h.DeleteVr)

	logger.Info(fmt.Sprintf("server is listening on %s", a.conf.ServiceAddress))
	server.Logger.Fatal(server.Start(a.conf.ServiceAddress))
}
