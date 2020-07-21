package getenv

import "os"

func getEnv(key, defaultValue string) string {
	value, err := os.Getenv(key)
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return defaultValue
	}
	return value, nil
}