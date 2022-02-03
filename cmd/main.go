package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/sachinagada/secretsanta"
	"github.com/sachinagada/secretsanta/send"
	"github.com/sachinagada/secretsanta/shuffle"

	"github.com/vimeo/dials"
	"github.com/vimeo/dials/env"
	"github.com/vimeo/dials/flag"
)

type Config struct {
	Port       int    `dialsdesc:"port where application should run"`
	Shuffler   string `dialsdesc:"type of shuffler to shuffle the participants (rand)"`
	MailConfig *send.Config
}

func defaultConfig() *Config {
	return &Config{
		Port:       3000,
		Shuffler:   "rand",
		MailConfig: send.DefaultConfig(),
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

	var shuf secretsanta.Shuffler
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

	srv := secretsanta.NewServer(shuf, mail)
	http.Handle("/santa", srv)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	server := http.Server{
		Addr: net.JoinHostPort("", strconv.Itoa(conf.Port)),
	}
	go func() {
		if srvErr := server.ListenAndServe(); srvErr != http.ErrServerClosed {
			log.Fatalf("unexpected error listening on server: %s", srvErr)
		}
	}()

	<-ctx.Done()
	stop()

	log.Print("initializing shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if shutErr := server.Shutdown(shutdownCtx); shutErr != nil {
		log.Fatalf("error shutting down server: %s", shutErr)
	}
}
