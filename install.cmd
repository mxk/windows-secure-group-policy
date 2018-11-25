@echo off
setlocal

if /i "%1" == "/y" goto :install
set /p "confirm=Install GPO as Local Computer Policy (y/[n])? "
if /i "%confirm%" neq "y" goto :eof

:install
set gpo=GPO\{07BDCD6A-3F72-473C-82B9-67BB69DBE54D}\DomainSysvol\GPO

pushd %~dp0
copy /Y %gpo%\GPT.INI %SystemRoot%\System32\GroupPolicy\
Tools\LGPO.exe /g GPO
copy /Y %gpo%\Machine\comment.cmtx %SystemRoot%\System32\GroupPolicy\Machine\
copy /Y %gpo%\User\comment.cmtx %SystemRoot%\System32\GroupPolicy\User\
popd

echo.
echo Restart computer to apply changes!
echo.

if /i "%1" neq "/y" pause
