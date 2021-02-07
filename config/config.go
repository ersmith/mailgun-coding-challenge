package config

import (
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"
)

// Interface for types which can provide DbConfiguration.
type DbConfigProvider interface {
	DbConfig() *DbConfig
}

type Config struct {
	Logger *zap.SugaredLogger
}

// Struct to make it easier to pass around database configuration
type DbConfig struct {
	Username    string
	Password    string
	Host        string
	Port        string
	Name        string
	MinPoolSize int32
	MaxPoolSize int32
}

// Creates a new
func (self *DbConfig) ConnnectionUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		self.Username,
		self.Password,
		self.Host,
		self.Port,
		self.Name)
}

// Returns the DbConfiguration based on the current config
func (self *Config) DbConfig() *DbConfig {
	return &DbConfig{
		Username:    self.DbUser(),
		Password:    self.DbPassword(),
		Host:        self.DbHost(),
		Port:        self.DbPort(),
		Name:        self.DbName(),
		MinPoolSize: self.DbMinPoolSize(),
		MaxPoolSize: self.DbMaxPoolSize(),
	}
}

// Returns the username that should be used when connecting to the database
func (self *Config) DbUser() string {
	return getEnv("DB_USER", "postgres")
}

// Returns the password that should be used when connecting to the database
func (self *Config) DbPassword() string {
	return getEnv("DB_PASSWORD", "postgres")
}

// Returns the hostname that should be used when connecting to the database
func (self *Config) DbHost() string {
	return getEnv("DB_HOST", "localhost")
}

// Returns the port that should be used when connecting to the database
func (self *Config) DbPort() string {
	return getEnv("DB_PORT", "15432")
}

// Returns the database name that should be used when connecting to the database
func (self *Config) DbName() string {
	return getEnv("DB_NAME", "mailgun_dev")
}

// Returns the starting pool size for the db connection pool
func (self *Config) DbMinPoolSize() int32 {
	num, err := strconv.ParseInt(getEnv("DB_MIN_POOL_SIZE", "10"), 10, 0)

	if err != nil {
		self.Logger.Fatalw("Failed to parse DB_MIN_POOL_SIZE.", zap.Error(err))
		os.Exit(1)
	}

	return int32(num)
}

// Returns the max pool size for the db connection pool
func (self *Config) DbMaxPoolSize() int32 {
	num, err := strconv.ParseInt(getEnv("DB_MAX_POOL_SIZE", "30"), 10, 0)

	if err != nil {
		self.Logger.Fatalw("Failed to parse DB_MAX_POOL_SIZE", zap.Error(err))
		os.Exit(1)
	}

	return int32(num)
}

// Returns the http port to listen on
func (self *Config) HttpPort() string {
	return getEnv("HTTP_PORT", "8080")
}

// Gets an environment variable or returns the default value if it isn't set.
func getEnv(name string, defaultValue string) string {
	value, set := os.LookupEnv(name)

	if set {
		return value
	}

	return defaultValue
}
