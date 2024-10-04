@echo off
:: Title for the terminal window
title Mz-Brute Compiler - MH-ProDev ItzK4sra

:: Set variables for color and formatting
color 0A

:: Display title and prompt
echo ========================================
echo         Mz-Brute Script
echo         MH-ProDev ItzK4sra
echo ========================================

:chooseOS
:: Prompt user for OS choice with validation
echo Choose the OS to compile for:
echo 1 - Linux
echo 2 - Windows
echo ========================================
set /p os_choice="Enter your choice (1 or 2): "

if "%os_choice%"=="1" (
    set "GOOS=linux"
    set "GOARCH=amd64"
    set "OUTPUT=Mz-Brute-linux"
    echo Compiling for Linux...
) else if "%os_choice%"=="2" (
    set "GOOS=windows"
    set "GOARCH=amd64"
    set "OUTPUT=Mz-Brute.exe"
    echo Compiling for Windows...
) else (
    echo Invalid choice. Please enter 1 or 2.
    goto chooseOS
)

:: Compile the Go program
echo ========================================
echo Running: go build -o %OUTPUT% Mz-Brute.go
go build -o %OUTPUT% Mz-Brute.go

:: Check if the build was successful
if "%ERRORLEVEL%"=="0" (
    echo ========================================
    echo Build successful!
    echo Output file: %OUTPUT%
    echo ========================================
) else (
    echo ========================================
    echo Build failed. Please check your Go file for errors.
    echo ========================================
)

pause
