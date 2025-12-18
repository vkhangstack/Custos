# Build script for Custos on Windows

Write-Host "Building Custos for Windows..."

# Ensure we are in the project root
$exclude = "-ldflags '-s -w'"
$platform = "windows/amd64"

# Clean build directory if needed (Wails handles this usually, but good to be sure if we want a fresh start)
# Remove-Item -Path "build/bin" -Recurse -ErrorAction SilentlyContinue

# various options: 
# -clean: cleans the build directory
# -platform: target platform
# -ldflags: linker flags to reduce binary size and strip debug info
# -nsis: generate installer
wails build -clean -platform $platform -ldflags "-s -w -H windowsgui" -nsis 

if ($LASTEXITCODE -eq 0) {
    # Get version from app.json
    $appConfig = Get-Content "app.json" | ConvertFrom-Json
    $version = $appConfig.version

    $installerName = "Custos-$version.exe"
    $generatedInstaller = "build/bin/Custos-amd64-installer.exe"
    
    if (Test-Path $generatedInstaller) {
        Rename-Item -Path $generatedInstaller -NewName $installerName -Force
        Write-Host "Build successful! Installer located at build/bin/$installerName" -ForegroundColor Green
    } else {
        Write-Host "Build successful! Executable located at build/bin/Custos.exe (Installer not found)" -ForegroundColor Yellow
    }
} else {
    Write-Host "Build failed." -ForegroundColor Red
}
