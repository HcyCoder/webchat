$services = @(
    @{Name="user";    Port="50051"; Exe="app\user\user.go";        Conf="app\user\etc\user.yaml"},
    @{Name="chat";    Port="50052"; Exe="app\chat\chat.go";        Conf="app\chat\etc\chat.yaml"},
    @{Name="group";   Port="50053"; Exe="app\group\group.go";      Conf="app\group\etc\group.yaml"},
    @{Name="media";   Port="50054"; Exe="app\media\media.go";      Conf="app\media\etc\media.yaml"},
    @{Name="gateway"; Port="8080";  Exe="app\gateway\gateway.go";  Conf="app\gateway\etc\gateway.yaml"}
)

$scriptDir = $PSScriptRoot
$tmpDir = "$env:TEMP\webchat"
New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null

$tabParts = @()

for ($i = 0; $i -lt $services.Count; $i++) {
    $svc = $services[$i]
    $runner = Join-Path $tmpDir "run_$($svc.Name).ps1"

@"
Write-Host "[$($svc.Name)] listening on :$($svc.Port)" -ForegroundColor Cyan
go run $($svc.Exe) -f $($svc.Conf)
"@ | Set-Content -Path $runner -Encoding UTF8

    $prefix = if ($i -eq 0) { "" } else { "new-tab " }
    $tabParts += "${prefix}--title webchat-$($svc.Name) pwsh -NoExit -File `"$runner`""
}

$wtCmd = "wt -d `"$scriptDir`" " + ($tabParts -join " ; ")

Write-Host "Starting all backend services in Windows Terminal..." -ForegroundColor Green
cmd /c $wtCmd
Write-Host "done." -ForegroundColor Green
