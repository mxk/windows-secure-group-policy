@echo off
setlocal

set GPO=GPO\{07BDCD6A-3F72-473C-82B9-67BB69DBE54D}\DomainSysvol\GPO

pushd %~dp0
Tools\LGPO.exe /q /r gpo-machine.txt /w %GPO%\Machine\Registry.pol && Tools\multisz-fix.exe %GPO%\Machine\Registry.pol
Tools\LGPO.exe /q /r gpo-user.txt /w %GPO%\User\Registry.pol
popd
