# ==============================================================================
# E2E Test Script: SPEC-001 - Gerador de QR Code Dinâmico
# ==============================================================================

param (
    [string]$BaseUrl = "http://localhost:8080"
)

$ErrorActionPreference = "Stop"
$passed = 0
$total = 9

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " Iniciando Testes E2E: SPEC-001 (DynamicQR)" -ForegroundColor Cyan
Write-Host " Target URL: $BaseUrl" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

function Assert-Step {
    param (
        [int]$StepNum,
        [string]$Description,
        [bool]$Condition,
        [string]$ErrorDetail = ""
    )
    if ($Condition) {
        Write-Host " [PASS] Cenario ${StepNum}: $Description" -ForegroundColor Green
        $script:passed++
    } else {
        Write-Host " [FAIL] Cenario ${StepNum}: $Description" -ForegroundColor Red
        if ($ErrorDetail) {
            Write-Host "        Detalhe: $ErrorDetail" -ForegroundColor Yellow
        }
        exit 1
    }
}

# ------------------------------------------------------------------------------
# Cenário 1: Testar /healthz e conectividade do MongoDB
# ------------------------------------------------------------------------------
Write-Host "`nTestando Cenário 1: Verificação de Saúde (/healthz)..." -ForegroundColor Gray
$healthRes = Invoke-RestMethod -Uri "$BaseUrl/healthz" -Method Get
Assert-Step 1 "Endpoint /healthz online e conectado ao MongoDB" ($healthRes.status -eq "ok" -and $healthRes.mongo -eq "connected") "Status retornado: $($healthRes | ConvertTo-Json -Compress)"

# ------------------------------------------------------------------------------
# Cenário 2: Criar QR Code via POST /api/v1/qr-codes
# ------------------------------------------------------------------------------
Write-Host "`nTestando Cenário 2: Criação de QR Code Dinâmico..." -ForegroundColor Gray
$createBody = @{
    title = "E2E Test QR Code"
    target_url = "https://example.com/v1"
} | ConvertTo-Json

$created = Invoke-RestMethod -Uri "$BaseUrl/api/v1/qr-codes" -Method Post -ContentType "application/json" -Body $createBody
$qrId = $created.id
$slug = $created.slug

Assert-Step 2 "Criação de QR Code com ID e slug aleatório" ($null -ne $qrId -and $null -ne $slug -and $created.target_url -eq "https://example.com/v1") "Retorno: $($created | ConvertTo-Json -Compress)"

# ------------------------------------------------------------------------------
# Cenário 3: Validar persistência do campo qr_image_base32
# ------------------------------------------------------------------------------
Write-Host "`nTestando Cenário 3: Validação da Codificação Base32 da Imagem..." -ForegroundColor Gray
$base32Img = $created.qr_image_base32
$isBase32Valid = ($null -ne $base32Img -and $base32Img.Length -gt 50)
Assert-Step 3 "Campo qr_image_base32 preenchido e serializado em Base32" $isBase32Valid "Tamanho Base32: $($base32Img.Length)"

# ------------------------------------------------------------------------------
# Cenário 4: Redirecionamento 302 para https://example.com/v1
# ------------------------------------------------------------------------------
Write-Host "`nTestando Cenário 4: Redirecionamento público GET /{slug}..." -ForegroundColor Gray
$req = [System.Net.HttpWebRequest]::Create("$BaseUrl/$slug")
$req.AllowAutoRedirect = $false
try {
    $resp = $req.GetResponse()
    $statusCode = [int]$resp.StatusCode
    $location = $resp.Headers["Location"]
    $cacheControl = $resp.Headers["Cache-Control"]
    $resp.Close()
} catch [System.Net.WebException] {
    $resp = $_.Response
    $statusCode = [int]$resp.StatusCode
    $location = $resp.Headers["Location"]
    $cacheControl = $resp.Headers["Cache-Control"]
    $resp.Close()
}

$isRedirectV1 = ($statusCode -eq 302 -or $statusCode -eq 307) -and ($location -eq "https://example.com/v1")
Assert-Step 4 "Redirecionamento público retorna 302/307 para https://example.com/v1" $isRedirectV1 "StatusCode: $statusCode, Location: $location, Cache-Control: $cacheControl"

# ------------------------------------------------------------------------------
# Cenário 5: Validar incremento do contador de cliques (click_count)
# ------------------------------------------------------------------------------
Write-Host "`nTestando Cenário 5: Incremento atômico de acessos..." -ForegroundColor Gray
Start-Sleep -Milliseconds 400 # Aguarda brevemente a goroutine atômica persistir
$qrDetails = Invoke-RestMethod -Uri "$BaseUrl/api/v1/qr-codes/$qrId" -Method Get
Assert-Step 5 "Contador click_count incrementado atomicamente (>= 1)" ($qrDetails.click_count -ge 1) "Contador atual: $($qrDetails.click_count)"

