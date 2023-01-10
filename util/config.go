package util

import "github.com/spf13/viper"

//cau hinh database

type Config struct {
	TOKEN      string `mapstructure:"TOKEN"`
	DBSource   string `mapstructure:"DB_SOURCE"`
	SERVER_GIT string `mapstructure:"SERVER_GIT"`
	GROUP_ID   string `mapstructure:"GROUP_ID"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return

}
