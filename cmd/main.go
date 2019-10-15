package main

import (
	"fmt"
	"log"
	"net/http"
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

	fmt.Println(viper.GetString("heLlo"), viper.GetString("cool"))
	// set up mgo config
	collFn, err := config.NewMgoDB()
	if err != nil {
		panic(err)
	}

	fmt.Println(viper.GetString("REDIS.ADDR"))

	redisClient := config.NewRedisClient()

	pollRepo := pollrepo.NewRepository(collFn)
	pollSvc := pollsvc.NewService(pollRepo)
	pollSet := pollendpoint.NewSet(pollSvc)

	voteRepo := voterepo.NewRepository(collFn)
	voteSvc := votesvc.NewService(voteRepo, pollRepo, redisClient)
	voteSet := voteendpoint.NewSet(voteSvc)

	r := pkg.NewRouter(pollSet, voteSet)
	port := viper.GetString("PORT")

	s := http.Server{
		Handler:      r,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		Addr:         port,
	}

	log.Println("server running on", port)
	log.Fatal(s.ListenAndServe(), "Server Stopped")

}
