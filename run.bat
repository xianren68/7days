@echo off
setlocal

:: Build the server binary
go build -o server.exe

:: Start servers in the background
start cmd /k  server.exe -port=8001
start cmd /k  server.exe -port=8002
start /b server.exe -port=8003 -api=1

:: Wait for a few seconds
timeout /t 2 /nobreak

:: Start the tests
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom"
echo \n
curl "http://localhost:9999/api?key=Tom"
echo \n
curl "http://localhost:9999/api?key=Tom"

timeout /t 10
:: Clean up the server binary
del server.exe