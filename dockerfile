# Utilise l’image officielle Go 1.24 (Alpine) comme builder
FROM golang:1.24-alpine AS builder

# Installe git pour le go mod download (si besoin de dépôts VCS) et les CA (pour certains modules)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# 1. Copie les fichiers de modules et télécharge les dépendances
COPY go.mod go.sum ./
RUN go mod download

# 2. Copie le code source et compile la version Linux statique
COPY . .
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags="-s -w" -o gotags

# --- image finale ---
FROM scratch

# Pour certains binaire statiques, il peut être utile de fournir les certificats
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

# Copie uniquement le binaire compilé depuis le builder
COPY --from=builder /app/gotags .

# Point d’entrée
ENTRYPOINT ["/app/gotags"]
