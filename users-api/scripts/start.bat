@echo off
REM Start script for users-api (Windows)
REM This script starts the users-api service in production mode

echo Starting users-api...

REM Check if .env exists and load it
if exist .env (
    echo Loading environment variables from .env
    for /f "usebackq tokens=*" %%a in (".env") do (
        set %%a
    )
)

REM Build the application
echo Building application...
go build -o bin\users-api.exe .\cmd\server\main.go

REM Run the application
echo Starting server on port 8080...
.\bin\users-api.exe
