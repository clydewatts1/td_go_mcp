package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/alexbrainman/odbc"
)

type Config struct {
	Driver           string `yaml:"driver"`
	ConnectionString string `yaml:"connection_string"`
	DSN              string `yaml:"dsn"`
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
	Database         string `yaml:"database"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
}

func LoadConfig() *Config {
	config := &Config{
		Driver: "odbc",
		DSN:    "teradw", // Default Teradata DSN
	}

	// Override from environment variables
	if driver := os.Getenv("DB_DRIVER"); driver != "" {
		config.Driver = driver
	}
	if dsn := os.Getenv("DB_DSN"); dsn != "" {
		config.DSN = dsn
	}
	if connStr := os.Getenv("DB_CONNECTION_STRING"); connStr != "" {
		config.ConnectionString = connStr
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.Port = port
	}
	if database := os.Getenv("DB_DATABASE"); database != "" {
		config.Database = database
	}
	if username := os.Getenv("DB_USERNAME"); username != "" {
		config.Username = username
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Password = password
	}

	return config
}

func (c *Config) GetConnectionString() string {
	if c.ConnectionString != "" {
		return c.ConnectionString
	}

	if c.DSN != "" {
		return fmt.Sprintf("dsn=%s", c.DSN)
	}

	// Build connection string from components
	var parts []string
	if c.Host != "" {
		parts = append(parts, fmt.Sprintf("server=%s", c.Host))
	}
	if c.Port != "" {
		parts = append(parts, fmt.Sprintf("port=%s", c.Port))
	}
	if c.Database != "" {
		parts = append(parts, fmt.Sprintf("database=%s", c.Database))
	}
	if c.Username != "" {
		parts = append(parts, fmt.Sprintf("uid=%s", c.Username))
	}
	if c.Password != "" {
		parts = append(parts, fmt.Sprintf("pwd=%s", c.Password))
	}

	return strings.Join(parts, ";")
}

type DB struct {
	conn   *sql.DB
	config *Config
}

func Connect(config *Config) (*DB, error) {
	connStr := config.GetConnectionString()
	conn, err := sql.Open(config.Driver, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn, config: config}, nil
}

func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func (db *DB) ExecuteQuery(query string) ([]map[string]interface{}, error) {
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}
