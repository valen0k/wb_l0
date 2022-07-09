package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
	Storage struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"storage"`
}

func NewConfig(configFile string) (*Config, error) {
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err = json.Unmarshal(file, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) NewDBConnection() (*pgx.Conn, error) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			c.Storage.User, c.Storage.Password,
			c.Storage.Host, c.Storage.Port,
			c.Storage.Database))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
