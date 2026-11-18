package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AuthTokenExpireTime int `yaml:"auth.token.expire_time"`

	ServerPort   int    `yaml:"server.port"`
	ServerDomain string `yaml:"server.domain"`

	AdminEmail    string `yaml:"admin.email"`
	AdminPassword string `yaml:"admin.password"`

	WorkspaceDefaultName string `yaml:"workspace.default.name"`

	DatabaseName     string `yaml:"database.name"`
	DatabaseHost     string `yaml:"database.host"`
	DatabaseSSLMode  string `yaml:"database.ssl_mode"`
	DatabasePort     int    `yaml:"database.port"`
	DatabaseUser     string `yaml:"database.user"`
	DatabasePassword string `yaml:"database.password"`

	RedisHost     string `yaml:"redis.host"`
	RedisPort     int    `yaml:"redis.port"`
	RedisPassword string `yaml:"redis.password"`
	RedisDB       int    `yaml:"redis.db"`

	KafkaBrokers []string `yaml:"kafka.brokers"`
	KafkaTopic   string   `yaml:"kafka.topic"`

	KafkaConsumerCount int `yaml:"kafka.consumer.count"`

	Hosts []string `yaml:"hosts"`

	Providers []ProviderConfig `yaml:",flow"`

	LogsDirectory    string `yaml:"logs.directory"`
	LogsLevel        string `yaml:"logs.level"`
	LogsFormat       string `yaml:"logs.format"`
	LogsRotationSize string `yaml:"logs.rotation.size"`
	LogsMaxFiles     int    `yaml:"logs.rotation.max_files"`
	LogsCompress     bool   `yaml:"logs.rotation.compress"`
	LogsErrorFile    string `yaml:"logs.error_file"`
	LogsAccessFile   string `yaml:"logs.access_file"`

	MetricsEnabled  bool   `yaml:"metrics.enabled"`
	MetricsEndpoint string `yaml:"metrics.endpoint"`
}

type ProviderConfig struct {
	Type           string `yaml:"type"`
	Name           string `yaml:"name"`
	SmtpHost       string `yaml:"smtp_host,omitempty"`
	SmtpPort       int    `yaml:"smtp_port,omitempty"`
	Username       string `yaml:"username,omitempty"`
	Password       string `yaml:"password,omitempty"`
	FromEmail      string `yaml:"from_email,omitempty"`
	Url            string `yaml:"url,omitempty"`
	Method         string `yaml:"method,omitempty"`
	AuthTokenKey   string `yaml:"auth_token_key,omitempty"`
	AuthTokenValue string `yaml:"auth_token_value,omitempty"`
	IsAuthInHeader bool   `yaml:"is_auth_in_header,omitempty"`
	BotToken       string `yaml:"bot_token,omitempty"`
	WebhookUrl     string `yaml:"webhook_url,omitempty"`
}

// NewConfig is the constructor that parses the config file
func NewConfig(appFlags AppFlags) Config {
	// Load .env file if exists
	loadEnvVariables()
	// Read the YAML file
	filePath := getConfigFile()

	if appFlags.ConfigPath != nil || *appFlags.ConfigPath != "" {
		filePath = *appFlags.ConfigPath
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("error reading config file: %v", err))
	}

	// Replace placeholders in the raw YAML data before unmarshalling
	data = replaceEnvInData(data)

	// Parse the YAML file
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(fmt.Sprintf("error unmarshalling config: %v", err))
	}

	// Replace placeholders with environment variables
	replaceEnvVariables(&config)

	if appFlags.LogLevel != nil {
		config.LogsLevel = *appFlags.LogLevel
	}

	return config
}

// replaceEnvVariables replaces ${ENV_VAR} placeholders with actual environment variable values
func replaceEnvVariables(config *Config) {
	// Auth
	config.DatabaseName = replaceEnv(config.DatabaseName)
	config.DatabaseHost = replaceEnv(config.DatabaseHost)
	config.DatabaseSSLMode = replaceEnv(config.DatabaseSSLMode)
	config.DatabasePort = replaceEnvInt(config.DatabasePort)
	config.DatabaseUser = replaceEnv(config.DatabaseUser)
	config.DatabasePassword = replaceEnv(config.DatabasePassword)

	// Redis
	config.RedisHost = replaceEnv(config.RedisHost)
	config.RedisPassword = replaceEnv(config.RedisPassword)
	config.RedisDB = replaceEnvInt(config.RedisDB)

	// Providers
	for i := range config.Providers {
		config.Providers[i].SmtpHost = replaceEnv(config.Providers[i].SmtpHost)
		config.Providers[i].Username = replaceEnv(config.Providers[i].Username)
		config.Providers[i].Password = replaceEnv(config.Providers[i].Password)
		config.Providers[i].FromEmail = replaceEnv(config.Providers[i].FromEmail)
		config.Providers[i].Url = replaceEnv(config.Providers[i].Url)
		config.Providers[i].BotToken = replaceEnv(config.Providers[i].BotToken)
		config.Providers[i].WebhookUrl = replaceEnv(config.Providers[i].WebhookUrl)
	}

	// Logs
	config.LogsDirectory = replaceEnv(config.LogsDirectory)
	config.LogsLevel = replaceEnv(config.LogsLevel)
	config.LogsFormat = replaceEnv(config.LogsFormat)
	config.LogsRotationSize = replaceEnv(config.LogsRotationSize)
	config.LogsErrorFile = replaceEnv(config.LogsErrorFile)
	config.LogsAccessFile = replaceEnv(config.LogsAccessFile)
}

func loadEnvVariables() {
	_ = godotenv.Load()
}

func getConfigFile() string {
	var filePaths = []string{
		"/etc/avan/config.yml",
		"/etc/avan/config.yaml",
		"config/config.yml",
		"config/config.yaml",
	}

	var filePath string
	for _, path := range filePaths {
		if _, err := os.Stat(path); err == nil {
			filePath = path
			break
		}
	}

	if filePath == "" {
		panic("config file not found in any of the specified paths")
	}

	return filePath
}

// replaceEnvInData replaces environment variable placeholders in the raw YAML data
func replaceEnvInData(data []byte) []byte {
	// Regular expression to find ${ENV_VAR} placeholders
	re := regexp.MustCompile(`\$\{([A-Za-z0-9_]+)\}`)

	// Replace each match with the value of the environment variable
	result := re.ReplaceAllFunc(data, func(match []byte) []byte {
		// Extract the environment variable name
		envVar := string(match[2 : len(match)-1]) // Remove ${ and }
		envValue := os.Getenv(envVar)             // Look up the environment variable

		// If the environment variable is not found, return the original placeholder
		if envValue == "" {
			return match
		}

		// If the environment variable exists, replace it in the data
		return []byte(envValue)
	})

	return result
}

// replaceEnv replaces a placeholder with the corresponding environment variable value
func replaceEnv(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		envVar := value[2 : len(value)-1]
		if val, exists := os.LookupEnv(envVar); exists {
			return val
		}
	}
	return value
}

// replaceEnvInt replaces an integer placeholder with the environment variable value
func replaceEnvInt(value int) int {
	return value // No actual replacement for integer in this case
}
