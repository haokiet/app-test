package config

import "github.com/spf13/viper"

type configEnv struct {
	App struct {
		Port string `json:"port"`
	} `json:"app"`

	Mongo struct {
		Uri      string `json:"uri"`
		Database string `json:"database"`
	} `json:"mongo"`
}

var C configEnv

func init() {
	viper.SetConfigFile(`config/config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		panic(err)
	}
}

func GetConfig() configEnv {
	return C
}
