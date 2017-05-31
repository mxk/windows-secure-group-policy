Windows 10 and Server 2016 Secure Baseline GPO
==============================================

This GPO contains security settings and administrative policies that should
apply by default to all domain and standalone (via Local Computer Policy)
Windows 10 and Server 2016 computers. These settings aim to provide maximum
privacy, security, and performance (in that order), with minimum network
communication. Security features that send data to Microsoft, such as
SmartScreen, are disabled.

**Use a separate GPO to override settings as needed! Do not use for
site-specific configuration!** This policy should apply to any Windows system,
especially laptops that get connected to insecure networks. Edit only to improve
privacy, security, and performance, or to disable non-essential features added
in new Windows releases.

This GPO overrides most, but not all, settings in the Default Domain
[Controllers] Policy. It also contains unmanaged (preference-type) registry
settings that are not disabled when the GPO goes out of scope. Use filters to
identify. Some of these settings are not covered by any ADMX template (see LGPO
section below).

Installation
------------

### Local Computer Policy

Run `install.cmd`. See LGPO documentation (in Tools) for more info.

Use `parse.cmd` and `build.cmd` to convert Registry.pol files to LGPO text
format and back.

ADMX templates for the local computer are in `C:\Windows\PolicyDefinitions`.
Change the owner to `Administrators` and replace owner on all subcontainers and
objects. After reopening advanced security settings, remove all non-inherited
permissions, enable inheritance, and replace all child permissions. Finally,
copy the PolicyDefinitions directory from the repository, overwriting all
existing files.

**Doing this makes sfc unhappy!** The solution is to run:

```
dism /online /cleanup-image /restorehealth
sfc /scannow
```

DISM requires working NCSI (see notes below) or you get error 0x800f0906. Since
the templates are only required to edit the policy, it's probably better to just
leave them alone.

### Active Directory

1. Create a blank GPO using Group Policy Management Console.
2. Copy the contents of GPO directory to
   `\\<domain>\SYSVOL\<domain>\Policies\<GUID>\`.
3. Make a single edit to security, audit, computer, and user policies. Simply
   toggle one of the configured settings off/on. This is needed to get AD to
   reload the actual settings and update version numbers (please tell me if
   there is a better way to do this - backup/restore doesn't work).
4. Copy PolicyDefinitions to `\\<domain>\SYSVOL\<domain>\Policies\`.

Required ADMX templates
-----------------------

These are contained in the PolicyDefinitions directory:

* Windows 10 and Server 2016 (1607)
* Windows 10 (1703) Creators Update
* MSS (Legacy) from Windows 10 and Server 2016 Security Baseline
* Windows Restricted Traffic Limited Functionality Baseline
* Fixed en-US\W32Time.adml file (v1703 is corrupt; bug reported to Microsoft)

Notes
-----

* Local Administrator is disabled and renamed to LocalAdmin. The password must
  be set before it can be enabled. Otherwise, you'll get error 7016 in
  `gpresult /h` report along with "Security has requested to process its policy
  settings again" message.
* Do not disable "Microsoft Account Sign-in Assistant" (aka "wlidsvc") service
  as suggested by the Restricted Traffic Baseline. Doing so results in error
  0x80070426 when running Windows Update.
* Disabling active Network Connectivity Status Indicator (NCSI) tests breaks
  Windows Update (0x800704cf), DISM (0x800f0906 or 0x800f081f), and probably
  other things. Passive-only configuration doesn't fix this, so best to leave
  both passive and active tests enabled. Windows Update error only happens once
  if Windows was installed without any network connectivity. After a single
  successful update, disabling NCSI doesn't seem to cause further problems, but
  DISM will still run into errors. Guess how many days it took to figure this
  out.
* All traffic to Microsoft servers except for NCSI tests and Windows, Root CA,
  and Defender definition updates is disabled.
* Everything related to Microsoft accounts, Store, Windows Apps, cloud,
  feedback, and customer ~~spying~~ experience improvement is disabled. Some of
  these settings are only supported by Enterprise and Education editions.
* Some functions that are only relevant to domain environments, like dynamic DNS
  updates, app virtualization, remote printing, etc., are disabled.
* NTFS 8dot3name creation is disabled for all volumes, but existing short names
  must be stripped manually (e.g. `fsutil 8dot3name strip /s C:`).
* NTP client is configured to use pool.ntp.org. Local NTP traffic should be
  intercepted/redirected by the firewall.

Suggestions to implement in a separate GPO
------------------------------------------

* Interactive logon: Machine inactivity limit
* Disable UAC on servers
* Firewall rules and logging (intentionally left unconfigured to allow local
  management on standalone systems)
* PKI
* Redirect NCSI DNS/HTTP tests to your own server
* Dynamic DNS updates
* Allow Print Spooler to accept client connections
* Power settings
* Reduce some of the IE 11 & Edge restrictions (these policies were set with the
  assumption that Firefox, Chrome, or another browser is the user's default, and
  IE/Edge is used only on servers or as an occasional backup)
* Remote Shell access
* Automatic Windows updates

LGPO
----

[LGPO.exe v2.0](https://blogs.technet.microsoft.com/secguide/2016/09/23/lgpo-exe-v2-0-pre-release-support-for-mlgpo-and-reg_qword/)
was used to set the following unmanaged registry settings not covered by any
ADMX template. This allows the settings to apply via Local Computer Policy,
which does not support Preferences (without hacks).

```
; Disable shared experiences
User
Software\Microsoft\Windows\CurrentVersion\CDP
RomeSdkChannelUserAuthzPolicy
DWORD:0

; Do not hide extensions for known file types
User
Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced
HideFileExt
DWORD:0

; Disable automatic proxy detection
User
Software\Microsoft\Windows\CurrentVersion\Internet Settings
AutoDetect
DWORD:0
```

References
----------

* https://docs.microsoft.com/en-us/windows/device-security/windows-security-baselines
* https://blogs.technet.microsoft.com/secguide/
* https://blogs.technet.microsoft.com/secguide/2016/10/02/the-mss-settings/
* https://blogs.technet.microsoft.com/secguide/2015/11/18/changes-from-the-windows-8-1-baseline-to-the-windows-10-th11507-baseline/
* https://docs.microsoft.com/en-us/windows/configuration/configure-windows-telemetry-in-your-organization
* https://docs.microsoft.com/en-us/windows/configuration/manage-connections-from-windows-operating-system-components-to-microsoft-services
* http://iase.disa.mil/stigs/os/windows/Pages/index.aspx
* https://www.stigviewer.com/stig/windows_10/
* https://blogs.technet.microsoft.com/josebda/2012/11/13/windows-server-2012-file-server-tip-disable-8-3-naming-and-strip-those-short-names-too/
* https://support.microsoft.com/en-us/help/2526083/disabling-user-account-control-uac-on-windows-server
