package configs

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	mapErr := godotenv.Load(".env")
	if mapErr != nil {
		fmt.Println(mapErr)
	}
}
