package main

// list_databases.yaml
type ListDatabasesArgs struct {
	DatabaseName string `json:"database_name" jsonschema:"required,description=Database name with wildcard (e.g. % for all databases)"`
	DatabaseType string `json:"database_type" jsonschema:"description=User type (D=Database, U=User)"`
}

// count_records.yaml
type CountRecordsArgs struct {
	TableName  string `json:"table_name" jsonschema:"required,description=Name of the table to count records from"`
	DateFilter string `json:"date_filter,omitempty" jsonschema:"description=Optional date filter (YYYY-MM-DD format)"`
	Status     string `json:"status,omitempty" jsonschema:"description=Optional status filter"`
}

// get_user_by_id.yaml
type GetUserByIDArgs struct {
	UserID         string `json:"user_id" jsonschema:"required,description=The unique identifier for the user"`
	IncludeDetails bool   `json:"include_details,omitempty" jsonschema:"description=Include detailed user information"`
}

// list_active_sessions.yaml
type ListActiveSessionsArgs struct {
	UserType string `json:"user_type,omitempty" jsonschema:"description=Filter by user type (admin, user, guest)"`
	Limit    int    `json:"limit,omitempty" jsonschema:"description=Maximum number of sessions to return"`
}

// list_of_tables.yaml
type ListOfTablesArgs struct {
	DatabaseName string `json:"database_name" jsonschema:"required,description=Name of the database to list tables from"`
}
