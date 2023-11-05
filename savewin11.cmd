@echo off
setlocal

pushd %~dp0
.\savelocal.cmd .\PolicyRules\Win11.PolicyRules "Windows 11 Secure Group Policy"
popd
