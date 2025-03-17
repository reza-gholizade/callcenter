package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT configuration
	JWTSecret     string
	JWTExpiration time.Duration

	// External API configuration
	NLPAPIKey            string
	NLPAPIURL            string
	MyTicketAPIKey       string
	MyTicketAPIURL       string
	FaranegarAPIKey      string
	FaranegarAPIURL      string
	PaymentGatewayAPIKey string
	PaymentGatewayAPIURL string

	// SMS/Email configuration
	SMSAPIKey   string
	SMSAPIURL   string
	EmailAPIKey string
	EmailAPIURL string

	// Telegram/WhatsApp configuration
	TelegramBotToken string
	WhatsAppAPIKey   string
	WhatsAppAPIURL   string

	// Logging configuration
	LogLevel string
	LogFile  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Server configuration
	config.Port = getEnvOrDefault("PORT", "8080")
	readTimeout, _ := strconv.Atoi(getEnvOrDefault("READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnvOrDefault("WRITE_TIMEOUT", "10"))
	config.ReadTimeout = time.Duration(readTimeout) * time.Second
	config.WriteTimeout = time.Duration(writeTimeout) * time.Second

	// Database configuration
	config.DBHost = getEnvOrDefault("DB_HOST", "localhost")
	config.DBPort = getEnvOrDefault("DB_PORT", "5432")
	config.DBUser = getEnvOrDefault("DB_USER", "postgres")
	config.DBPassword = getEnvOrDefault("DB_PASSWORD", "postgres")
	config.DBName = getEnvOrDefault("DB_NAME", "callcenter")
	config.DBSSLMode = getEnvOrDefault("DB_SSL_MODE", "disable")

	// JWT configuration
	config.JWTSecret = getEnvOrDefault("JWT_SECRET", "your-secret-key")
	jwtExpiration, _ := strconv.Atoi(getEnvOrDefault("JWT_EXPIRATION", "24"))
	config.JWTExpiration = time.Duration(jwtExpiration) * time.Hour

	// External API configuration
	config.NLPAPIKey = getEnvOrDefault("NLP_API_KEY", "")
	config.NLPAPIURL = getEnvOrDefault("NLP_API_URL", "")
	config.MyTicketAPIKey = getEnvOrDefault("MYTICKET_API_KEY", "")
	config.MyTicketAPIURL = getEnvOrDefault("MYTICKET_API_URL", "")
	config.FaranegarAPIKey = getEnvOrDefault("FARANEGAR_API_KEY", "")
	config.FaranegarAPIURL = getEnvOrDefault("FARANEGAR_API_URL", "")
	config.PaymentGatewayAPIKey = getEnvOrDefault("PAYMENT_GATEWAY_API_KEY", "")
	config.PaymentGatewayAPIURL = getEnvOrDefault("PAYMENT_GATEWAY_API_URL", "")

	// SMS/Email configuration
	config.SMSAPIKey = getEnvOrDefault("SMS_API_KEY", "")
	config.SMSAPIURL = getEnvOrDefault("SMS_API_URL", "")
	config.EmailAPIKey = getEnvOrDefault("EMAIL_API_KEY", "")
	config.EmailAPIURL = getEnvOrDefault("EMAIL_API_URL", "")

	// Telegram/WhatsApp configuration
	config.TelegramBotToken = getEnvOrDefault("TELEGRAM_BOT_TOKEN", "")
	config.WhatsAppAPIKey = getEnvOrDefault("WHATSAPP_API_KEY", "")
	config.WhatsAppAPIURL = getEnvOrDefault("WHATSAPP_API_URL", "")

	// Logging configuration
	config.LogLevel = getEnvOrDefault("LOG_LEVEL", "info")
	config.LogFile = getEnvOrDefault("LOG_FILE", "app.log")

	// Validate required configuration
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate checks if all required configuration is present
func (c *Config) validate() error {
	required := map[string]string{
		"DB_PASSWORD": c.DBPassword,
		"JWT_SECRET":  c.JWTSecret,
	}

	for name, value := range required {
		if value == "" {
			return fmt.Errorf("required configuration %s is missing", name)
		}
	}

	return nil
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}
