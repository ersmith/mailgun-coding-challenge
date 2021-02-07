package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func createConfig() Config {
	return Config{
		Logger: zap.NewNop().Sugar(),
	}
}

func TestDbUserDefault(t *testing.T) {
	config := createConfig()
	os.Unsetenv("DB_USER")
	assert.Equal(t, "postgres", config.DbUser())
}

func TestDbUserNotDefault(t *testing.T) {
	config := createConfig()
	var dbUser = "test_user"
	os.Setenv("DB_USER", dbUser)
	assert.Equal(t, dbUser, config.DbUser())
	os.Unsetenv("DB_USER")
}

func TestDbPasswordDefault(t *testing.T) {
	config := createConfig()
	os.Unsetenv("DB_PASSWORD")
	assert.Equal(t, "postgres", config.DbPassword())
}

func TestDbPasswordNotDefault(t *testing.T) {
	config := createConfig()
	var dbPassword = "testPassword1234!"
	os.Setenv("DB_PASSWORD", dbPassword)
	assert.Equal(t, dbPassword, config.DbPassword())
	os.Unsetenv("DB_PASSWORD")
}

func TestDbHostDefault(t *testing.T) {
	config := createConfig()
	os.Unsetenv("DB_HOST")
	assert.Equal(t, "localhost", config.DbHost())
}

func TestDbHostNotDefault(t *testing.T) {
	config := createConfig()
	var dbHost = "example.com"
	os.Setenv("DB_HOST", dbHost)
	assert.Equal(t, dbHost, config.DbHost())
	os.Unsetenv("DB_HOST")
}

func TestDbPortDefault(t *testing.T) {
	config := createConfig()
	os.Unsetenv("DB_PORT")
	assert.Equal(t, "15432", config.DbPort())
}

func TestDbPortNotDefault(t *testing.T) {
	config := createConfig()
	var dbPort = "98765"
	os.Setenv("DB_PORT", dbPort)
	assert.Equal(t, dbPort, config.DbPort())
	os.Unsetenv("DB_PORT")
}

func TestDbNameDefault(t *testing.T) {
	config := createConfig()
	os.Unsetenv("DB_NAME")
	assert.Equal(t, "mailgun_dev", config.DbName())
}

func TestDbNameNotDefault(t *testing.T) {
	config := createConfig()
	var dbName = "test_db_name"
	os.Setenv("DB_NAME", dbName)
	assert.Equal(t, dbName, config.DbName())
	os.Unsetenv("DB_NAME")
}

func TestDbMinPoolSizeDefault(t *testing.T) {
	config := createConfig()
	os.Unsetenv("DB_MIN_POOL_SIZE")
	assert.Equal(t, int32(10), config.DbMinPoolSize())
}

func TestDbMinPoolSizeNotDefault(t *testing.T) {
	config := createConfig()
	var minPoolSize = int32(1000)
	os.Setenv("DB_MIN_POOL_SIZE", fmt.Sprint(minPoolSize))
	assert.Equal(t, minPoolSize, config.DbMinPoolSize())
	os.Unsetenv("DB_MIN_POOL_SIZE")
}

func TestDbMaxPoolSizeDefault(t *testing.T) {
	config := createConfig()
	os.Unsetenv("DB_MAX_POOL_SIZE")
	assert.Equal(t, int32(30), config.DbMaxPoolSize())
}

func TestDbMaxPoolSizeNotDefault(t *testing.T) {
	config := createConfig()
	var maxPoolSize = int32(9999)
	os.Setenv("DB_MAX_POOL_SIZE", fmt.Sprint(maxPoolSize))
	assert.Equal(t, maxPoolSize, config.DbMaxPoolSize())
	os.Unsetenv("DB_MAX_POOL_SIZE")
}

func TestHttpPortDefault(t *testing.T) {
	config := createConfig()
	os.Unsetenv("HTTP_PORT")
	assert.Equal(t, "8080", config.HttpPort())
}

func TestHttpPortNotDefault(t *testing.T) {
	config := createConfig()
	var httpPort = "18181"
	os.Setenv("HTTP_PORT", httpPort)
	assert.Equal(t, httpPort, config.HttpPort())
	os.Unsetenv("HTTP_PORT")
}

func TestGetEnvUnset(t *testing.T) {
	var varName = "RANDOM_TEST_VAR_12345"
	var defaultValue = "some_random_value"
	assert.Equal(t, defaultValue, getEnv(varName, defaultValue))
	os.Unsetenv(varName)
}

func TestGetEnvSet(t *testing.T) {
	var varName = "RANDOM_TEST_VAR_12345"
	var defaultValue = "some_random_value"
	var setValue = "123456"
	os.Setenv(varName, setValue)
	assert.Equal(t, setValue, getEnv(varName, defaultValue))
	os.Unsetenv(varName)
}
