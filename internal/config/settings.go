package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Filter string

const (
	defaultListenAddress = "0.0.0.0"
	defaultListenPort    = 8080
	localPath            = "."

	defaultLDAPWorkersCount = 10
	defaultLDAPTimeout      = 10 * time.Second

	SiteFilter       Filter = "site"
	SiteGroupFilter  Filter = "site_group"
	SiteRegionFilter Filter = "region"
)

var (
	Cfg Config
)

type Config struct {
	Authentication AuthConfig
	NetBox         struct {
		URL                 string
		APIKey              string
		DatacenterFilterKey Filter
	}
	Log struct {
		Level  string
		Pretty bool
	}
	Datacenter string
	API        struct {
		ListenAddress string
		ListenPort    int
	}
	Build struct {
		Interval            time.Duration
		AllDevicesMustBuild bool
	}
}

type AuthConfig struct {
	LDAP *LDAPConfig
}

type LDAPConfig struct {
	URL                   string
	BaseDN                string
	BindDN                string
	Password              string
	Timeout               time.Duration
	MaxConnectionLifetime time.Duration
	InsecureSkipVerify    bool
	WorkersCount          int
}

func setDefaults() {
	viper.SetDefault("Datacenter", "")
	viper.SetDefault("Log.Level", "info")
	viper.SetDefault("Log.Pretty", false)

	viper.SetDefault("API.ListenAddress", defaultListenAddress)
	viper.SetDefault("API.ListenPort", defaultListenPort)

	viper.SetDefault("NetBox.URL", "")
	viper.SetDefault("NetBox.APIKey", "")
	viper.SetDefault("NetBox.DatacenterFilterKey", SiteFilter)

	viper.SetDefault("Build.Interval", time.Minute)
	viper.SetDefault("Build.AllDevicesMustBuild", false)

	viper.SetDefault("Authentication.LDAP.URL", "")
	viper.SetDefault("Authentication.LDAP.BaseDN", "")
	viper.SetDefault("Authentication.LDAP.BindDN", "")
	viper.SetDefault("Authentication.LDAP.Password", "")
	viper.SetDefault("Authentication.LDAP.InsecureSkipVerify", false)
	viper.SetDefault("Authentication.LDAP.WorkersCount", defaultLDAPWorkersCount)
	viper.SetDefault("Authentication.LDAP.Timeout", defaultLDAPTimeout)
	viper.SetDefault("Authentication.LDAP.MaxConnectionLifetime", time.Minute)
}

func LoadConfig() error {
	viper.SetConfigName("settings")
	viper.AddConfigPath(localPath)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("DAAPI")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	setDefaults()
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok { //nolint:errorlint  // failing only if invalid file
			return fmt.Errorf("invalid config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		return err
	}

	return nil
}
