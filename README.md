# Windows 11 Secure Group Policy

A local group policy intended for standalone Windows 11 devices. It aims to improve privacy, security, and performance, in that order.

All settings are maintained in a single `PolicyRules` file that is applied using [LGPO]. Security features that send data to Microsoft, such as SmartScreen, are disabled. This deviates from [Microsoft's Security Baseline][Baseline]. Some settings are only effective on the Enterprise edition.

The target Feature Update version is **Windows 11 23H2**. This prevents automatic updates to the next release before the policy is updated.

[LGPO]: https://learn.microsoft.com/en-us/windows/security/operating-system-security/device-management/windows-security-configuration-framework/security-compliance-toolkit-10#what-is-the-local-group-policy-object-lgpo-tool
[Baseline]: https://learn.microsoft.com/en-us/windows/security/operating-system-security/device-management/windows-security-configuration-framework/windows-security-baselines

## Installation

Run `install.cmd` as an Administrator and restart the computer. Some settings may only be effective if applied immediately after a clean install.

## Templates

The `PolicyDefinitions` directory contains the following templates:

1. Windows 11 Enterprise 23H2 ISO (22621.2428)
2. [Windows 11 v23H2 Security Baseline][SCT]
3. [Windows Restricted Traffic Limited Functionality Baseline - Windows 11 23H2][RTLFB]
4. [Microsoft Edge (119.0.2151.44)][Edge]

[SCT]: https://www.microsoft.com/en-us/download/details.aspx?id=55319
[RTLFB]: https://learn.microsoft.com/en-us/windows/privacy/manage-connections-from-windows-operating-system-components-to-microsoft-services
[Edge]: https://www.microsoft.com/en-us/edge/business/download

Before editing the policy, copy the templates to `C:\Windows\PolicyDefinitions`. Overwriting existing files is not recommended because it requires ownership changes, which makes SFC unhappy, which may break Windows Update. In general, it's better to start with a VM running a matching version of Windows. For each new release, the `PolicyDefinitions` directory should be rebuilt from scratch by copying the templates over in the listed order to ensure removal of outdated templates.

To extract `PolicyDefinitions` from a Windows ISO:

1. Mount the ISO file.
2. Open `sources\install.wim` with 7-Zip.
3. Check `[1].xml` for the appropriate image index and build version.
4. Extract `\<N>\Windows\PolicyDefinitions`.

## Saving local policy

Run `.\savelocal.cmd <out-file> <policy-name>` or `.\savewin11.cmd` (creates `Win11-Local.PolicyRules`) as an Administrator to save the local group policy as a `PolicyRules` file.

*Warning:* This will overwrite the contents of `C:\GPO`.

## Updating policy

When `LGPO.exe` and `GPO2PolicyRules.exe` export the local policy, they include many default settings that shouldn't be overwritten when applying the resulting `PolicyRules` file. There is also a bug in handling `(Default)` registry values. This is an annoyance that prevents a clean install/save roundtrip and adds noise when comparing against Microsoft's Security Baseline. Default settings were manually removed from `Win11.PolicyRules` by doing a three-way comparison between it, `MSFT-Win11.PolicyRules`, and `Win11-CleanInstall.PolicyRules` with the Policy Analyzer. To avoid reverting these edits, any updates to the policy must be merged in manually:

1. Use `gpedit.msc` to modify the local policy.
2. Run `.\savewin11.cmd` to create `Win11-Local.PolicyRules` file.
3. Diff and copy the relevant settings to `Win11.PolicyRules`.

## Notes

* Local Administrator is disabled and renamed to LocalAdmin. The password must be set before it can be enabled or you'll get error 7016 in `gpresult /h` report along with "Security has requested to process its policy settings again" message.
* Microsoft Account sign-in is disabled. Do not disable "Microsoft Account Sign-in Assistant" service (aka "wlidsvc") as suggested by the Restricted Traffic Baseline. Doing so results in error 0x80070426 when running Windows Update.
* Disabling active Network Connectivity Status Indicator (NCSI) tests breaks Windows Update (0x800704cf), DISM (0x800f0906 or 0x800f081f), and probably other things. Passive-only configuration doesn't fix this, so it's best to leave both passive and active tests enabled. Windows Update error only happens once if Windows was installed without any network connectivity. After a single successful update, disabling NCSI doesn't seem to cause further problems, but DISM will still run into errors. Guess how many days it took to figure this out.
* "Turn off all Windows spotlight features" policy must be applied within 15 minutes after Windows is installed (was true for Windows 10, no longer the case for Windows 11?).

### Settings without a template

The following registry entries do not have an associated template and are treated as preference-type settings that do not get removed automatically when no longer applied by the policy:

* `DisableWpad=1` and `AutoDetect=0` disable [Automatic proxy detection (WPAD)][WPAD]. Disabling "WinHTTP Web Proxy Auto-Discovery Service" (aka "WinHttpAutoProxySvc") [will break things][break].
* `HKCU\Software\Classes\CLSID\{86CA1AA0-34AA-4E8B-A509-50C905BAE2A2}\InprocServer32` key restores classic File Explorer context menus.
* `HideFileExt=0` shows all file extensions in File Explorer.
* `ScoobeSystemSettingEnabled=0` disables "Let's finish setting up your device" notification (Settings > System > Notifications > Additional settings > Suggest ways to get the most out of Windows and finish setting up this device).

[WPAD]: https://learn.microsoft.com/en-us/troubleshoot/windows-server/networking/disable-http-proxy-auth-features#how-to-disable-wpad
[break]: https://github.com/MicrosoftDocs/windows-itpro-docs/issues/2965#issuecomment-475441420
