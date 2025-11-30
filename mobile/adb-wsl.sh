#!/bin/bash
# Script para usar adb.exe do Windows no WSL

ADB_PATH="/mnt/c/Users/micae/AppData/Local/Android/Sdk/platform-tools/adb.exe"

# Executar adb.exe do Windows com os argumentos passados
"$ADB_PATH" "$@"
