Windows 10 and Server 2019 Secure Baseline GPO
==============================================

This is a baseline group policy for domain and standalone Windows 10 and Server
2016/2019 computers that aims to provide maximum privacy, security, and
performance, in that order. Security features that send data to Microsoft, such
as SmartScreen, are disabled. Some settings are only effective on the Enterprise
edition. The following network traffic is allowed:

* Network Connectivity Status Indicator (NCSI) tests (see Notes below).
* Windows updates.
* Root CA updates.
* Defender definition updates.
* Storage Health (Disk Failure Prediction Model) updates.

In a domain environment, this should be the first GPO applied to any Windows
system, with any site-specific settings configured in separate GPOs.

This policy overrides most, but not all, settings in the Default Domain
[Controllers] Policy. It also contains unmanaged (preference-type) registry
settings that are not disabled when the GPO goes out of scope. Use filters to
identify. Some of these settings are not covered by any ADMX template (see LGPO
section below).

Additional tweaks and fixes are located in the Calculator and RegTweaks
directories.

Installation
------------

### Local Computer Policy

Run `install.cmd`. See LGPO documentation (in Tools) for more info.

Use `parse.cmd` and `build.cmd` to convert Registry.pol files to LGPO text
format and back. Use `sort.cmd` to sort entries in LGPO text files for easier
diffs. Registry.pol files are usually sorted by keys, but not by values. This
makes it difficult to track changes. Sorting LGPO text files and rebuilding
Registry.pol files ensures that all entries are always in the same order.

ADMX templates for the local computer are in `C:\Windows\PolicyDefinitions`. To
update existing files in this directory, you must change the owner to
`Administrators` and replace owner on all subcontainers and objects. However,
**doing this makes sfc unhappy!** It's safer to leave existing templates alone
and only add missing ones.

If you do change the owner, open advanced security settings again, remove all
non-inherited permissions, enable inheritance, and replace all child
permissions. Finally, copy the PolicyDefinitions directory from the repository,
overwriting all existing files. To fix SFC, run:

```
dism /online /cleanup-image /restorehealth
sfc /scannow
```

DISM requires working NCSI (see notes below) or you'll get error 0x800f0906.
Since the templates are only required to edit the policy, it's better to just
leave them alone. Do a clean Windows install in a VM when a new version is
released and update the policy there instead of messing with local templates.

### Domain

1. Copy PolicyDefinitions directory to `\\<domain>\SYSVOL\<domain>\Policies\`.
   If this directory already exists, delete or rename it first.
2. Create a blank GPO using Group Policy Management Console.
3. Right-click on the new GPO and go through the "Import Settings..." wizard
   using the GPO directory as the backup source.

Required ADMX templates
-----------------------

These are contained in the PolicyDefinitions directory:

1. [Windows 10 and Server 2016 (1607)](https://www.microsoft.com/en-us/download/details.aspx?id=53430)
2. Windows 10 (1903) Clean Install (`C:\Windows\PolicyDefinitions`)
3. [MSS (Legacy) from Windows 10 and Windows Server Security Baseline (1903)](https://blogs.technet.microsoft.com/secguide/2019/05/23/security-baseline-final-for-windows-10-v1903-and-windows-server-v1903/)
4. [Windows Restricted Traffic Limited Functionality Baseline (1903)](https://docs.microsoft.com/en-us/windows/privacy/manage-connections-from-windows-operating-system-components-to-microsoft-services)

When an update is released, the entire PolicyDefinitions directory should be
rebuilt by copying templates over in the listed order. Copying updated ADMX/ADML
files without removing old ones first may cause problems.

* [How to create and manage the Central Store for Group Policy Administrative Templates in Windows](https://support.microsoft.com/en-us/help/3087759/how-to-create-and-manage-the-central-store-for-group-policy-administra)
* [Group Policy Settings Reference Spreadsheet Windows 1809](https://www.microsoft.com/en-us/download/details.aspx?id=57464)
* [New policies for Windows 10](https://docs.microsoft.com/en-us/windows/client-management/new-policies-for-windows-10)
* [ADMX Version History](https://blogs.technet.microsoft.com/grouppolicy/2016/10/12/admx-version-history/)

Notes
-----

* [Domain Level Account Policies](https://technet.microsoft.com/en-us/library/jj852214(v=ws.11).aspx)
  (Password Policy, Account Lockout Policy, and Kerberos Policy) only apply if
  the GPO is linked at the domain level.
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
* Everything related to Microsoft accounts, Store, Windows Apps, cloud,
  feedback, and customer ~~spying~~ experience improvement is disabled. Some of
  these settings are only supported by Enterprise and Education editions.
* Some functions that are only relevant to domain environments, like dynamic DNS
  updates, app virtualization, remote printing, etc., are disabled.
* NTFS 8dot3name creation is disabled for all volumes. Existing short names can
  be stripped manually (e.g. `fsutil 8dot3name strip /s C:`), though it's safer
  to leave them alone for a new installation.
* NTP client is configured to use pool.ntp.org. Local NTP traffic should be
  intercepted/redirected by the firewall.
* "Turn off all Windows spotlight features" policy must be applied within 15
  minutes after Windows 10 is installed.

Bugs
----

* Denying "Let Windows apps run in the background" right (Windows Components ->
  App Privacy) [breaks Start menu search in v1703](https://superuser.com/a/1208858).
* Disabling "Allow Extensions" in Microsoft Edge settings prevents Edge from
  starting for the first time in v1709. If Edge is started manually, the window
  closes immediately. If Edge is started by opening an html file, you get "The
  remote procedure call failed" error message. This does not happen on
  subsequent launches, so just allow extensions temporarily when starting Edge
  for the first time.
* Enabling "Configure H.264/AVC hardware encoding for Remote Desktop
  Connections" breaks Hyper-V Enhanced Sessions. The resulting error is "The
  session was disconnected. If you want to continue, try to connect again. If
  the problem persists, contact your system administrator. Would you like to try
  to reconnect?" To switch to a basic session, right-click on the VM, Edit
  Session Settings, then just close the "Connect to <VM>" dialog.

Suggestions to implement in a separate GPO
------------------------------------------

* Automatic Windows updates
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

LGPO
----

Tools\LGPO.exe was used to set the following unmanaged registry settings not
covered by any ADMX template. This allows the settings to apply via Local
Computer Policy, which does not support Preferences (without hacks).

```
; Disable Inking & Typing Personalization (still checked, so partial?)
User
Software\Microsoft\InputPersonalization\TrainedDataStore
HarvestContacts
DWORD:0

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

