package main

import (
	"flag"

	"github.com/linshenqi/payground/src/services/payment"
	"github.com/linshenqi/sptty"
)

func main() {
	sptty.SetTag("payment")

	cfg := flag.String("config", "./config.yml", "--config")
	flag.Parse()

	app := sptty.GetApp()
	app.ConfFromFile(*cfg)

	services := sptty.Services{
		&payment.Service{},
	}

	configs := sptty.Configs{
		&payment.Config{},
	}

	app.AddServices(services)
	app.AddConfigs(configs)

	app.Sptting()
}
