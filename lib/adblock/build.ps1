# Build script for libadblock on Windows
Write-Host "Building libadblock Rust library..."

pushd $PSScriptRoot
cargo build --release
popd

if ($LASTEXITCODE -eq 0) {
    Write-Host "libadblock build successful!" -ForegroundColor Green
} else {
    Write-Host "libadblock build failed." -ForegroundColor Red
    exit 1
}
