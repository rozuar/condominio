#!/bin/bash

# Script para ejecutar linting y formateo

set -e

echo "=========================================="
echo "  Condominio App - Lint & Format"
echo "=========================================="

cd "$(dirname "$0")/.."

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_status() {
    echo -e "${GREEN}[OK]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Ejecutar Detekt con auto-correct
echo ""
echo "Ejecutando Detekt con auto-correct..."
if ./gradlew detekt --auto-correct; then
    print_status "Detekt completado"
else
    print_warning "Se encontraron problemas que requieren correccion manual"
fi

# Android Lint
echo ""
echo "Ejecutando Android Lint..."
if ./gradlew lintDebug; then
    print_status "Android Lint completado"
else
    print_warning "Se encontraron advertencias de lint"
fi

echo ""
echo "=========================================="
echo "  Reportes generados:"
echo "=========================================="
echo "  - Detekt: app/build/reports/detekt/detekt.html"
echo "  - Android Lint: app/build/reports/lint-results-debug.html"
echo ""
