package main

import (
	"fmt"
	"os"

	"github.com/chonla/oddsvr-api/api"
	"github.com/chonla/oddsvr-api/logger"
)

var AppVersion string

func main() {

	dbServer := env("ODDSVR_DB", "127.0.0.1:27017", "Database address", false)
	boundAddress := env("ODDSVR_ADDR", ":1323", "Application address", false)
	frontBase := env("ODDSVR_FRONT_BASEURL", "http://localhost", "Front end web address", false)
	stravaClientID := env("ODDSVR_STRAVA_CLIENT_ID", "", "Strava client id", true)
	stravaClientSecret := env("ODDSVR_STRAVA_CLIENT_SECRET", "", "Strava client secret", true)
	jwtSecret := env("ODDSVR_JWT_SECRET", "", "JWT secret", true)

	conf := &api.Conf{
		AppVersion:         AppVersion,
		DBConnection:       dbServer,
		ServiceAddress:     boundAddress,
		FrontBaseURL:       frontBase,
		JWTSecret:          jwtSecret,
		StravaClientID:     stravaClientID,
		StravaClientSecret: stravaClientSecret,
	}

	api := api.NewAPI(conf)

	api.Start()
}

func env(key, defaultValue, name string, errorIfMissing bool) string {
	value, found := os.LookupEnv(key)
	if !found || value == "" {
		if errorIfMissing {
			logger.Error(fmt.Errorf("It seems like %s (%s) is missing from environment variables", name, key).Error())
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("%s is set to default value", name))
		logger.Info(fmt.Sprintf("You can override this using %s environment variable", key))
		value = defaultValue
	}
	return value
}
