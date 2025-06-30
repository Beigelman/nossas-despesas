#!/bin/bash

# Extrair a cobertura total
gocov_total=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

# Determinar emoji baseado na cobertura
declare EMOJI STATUS
if (( $(echo "$gocov_total >= 80" | bc -l) )); then
  EMOJI="🟢"
  STATUS="Excelente"
elif (( $(echo "$gocov_total >= 60" | bc -l) )); then
  EMOJI="🟡"
  STATUS="Boa"
else
  EMOJI="🔴"
  STATUS="Baixa"
fi

# Gerar relatório condensado
echo "## ${EMOJI} Relatório de Cobertura de Código" > coverage_report.md
echo "" >> coverage_report.md
echo "### 📊 **Cobertura Total: ${gocov_total}%** (${STATUS})" >> coverage_report.md
echo "" >> coverage_report.md

echo "### 📦 Cobertura por Pacote:" >> coverage_report.md
echo '```' >> coverage_report.md
go test -cover ./internal/... | grep 'coverage:' | sed -E 's/ok[[:space:]]+//;s/\?[[:space:]]+//' | sed -E 's/(.*)[[:space:]]+coverage: ([0-9.]+)% of statements/\1 \2/' | while read pkg cov; do
  cov_num=$(echo $cov | sed 's/%//')
  if (( $(echo "$cov_num >= 80" | bc -l) )); then
    emoji="🟢"
  elif (( $(echo "$cov_num >= 60" | bc -l) )); then
    emoji="🟡"
  else
    emoji="🔴"
  fi
  echo "$emoji $pkg: $cov%" >> coverage_report.md
done
echo '```' >> coverage_report.md
echo "" >> coverage_report.md
echo "---" >> coverage_report.md
echo "*Relatório gerado automaticamente pelo GitHub Actions*" >> coverage_report.md

echo "coverage=${gocov_total}" >> $GITHUB_OUTPUT 