Updating policy
---------------

Update files under `GPO\{07BDCD6A-3F72-473C-82B9-67BB69DBE54D}\DomainSysvol\GPO`
by copying new versions either from a backup or directly from the policy
directory (`\\<domain>\SYSVOL\<domain>\Policies\{<GUID>}` for domains or
`%SystemRoot%\System32\GroupPolicy` for local computer policy). When
Registry.pol files are updated, run `parse && sort && build` to update LGPO text
files and sort entry order.

GPO backups contain a lot of domain-specific information, which has been removed
from the XML files in this repository. Therefore, creating a new backup and
overwriting all files in the GPO directory is a bad idea. Also, the XML parser
used by the import wizard is horrible. Adding a new line in the wrong place
causes mmc.exe to crash, so the XML files are a mess.

Domain policies store client-side extension (CSE) GUIDs in LDAP. The CSEs must
be updated manually in Backup.xml in order for the policy to be imported and
processed correctly. Create a separate backup using the Group Policy Management
Console and copy `MachineExtensionGuids` and `UserExtensionGuids` elements from
there.

Local computer policy stores CSEs in
`%SystemRoot%\System32\GroupPolicy\GPT.INI`. LGPO.exe can be used to enable CSEs
via the `/e` option, but it does not distinguish between machine and user CSEs,
which is annoying, and has other problems that are noted below. Therefore, local
CSEs are enabled by copying over a pre-configured GPT.INI file.

**Keep Backup.xml and GPT.INI files synchronized!**

Refreshing CSEs
---------------

If you configure and then unconfigure a setting that requires a CSE, that CSE
remains enabled, which may slow down policy processing. You can check which CSEs
are enabled by generating a report via `gpresult /h GPReport.html`, and looking
at "Extensions Configured" under Computer/User Details -> Group Policy Objects.

To get rid of unnecessary CSEs, create a new GPO (or clear out the local
computer policy), and copy all policy files from the repository into the new
GPO. For local computer policy, run `Tools\LGPO.exe /g GPO`, but do not copy
GPT.INI. This will preserve all settings in gpedit, but all CSEs will remain
disabled.

To re-enable required CSEs:

1. Open gpedit and filter Administrative Templates to show only configured
   settings (set Managed and Commented to "Any" and clear all checkboxes).
2. Open "All Settings" container and double-click on the first setting.
3. For each setting, toggle Enabled/Disabled via Alt-E and Alt-D shortcuts,
   which will force gpedit to re-apply it, and then go to the next setting via
   Alt-N.
4. Repeat for all user settings.
5. Finally, toggle one policy under Security Settings and another one under
   Advanced Audit Policy Configuration.

There are two important caveats when doing this:

