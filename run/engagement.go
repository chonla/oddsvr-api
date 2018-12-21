package run

import (
	"github.com/labstack/echo"
)

type VrEngagement struct{}

type Engagement struct {
	AthleteID uint32  `json:"athlete_id" bson:"athlete_id"`
	Distance  float64 `json:"distance" bson:"distance"`
}

func (v *VirtualRun) NewEngagement() *VrEngagement {
	return &VrEngagement{}
}

func (v *VrEngagement) FromContext(c echo.Context) (*Engagement, error) {
	eng := new(Engagement)
	if err := c.Bind(eng); err != nil {
		return nil, err
	}
	return eng, nil
}
