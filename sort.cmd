@echo off
setlocal

pushd %~dp0
Tools\lgpo-sort.exe gpo-machine.txt gpo-user.txt
popd
