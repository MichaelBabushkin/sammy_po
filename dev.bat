@echo off
REM filepath: c:\Users\misha\OneDrive\Desktop\sammy-po\dev.bat
ECHO Setting up development environment with Air for hot reloading

ECHO Installing Air if not already installed...
go install github.com/air-verse/air@latest

ECHO Ensuring Go modules are tidy...
go mod tidy

ECHO Starting Go server with hot reloading...
air -c .air.toml