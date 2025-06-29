@echo off
echo Building go_ImagePreviewer for Windows (no console)...

REM Windows用にコンソールなしでビルド
go build -ldflags="-s -w -H=windowsgui" -trimpath -o go_ImagePreviewer.exe

if %ERRORLEVEL% EQU 0 (
    echo Build successful: go_ImagePreviewer.exe
) else (
    echo Build failed
)

pause
