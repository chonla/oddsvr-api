package run

import (
	"github.com/chonla/oddsvr-api/database"
	"github.com/chonla/rnd"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

type VirtualRun struct {
	db *database.Database
}

type Vr struct {
	ID              bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedBy       uint32        `json:"created_by" bson:"created_by"`
	CreatedDateTime string        `json:"created_datetime" bson:"created_datetime"`
	Title           string        `json:"title" bson:"title"`
	Detail          string        `json:"detail" bson:"detail"`
	Period          []string      `json:"period" bson:"period"`
	Link            string        `json:"link" bson:"link"`
	Engagements     []Engagement  `json:"engagements" bson:"engagements"`
}

func NewVirtualRun(db *database.Database) *VirtualRun {
	return &VirtualRun{
		db: db,
	}
}

func (v *VirtualRun) FromContext(c echo.Context) (*Vr, error) {
	vr := new(Vr)
	if err := c.Bind(vr); err != nil {
		return nil, err
	}
	return vr, nil
}

func (v *VirtualRun) AllAthleteCredentials() []AthleteCredential {
	creds := []AthleteCredential{}
	tokens := []InvertedToken{}
	e := v.db.List("athlete", bson.M{}, []string{"_id"}, &tokens)
	if e == nil {
		for _, t := range tokens {
			creds = append(creds, AthleteCredential{
				ID:           t.ID,
				AccessToken:  t.AccessToken,
				RefreshToken: t.RefreshToken,
				Expiry:       t.Expiry,
			})
		}
	}
	return creds
}

func (v *VirtualRun) StampLastSync(id uint32, stamp int64) {
	v.db.Upsert("sync", bson.M{
		"_id": id,
	}, bson.M{
		"_id":   id,
		"stamp": stamp,
	})
}

func (v *VirtualRun) GetLastSync(id uint32) int64 {
	output := map[string]interface{}{}
	e := v.db.Get("sync", bson.M{
		"_id": id,
	}, output)
	if e != nil {
		return 0
	}
	return output["stamp"].(int64)
}

func (v *VirtualRun) UpdateToken(token AthleteCredential) error {
	invToken := InvertedToken{}
	e := v.db.Get("athlete", bson.M{
		"_id": token.ID,
	}, &invToken)

	if e != nil {
		return e
	}

	invToken.Token.Expiry = token.Expiry
	invToken.Token.AccessToken = token.AccessToken
	invToken.Token.RefreshToken = token.RefreshToken

	return v.db.Replace("athlete", bson.M{
		"_id": token.ID,
	}, invToken)
}

func (v *VirtualRun) SaveToken(token *Token) error {
	invToken := InvertedToken{
		ID:    token.ID,
		Token: token,
	}

	return v.db.Upsert("athlete", bson.M{
		"_id": token.ID,
	}, invToken)
}

func (v *VirtualRun) SaveVr(vr *Vr) error {
	return v.db.Insert("virtualrun", vr)
}

func (v *VirtualRun) Exists(link string) bool {
	return v.db.Has("virtualrun", bson.M{
		"link": link,
	})
}

func (v *VirtualRun) FromLink(link string) (Vr, error) {
	vr := Vr{}
	e := v.db.Get("virtualrun", bson.M{
		"link": link,
	}, &vr)
	return vr, e
}

func (v *VirtualRun) UnexpiredRuns() ([]Vr, error) {
	vrs := []Vr{}
	e := v.db.List("virtualrun", bson.M{}, []string{"-startdate"}, &vrs)
	return vrs, e
}

func (v *VirtualRun) Joined(id uint32) ([]Vr, error) {
	vrs := []Vr{}
	e := v.db.List("virtualrun", bson.M{
		"engagements": bson.M{
			"$elemMatch": bson.M{
				"athlete_id": id,
			},
		},
	}, []string{"-startdate"}, &vrs)
	return vrs, e
}

func (v *VirtualRun) CreateSafeVrLink() string {
	link := rnd.Alphanum(12)
	return link
}

func (v *VirtualRun) HasJoined(id string, athleteID uint32) bool {
	return v.db.Has("virtualrun", bson.M{
		"link": id,
		"engagements": bson.M{
			"$elemMatch": bson.M{
				"athlete_id": athleteID,
			},
		},
	})
}

func (v *VirtualRun) Join(id string, eng *Engagement) error {
	if v.HasJoined(id, eng.AthleteID) {
		return nil
	}

	vr, e := v.FromLink(id)
	if e != nil {
		return e
	}

	activities, e := v.AthleteActivities(eng.AthleteID, vr.Period)
	if e != nil {
		return e
	}

	taken_distance := 0.0
	for _, activity := range activities {
		taken_distance += activity.Distance
	}

	eng.TakenDistance = taken_distance

	return v.db.Push("virtualrun", bson.M{
		"link": id,
	}, bson.M{
		"engagements": eng,
	})
}
