package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/johnrirwin/mcp-news-feed/internal/aggregator"
	"github.com/johnrirwin/mcp-news-feed/internal/cache"
	"github.com/johnrirwin/mcp-news-feed/internal/database"
	"github.com/johnrirwin/mcp-news-feed/internal/equipment"
	"github.com/johnrirwin/mcp-news-feed/internal/httpapi"
	"github.com/johnrirwin/mcp-news-feed/internal/inventory"
	"github.com/johnrirwin/mcp-news-feed/internal/logging"
	"github.com/johnrirwin/mcp-news-feed/internal/mcp"
	"github.com/johnrirwin/mcp-news-feed/internal/ratelimit"
	"github.com/johnrirwin/mcp-news-feed/internal/sellers"
	"github.com/johnrirwin/mcp-news-feed/internal/sources"
	"github.com/johnrirwin/mcp-news-feed/internal/tagging"
)

func main() {
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

	// Environment variable overrides
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
	if v := os.Getenv("DB_USER"); v != "" {
		*dbUser = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		*dbPassword = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		*dbName = v
	}

	level := logging.LevelInfo
	switch *logLevel {
	case "debug":
		level = logging.LevelDebug
	case "warn":
		level = logging.LevelWarn
	case "error":
		level = logging.LevelError
	}
	logger := logging.New(level)

	// Initialize cache backend
	var feedCache cache.Cache
	switch *cacheBackend {
	case "redis":
		logger.Info("Using Redis cache backend", logging.WithField("addr", *redisAddr))
		redisCache, err := cache.NewRedis(cache.RedisConfig{
			Addr:   *redisAddr,
			Prefix: "mcp-news:",
		}, *cacheTTL)
		if err != nil {
			logger.Error("Failed to connect to Redis, falling back to memory cache", logging.WithField("error", err.Error()))
			feedCache = cache.NewMemory(*cacheTTL)
		} else {
			feedCache = redisCache
		}
	default:
		logger.Info("Using in-memory cache backend")
		feedCache = cache.NewMemory(*cacheTTL)
	}

	limiter := ratelimit.New(*rateLimitDur)
	tagger := tagging.New()

	config := sources.DefaultConfig()
	var fetchers []sources.Fetcher
	fetchers = append(fetchers, sources.CreateDroneRSSFetchers(limiter, config)...)
	fetchers = append(fetchers, sources.CreateDroneRedditFetchers(limiter, config)...)

	agg := aggregator.New(fetchers, feedCache, tagger, logger)

	// Initialize seller registry with adapters
	sellerRegistry := sellers.NewRegistry()
	sellerRegistry.Register(sellers.NewRaceDayQuads(limiter, feedCache))
	sellerRegistry.Register(sellers.NewGetFPV(limiter, feedCache))
	logger.Info("Registered seller adapters", logging.WithField("count", len(sellerRegistry.List())))

	// Initialize equipment service
	equipmentSvc := equipment.NewService(sellerRegistry, feedCache, logger)

	// Initialize inventory service (in-memory for now, can upgrade to PostgreSQL)
	var inventorySvc inventory.InventoryManager

	// Try to connect to PostgreSQL for inventory persistence
	dbConfig := database.Config{
		Host:     *dbHost,
		Port:     *dbPort,
		User:     *dbUser,
		Password: *dbPassword,
		Database: *dbName,
	}

	db, err := database.New(dbConfig)
	if err != nil {
		logger.Warn("Failed to connect to PostgreSQL, using in-memory inventory", logging.WithField("error", err.Error()))
		inventorySvc = inventory.NewInMemoryService(logger)
	} else {
		logger.Info("Connected to PostgreSQL for inventory persistence")
		if err := db.Migrate(context.Background()); err != nil {
			logger.Warn("Failed to run migrations, using in-memory inventory", logging.WithField("error", err.Error()))
			inventorySvc = inventory.NewInMemoryService(logger)
		} else {
			inventoryStore := database.NewInventoryStore(db)
			inventorySvc = inventory.NewService(inventoryStore, logger)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Info("Shutting down...")
		cancel()
	}()

	if *mcpMode {
		logger.Info("Starting MCP server in stdio mode")
		mcpHandler := mcp.NewHandler(agg, equipmentSvc, inventorySvc, logger)
		mcpServer := mcp.NewServer(mcpHandler, logger)

		logger.Info("Pre-fetching feeds...")
		if err := agg.Refresh(ctx); err != nil {
			logger.Warn("Initial fetch had errors", logging.WithField("error", err.Error()))
		}

		if err := mcpServer.Run(ctx); err != nil && err != context.Canceled {
			logger.Error("MCP server error", logging.WithField("error", err.Error()))
			os.Exit(1)
		}
	} else {
		logger.Info("Starting HTTP server", logging.WithField("addr", *httpAddr))

		httpServer := httpapi.New(agg, equipmentSvc, inventorySvc, logger)

		go func() {
			logger.Info("Pre-fetching feeds in background...")
			if err := agg.Refresh(ctx); err != nil {
				logger.Warn("Initial fetch had errors", logging.WithField("error", err.Error()))
			}
			logger.Info("Initial fetch complete")
		}()

		if err := httpServer.Start(*httpAddr); err != nil {
			logger.Error("HTTP server error", logging.WithField("error", err.Error()))
			os.Exit(1)
		}
	}
}
