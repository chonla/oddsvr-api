package handler

import (
	"github.com/chonla/oddsvr-api/jwt"
	"github.com/chonla/oddsvr-api/run"
	"github.com/chonla/oddsvr-api/strava"
)

type Handler struct {
	strava *strava.Strava
	vr     *run.VirtualRun
	jwt    *jwt.JWT
	conf   *Conf
}

type Conf struct {
	FrontBaseURL string
}

func NewHandler(s *strava.Strava, vr *run.VirtualRun, jwt *jwt.JWT, conf *Conf) *Handler {
	return &Handler{
		strava: s,
		vr:     vr,
		jwt:    jwt,
		conf:   conf,
	}
}
