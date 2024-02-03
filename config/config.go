package config

import "sync"

var cfg *Config

var lock = &sync.Mutex{}

func GetInstance() *Config {
	if cfg == nil {
		lock.Lock()
		defer lock.Unlock()
		if cfg == nil {
			cfg = &Config{}
		}
	}
	return cfg
}

type Config struct {
	DNS     []ServerGroup
	General General
}

type ServerGroup struct {
	Name             string
	Servers          []string
	FilterDuplicates bool
}

type General struct {
	Cors Cors
}

type Cors struct {
	Origins []string
	Methods []string
	Headers []string
}
