package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/arizon-dread/webdig-backend/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	readConfig()
	cfg := config.GetInstance()
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(cors.New(cors.Config{
		AllowOrigins: cfg.General.Cors.Origins,
		AllowHeaders: cfg.General.Cors.Headers,
		AllowMethods: cfg.General.Cors.Methods,
		AllowOriginFunc: func(origin string) bool {
			return slices.Contains(cfg.General.Cors.Origins, origin)
		},
	}))
	router.POST("/api/dig", lookup)
	router.GET("/healthz", healthz)
	router.Run(":8080")
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
	atLeastOneUnfiltered := false
	for _, v := range cfg.DNS {
		if v.FilterDuplicates == false {
			atLeastOneUnfiltered = true
		}
	}
	if !atLeastOneUnfiltered {
		panic("At least one dns entry must be set to FilterDuplicates: false.")
	}
}

func breakOnNoConfig(err error) {
	fmt.Printf("error when reading config, %v\n", err)
	panic("Failed to read config")
}
