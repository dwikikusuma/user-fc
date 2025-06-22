package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func LoadConfig() Config {
	var cfg Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./files/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Error reading config file")
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic("Error unmarshalling config file")
	}
	configDetails := fmt.Sprintf("Config loaded successfully: %v", cfg)
	fmt.Println(configDetails)

	return cfg
}
