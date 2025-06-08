package utils

import "os"

const configFile = "./config.yaml"

func LoadConfig() ([]byte, error) {
	return os.ReadFile(configFile)
}
