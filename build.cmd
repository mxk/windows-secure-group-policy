@echo off
setlocal

set gpo=GPO\{07BDCD6A-3F72-473C-82B9-67BB69DBE54D}\DomainSysvol\GPO
set machine=%gpo%\Machine\Registry.pol
set user=%gpo%\User\Registry.pol

pushd %~dp0
Tools\LGPO.exe /q /r gpo-machine.txt /w %machine%
Tools\LGPO.exe /q /r gpo-user.txt /w %user%
Tools\gpx.exe polfix %machine% %user%
popd
