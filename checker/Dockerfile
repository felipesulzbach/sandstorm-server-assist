# Verificador de portas do Sandstorm Server Assist (ADR-021). Roda atrás do Caddy do servidor da página de
# download, que monta a imagem direto do GitHub (deploy/compose.prod.yaml do game-in).
# Build local: docker build -t ssa-checker checker

FROM golang:1.27.4-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
# Binário estático (sem cgo), sem informações de depuração. Só a biblioteca padrão do Go.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/checker .

# Imagem final mínima; roda como usuário sem privilégios.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/checker /checker
EXPOSE 8095
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 CMD ["/checker", "healthcheck"]
ENTRYPOINT ["/checker"]
