@echo off

pushd %~dp0
Tools\LGPO.exe /q /parse /m GPO\Machine\Registry.pol > gpo-machine.txt
Tools\LGPO.exe /q /parse /u GPO\User\Registry.pol > gpo-user.txt
popd
