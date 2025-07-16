// Author: Bontus
// Email: bontus.doku@gmail.com
// License: MIT
// Created: 5/7/2025

package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// config TODO:

type Config struct {
	DBUsername     string `mapstructure:"DB_USERNAME"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBDriver       string `mapstructure:"DB_DRIVER"`
	DBName         string `mapstructure:"DB_NAME"`
	SSLMode        string `mapstructure:"SSLMODE"`
	ServerPort     int    `mapstructure:"SERVER_PORT"`
	JwtSecret      string `mapstructure:"JWT_SECRET"`
	AdminEmail     string `mapstructure:"ADMIN_EMAIL"`
	AdminPassword  string `mapstructure:"ADMIN_PASSWORD"`
	PlunkBaseUrl   string `mapstructure:"PLUNK_BASE_URL"`
	PlunkSecretKey string `mapstructure:"PLUNK_SECRET_KEY"`
	BaseURL        string `mapstructure:"BaseURL"`
}

func LoadConfig(path string) (*Config, error) {
	// Validate that the path is not empty
	if path == "" {
		path = "."
	}

	// Create a new Viper instance to avoid global state
	v := viper.New()

	// Disable environment variable prefix
	v.SetEnvPrefix("")
	v.AutomaticEnv()

	// Configure config file
	v.AddConfigPath(path)
	v.SetConfigName(".env")
	v.SetConfigType("env")

	// Read a config file
	if err := v.ReadInConfig(); err != nil {
		// Log the error but don't fail entirely
		log.Printf("Warning: Unable to read config file: %v", err)
	}

	_ = v.BindEnv("DB_USERNAME")
	_ = v.BindEnv("DB_PASSWORD")
	_ = v.BindEnv("DB_HOST")
	_ = v.BindEnv("DB_NAME")
	_ = v.BindEnv("DB_PORT")
	_ = v.BindEnv("SSLMODE")
	_ = v.BindEnv("DB_DRIVER")
	_ = v.BindEnv("SERVER_PORT")
	_ = v.BindEnv("ADMIN_EMAIL")
	_ = v.BindEnv("ADMIN_PASSWORD")
	_ = v.BindEnv("PLUNK_BASE_URL")
	_ = v.BindEnv("PLUNK_SECRET_KEY")
	_ = v.BindEnv("BaseURL")

	// Create config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Additional security: Validate critical configurations
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	// Add validation for critical configurations
	if config.ServerPort == 0 {
		return fmt.Errorf("server port must be specified")
	}

	// Add more validation as needed
	if config.DBUsername == "" || config.DBPassword == "" {
		return fmt.Errorf("database credentials must be provided")
	}

	return nil
}
