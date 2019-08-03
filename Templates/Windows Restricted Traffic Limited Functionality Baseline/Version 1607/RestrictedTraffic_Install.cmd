@echo off
ECHO ===============================================================================
ECHO ===============================================================================
ECHO.
ECHO This script installs restricted traffic baselines into local policy for Windows 10.
ECHO.
ECHO Press Ctrl+C to stop the installation, or press any other key to continue...
ECHO.
ECHO ===============================================================================
ECHO ===============================================================================
PAUSE > nul

:: Make the directory where this script lives the current dir.
PUSHD %~dp0
SET RTGUIDE=%CD%
SET RTGUIDELOGS=%RTGUIDE%\LOGS
MD "%RTGUIDELOGS%" 2> nul

ECHO RestrictedTraffic-install.log > "%RTGUIDELOGS%%\RestrictedTraffic-install.log"
ECHO. >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"

ECHO Checking if LGPO.exe exists in Tools folder ...
IF NOT EXIST .\Tools\LGPO.exe (
	echo.
	ECHO LGPO.exe is not found in .\Tools folder. Failed to apply 'Windows Restricted Traffic Limited Functionality Baseline'.
	ECHO Please check Tools\readme.txt to install the tool and retry.
	ECHO .
	ECHO LGPO.exe is not found in .\Tools folder. Failed to apply 'Windows Restricted Traffic Limited Functionality Baseline'. >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"
	ECHO Please check Tools\readme.txt to install the tool and retry. >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"
	ECHO. >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"
	EXIT /B
)

ECHO Installing Windows 10 Restricted Traffic settings and policies...
:: Apply Windows 10 Restricted Traffic
.\Tools\LGPO.exe /g .\GPO >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log" 2>>&1
echo Windows 10 Local Policy Applied


:: Copy Custom Administrative Templates
ECHO Copying custom administrative templates
copy /y Templates\*.admx %SystemRoot%\PolicyDefinitions >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"
copy /y Templates\*.adml %SystemRoot%\PolicyDefinitions\en-US >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"


:: Uninstall Windows Media Player
:: ECHO Uninstalling Windows Media Player
:: DISM.exe /Online /Disable-Feature /FeatureName:WindowsMediaPlayer >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"


:: Uninstall Modern Applications
ECHO Uninstalling Modern Applications
PowerShell.exe -ExecutionPolicy Bypass -File RemoveMA.ps1 >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"


::Display Notifications
ECHO.
ECHO.
ECHO ===============================================================================
ECHO ===============================================================================
ECHO.
ECHO In order to test properly, reboot and login with current account.
ECHO.
ECHO Additionally, check log files located in this directory:
ECHO.    %RTGUIDELOGS%
ECHO.
ECHO ===============================================================================
ECHO ===============================================================================
ECHO.
POPD

