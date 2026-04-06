package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AI struct {
		Key     string `mapstructure:"key"`
		BaseURL string `mapstructure:"base_url"`
	} `mapstructure:"ai"`
}

var AppConfig Config

func Load() error {
	_ = godotenv.Load(".env")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("TMK")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("ai.base_url", "https://api.siliconflow.cn/v1")

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}
	if osKey := os.Getenv("TMK_AI_KEY"); osKey != "" {
		AppConfig.AI.Key = osKey
	}
	return nil
}
