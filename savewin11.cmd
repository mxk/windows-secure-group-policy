@echo off
setlocal

pushd %~dp0
.\save.cmd .\PolicyRules\Win11-Local.PolicyRules "Windows 11 Secure Group Policy"
popd
