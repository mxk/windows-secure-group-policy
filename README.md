# Windows 11 Secure Group Policy

A local group policy intended for standalone Windows 11 devices. It aims to improve privacy, security, and performance, in that order.

All settings are maintained in a single [`Win11.PolicyRules`] file that is applied with [LGPO]. Security features that send data to Microsoft, such as SmartScreen, are disabled, deviating from [Microsoft's Security Baseline][Baseline]. Some settings are only effective on the Enterprise and/or Education editions.

The target Feature Update version is **Windows 11 25H2**. This prevents automatic updates to the next release before the policy is updated with new settings.

Releases are tagged using SemVer 2.0 format where `MAJOR.MINOR` is the target Windows feature update version (e.g. `24H2` is `24.2`) and `PATCH` is the policy revision.

[`Win11.PolicyRules`]: ./PolicyRules/Win11.PolicyRules
[LGPO]: https://learn.microsoft.com/en-us/windows/security/operating-system-security/device-management/windows-security-configuration-framework/security-compliance-toolkit-10#what-is-the-local-group-policy-object-lgpo-tool
[Baseline]: https://learn.microsoft.com/en-us/windows/security/operating-system-security/device-management/windows-security-configuration-framework/windows-security-baselines

## Usage

*All scripts must be run as an Administrator.*

Run `install.cmd` and restart the computer to install [`Win11.PolicyRules`] as the local group policy. This does not clear out any existing settings. Run `reset.cmd` first to start with a clean slate.

Run `save.cmd <out-file> <policy-name>` or `savewin11.cmd` (creates `Win11-Local.PolicyRules`) to save the local group policy as a `PolicyRules` file. *Warning: This will overwrite the contents of `C:\GPO`.*

Run `reset.cmd` followed by `gpupdate /force` or a restart to clear out the local group policy.

## Comparing policies

Download and use the [Policy Analyzer][SCT] to compare `PolicyRules` files. Make sure you configure it to use the repository's `PolicyDefinitions` directory rather than `C:\Windows\PolicyDefinitions`.

[SCT]: https://www.microsoft.com/en-us/download/details.aspx?id=55319

## Updating policy

When `LGPO.exe` and `GPO2PolicyRules.exe` export the local policy, they include many default settings that shouldn't be overwritten when applying the resulting `PolicyRules` file. There is also a bug in handling `(Default)` registry values. These are annoyances that prevent a clean install/save roundtrip and add noise when comparing against Microsoft's Security Baseline. Default settings were manually removed from `Win11.PolicyRules` by doing a three-way comparison between it, `MSFT-Win11-*.PolicyRules`, and `Win11-Clean-*.PolicyRules`. The easiest way to do this is by removing all `ComputerConfig` and `UserConfig` entries, and then comparing all three files with Policy Analyzer. To avoid reverting these edits, any updates to the policy must be merged in manually:

1. Use `gpedit.msc` to modify the local policy.
2. Run `savewin11.cmd` to create `Win11-Local.PolicyRules` file (not version-controlled).
3. Diff and copy the relevant settings to `Win11.PolicyRules`. Ignore `SecurityTemplate` changes unless those were actually modified. Be sure to check for new `CSE-Machine` and `CSE-User` entries.

To update the policy for a new Windows feature release:

1. Download the latest ISO.
2. Check for updates to [LGPO and Policy Analyzer tools][SCT].
3. Rebuild `PolicyDefinitions` directory as described in the next section and update README.md with new version information.
4. Create a new Hyper-V VM and do a clean install.
5. Configure VM enhanced session settings to redirect the host drive that contains this repository.
6. Run `cmd.exe` as an Administrator in the VM.
7. Map host drive for easier access: `net use Z: \\tsclient\<drive>\<path-to-repo>`
8. Run `Z:\savewin11.cmd` and rename `Z:\PolicyRules\Win11-Local.PolicyRules` to `Z:\PolicyRules\Win11-Clean-vXXH2.PolicyRules`.
9. Run `Z:\install.cmd`, restart, and remap `Z:`.
10. Copy updated `Z:\PolicyDefinitions` directory to the VM, skipping all existing files.
11. Use Policy Analyzer to view differences between the old and new Security Baselines.
12. Follow the steps above to update the policy via `gpedit.msc`.
13. Update `Computer Configuration → Windows Components → Windows Update → Manage updates offered from Windows Update → Select the target Feature Update version` and README.md.
14. Use Policy Analyzer to view and resolve any additional differences between `MSFT-Win11-vXXH2.PolicyRules` and `Win11.PolicyRules`.
    * Do not copy settings directly. Always use `gpedit.msc` to make changes, followed by `savewin11.cmd`, and merge from `Win11-Local.PolicyRules`.
    * In general, always set a value for any setting in the Security Baseline, even if it's the default. Conflicts are easier to review, whereas if a setting is missing, it's not clear whether it is new or was omitted intentionally. Exceptions are Internet Explorer, LAPS, and Attack Surface Reduction policies, which are not used, and a few other settings that can't be set.
15. Reapply the updated policy on the VM and make any additional changes as needed. See Security Baseline `Documentation\Windows 11 XXH2 to XXH2 Delta.xlsx` for information about new settings.

## Templates

Templates contained in the `PolicyDefinitions` directory:

1. Windows 11 Enterprise 25H2 ISO (26200.6584)
   * SHA256: `2B65DF49334B64E9341DC404E9C527BF1B2A9A105E95314A347FD29AC9900581`
2. [Windows 11 v25H2 Security Baseline][SCT]
3. [Windows Restricted Traffic Limited Functionality Baseline - Windows 11 23H2][RTLFB]
4. [Microsoft Edge (140.0.3485.94)][Edge]
5. [Mozilla Firefox (7.3)][Firefox]

[RTLFB]: https://learn.microsoft.com/en-us/windows/privacy/manage-connections-from-windows-operating-system-components-to-microsoft-services
[Edge]: https://www.microsoft.com/en-us/edge/business/download
[Firefox]: https://github.com/mozilla/policy-templates/releases

Before editing the policy with `gpedit.msc`, copy the templates to `C:\Windows\PolicyDefinitions`. Overwriting existing files is not recommended because it requires ownership changes, which makes SFC unhappy, which may break Windows Update. In general, it's better to start with a VM running a matching version of Windows. For each new release, the `PolicyDefinitions` directory should be rebuilt from scratch by copying the templates over in the listed order to ensure removal of outdated templates.

To extract `PolicyDefinitions` from a Windows ISO:

1. Mount the ISO file.
2. Open `sources\install.wim` with 7-Zip.
3. Check `[1].xml` for the appropriate image index and build version.
4. Extract `\<N>\Windows\PolicyDefinitions`.
5. Record `BUILD` and `SPBUILD` versions.

## Notes

* Local Administrator is disabled and renamed to LocalAdmin. The password must be set before it can be enabled or you'll get error 7016 in `gpresult /h` report along with "Security has requested to process its policy settings again" message.
* Microsoft Account sign-in is disabled. Do not disable "Microsoft Account Sign-in Assistant" service (aka "wlidsvc") as suggested by the Restricted Traffic Baseline. Doing so results in error 0x80070426 when running Windows Update.
* Disabling active Network Connectivity Status Indicator (NCSI) tests breaks Windows Update (0x800704cf), DISM (0x800f0906 or 0x800f081f), and probably other things. Passive-only configuration doesn't fix this, so it's best to leave both passive and active tests enabled. Windows Update error only happens once if Windows was installed without any network connectivity. After a single successful update, disabling NCSI doesn't seem to cause further problems, but DISM will still run into errors. Guess how many days it took to figure this out.
* "Turn off all Windows spotlight features" policy must be applied within 15 minutes after Windows is installed (was true for Windows 10, no longer the case for Windows 11?).
* `LGPO.exe` has a [known bug] with importing `REG_MULTI_SZ` values literally rather than converting `\0` escape sequences into a separator, so these values cannot be added to [`Win11.PolicyRules`] file.

[known bug]: https://techcommunity.microsoft.com/discussions/security-baselines/reg-multi-sz-are-not-imported-properly/2306822

### Settings without a template

The following registry entries do not have an associated template and are treated as preference-type settings that are not removed automatically when no longer applied by the policy:

* `DisableWpad=1` and `AutoDetect=0` disable [automatic proxy detection (WPAD)][WPAD]. Do not disable "WinHTTP Web Proxy Auto-Discovery Service" (aka "WinHttpAutoProxySvc") - doing so [will break things][WinHTTP].
* `HKCU\Software\Classes\CLSID\{86CA1AA0-34AA-4E8B-A509-50C905BAE2A2}\InprocServer32` key restores classic File Explorer context menus.
* `HideFileExt=0` shows all file extensions in File Explorer.
* `ShowSyncProviderNotifications=0` disables sync provider notifications, which are used to show Microsoft ads in File Explorer.
* `ScoobeSystemSettingEnabled=0` disables "Let's finish setting up your device" notification (Settings > System > Notifications > Additional settings > Suggest ways to get the most out of Windows and finish setting up this device).
* `IsResumeAllowed=0` disables cross device resume.
* `Start_IrisRecommendations=0` disables "Show recommendations for tips, shortcuts, new apps, and more" Start setting.

[WPAD]: https://learn.microsoft.com/en-us/troubleshoot/windows-server/networking/disable-http-proxy-auth-features#how-to-disable-wpad
[WinHTTP]: https://github.com/MicrosoftDocs/windows-itpro-docs/issues/2965#issuecomment-475441420

### Firefox

Firefox is configured using a combination of managed policies and [JSON preferences](./firefox-prefs.json) (copied to the [Preferences (JSON on one line)][Preferences] policy setting in minified format). See the following links for detailed setting information:

* [Firefox Policy Templates]
* [How to stop Firefox from making automatic connections][Firefox Connections]
* [brainfucksec/user.js]
* [arkenfox/user.js]

Preferences setting does not support the `app.*` prefix, so `false` values for `app.normandy.enabled` and `app.shield.optoutstudies.enabled` are not applied. Also, the one line minified `REG_SZ` format must be used because `LGPO.exe` has a bug in handling `\0` separators in `REG_MULTI_SZ` values.

[Preferences]: https://mozilla.github.io/policy-templates/#preferences
[Firefox Policy Templates]: https://mozilla.github.io/policy-templates/
[Firefox Connections]: https://support.mozilla.org/en-US/kb/how-stop-firefox-making-automatic-connections
[brainfucksec/user.js]: https://gist.github.com/brainfucksec/68e79da1c965aeaa4782914afd8f7fa2
[arkenfox/user.js]: https://github.com/arkenfox/user.js

### Client-Side Extensions (CSEs)

`LGPO.exe` writes `CSE-Machine` and `CSE-User` entries to `C:\Windows\System32\GroupPolicy\GPT.INI`. The format of the resulting `gPCMachineExtensionNames` and `gPCUserExtensionNames` values is `[{CSE1}{Tool1}{Tool2}...][{CSE2}...]...`. It's not clear whether tool extension GUIDs are important for correct policy enforcement (if you know, please message me). `LGPO.exe` only adds an undocumented `DF3DC19F-F72C-4030-940E-4C2A65A6B612` tool GUID for all entries when applying a PolicyRules file. The following table lists tool GUIDs that are normally set when the policy is configured via `gpedit.msc`:

| CSE                                      | Tool GUID                                | Added By                                         |
|------------------------------------------|------------------------------------------|--------------------------------------------------|
| `{35378EAC-683F-11D2-A89A-00C04FBBCFA2}` | `{B05566AC-FE9C-4368-BE01-7A4CBB6CBA11}` | Windows Defender Firewall with Advanced Security |
| All registry computer CSEs               | `{D02B1F72-3407-48AE-BA88-E8213C6761F1}` | Computer Administrative Templates                |
| All registry user CSEs                   | `{D02B1F73-3407-48AE-BA88-E8213C6761F1}` | User Administrative Templates                    |
| `{F3CCC681-B74C-4060-9F26-CD84525DCA2A}` | `{0F3F3735-573D-9804-99E4-AB2A69BA5FD4}` | Advanced Audit Policy Configuration              |

See [\[MS-GPSO\].pdf] for additional CSE and Tool Extension GUIDs.

To refresh CSE configuration, removing any extensions that are no longer needed by the current policy:

1. Delete `C:\Windows\System32\GroupPolicy\GPT.INI` file.
2. Open `gpedit.msc` and filter Computer Administrative Templates to show only Configured settings (set Managed and Commented to "Any" and clear all checkboxes).
3. Open "All Settings" container and double-click on the first setting.
4. For each setting, toggle Enabled/Disabled via Alt-E/Alt-D shortcuts, which will force `gpedit.msc` to re-apply it, and then go to the next setting via Alt-N.
   * You may occasionally see "The process cannot access the file because it is being used by another process. (Exception from HRESULT: 0x80070020)" error. Just repeat the operation for the current setting and keep going.
   * "Join Microsoft MAPS" and "Limit optional diagnostic data for Desktop Analytics" show up as Disabled, but must actually be set to Enabled and then configured via the drop-down control.
5. Repeat for all User settings.
6. Expand Windows Settings → Security Settings and toggle one setting under Security Options, Windows Defender Firewall with Advanced Security, and Advanced Audit Policy Configuration.
7. Extract the resulting CSEs from `GPT.INI`.

[\[MS-GPSO\].pdf]: https://download.microsoft.com/download/5/0/1/501ED102-E53F-4CE0-AA6B-B0F93629DDC6/Windows/[MS-GPSO].pdf
