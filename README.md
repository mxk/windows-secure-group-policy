# Windows 11 Secure Group Policy

A local group policy intended for standalone Windows 11 devices. It aims to improve privacy, security, and performance, in that order.

All settings are maintained in a single `PolicyRules` file that is applied using `LGPO.exe`. Security features that send data to Microsoft, such as SmartScreen, are disabled. Some settings are only effective on the Enterprise edition. The following network traffic is allowed:

* Network Connectivity Status Indicator (NCSI) tests (see Notes below).
* Windows updates.
* Root CA updates.
* Defender definition updates.
* Storage Health (Disk Failure Prediction Model) updates.

## Installation

Run `install.cmd` as an Administrator and restart the computer. Some settings may only be effective if applied immediately after a clean install.

## Templates

`PolicyDefinitions` directory contains the following templates:

1. Windows 11 Enterprise 23H2 ISO (22621.2428)
2. [Windows 11 v23H2 Security Baseline][SCT]
3. [Windows Restricted Traffic Limited Functionality Baseline - Windows 11 23H2][RTLFB]
4. [Microsoft Edge (119.0.2151.44)][Edge]

[SCT]: https://www.microsoft.com/en-us/download/details.aspx?id=55319
[RTLFB]: https://learn.microsoft.com/en-us/windows/privacy/manage-connections-from-windows-operating-system-components-to-microsoft-services
[Edge]: https://www.microsoft.com/en-us/edge/business/download

Before editing the policy, copy the templates to `C:\Windows\PolicyDefinitions`. Overwriting existing files is not recommended because it requires ownership changes, which makes SFC unhappy, which may break Windows Update. In general, it's better to use a VM running a matching version of Windows. For each new release, the `PolicyDefinitions` directory should be rebuilt from scratch by copying templates over in the listed order to ensure removal of outdated templates.

To extract `PolicyDefinitions` from a Windows ISO:

1. Mount the ISO file.
2. Open `sources\install.wim` with 7-Zip.
3. Check `[1].xml` for the appropriate image index and build version.
4. Extract `\<N>\Windows\PolicyDefinitions`.

## Local policy

The current local policy can be obtained with the following command in a Command Prompt running as an Administrator:

```
rmdir /s /q LGPO & mkdir LGPO & LGPO.exe /b %cd%\LGPO & GPO2PolicyRules.exe .\LGPO LGPO.PolicyRules
```
