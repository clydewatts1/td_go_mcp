# E2E Test Script for td-go-mcp using Gemini CLI

# Exit on first error
$ErrorActionPreference = 'Stop'

# --- Test Configuration ---
$ServerName = "td-go-mcp"
$GeminiExe = "gemini" # Assuming 'gemini' is in the PATH

# --- Helper Functions ---
function Invoke-GeminiTool {
    param(
        [string]$ToolName,
        [string]$Parameters
    )
    $Command = "$GeminiExe @$ServerName $ToolName $Parameters"
    Write-Host "Executing: $Command"
    try {
        $output = Invoke-Expression $Command
        return $output | ConvertFrom-Json
    }
    catch {
        Write-Error "Failed to execute Gemini CLI command: $_"
        return $null
    }
}

function Assert-Condition {
    param(
        [string]$TestName,
        [bool]$Condition
    )
    if ($Condition) {
        Write-Host "✅ PASS: $TestName"
    }
    else {
        Write-Error "❌ FAIL: $TestName"
    }
}

# --- Test Cases ---

# Test 1: List Databases
Write-Host "`n--- Testing 'list_databases' ---"
$result = Invoke-GeminiTool -ToolName "list_databases" -Parameters "--database_name %"
Assert-Condition -TestName "list_databases should return a non-null result" -Condition ($null -ne $result)
# Add more specific assertions here based on expected output

# Test 2: Count Records (Example)
# Note: This test requires a valid table name and will likely fail without a real database connection.
Write-Host "`n--- Testing 'count_records' ---"
$result = Invoke-GeminiTool -ToolName "count_records" -Parameters "--table_name dbc.dbcinfo"
Assert-Condition -TestName "count_records should return a non-null result" -Condition ($null -ne $result)
if ($null -ne $result) {
    Assert-Condition -TestName "count_records result should have a 'count' property" -Condition ($result.PSObject.Properties.Name -contains 'count')
}

# Test 3: Get User by ID (Example)
Write-Host "`n--- Testing 'get_user_by_id' ---"
$result = Invoke-GeminiTool -ToolName "get_user_by_id" -Parameters "--user_id 123"
Assert-Condition -TestName "get_user_by_id should return a non-null result" -Condition ($null -ne $result)
if ($null -ne $result) {
    Assert-Condition -TestName "get_user_by_id result should have a 'user' property" -Condition ($result.PSObject.Properties.Name -contains 'user')
}


# Test 4: List Active Sessions (Example)
Write-Host "`n--- Testing 'list_active_sessions' ---"
$result = Invoke-GeminiTool -ToolName "list_active_sessions" -Parameters ""
Assert-Condition -TestName "list_active_sessions should return a non-null result" -Condition ($null -ne $result)
if ($null -ne $result) {
    Assert-Condition -TestName "list_active_sessions result should have a 'sessions' property" -Condition ($result.PSObject.Properties.Name -contains 'sessions')
}


Write-Host "`nAll tests completed."
