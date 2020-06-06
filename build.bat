@echo off

REM Clean up old build artifacts.
del stats.exe stats

REM Build the server for Windows.
set GOOS=windows
set GOARCH=amd64
go build -o stats.exe

REM Build the server for Linux.
set GOOS=linux
set GOARCH=amd64
go build -o stats

REM Back to Windows for next build.
set GOOS=windows