package main

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	envFileName = "./.env"
)

type enver interface {
	loadEnv() error
}

func newEnv() *env {
	return &env{}
}

type env struct {
}

func (*env) loadEnv() error {
	_, err := os.Stat(envFileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return godotenv.Load()
}
