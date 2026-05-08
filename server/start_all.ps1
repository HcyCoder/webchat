$services = @(
    @{Name="user";    Port="50051"; Exe="app\user\user.go";        Conf="app\user\etc\user.yaml"},
    @{Name="chat";    Port="50052"; Exe="app\chat\chat.go";        Conf="app\chat\etc\chat.yaml"},
    @{Name="group";   Port="50053"; Exe="app\group\group.go";      Conf="app\group\etc\group.yaml"},
    @{Name="media";   Port="50054"; Exe="app\media\media.go";      Conf="app\media\etc\media.yaml"},
    @{Name="gateway"; Port="8080";  Exe="app\gateway\gateway.go";  Conf="app\gateway\etc\gateway.yaml"}
)

Write-Host "Starting all backend services..." -ForegroundColor Green

Push-Location $PSScriptRoot

foreach ($svc in $services) {
    $title = "webchat-$($svc.Name) :$($svc.Port)"
    Start-Process pwsh -ArgumentList "-NoExit", "-Command",
        "Write-Host '[$($svc.Name)] listening on :$($svc.Port)' -ForegroundColor Cyan; go run $($svc.Exe) -f $($svc.Conf)" `
        -WindowStyle Normal
    Write-Host "  started $title"
}

Pop-Location
Write-Host "done. close each window to stop the service." -ForegroundColor Green
