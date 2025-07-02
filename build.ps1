# Define as plataformas de destino [GOOS, GOARCH]
$platforms = @(
    @{goos="windows"; goarch="amd64"},
    @{goos="linux"; goarch="amd64"},
    @{goos="darwin"; goarch="amd64"} # macOS
)

# Diretório de saída para os binários
$outputDir = "release"
if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
}

# Ficheiro de entrada principal da aplicação
$sourceFile = "cmd/server/main/main.go"

# Itera sobre cada plataforma e compila o binário
foreach ($p in $platforms) {
    $env:GOOS = $p.goos
    $env:GOARCH = $p.goarch
    $ext = ""
    if ($env:GOOS -eq "windows") {
        $ext = ".exe"
    }
    $outputName = "./$outputDir/main-$($env:GOOS)-$($env:GOARCH)$ext"
    
    Write-Host "A compilar para $($env:GOOS)/$($env:GOARCH)..."
    go build -ldflags "-s -w" -o $outputName $sourceFile
    if ($LASTEXITCODE -ne 0) {
        Write-Error "A compilação falhou para $($env:GOOS)/$($env:GOARCH)"
        exit 1
    }
}

Write-Host "Compilação para todas as plataformas concluída com sucesso! Os ficheiros estão em '$outputDir'." 