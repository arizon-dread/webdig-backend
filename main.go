package main

import (
	"fmt"
	"strings"

	"github.com/arizon-dread/webdig-backend/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	router := gin.Default()
	router.POST("/api/dig", lookup)
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
}

func breakOnNoConfig(err error) {
	fmt.Printf("error when reading config, %v\n", err)
	panic("Failed to read config")
}
