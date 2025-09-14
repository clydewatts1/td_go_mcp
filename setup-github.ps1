# GitHub Setup Script
# Run this after creating your GitHub repository

# Instructions:
# 1. Go to https://github.com and create a new repository named 'td_go_mcp'
# 2. Don't initialize with README, .gitignore, or license
# 3. Replace YOUR_USERNAME below with your actual GitHub username
# 4. Run this script

param(
    [Parameter(Mandatory=$true)]
    [string]$GitHubUsername
)

$repoUrl = "https://github.com/$GitHubUsername/td_go_mcp.git"

Write-Host "Setting up GitHub repository..." -ForegroundColor Green
Write-Host "Repository URL: $repoUrl" -ForegroundColor Cyan

try {
    # Add remote origin
    git remote add origin $repoUrl
    Write-Host "✓ Added remote origin" -ForegroundColor Green
    
    # Rename branch to main
    git branch -M main
    Write-Host "✓ Renamed branch to main" -ForegroundColor Green
    
    # Push to GitHub
    git push -u origin main
    Write-Host "✓ Pushed to GitHub" -ForegroundColor Green
    
    Write-Host "" 
    Write-Host "Repository successfully created and pushed to GitHub!" -ForegroundColor Green
    Write-Host "View your repository at: https://github.com/$GitHubUsername/td_go_mcp" -ForegroundColor Cyan
}
catch {
    Write-Error "Failed to set up GitHub repository: $_"
    Write-Host "Please check that:"
    Write-Host "1. You created the repository on GitHub first"
    Write-Host "2. Your GitHub username is correct"
    Write-Host "3. You have push permissions to the repository"
}