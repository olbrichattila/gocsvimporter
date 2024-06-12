package main

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	envFileName = "./.env"
)

func NewEnv() *env {
	return &env{}
}

type env struct {
}

func (*env) LoadEnv() error {
	_, err := os.Stat(envFileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil
		}

		return err
	}

	return godotenv.Load()
}
