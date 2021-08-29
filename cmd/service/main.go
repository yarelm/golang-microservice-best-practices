package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yarelm/golang-microservice-best-practices/internal/api"
)

func main() {
	flagSet := flag.NewFlagSet("main", flag.ExitOnError)
	debug := flagSet.Bool("debug", false, "sets log level to debug")
	port := flagSet.String("port", "8080", "port to listen")

	if err := viper.BindPFlags(flagSet); err != nil {
		log.Fatal().Err(err).Msg("failed to bind viper flags")
	}
	viper.AutomaticEnv()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Msg("Service is booting up...")
	defer log.Info().Msg("Service is going down...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	server := api.NewServer(fmt.Sprintf(":%v", *port))
	if err := server.Serve(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to serve API")
	}
}
