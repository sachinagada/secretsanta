package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sachinagada/secretsanta/monitor"
	"github.com/sachinagada/secretsanta/pick"
	"github.com/sachinagada/secretsanta/send"
	"github.com/sachinagada/secretsanta/shuffle"

	"github.com/vimeo/dials"
	"github.com/vimeo/dials/env"
	"github.com/vimeo/dials/flag"
)

type Config struct {
	Port                string `dialsdesc:"port where application should run"`
	Shuffler            string `dialsdesc:"type of shuffler to shuffle the participants (rand)"`
	MailConfig          *send.Config
	MetricsPort         string        `dialsdesc:"port for the metrics server"`
	ShutdownGracePeriod time.Duration `dialsdesc:"amount of time to give servers to gracefully shutdown"`
}

func defaultConfig() *Config {
	return &Config{
		Port:                "3000",
		MetricsPort:         "8080",
		Shuffler:            "rand",
		MailConfig:          send.DefaultConfig(),
		ShutdownGracePeriod: 30 * time.Second,
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config := defaultConfig()

	flagSet, flagErr := flag.NewCmdLineSet(flag.DefaultFlagNameConfig(), config)
	if flagErr != nil {
		log.Fatalf("error configuring command line flag: %s", flagErr)
	}

	d, confErr := dials.Config(context.Background(), config, &env.Source{}, flagSet)
	if confErr != nil {
		log.Fatalf("error getting configuration: %s", confErr)
	}

	conf := d.View().(*Config)

	var shuf pick.Shuffler
	switch conf.Shuffler {
	case "rand":
		shuf = &shuffle.Rand{}
	default:
		log.Fatalf("unsupported type of shuffler %q", conf.Shuffler)
	}

	mail, mailErr := send.NewMail(conf.MailConfig)
	if mailErr != nil {
		log.Fatalf("error initializing mail: %s", mailErr)
	}

	monSvr, monErr := monitor.NewServer(conf.MetricsPort)
	if monErr != nil {
		log.Fatalf("error initializing monitoring server: %s", monErr)
	}

	go func() {
		if monSrvErr := monSvr.ListenAndServe(); !errors.Is(monSrvErr, http.ErrServerClosed) {
			log.Fatalf("unexpected error listening on server: %s", monSrvErr)
		}
	}()

	server := pick.NewServer(shuf, mail, conf.Port)

	go func() {
		if srvErr := server.ListenAndServe(); !errors.Is(srvErr, http.ErrServerClosed) {
			log.Fatalf("unexpected error listening on server: %s", srvErr)
		}
	}()

	<-ctx.Done()
	stop()

	log.Print("initializing shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), conf.ShutdownGracePeriod)
	defer cancel()

	if shutErr := server.Shutdown(shutdownCtx); shutErr != nil {
		log.Printf("error shutting down server: %s", shutErr)
	}

	// shutdown the metrics server last so it can be scraped one more time
	// by prometheus
	monitorSDCtx, monCancel := context.WithTimeout(context.Background(), conf.ShutdownGracePeriod)
	defer monCancel()

	if monShutErr := monSvr.Shutdown(monitorSDCtx); monShutErr != nil {
		log.Printf("error shutting down monitoring server: %s", monShutErr)
	}
}
