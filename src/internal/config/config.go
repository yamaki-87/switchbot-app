package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	App       AppConfig       `toml:"app"`
	SwitchBot SwitchBotConfig `toml:"switchbot"`
	Database  DatabaseConfig  `toml:"database"`
}

type AppConfig struct {
	IntervalTimeout int `toml:"interval_timeout"`
}

type SwitchBotConfig struct {
	DeviceIDs []string `toml:"device_ids"`
}

type DatabaseConfig struct {
	URL string `toml:"url"`
}

var config *Config

func GetConfig() *Config {
	return config
}

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, &config)
}

type SecretEnv struct {
	Token  string
	Secret string
}

var secretEnv *SecretEnv

func GetSecretEnv() *SecretEnv {
	return secretEnv
}

func LoadSecretEnv() error {
	token := os.Getenv("SWITCHBOT_TOKEN")
	secret := os.Getenv("SWITCHBOT_SECRET")
	if token == "" || secret == "" {
		return fmt.Errorf("token or secret is empty")
	}

	secretEnv = &SecretEnv{
		Token:  token,
		Secret: secret,
	}
	return nil
}

var settings *Settings

type Settings struct {
	Config *Config
	Secret *SecretEnv
}

func LoadSettings(path string) error {
	if err := LoadConfig(path); err != nil {
		return err
	}
	if err := LoadSecretEnv(); err != nil {
		return err
	}
	settings = &Settings{
		Config: config,
		Secret: secretEnv,
	}
	return nil
}

func GetSettings() *Settings {
	return settings
}