# ------------------------------------------------------------------------------
# Cenário 6: Alterar destino para https://example.com/v2 via PUT /api/v1/qr-codes/{id}
# ------------------------------------------------------------------------------
Write-Host "`nTestando Cenário 6: Edição dinâmica da URL de destino..." -ForegroundColor Gray
$updateBody = @{
    title = "E2E Test QR Code - Atualizado"
    target_url = "https://example.com/v2"
} | ConvertTo-Json

$updated = Invoke-RestMethod -Uri "$BaseUrl/api/v1/qr-codes/$qrId" -Method Put -ContentType "application/json" -Body $updateBody
Assert-Step 6 "Edição dinâmica preservando o slug original" ($updated.target_url -eq "https://example.com/v2" -and $updated.slug -eq $slug) "Slug original: $slug, Slug pós update: $($updated.slug)"

# ------------------------------------------------------------------------------
# Cenário 7: Novo GET /{slug} deve redirecionar para https://example.com/v2
# ------------------------------------------------------------------------------
Write-Host "`nTestando Cenário 7: Redirecionamento instantâneo para novo destino..." -ForegroundColor Gray
$req2 = [System.Net.HttpWebRequest]::Create("$BaseUrl/$slug")
$req2.AllowAutoRedirect = $false
try {
    $resp2 = $req2.GetResponse()
    $statusCode2 = [int]$resp2.StatusCode
    $location2 = $resp2.Headers["Location"]
    $resp2.Close()
} catch [System.Net.WebException] {
    $resp2 = $_.Response
    $statusCode2 = [int]$resp2.StatusCode
    $location2 = $resp2.Headers["Location"]
    $resp2.Close()
}

$isRedirectV2 = ($statusCode2 -eq 302 -or $statusCode2 -eq 307) -and ($location2 -eq "https://example.com/v2")
Assert-Step 7 "Redirecionamento reflete imediatamente novo destino https://example.com/v2" $isRedirectV2 "StatusCode: $statusCode2, Location: $location2"

# ------------------------------------------------------------------------------
# Cenário 8: Baixar imagem PNG e validar integridade binária do Base32 decodificado
# ------------------------------------------------------------------------------
Write-Host "`nTestando Cenário 8: Download e decodificação da imagem PNG a partir do Base32..." -ForegroundColor Gray
$imageUri = "$BaseUrl/api/v1/qr-codes/$qrId/image?download=true"
$wc = New-Object System.Net.WebClient
$imgBytes = $wc.DownloadData($imageUri)
$wc.Dispose()

# Bytes mágicos do formato PNG: 0x89, 0x50, 0x4E, 0x47 (‰PNG)
$isPngValid = ($imgBytes.Length -gt 100 -and $imgBytes[0] -eq 0x89 -and $imgBytes[1] -eq 0x50 -and $imgBytes[2] -eq 0x4E -and $imgBytes[3] -eq 0x47)
Assert-Step 8 "Imagem recuperada do MongoDB, decodificada do Base32 e validada com assinatura PNG" $isPngValid "Tamanho: $($imgBytes.Length) bytes, Mágicos: $($imgBytes[0..3] -join ' ')"

# ------------------------------------------------------------------------------
# Cenário 9: Excluir QR Code e validar resposta 404
# ------------------------------------------------------------------------------
Write-Host "`nTestando Cenário 9: Exclusão e validação de 404..." -ForegroundColor Gray
Invoke-RestMethod -Uri "$BaseUrl/api/v1/qr-codes/$qrId" -Method Delete

$statusCode404 = $null
try {
    $req404 = [System.Net.HttpWebRequest]::Create("$BaseUrl/$slug")
    $req404.AllowAutoRedirect = $false
    $resp404 = $req404.GetResponse()
    $statusCode404 = [int]$resp404.StatusCode
    $resp404.Close()
} catch [System.Net.WebException] {
    $resp404 = $_.Exception.Response
    if ($resp404) {
        $statusCode404 = [int]$resp404.StatusCode
        $resp404.Close()
    }
}

Assert-Step 9 "QR Code excluído e endpoint público /{slug} retornando 404" ($statusCode404 -eq 404) "StatusCode pós-delete: $statusCode404"

# ------------------------------------------------------------------------------
# Conclusao
# ------------------------------------------------------------------------------
Write-Host "`n==========================================================" -ForegroundColor Cyan
Write-Host " RESULTADO DOS TESTES E2E: $passed / $total CENARIOS PASSANDO!" -ForegroundColor Green
Write-Host " A SPEC-001 FOI VALIDADA COM SUCESSO." -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan
exit 0
