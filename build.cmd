@echo off

pushd %~dp0
Tools\LGPO.exe /q /r gpo-machine.txt /w GPO\Machine\Registry.pol && Tools\multisz-fix.exe GPO\Machine\Registry.pol
Tools\LGPO.exe /q /r gpo-user.txt /w GPO\User\Registry.pol
popd
