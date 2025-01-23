package global

import (
	"github.com/spf13/viper"
	"log"
	"short-url/config"
	"short-url/internal/common/enum"
)

var Conf *config.Config

func InitConfig() {
	viper.SetConfigFile(enum.ConfigFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Read in config error: ", err.Error())
	}

	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Fatal("Init config error: ", err.Error())
	}
}
