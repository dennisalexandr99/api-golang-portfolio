package config

import "github.com/tkanos/gonfig"

type Configuration struct {
	DB_USERNAME      string
	DB_PASSWORD      string
	DB_HOST          string
	DB_PORT          string
	DB_NAME          string
	HMAC256_SECRET   string
	JWT_SECRET       string
	JWT_TIME_MINUTES int
}

func GetConfig() Configuration {
	conf := Configuration{}
	gonfig.GetConf("config/config.json", &conf)
	return conf
}
