package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration for our application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	App      AppConfig      `mapstructure:"app"`
}

type ServerConfig struct {
	Port         string `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// Load reads configuration from file and environment variables
func Load() (*Config, error) {
	config := &Config{}

	// Set the file name of the configurations file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Set the path to look for the configurations file
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("./")
	viper.AddConfigPath("/etc/chat-agent/")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal the config into our struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.idle_timeout", 60)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "chat_agent")
	viper.SetDefault("database.sslmode", "disable")

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	// App defaults
	viper.SetDefault("app.name", "Chat Agent")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", getEnv("ENVIRONMENT", "development"))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetConfigPath returns the path to config file
func GetConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs"
	}
	return configPath
}

// CreateConfigFile creates a sample config file
func CreateConfigFile() error {
	configPath := GetConfigPath()
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}

	configFile := filepath.Join(configPath, "config.yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		sampleConfig := `# Chat Agent Configuration
server:
  port: "8080"
  host: "0.0.0.0"
  read_timeout: 10
  write_timeout: 10
  idle_timeout: 60

database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: ""
  dbname: "chat_agent"
  sslmode: "disable"

log:
  level: "info"
  format: "json"

app:
  name: "Chat Agent"
  version: "1.0.0"
  environment: "development"
`
		return os.WriteFile(configFile, []byte(sampleConfig), 0644)
	}
	return nil
}

