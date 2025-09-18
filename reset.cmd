@echo off
setlocal

if /i "%1" == "/y" goto :reset
set /p "confirm=Reset local group policy (y/[n])? "
if /i "%confirm%" neq "y" goto :eof

:reset
rmdir /s /q %SystemRoot%\System32\GroupPolicy\User
rmdir /s /q %SystemRoot%\System32\GroupPolicy\Machine

echo.
echo Run 'gpupdate /force' to apply.
echo.

if /i "%1" neq "/y" pause
