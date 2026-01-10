#!/bin/bash

# Script para ejecutar tests y generar reportes

set -e

echo "=========================================="
echo "  Condominio App - Test Runner"
echo "=========================================="

cd "$(dirname "$0")/.."

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Funcion para imprimir con color
print_status() {
    echo -e "${GREEN}[OK]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Limpiar build anterior
echo ""
echo "Limpiando build anterior..."
./gradlew clean

# Ejecutar tests unitarios
echo ""
echo "Ejecutando tests unitarios..."
if ./gradlew testDebugUnitTest; then
    print_status "Tests unitarios completados"
else
    print_error "Algunos tests fallaron"
    exit 1
fi

# Generar reporte de coverage
echo ""
echo "Generando reporte de coverage..."
if ./gradlew koverHtmlReportDebug; then
    print_status "Reporte de coverage generado"
    echo "  Ver: app/build/reports/kover/htmlDebug/index.html"
else
    print_warning "Error generando reporte de coverage"
fi

# Ejecutar Detekt
echo ""
echo "Ejecutando analisis estatico (Detekt)..."
if ./gradlew detekt; then
    print_status "Analisis estatico completado"
else
    print_warning "Se encontraron problemas de estilo"
fi

echo ""
echo "=========================================="
echo "  Reportes generados:"
echo "=========================================="
echo "  - Tests: app/build/reports/tests/testDebugUnitTest/index.html"
echo "  - Coverage: app/build/reports/kover/htmlDebug/index.html"
echo "  - Detekt: app/build/reports/detekt/detekt.html"
echo ""
