@echo off
setlocal

set vbs=%TEMP%\win32calc.tmp.vbs

pushd %~dp0
copy /Y win32calc.exe %SystemRoot%\System32\
copy /Y win32calc.exe.mui %SystemRoot%\System32\en-US\
regedit /S win32calc.reg

echo Set ws = WScript.CreateObject("WScript.Shell") > %vbs%
echo Set lnk = ws.CreateShortcut("%ProgramData%\Microsoft\Windows\Start Menu\Programs\Accessories\Calculator.lnk") >> %vbs%
echo lnk.TargetPath = "%SystemRoot%\System32\win32calc.exe" >> %vbs%
echo lnk.Description = "Performs basic arithmetic tasks with an on-screen calculator." >> %vbs%
echo lnk.IconLocation = "%SystemRoot%\System32\win32calc.exe,0" >> %vbs%
echo lnk.Save >> %vbs%
cscript %vbs%
del %vbs%

popd
pause
