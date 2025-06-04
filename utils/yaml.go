package utils

import "os"

const configFile = "E:\\go代码\\parking-system-go\\parking-system-go\\config.yaml"

func LoadConfig() ([]byte, error) {
	return os.ReadFile(configFile)
}
