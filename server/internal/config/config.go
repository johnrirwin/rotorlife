package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Cache    CacheConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

// ServerConfig holds HTTP/MCP server configuration
type ServerConfig struct {
	HTTPAddr     string
	MCPMode      bool
	RateLimitDur time.Duration
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Backend   string // "memory" or "redis"
	TTL       time.Duration
	RedisAddr string
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string
}

// Load parses flags and environment variables to build configuration
func Load() *Config {
	cfg := &Config{}

	// Define flags with defaults
	httpAddr := flag.String("http", ":8080", "HTTP server address")
	mcpMode := flag.Bool("mcp", false, "Run in MCP stdio mode")
	cacheTTL := flag.Duration("cache-ttl", 5*time.Minute, "Cache TTL for feed items")
	cacheBackend := flag.String("cache-backend", "memory", "Cache backend: memory or redis")
	redisAddr := flag.String("redis-addr", "localhost:6379", "Redis server address")
	rateLimitDur := flag.Duration("rate-limit", time.Second, "Minimum delay between requests to same host")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	dbHost := flag.String("db-host", "localhost", "PostgreSQL host")
	dbPort := flag.Int("db-port", 5432, "PostgreSQL port")
	dbUser := flag.String("db-user", "postgres", "PostgreSQL user")
	dbPassword := flag.String("db-password", "postgres", "PostgreSQL password")
	dbName := flag.String("db-name", "drone_inventory", "PostgreSQL database name")

	flag.Parse()

	// Apply environment variable overrides
	applyEnvOverrides(httpAddr, mcpMode, cacheTTL, cacheBackend, redisAddr, rateLimitDur, logLevel, dbHost, dbPort, dbUser, dbPassword, dbName)

	// Build config struct
	cfg.Server = ServerConfig{
		HTTPAddr:     *httpAddr,
		MCPMode:      *mcpMode,
		RateLimitDur: *rateLimitDur,
	}

	cfg.Cache = CacheConfig{
		Backend:   *cacheBackend,
		TTL:       *cacheTTL,
		RedisAddr: *redisAddr,
	}

	cfg.Database = DatabaseConfig{
		Host:     *dbHost,
		Port:     *dbPort,
		User:     *dbUser,
		Password: *dbPassword,
		Database: *dbName,
	}

	cfg.Logging = LoggingConfig{
		Level: *logLevel,
	}

	return cfg
}

func applyEnvOverrides(
	httpAddr *string,
	mcpMode *bool,
	cacheTTL *time.Duration,
	cacheBackend *string,
	redisAddr *string,
	rateLimitDur *time.Duration,
	logLevel *string,
	dbHost *string,
	dbPort *int,
	dbUser *string,
	dbPassword *string,
	dbName *string,
) {
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		*httpAddr = v
	}
	if v := os.Getenv("MCP_MODE"); v == "true" || v == "1" {
		*mcpMode = true
	}
	if v := os.Getenv("CACHE_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			*cacheTTL = d
		}
	}
	if v := os.Getenv("CACHE_BACKEND"); v != "" {
		*cacheBackend = v
	}
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		*redisAddr = v
	}
	if v := os.Getenv("RATE_LIMIT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			*rateLimitDur = d
		}
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		*logLevel = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		*dbHost = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			*dbPort = p
		}
	}
	if v := os.Getenv("DB_USER"); v != "" {
		*dbUser = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		*dbPassword = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		*dbName = v
	}
}
