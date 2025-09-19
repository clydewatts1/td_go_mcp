package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"td_go_mcp/internal/db"
	"td_go_mcp/internal/tools"

	"golang.org/x/exp/slog"
)

var (
	loadedTools      []tools.ToolDefinition
	loadedPrompts    []tools.PromptDefinition
	loadedGlossaries []tools.GlossaryResource
	processors       map[string]*tools.SQLProcessor
	database         *db.DB
	logger           *slog.Logger
	basePath         string
)

func init() {
	// Determine the base path of the executable
	exePath, err := os.Executable()
	if err != nil {
		slog.Error("Failed to get executable path", "err", err)
		os.Exit(1)
	}
	// Check if running with "go run"
	if strings.Contains(exePath, "go-build") || strings.Contains(exePath, "exe\\main.exe") {
		// Likely running with "go run", use current working directory
		wd, err := os.Getwd()
		if err != nil {
			slog.Error("Failed to get working directory", "err", err)
			os.Exit(1)
		}
		basePath = wd
	} else {
		// Running as a compiled binary, use the executable's directory
		basePath = filepath.Dir(exePath)
	}

	// Set up logging directory and slog logger
	logDir := filepath.Join(basePath, "logging")
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

	// Load glossary resources from tools
	glossaryPath := filepath.Join(basePath, "tools")
	loadedGlossaries, err = tools.LoadGlossaryResourcesFromDirectory(glossaryPath)
	if err != nil {
		logger.Error("Error loading glossary resources", "err", err, "path", glossaryPath)
		loadedGlossaries = []tools.GlossaryResource{}
	}
	logger.Info("Loaded glossary resources", "glossaries", len(loadedGlossaries))

	toolsPath := filepath.Join(basePath, "tools")
	loadedTools, err = tools.LoadToolsFromDirectory(toolsPath)
	if err != nil {
		logger.Error("Error loading tools", "err", err, "path", toolsPath)
		loadedTools = []tools.ToolDefinition{} // Continue with empty tools
	}

	processors = make(map[string]*tools.SQLProcessor)
	for i := range loadedTools {
		processors[loadedTools[i].Name] = tools.NewSQLProcessor(loadedTools[i])
	}

	// Load prompts from YAML files
	promptsPath := filepath.Join(basePath, "prompts")
	loadedPrompts, err = tools.LoadPromptsFromDirectory(promptsPath)
	if err != nil {
		logger.Error("Error loading prompts", "err", err, "path", promptsPath)
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
