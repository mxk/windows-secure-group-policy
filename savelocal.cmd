@echo off
setlocal

if "%~1" == "" goto :usage
if "%~2" == "" goto :usage
goto :main

:usage
echo usage: %0 out-file policy-name
echo.
goto :eof

:main
pushd %~dp0
rmdir /s /q C:\GPO
mkdir C:\GPO
.\LGPO\LGPO.exe /b C:\GPO /n "%~2"
move C:\GPO\{*} C:\GPO\{00000000-0000-0000-0000-000000000000}
secedit /export /cfg "C:\GPO\{00000000-0000-0000-0000-000000000000}\DomainSysvol\GPO\Machine\microsoft\windows nt\SecEdit\GptTmpl.inf"
copy /y "%SystemRoot%\System32\GroupPolicy\Machine\Microsoft\Windows NT\Audit\audit.csv" "C:\GPO\{00000000-0000-0000-0000-000000000000}\DomainSysvol\GPO\Machine\microsoft\windows nt\Audit\"
.\PolicyAnalyzer\GPO2PolicyRules.exe C:\GPO "%~1"
popd
