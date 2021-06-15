package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	MysqlDSN string         `yaml:"mysqlDSN"`
	Port     string         `yaml:"prepare"`
	MVS      []VolumeStatus `multiple volume status. yaml:"mvs"`
}

// UnmarshalConfFile parses the CONFIG yaml file.
func UnmarshalConfFile(confFile string) (*Config, error) {
	conf := &Config{}
	viper.SetConfigFile(confFile)
	viper.ReadInConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, fmt.Errorf("read CONFIG file %s error: %q", confFile, err)
	}
	return conf, nil
}

// ParseConfFile parses CONFIG yaml file and flags to *Conf.
func ParseConfFile(confFile string, port string) *Config {
	c := loadConfFile(confFile)
	if c == nil {
		c = &Config{}
	}

	if port != "" {
		c.Port = port
	}

	confJSON, _ := json.Marshal(c)
	log.Printf("Configuration: %s", confJSON)

	return c
}

func loadConfFile(confFile string) *Config {
	if confFile == "" {
		if s, err := os.Stat(defaultConfFile); err != nil || s.IsDir() {
			return nil // not exists
		}
		confFile = defaultConfFile
	}

	c, err := UnmarshalConfFile(confFile)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
