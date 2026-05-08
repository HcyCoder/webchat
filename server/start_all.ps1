$services = @(
    @{Name="user";    Port="50051"; Exe="app\user\user.go";        Conf="app\user\etc\user.yaml"},
    @{Name="chat";    Port="50052"; Exe="app\chat\chat.go";        Conf="app\chat\etc\chat.yaml"},
    @{Name="group";   Port="50053"; Exe="app\group\group.go";      Conf="app\group\etc\group.yaml"},
    @{Name="media";   Port="50054"; Exe="app\media\media.go";      Conf="app\media\etc\media.yaml"},
    @{Name="gateway"; Port="8080";  Exe="app\gateway\gateway.go";  Conf="app\gateway\etc\gateway.yaml"}
)

$dir = Resolve-Path $PSScriptRoot
$tabCmds = @()

for ($i = 0; $i -lt $services.Count; $i++) {
    $svc = $services[$i]
    $cmd = "Write-Host '[$($svc.Name)] listening on :$($svc.Port)' -ForegroundColor Cyan; go run $($svc.Exe) -f $($svc.Conf)"
    if ($i -eq 0) {
        $tabCmds += "--title webchat-$($svc.Name) pwsh -NoExit -Command `"$cmd`""
    } else {
        $tabCmds += "new-tab --title webchat-$($svc.Name) pwsh -NoExit -Command `"$cmd`""
    }
}

$args = $tabCmds -join " ; "
$fullCmd = "wt -d `"$dir`" $args"

Write-Host "Starting all backend services in Windows Terminal..." -ForegroundColor Green
Invoke-Expression $fullCmd
Write-Host "done." -ForegroundColor Green
