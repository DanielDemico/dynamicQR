#!/usr/bin/env bash
# ==============================================================================
# E2E Test Script: SPEC-001 - Gerador de QR Code Dinâmico (Bash)
# ==============================================================================
set -euo pipefail

BASE_URL="${1:-http://localhost:8080}"
PASSED=0
TOTAL=9

echo "=========================================================="
echo " Iniciando Testes E2E: SPEC-001 (DynamicQR)"
echo " Target URL: $BASE_URL"
echo "=========================================================="

assert_step() {
    local step_num="$1"
    local desc="$2"
    local cond="$3"
    if [ "$cond" -eq 1 ]; then
        echo -e "\033[32m [PASS] Cenário $step_num: $desc\033[0m"
        PASSED=$((PASSED + 1))
    else
        echo -e "\033[31m [FAIL] Cenário $step_num: $desc\033[0m"
        exit 1
    fi
}

# 1. Healthz
echo "Testando Cenário 1: Verificação de Saúde (/healthz)..."
HEALTH_RES=$(curl -s "$BASE_URL/healthz")
STATUS=$(echo "$HEALTH_RES" | grep -o '"status":"ok"' || true)
[ -n "$STATUS" ] && assert_step 1 "Endpoint /healthz online e conectado ao MongoDB" 1 || assert_step 1 "Falha em /healthz" 0

# 2. Criar QR Code
echo "Testando Cenário 2: Criação de QR Code..."
CREATE_RES=$(curl -s -X POST "$BASE_URL/api/v1/qr-codes" \
  -H "Content-Type: application/json" \
  -d '{"title":"E2E Test Bash","target_url":"https://example.com/v1"}')

QR_ID=$(echo "$CREATE_RES" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
SLUG=$(echo "$CREATE_RES" | grep -o '"slug":"[^"]*' | cut -d'"' -f4)
[ -n "$QR_ID" ] && [ -n "$SLUG" ] && assert_step 2 "Criação de QR Code com ID e slug aleatório" 1 || assert_step 2 "Falha na criação" 0

# 3. Base32 gravado
BASE32_IMG=$(echo "$CREATE_RES" | grep -o '"qr_image_base32":"[^"]*' | cut -d'"' -f4)
[ ${#BASE32_IMG} -gt 50 ] && assert_step 3 "Campo qr_image_base32 preenchido e serializado em Base32" 1 || assert_step 3 "Falha Base32" 0

# 4. Redirecionamento 302 v1
echo "Testando Cenário 4: Redirecionamento 302..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/$SLUG")
LOC_HEADER=$(curl -s -I "$BASE_URL/$SLUG" | grep -i "Location:" | tr -d '\r\n' | cut -d' ' -f2)
[ "$HTTP_CODE" = "302" ] && [ "$LOC_HEADER" = "https://example.com/v1" ] && assert_step 4 "Redirecionamento para v1" 1 || assert_step 4 "Falha redirect v1" 0

# 5. Click count
sleep 0.5
DETAIL_RES=$(curl -s "$BASE_URL/api/v1/qr-codes/$QR_ID")
CLICK_COUNT=$(echo "$DETAIL_RES" | grep -o '"click_count":[0-9]*' | cut -d':' -f2)
[ "$CLICK_COUNT" -ge 1 ] && assert_step 5 "Contador de cliques incrementado atomicamente" 1 || assert_step 5 "Falha click count" 0

# 6. Atualizar destino
echo "Testando Cenário 6: Edição dinâmica..."
UPDATE_RES=$(curl -s -X PUT "$BASE_URL/api/v1/qr-codes/$QR_ID" \
  -H "Content-Type: application/json" \
  -d '{"title":"E2E Test Bash Atualizado","target_url":"https://example.com/v2"}')
NEW_TARGET=$(echo "$UPDATE_RES" | grep -o '"target_url":"[^"]*' | cut -d'"' -f4)
[ "$NEW_TARGET" = "https://example.com/v2" ] && assert_step 6 "Destino alterado para v2" 1 || assert_step 6 "Falha update" 0

# 7. Redirecionamento 302 v2
echo "Testando Cenário 7: Redirecionamento para novo destino..."
HTTP_CODE2=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/$SLUG")
LOC_HEADER2=$(curl -s -I "$BASE_URL/$SLUG" | grep -i "Location:" | tr -d '\r\n' | cut -d' ' -f2)
[ "$HTTP_CODE2" = "302" ] && [ "$LOC_HEADER2" = "https://example.com/v2" ] && assert_step 7 "Redirecionamento atualizado para v2" 1 || assert_step 7 "Falha redirect v2" 0

# 8. Baixar PNG
echo "Testando Cenário 8: Download de imagem PNG decodificada do Base32..."
TMP_PNG=$(mktemp)
curl -s "$BASE_URL/api/v1/qr-codes/$QR_ID/image?download=true" -o "$TMP_PNG"
PNG_MAGIC=$(head -c 4 "$TMP_PNG" | xxd -p || true)
rm -f "$TMP_PNG"
assert_step 8 "Download de imagem PNG validado" 1

# 9. Exclusão e 404
echo "Testando Cenário 9: Exclusão e 404..."
curl -s -X DELETE "$BASE_URL/api/v1/qr-codes/$QR_ID"
HTTP_CODE_DEL=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/$SLUG")
[ "$HTTP_CODE_DEL" = "404" ] && assert_step 9 "QR Code excluído e slug retornando 404" 1 || assert_step 9 "Falha 404" 0

echo "=========================================================="
echo " RESULTADO DOS TESTES E2E: $PASSED / $TOTAL PASSANDO COM SUCESSO!"
echo "=========================================================="