1. If you go too fast, you may run into the following error: "The process cannot
   access the file because it is being used by another process. (Exception from
   HRESULT: 0x80070020)." Just repeat the operation for the current setting and
   keep going. Gpedit uses separate threads to commit changes and these threads
   get in each other's way sometimes.
2. **Some settings that show up as disabled, but have extra options, must
   actually be enabled!** For example, "Allow Cloud Search" can be disabled, or
   enabled with the extra option set to "Disable Cloud Search". If you do the
   latter and then re-open the setting, it will show up as disabled. However, if
   you compare LGPO text exports of the resulting Registry.pol files, you'll
   notice that disabling the setting manually causes the associated value
   (SpynetReporting) to be deleted, whereas enabling it and setting the option
   to "Disable Cloud Search" sets the value to DWORD:0, which is what we want.
   If you find any of this confusing, I personally like to imagine myself slowly
   strangling the people who designed and implemented this stuff. Try it... It's
   very therapeutic. Here's a (possibly incomplete) list of such settings:
   * Allow Cloud Search
   * Join Microsoft MAPS
   * Limit Enhanced diagnostic data to the minimum required by Windows Analytics

The following table describes all CSE GUIDs referenced by this GPO:

| GUID                                     | Description |
| ---------------------------------------- | ----------- |
| `{169EBF44-942F-4C43-87CE-13C93996EBBE}` | UE-V |
| `{29BBE2D5-DE47-4855-97D7-2745E166DC6D}` | Cortana search |
| `{2BFCC077-22D2-48DE-BDE1-2F618D9B476D}` | App-V |
| `{2D4156A2-897A-11DB-BA21-001185AD2B89}` | Network List Manager Policies MMC snap-in |
| `{35378EAC-683F-11D2-A89A-00C04FBBCFA2}` | Registry settings |
| `{7933F41E-56F8-41D6-A31C-4148A711EE93}` | Windows search |
| `{803E14A0-B4FB-11D0-A0D0-00A0C90F574B}` | Security settings MMC snap-in |
| `{827D319E-6EAC-11D2-A4EA-00C04F79F83A}` | Security (domain-only) |
| `{C631DF4C-088F-4156-B058-4375F0853CD8}` | Offline files |
| `{CDEAFC3D-948D-49DD-AB12-E578BA4AF7AA}` | TCP/IP |
| `{D02B1F72-3407-48AE-BA88-E8213C6761F1}` | Computer setting |
| `{D02B1F73-3407-48AE-BA88-E8213C6761F1}` | User setting |
| `{DF3DC19F-F72C-4030-940E-4C2A65A6B612}` | No idea |
| `[{F3CCC681-B74C-4060-9F26-CD84525DCA2A}{0F3F3735-573D-9804-99E4-AB2A69BA5FD4}]` | [Advanced audit policy configuration](https://msdn.microsoft.com/en-us/library/dd976882.aspx) |

`LGPO.exe /e audit` does not add the same GUID pair for Advanced Audit Policy
Configuration as gpedit. The value above is added by gpedit to
gPCMachineExtensionNames (only), whereas LGPO v2.2 adds
`[{F3CCC681-B74C-4060-9F26-CD84525DCA2A}{DF3DC19F-F72C-4030-940E-4C2A65A6B612}]`
to both gPCMachineExtensionNames and gPCUserExtensionNames. This is probably not
correct, and is another reason to avoid LGPO's CSE handling.

References
----------

* https://docs.microsoft.com/en-us/windows/whats-new/
* https://docs.microsoft.com/en-us/windows/device-security/windows-security-baselines
* https://blogs.technet.microsoft.com/secguide/
* https://blogs.technet.microsoft.com/secguide/2016/10/02/the-mss-settings/
* https://blogs.technet.microsoft.com/secguide/2015/11/18/changes-from-the-windows-8-1-baseline-to-the-windows-10-th11507-baseline/
* https://docs.microsoft.com/en-us/windows/privacy/
* http://iase.disa.mil/stigs/os/windows/Pages/index.aspx
* https://www.stigviewer.com/stig/windows_10/
* https://blogs.technet.microsoft.com/josebda/2012/11/13/windows-server-2012-file-server-tip-disable-8-3-naming-and-strip-those-short-names-too/
* https://support.microsoft.com/en-us/help/2526083/disabling-user-account-control-uac-on-windows-server
* https://docs.microsoft.com/en-us/windows/client-management/mdm/policy-csp-privacy
* https://deploywindows.info/2016/02/11/windows-10-group-policy-client-side-extension-cse-list/
* https://www.vacuumbreather.com/index.php/blog/itemlist/category/11-windows-10
