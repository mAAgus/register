package main

import (
	"flag"
	"log"
	"register/internal/app/api"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "path", "configs/api.toml", "path to config file in .tomal format")
}

func main() {
	flag.Parse()
	log.Println("It is works")
	config := api.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Println("Can not find configs file. using default values:", err)
	}
	server := api.New(config)

	log.Fatal(server.Start())
}
