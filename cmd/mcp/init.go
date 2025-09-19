package main

import (
	"os"
	"path/filepath"
	"time"

	"td_go_mcp/internal/db"
	"td_go_mcp/internal/tools"

	"golang.org/x/exp/slog"
)

var (
	loadedTools   []tools.ToolDefinition
	loadedPrompts []tools.PromptDefinition
	processors    map[string]*tools.SQLProcessor
	database      *db.DB
	logger        *slog.Logger
)

func init() {
	// Set up logging directory and slog logger
	logDir := "logging"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		slog.Error("Failed to create logging directory", "err", err)
		os.Exit(1)
	}
	logFile := filepath.Join(logDir, time.Now().Format("2006-01-02_15-04-05")+".log")
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("Failed to open log file", "err", err)
		os.Exit(1)
	}
	logger = slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(logger)

	loadedTools, err = tools.LoadToolsFromDirectory("tools")
	if err != nil {
		logger.Error("Error loading tools", "err", err)
		loadedTools = []tools.ToolDefinition{} // Continue with empty tools
	}

	processors = make(map[string]*tools.SQLProcessor)
	for i := range loadedTools {
		processors[loadedTools[i].Name] = tools.NewSQLProcessor(loadedTools[i])
	}

	// Load prompts from YAML files
	loadedPrompts, err = tools.LoadPromptsFromDirectory("tools")
	if err != nil {
		logger.Error("Error loading prompts", "err", err)
		loadedPrompts = []tools.PromptDefinition{} // Continue with empty prompts
	}

	logger.Info("Loaded tools and prompts", "tools", len(loadedTools), "prompts", len(loadedPrompts))

	// Initialize database connection
	dbConfig := db.LoadConfig()
	database, err = db.Connect(dbConfig)
	if err != nil {
		logger.Error("Database connection failed", "err", err)
		logger.Warn("Continuing without database - SQL preview mode only")
		database = nil
	} else {
		logger.Info("Database connection established successfully")
	}
}
