@echo off
setlocal

set /p "confirm=Install GPO as Local Computer Policy (y/[n])? "
if /i "%confirm%" neq "y" goto :eof

pushd %~dp0
Tools\LGPO.exe /g GPO
popd
pause
