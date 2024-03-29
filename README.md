# Windows 11 Secure Group Policy

A local group policy intended for standalone Windows 11 devices. It aims to improve privacy, security, and performance, in that order.

All settings are maintained in a single `PolicyRules` file that is applied with [LGPO]. Security features that send data to Microsoft, such as SmartScreen, are disabled, deviating from [Microsoft's Security Baseline][Baseline]. Some settings are only effective on the Enterprise edition.

The target Feature Update version is **Windows 11 23H2**. This prevents automatic updates to the next release before the policy is updated with new settings.

[LGPO]: https://learn.microsoft.com/en-us/windows/security/operating-system-security/device-management/windows-security-configuration-framework/security-compliance-toolkit-10#what-is-the-local-group-policy-object-lgpo-tool
[Baseline]: https://learn.microsoft.com/en-us/windows/security/operating-system-security/device-management/windows-security-configuration-framework/windows-security-baselines

## Installation

Run `install.cmd` as an Administrator and restart the computer.

## Saving local policy

Run `savelocal.cmd <out-file> <policy-name>` or `savewin11.cmd` (creates `Win11-Local.PolicyRules`) as an Administrator to save the local group policy as a `PolicyRules` file.

*Warning:* This will overwrite the contents of `C:\GPO`.

## Comparing policies

Download and use the [Policy Analyzer][SCT] to compare `PolicyRules` files. Make sure you configure it to use the repository's `PolicyDefinitions` directory rather than `C:\Windows\PolicyDefinitions`.

[SCT]: https://www.microsoft.com/en-us/download/details.aspx?id=55319

## Updating policy

When `LGPO.exe` and `GPO2PolicyRules.exe` export the local policy, they include many default settings that shouldn't be overwritten when applying the resulting `PolicyRules` file. There is also a bug in handling `(Default)` registry values. These are annoyances that prevent a clean install/save roundtrip and add noise when comparing against Microsoft's Security Baseline. Default settings were manually removed from `Win11.PolicyRules` by doing a three-way comparison between it, `MSFT-Win11.PolicyRules`, and `Win11-CleanInstall.PolicyRules`. To avoid reverting these edits, any updates to the policy must be merged in manually:

1. Use `gpedit.msc` to modify the local policy.
2. Run `savewin11.cmd` to create `Win11-Local.PolicyRules` file (not version-controlled).
3. Diff and copy the relevant settings to `Win11.PolicyRules`.

To update the policy for a new Windows feature release:

1. Download and install the new ISO in a Hyper-V VM.
2. Configure VM enhanced session settings to redirect the host drive that contains this repository.
3. Run `cmd.exe` as an Administrator in the VM.
4. Map host drive for easier access: `net use Z: \\tsclient\<drive>\<path-to-repo>`
5. Update and copy over the templates as described in the next section. Update this file with new version information.
6. Save an updated `Win11-CleanInstall.PolicyRules` file.
7. Install the current policy and restart.
8. Follow the steps above to update the policy, comparing it against the new security baselines.

## Templates

Templates contained in the `PolicyDefinitions` directory:

1. Windows 11 Enterprise 23H2 ISO (22621.2428)
2. [Windows 11 v23H2 Security Baseline][SCT]
3. [Windows Restricted Traffic Limited Functionality Baseline - Windows 11 23H2][RTLFB]
4. [Microsoft Edge (119.0.2151.72)][Edge]
5. [Mozilla Firefox (5.4)][Firefox]

[RTLFB]: https://learn.microsoft.com/en-us/windows/privacy/manage-connections-from-windows-operating-system-components-to-microsoft-services
[Edge]: https://www.microsoft.com/en-us/edge/business/download
[Firefox]: https://github.com/mozilla/policy-templates/releases

Before editing the policy with `gpedit.msc`, copy the templates to `C:\Windows\PolicyDefinitions`. Overwriting existing files is not recommended because it requires ownership changes, which makes SFC unhappy, which may break Windows Update. In general, it's better to start with a VM running a matching version of Windows. For each new release, the `PolicyDefinitions` directory should be rebuilt from scratch by copying the templates over in the listed order to ensure removal of outdated templates.

To extract `PolicyDefinitions` from a Windows ISO:

1. Mount the ISO file.
2. Open `sources\install.wim` with 7-Zip.
3. Check `[1].xml` for the appropriate image index and build version.
4. Extract `\<N>\Windows\PolicyDefinitions`.

## Notes

* Local Administrator is disabled and renamed to LocalAdmin. The password must be set before it can be enabled or you'll get error 7016 in `gpresult /h` report along with "Security has requested to process its policy settings again" message.
* Microsoft Account sign-in is disabled. Do not disable "Microsoft Account Sign-in Assistant" service (aka "wlidsvc") as suggested by the Restricted Traffic Baseline. Doing so results in error 0x80070426 when running Windows Update.
* Disabling active Network Connectivity Status Indicator (NCSI) tests breaks Windows Update (0x800704cf), DISM (0x800f0906 or 0x800f081f), and probably other things. Passive-only configuration doesn't fix this, so it's best to leave both passive and active tests enabled. Windows Update error only happens once if Windows was installed without any network connectivity. After a single successful update, disabling NCSI doesn't seem to cause further problems, but DISM will still run into errors. Guess how many days it took to figure this out.
* "Turn off all Windows spotlight features" policy must be applied within 15 minutes after Windows is installed (was true for Windows 10, no longer the case for Windows 11?).

### Settings without a template

The following registry entries do not have an associated template and are treated as preference-type settings that are not removed automatically when no longer applied by the policy:

* `DisableWpad=1` and `AutoDetect=0` disable [automatic proxy detection (WPAD)][WPAD]. Do not disable "WinHTTP Web Proxy Auto-Discovery Service" (aka "WinHttpAutoProxySvc") - doing so [will break things][WinHTTP].
* `HKCU\Software\Classes\CLSID\{86CA1AA0-34AA-4E8B-A509-50C905BAE2A2}\InprocServer32` key restores classic File Explorer context menus.
* `HideFileExt=0` shows all file extensions in File Explorer.
* `ShowSyncProviderNotifications=0` disables sync provider notifications, which are used to show Microsoft ads in File Explorer.
* `ScoobeSystemSettingEnabled=0` disables "Let's finish setting up your device" notification (Settings > System > Notifications > Additional settings > Suggest ways to get the most out of Windows and finish setting up this device).

[WPAD]: https://learn.microsoft.com/en-us/troubleshoot/windows-server/networking/disable-http-proxy-auth-features#how-to-disable-wpad
[WinHTTP]: https://github.com/MicrosoftDocs/windows-itpro-docs/issues/2965#issuecomment-475441420
