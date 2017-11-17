@echo off
setlocal

set gpo=GPO\{07BDCD6A-3F72-473C-82B9-67BB69DBE54D}\DomainSysvol\GPO

pushd %~dp0
Tools\LGPO.exe /q /parse /m %gpo%\Machine\Registry.pol > gpo-machine.txt
Tools\LGPO.exe /q /parse /u %gpo%\User\Registry.pol > gpo-user.txt
popd
