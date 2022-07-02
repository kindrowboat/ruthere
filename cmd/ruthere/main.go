package main

import (
	"net/mail"
	"os"
	"time"

	"github.com/motevets/ruthere/pkg/ruthere"
	"gopkg.in/yaml.v3"
)

type configYaml struct {
	FromName    string   `yaml:"from_name"`
	FromAddress string   `yaml:"from_address"`
	ToName      string   `yaml:"to_name"`
	ToAddress   string   `yaml:"to_address"`
	CheckSites  []string `yaml:"check_sites"`
	SleepTime   string   `yaml:"sleep_time"`
}

func main() {
	var err error

	var config configYaml
	rawYaml, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(rawYaml, &config)
	if err != nil {
		panic(err)
	}

	fromAddress := mail.Address{
		Name:    config.FromName,
		Address: config.FromAddress,
	}
	toAddresses := []mail.Address{{
		Name:    config.ToName,
		Address: config.ToAddress,
	}}

	sleepTime, err := time.ParseDuration(config.SleepTime)
	if err != nil {
		panic(err)
	}

	checker := ruthere.NewHttpChecker(fromAddress, toAddresses, config.CheckSites, sleepTime)
	checker.Run()
	return
}
