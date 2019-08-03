Last updated:5/29/2019

'Version 1607' can be applied to version 1607 of Windows 10 Enterprise and Windows Server 2016 Datacenter except Nano Server.
'Version 1703' can be applied to version 1703 of Windows 10 Enterprise.
'Version 1709' can be applied to version 1709 of Windows 10 Enterprise.
'Version 1803' can be applied to version 1803 of Windows 10 Enterprise.
'Version 1809' can be applied to version 1809 of Windows 10 Enterprise, Windows Server  Datacenter and Windows Server Standard.
'Version 1903' can be applied to version 1903 of Windows 10 Enterprise

==============================================================================

LGPO.exe is required to install the Restricted Traffic Limited Functionality package, so please follow the below instructions to download LGPO.exe:  

1. Download LGPO.zip from https://www.microsoft.com/en-us/download/details.aspx?id=55319
2. Extract LGPO.zip into the "Tools" folder under respective version package (e.g. “\Version 1809\Enterprise\Tools\” )
3. Run a command prompt with administrator privilege.
4. [For Version 1903 and above] Run "\Version 1903\Enterprise\RestrictedTraffic_ClientEnt_Install.cmd for Enterprise OR \Version 1903\Server\RestrictedTraffic_ClientyServer_Install.cmd” for Server SKU 
5. Check the “\Logs” folder to make sure execution was successful 
6. Restart Windows 

Note: For Version 1607 Enterprise, background apps under path "Setting-Privacy-Background apps" need to be turned OFF manually.
