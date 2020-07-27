package cmd

import (
	"context"
	"fmt"
	"github.com/djeebus/gpsctf/Backend/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gpsctf",
	Short: "GPSCTF",

	PreRun: func(cmd *cobra.Command, args []string) {
		db.OpenDatabase(viper.Get("DatabasePath"))
	},

	Run: func(cmd *cobra.Command, args []string) {
		sigIntCh := make(chan os.Signal, 1)
		signal.Notify(sigIntCh, os.Interrupt)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			oscall := <-sigIntCh
			log.Printf("system call: %+v", oscall)
			cancel()
		}()

		go runWebServer()

		<-ctx.Done()
		log.Printf("server stopped")
	},

	PostRun: func(cmd *cobra.Command, args []string) {
		err := db.CloseDatabase()
		if err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .gpsctf)")

	rand.Seed(time.Now().UnixNano())
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
