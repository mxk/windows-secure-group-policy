Restore Win32 Calculator
========================

Server 2016 still has the old calculator installed as a package, but it is
renamed to win32calc.exe. Do not replace calc.exe on Windows 10. Doing so will
cause sfc to report file corruption, which may cause problems for Windows
Update.

To install locally, run `install.cmd`. For a domain, import Win32Calculator.xml
as a computer registry preference, and use `\\<domain>\NETLOGON\win32calc` share
to distribute .exe and .mui files to client computers.
