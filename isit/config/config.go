package config

import (
	"github.com/spf13/viper"
	"sync"
)

var c Config
var mutex sync.Mutex

// C returning config copy
func C() *Config {
	mutex.Lock()
	defer mutex.Unlock()
	configCopy := c
	return &configCopy
}

type Config struct {
	Port       int
	Postgres   Postgres
	Redis      Redis
	TTLSession int // on hours
}

type Redis struct {
	Host     string
	Post     int
	Password string
	DB       int
}

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

func LoadConfig() (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath("config")
	v.SetConfigName("config")
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) error {
	err := v.Unmarshal(&c)
	if err != nil {
		return err
	}
	return nil
}
