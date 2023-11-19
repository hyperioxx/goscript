Write-Output "Installing goscript..."
Move-Item -Path .\goscript.exe -Destination "C:\Program Files\goscript\goscript.exe"

Write-Output "Creating modules directory..."
New-Item -ItemType Directory -Force -Path "C:\Program Files\goscript\modules"

$acl = Get-Acl "C:\Program Files\goscript\modules"
$acl.SetAccessRuleProtection($True, $False)
$rule = New-Object System.Security.AccessControl.FileSystemAccessRule("Everyone","FullControl","ContainerInherit, ObjectInherit", "None", "Allow")
$acl.AddAccessRule($rule)
Set-Acl "C:\Program Files\goscript\modules" $acl

Write-Output "Setting GS_MODULES environment variable..."
[Environment]::SetEnvironmentVariable("GS_MODULES", "C:\Program Files\goscript\modules", "Machine")

Write-Output "Installation completed."
Write-Output "Please restart your terminal or computer for the changes to take effect."
