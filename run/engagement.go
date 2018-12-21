package run

import (
	"github.com/labstack/echo"
)

type VrEngagement struct{}

type Engagement struct {
	AthleteID     uint32  `json:"athlete_id" bson:"athlete_id"`
	AthleteName   string  `json:"athlete_name" bson:"athlete_name"`
	Distance      float64 `json:"distance" bson:"distance"`
	TakenDistance float64 `json:"taken_distance" bson:"taken_distance"`
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
