Restore Win32 Calculator
========================

Server 2016 still has the old calculator installed as a package, but it is
renamed to win32calc.exe. Do not replace calc.exe on Windows 10. Doing so will
cause sfc to report file corruption, which may cause problems for Windows
Update.

1. Copy win32calc.exe to `%SystemRoot%\System32\`.
2. Copy win32calc.exe.mui to `%SystemRoot%\System32\en-US`.
3. Run the .reg file or import the .xml file as a computer group policy registry
   preference.
4. Create a shortcut in
   `%ProgramData%\Microsoft\Windows\Start Menu\Programs\Accessories`.

Use `\\<domain>\NETLOGON\win32calc` share to distribute .exe and .mui to domain
clients.
