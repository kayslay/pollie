package main

import (
	_log "log"
	"pollie"

	"github.com/go-kit/kit/log"

	"net/http"
	"os"
	"pollie/config"
	"pollie/pkg"
	"pollie/pkg/poll/pollendpoint"
	"pollie/pkg/poll/pollrepo"
	"pollie/pkg/poll/pollsvc"
	"pollie/pkg/vote/voteendpoint"
	"pollie/pkg/vote/voterepo"
	"pollie/pkg/vote/votesvc"
	"time"

	"github.com/spf13/viper"
)

func main() {
	// set viper config
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetDefault("port", ":6000")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// set up mgo config
	collFn, err := config.NewMgoDB()
	if err != nil {
		panic(err)
	}

	// init mgo IP interfacer
	pollie.InitMgoIPer(collFn)

	redisClient := config.NewRedisClient()

	logger := log.NewJSONLogger(os.Stderr)

	pollRepo := pollrepo.NewRepository(collFn)
	pollSvc := pollsvc.NewService(pollRepo)
	pollSet := pollendpoint.NewSet(pollSvc, logger)

	voteRepo := voterepo.NewRepository(collFn)
	voteSvc := votesvc.NewService(voteRepo, pollRepo, redisClient)
	voteSet := voteendpoint.NewSet(voteSvc, logger)

	r := pkg.NewRouter(pollSet, voteSet)
	port := viper.GetString("PORT")

	s := http.Server{
		Handler:      r,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		Addr:         port,
	}

	_log.Println("server running on", port)
	_log.Fatal(s.ListenAndServe(), "Server Stopped")

}
