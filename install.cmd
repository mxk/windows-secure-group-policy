@echo off
setlocal

if /i "%1" == "/y" goto :install
set /p "confirm=Install GPO as Local Computer Policy (y/[n])? "
if /i "%confirm%" neq "y" goto :eof

:install
pushd %~dp0
.\LGPO\LGPO.exe /p .\PolicyRules\Win11.PolicyRules /v
popd

echo.
echo Restart computer to apply changes!
echo.

if /i "%1" neq "/y" pause
