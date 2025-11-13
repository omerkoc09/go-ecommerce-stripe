#!/bin/bash
# .env dosyasını yükle ve make komutlarını çalıştır

# .env dosyasını yükle
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
    echo "Environment variables loaded from .env"
else
    echo "Error: .env file not found!"
    exit 1
fi

# Make komutunu çalıştır
make "$@"
