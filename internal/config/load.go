package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (c *Config) SetConfigPath(path string) {
	extension := filepath.Ext(path)                       // eg: .yml
	filename := filepath.Base(path)                       // eg: config.yml
	configName := strings.TrimSuffix(filename, extension) // eg: config
	c.loader.SetConfigName(configName)                    // viper takes filename without extension

	if len(extension) > 1 {
		configType := extension[1:]
		c.loader.SetConfigType(configType)
	}

	configDir := filepath.Dir(path) // eg: /app or .
	c.loader.AddConfigPath(configDir)
}

func (c *Config) LoadConfig() error {
	if err := c.loader.ReadInConfig(); err != nil {
		return fmt.Errorf("loader.ReadInConfig: %w", err)
	}

	envLoader := c.loader.Sub(c.Env.String())

	if err := bindEnv(envLoader, "PORT"); err != nil {
		return fmt.Errorf("bindEnv: %w", err)
	}
	if err := bindEnv(envLoader, "LOG_LEVEL"); err != nil {
		return fmt.Errorf("bindEnv: %w", err)
	}
	if err := bindStructEnv(envLoader, c.Db); err != nil {
		return fmt.Errorf("bindStructEnv: %w", err)
	}
	if err := envLoader.Unmarshal(&c); err != nil {
		return fmt.Errorf("envLoader.Unmarshal: %w", err)
	}

	return nil
}

func (c *Config) DBConnectionString() string {
	if c.Db.ConnectionString != "" {
		return c.Db.ConnectionString
	} else if c.Db.Type == "postgres" {
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", c.Db.Host, c.Db.Port, c.Db.User, c.Db.Name, c.Db.Password)
	} else if c.Db.Type == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Db.User, c.Db.Password, c.Db.Host, c.Db.Port, c.Db.Name)
	}

	panic("Invalid database type")
}
