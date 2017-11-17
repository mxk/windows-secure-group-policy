@echo off
setlocal

set dir=%~dp0
set dir=%dir:~,-1%
set build=cmd /c go build -gcflags "-trimpath=%dir:\=/%" -ldflags "-s -w"

pushd %~dp0
%build% lgpo-sort.go
%build% multisz-fix.go
popd
