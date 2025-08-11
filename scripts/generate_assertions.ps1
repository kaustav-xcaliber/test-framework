# Assertion Generator Script Wrapper for PowerShell
# Usage: .\scripts\generate_assertions.ps1 <json_file> [status_code]

param(
    [Parameter(Mandatory=$true)]
    [string]$JsonFile,
    
    [Parameter(Mandatory=$false)]
    [int]$StatusCode
)

# Function to print colored output
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Check if Go is installed
try {
    $goVersion = go version 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "Go not found"
    }
} catch {
    Write-Error "Go is not installed. Please install Go to use this script."
    exit 1
}

# Check if JSON file exists
if (-not (Test-Path $JsonFile)) {
    Write-Error "JSON file '$JsonFile' not found."
    exit 1
}

Write-Info "Generating assertions from '$JsonFile'"
if ($StatusCode) {
    Write-Info "Status code: $StatusCode"
}

Write-Host ""

# Run the Go script
if ($StatusCode) {
    go run scripts/generate_assertions.go $JsonFile $StatusCode
} else {
    go run scripts/generate_assertions.go $JsonFile
}

Write-Host ""
Write-Success "Assertion generation completed!"
Write-Info "Copy the JSON array above and use it in your test spec."
