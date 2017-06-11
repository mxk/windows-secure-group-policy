@echo off
setlocal

set /p "confirm=Install GPO as Local Computer Policy (y/[n])? "
if /i "%confirm%" neq "y" goto :eof

:: Extensions should be kept in sync with Backup.xml
pushd %~dp0
Tools\LGPO.exe /g GPO ^
	/e {169EBF44-942F-4C43-87CE-13C93996EBBE} ^
	/e {29BBE2D5-DE47-4855-97D7-2745E166DC6D} ^
	/e {2BFCC077-22D2-48DE-BDE1-2F618D9B476D} ^
	/e {426031C0-0B47-4852-B0CA-AC3D37BFCB39} ^
	/e {7933F41E-56F8-41D6-A31C-4148A711EE93} ^
	/e {C631DF4C-088F-4156-B058-4375F0853CD8} ^
	/e {CDEAFC3D-948D-49DD-AB12-E578BA4AF7AA} ^
	/e audit
popd
pause
