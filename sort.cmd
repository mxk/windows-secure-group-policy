@echo off
setlocal

pushd %~dp0
Tools\gpx.exe sort gpo-machine.txt gpo-user.txt
popd
