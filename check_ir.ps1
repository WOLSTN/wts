param (
    [Parameter(Mandatory=$true)]
    [string]$WirFile
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $WirFile)) {
    Write-Error "File not found: $WirFile"
    exit 1
}

Write-Host "=== Checking IR: $WirFile ===" -ForegroundColor Cyan
Write-Host ""

$json = Get-Content $WirFile -Raw | ConvertFrom-Json
$ok = $true

$types = @{}

$json.types | ForEach-Object { $types[$_.id] = $_ }

foreach ($func in $json.functions) {
    if (-not $func.body) { continue }
    if (-not $func.body.blocks) { continue }

    $exprCount = 0
    $badBinaryCount = 0
    $totalInstrs = 0
    $exprs = @()
    $badBinaries = @()

    foreach ($block in $func.body.blocks) {
        if (-not $block.instrs) { continue }
        foreach ($instr in $block.instrs) {
            $totalInstrs++

            if ($instr.opcode -eq "expr") {
                $exprCount++
                $exprs += "$($block.label)::$($instr.id) value=$($instr.value)"
            }

            if ($instr.opcode -eq "binary" -and (
                $instr.value -eq "KindEqualsToken" -or
                $instr.value -eq "KindPlusEqualsToken" -or
                $instr.value -eq "KindMinusEqualsToken" -or
                $instr.value -eq "KindAsteriskEqualsToken" -or
                $instr.value -eq "KindSlashEqualsToken"
            )) {
                $badBinaryCount++
                $badBinaries += "$($block.label)::$($instr.id) op=$($instr.value)"
            }
        }
    }

    if ($exprCount -gt 0 -or $badBinaryCount -gt 0) {
        $ok = $false
        Write-Host "  [!] $($func.name): $exprCount expr, $badBinaryCount bad binary (total $totalInstrs instrs)" -ForegroundColor Red

        if ($exprs.Count -gt 0) {
            Write-Host "      Unhandled AST nodes (expr):" -ForegroundColor Yellow
            foreach ($e in $exprs) {
                Write-Host "        $e" -ForegroundColor Gray
            }
        }
        if ($badBinaries.Count -gt 0) {
            Write-Host "      Suspect binary ops:" -ForegroundColor Yellow
            foreach ($b in $badBinaries) {
                Write-Host "        $b" -ForegroundColor Gray
            }
        }
    } else {
        Write-Host "  [√] $($func.name): clean ($totalInstrs instrs)" -ForegroundColor Green
    }
}

Write-Host ""

# Check variables
$untypedVars = 0
foreach ($v in $json.variables) {
    if (-not $v.type) { $untypedVars++ }
}

if ($untypedVars -gt 0) {
    Write-Host "  [!] $untypedVars variables without type" -ForegroundColor Yellow
} else {
    Write-Host "  [√] All variables have types" -ForegroundColor Green
}

Write-Host ""

if ($ok) {
    Write-Host "=== PASS: IR is clean ===" -ForegroundColor Green
    exit 0
} else {
    Write-Host "=== FAIL: IR has issues ===" -ForegroundColor Red
    exit 1
}
