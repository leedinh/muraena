# Windows PowerShell script to modify hosts file
$domain = "github.local"
$ip = "127.0.0.1"  # For Windows, we should use 127.0.0.1 to connect to WSL
$hostsFile = "C:\Windows\System32\drivers\etc\hosts"

Write-Host "`nTo add the host entry, run PowerShell as Administrator and execute:"
Write-Host "Add-Content -Path $hostsFile -Value '$ip`t$domain' -Force"

Write-Host "`nTo remove the entry later, run PowerShell as Administrator and execute:"
Write-Host "(Get-Content $hostsFile) | Where-Object { `$_ -notmatch '$domain' } | Set-Content $hostsFile"

Write-Host "`nCurrent entries in hosts file for $domain:"
Get-Content $hostsFile | Select-String $domain 