package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/arizon-dread/webdig-backend/config"
	"github.com/arizon-dread/webdig-backend/pkg/api"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func main() {

	readConfig()
	cfg := config.GetInstance()
	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins: cfg.General.Cors.Origins,
		AllowedMethods: cfg.General.Cors.Methods,
		AllowedHeaders: cfg.General.Cors.Headers,
	})

	mux.HandleFunc("POST /api/dig", api.Lookup)
	mux.HandleFunc("GET /healthz", api.Healthz)
	handler := c.Handler(mux)

	var protos http.Protocols
	protos.SetHTTP1(true)
	protos.SetHTTP2(true)
	httpServer := &http.Server{
		Addr:      ":8080",
		Handler:   handler,
		Protocols: &protos,
	}
	log.Fatal(httpServer.ListenAndServe())
}

func readConfig() {

	cfg := config.GetInstance()
	//We use a dedicated folder for the config file to ease the configMap volume mount.
	viper.SetConfigFile("./confFile/config.yaml")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		breakOnNoConfig(err)
	}
	//all keys that are read from config will get overwritten by their env equivalents, as long as they exist in config (empty or not)
	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		viper.Set(key, val)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		breakOnNoConfig(err)
	}
	//Perhaps we need to check that they don't circularly reference eachother...
	atLeastOneUnfiltered := false
	for _, v := range cfg.DNS {
		if len(v.FilterDuplicates) == 0 {
			atLeastOneUnfiltered = true
		}
	}
	if !atLeastOneUnfiltered {
		fmt.Println("[WARNING] You have set filterDuplicates on all server groups. If you create a circular filter for all groups, you will get an empty response every time.")
		panic("Example circular filter that would yield an empty result on every request: A -> B -> C -> A. ")
	}
}

func breakOnNoConfig(err error) {
	fmt.Printf("error when reading config, %v\n", err)
	panic("Failed to read config")
}
