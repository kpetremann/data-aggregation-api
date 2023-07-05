package config

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

const (
	defaultListenAddress = "0.0.0.0"
	defaultListenPort    = 8080
)

var (
	Cfg Config
)

type Config struct {
	Datacenter             string
	NetboxURL              string
	NetboxAPIKey           string
	LdapURL                string
	LdapBindDN             string
	LdapPassword           string
	LdapBaseDN             string
	LogLevel               string
	ListenAddress          string
	ListenPort             int
	BuildInterval          time.Duration
	LdapInsecureSkipVerify bool
	AllDevicesMustBuild    bool
}

func setDefaults() {
	viper.SetDefault("ListenAddress", defaultListenAddress)
	viper.SetDefault("ListenPort", defaultListenPort)
	viper.SetDefault("Datacenter", "")
	viper.SetDefault("NetboxURL", "")
	viper.SetDefault("NetboxAPIKey", "")
	viper.SetDefault("BuildInterval", time.Minute)
	viper.SetDefault("LdapURL", "")
	viper.SetDefault("LdapBindDN", "")
	viper.SetDefault("LdapPassword", "")
	viper.SetDefault("LdapBaseDN", "")
	viper.SetDefault("LdapInsecureSkipVerify", false)
	viper.SetDefault("AllDevicesMustBuild", false)
	viper.SetDefault("LogLevel", "info")
}

func LoadConfig() error {
	viper.SetConfigName("settings")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("DAAPI")

	setDefaults()
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Send()
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		return err
	}

	return nil
}
