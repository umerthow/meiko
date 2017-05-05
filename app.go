package main

import (
	"flag"
	"runtime"

	"log"

	"github.com/melodiez14/meiko/src/cron"
	"github.com/melodiez14/meiko/src/email"
	"github.com/melodiez14/meiko/src/util/conn"
	"github.com/melodiez14/meiko/src/util/env"
	"github.com/melodiez14/meiko/src/util/jsonconfig"
	"github.com/melodiez14/meiko/src/webserver"
)

type configuration struct {
	Database  conn.DatabaseConfig `json:"database"`
	Redis     conn.RedisConfig    `json:"redis"`
	Webserver webserver.Config    `json:"webserver"`
	Email     email.Config        `json:"email"`
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	flag.Parse()

	// load configuration
	cfgenv := env.Get()
	config := &configuration{}
	isLoaded := jsonconfig.Load(&config, "/etc/meiko", cfgenv) || jsonconfig.Load(&config, "./etc/meiko", cfgenv)
	if !isLoaded {
		log.Fatal("Failed to load configuration")
	}

	// initiate instance
	conn.InitDB(config.Database)
	conn.InitRedis(config.Redis)
	cron.Init()
	email.Init(config.Email)
	webserver.Start(config.Webserver)
}