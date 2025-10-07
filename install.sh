#!/bin/bash

BIN_NAME="pkm"

INSTALL_DIR="/usr/local/bin"

echo "🔨 Compilando proyecto Go..."
go build -o $BIN_NAME


if [ ! -f "$BIN_NAME" ]; then
  echo "❌ Error: no se pudo compilar el binario."
  exit 1
fi

echo "✅ Compilación exitosa."

echo "🔑 Dando permisos de ejecución..."
chmod +x $BIN_NAME

echo "📦 Moviendo $BIN_NAME a $INSTALL_DIR (puede pedir contraseña)..."
sudo mv $BIN_NAME $INSTALL_DIR/

echo "🔍 Verificando instalación..."
if command -v $BIN_NAME > /dev/null 2>&1; then
  echo "🎉 Instalación completa. Puedes usar tu CLI con: $BIN_NAME"
  $BIN_NAME version || echo "ℹ️ Ejecuta '$BIN_NAME' para probarlo."
else
  echo "⚠️ Algo salió mal, $BIN_NAME no está en tu PATH."
fi
