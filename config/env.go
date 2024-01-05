package config

import (
	"errors"
	"fmt"

	"github.com/joho/godotenv"
)

func loadenv() error {
	err := godotenv.Load()

	if err != nil {
		return errors.New("error loading .env file")
	}
	return nil
}
func Loadenv() {
	load := loadenv()
	if load == nil {
		fmt.Println(".env Load Success")
	}
}
