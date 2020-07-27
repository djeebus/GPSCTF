package cmd

import (
	"fmt"
	"github.com/djeebus/gpsctf/Backend/api"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

func runWebServer() {
	r := api.NewRouter()
	listen := fmt.Sprintf("%s:%d", viper.Get("Host"), viper.Get("Port"))
	fmt.Println("Listening on", listen)

	srv := &http.Server{
		Addr:         listen,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 15,
		Handler:      r,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}
}
