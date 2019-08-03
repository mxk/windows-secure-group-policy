@echo off

ECHO ===============================================================================
ECHO ===============================================================================
ECHO.
ECHO This script installs restricted traffic baselines into local policy for Windows 10.
ECHO.
ECHO Press Ctrl+C to stop the installation, or press any other key to continue...
PAUSE > nul

ECHO.
ECHO You are about to apply the Windows Restricted Traffic Limited Functionality settings on this device.  For details on what settings are applied please refer to this online article (https://review.docs.microsoft.com/en-us/windows/privacy/manage-connections-from-windows-operating-system-components-to-microsoft-services).
ECHO.
ECHO Do you agree to apply these settings? 
ECHO [Y] Yes  [N] No (default is 'N'):
SET /P reply=
IF /I not "%reply%" == "y" GOTO :End   

:: Make the directory where this script lives the current dir.
PUSHD %~dp0
SET RTGUIDE=%CD%
SET RTGUIDELOGS=%RTGUIDE%\LOGS
MD "%RTGUIDELOGS%" 2> nul

ECHO RestrictedTraffic-install.log > "%RTGUIDELOGS%%\RestrictedTraffic-install.log"
ECHO. >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"
ECHO User agreed to apply the settings >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"

ECHO Checking if LGPO.exe exists in Tools folder ...
ECHO Checking if LGPO.exe exists in Tools folder ...>>"%RTGUIDELOGS%%\RestrictedTraffic-install.log"
IF NOT EXIST .\Tools\LGPO.exe (
	echo.
	ECHO LGPO.exe is not found in .\Tools folder. Failed to apply 'Windows Restricted Traffic Limited Functionality Baseline'.
	ECHO Please check '.\.\Windows Restricted Traffic Limited Functionality Baseline\readme.txt' to install the tool and retry.
	ECHO .
	ECHO LGPO.exe is not found in .\Tools folder. Failed to apply 'Windows Restricted Traffic Limited Functionality Baseline'. >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"
	ECHO Please check '.\.\Windows Restricted Traffic Limited Functionality Baseline\readme.txt' to install the tool and retry. >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"
	ECHO. >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"
	EXIT /B
)

ECHO Installing Windows Server Restricted Traffic settings and policies...
:: Apply Windows Server Restricted Traffic
.\Tools\LGPO.exe /g .\GPO >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log" 2>>&1
echo Windows Server Local Policy Applied

:: Copy Custom Administrative Templates
ECHO Copying custom administrative templates...
copy /y Templates\*.admx %SystemRoot%\PolicyDefinitions >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"
copy /y Templates\*.adml %SystemRoot%\PolicyDefinitions\en-US >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"

:: Uninstall Windows Media Player
:: ECHO Uninstalling Windows Media Player
:: DISM.exe /Online /Disable-Feature /FeatureName:WindowsMediaPlayer >> "%RTGUIDELOGS%%\RestrictedTraffic-install.log"

::Display Notifications
ECHO.
ECHO ===============================================================================
ECHO ===============================================================================
ECHO.
ECHO The Restricted Traffic Limited Functionality settings have been applied successfully
ECHO Please reboot and login with current account.
ECHO.
ECHO Additionally, check log files located in this directory:
ECHO.    %RTGUIDELOGS%
ECHO.
ECHO ===============================================================================
ECHO ===============================================================================
ECHO.
POPD
:End