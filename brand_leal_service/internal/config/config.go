package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	DBHost              string
	DBPort              int
	DBUser              string
	DBPassword          string
	DBName              string
	KafkaBrokers        []string
	MsgPurchaseTopic    string
	MsgApplyPointsTopic string
	BrandGroup          string
	HTTPServerPort      string
}

var (
	configInstance *Config
	once           sync.Once
)

// LoadConfig initializes the configuration singleton
func LoadConfig() (*Config, error) {
	var loadErr error
	once.Do(func() {
		dbPort, err := strconv.Atoi(getEnv("DB_PORT"))
		if err != nil {
			loadErr = err
			return
		}

		configInstance = &Config{
			DBHost:              getEnv("DB_HOST"),
			DBPort:              dbPort,
			DBUser:              getEnv("DB_USER"),
			DBPassword:          getEnv("DB_PASSWORD"),
			DBName:              getEnv("DB_NAME"),
			KafkaBrokers:        []string{getEnv("MSG_BROKER_ADDRESS")},
			MsgPurchaseTopic:    getEnv("MSG_PURCHASE"),
			MsgApplyPointsTopic: getEnv("MSG_APPLY_POINTS"),
			BrandGroup:          getEnv("BRAND_GROUP_NAME"),
			HTTPServerPort:      getEnv("HTTP_SERVER_PORT"),
		}
	})

	return configInstance, loadErr
}

// GetConfig returns the singleton instance of the configuration
func GetConfig() *Config {
	return configInstance
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}
