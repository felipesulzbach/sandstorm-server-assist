# Página de download e publicação (ADR-016)

A página fica no mesmo servidor do game-in e do midia-multi-tools, em `https://game-in.duckdns.org/sandstorm/`. O app
instalado lê `https://game-in.duckdns.org/sandstorm/downloads/latest.json` para se atualizar (ADR-014).

## Arquivos

- `site/`: modelo da página (`index.html` com marcadores `{{...}}`, `style.css`, `lang.js` com os 5 idiomas,
  `favicon.svg`). Sem script nem estilo inline (CSP do Caddy).
- `scripts/publish-site.mjs`: monta a página com versão, tamanhos e SHA-256, gera `terms.html` e `privacy.html` a
  partir de `TERMS.md`/`PRIVACY.md`, assina o `latest.json` e envia tudo por SSH.

## Preparar o servidor (uma vez, no repositório do game-in)

Já aplicado no checkout local do game-in em 09/11/2026 (falta commit e levar ao servidor). O trecho fica aqui como
referência.

1. Em `deploy/compose.prod.yaml`, no serviço do Caddy, montar a pasta (ao lado da do midia):

   ```yaml
   # Página e instalador do Sandstorm Server Assist (/sandstorm/), enviados por aquele projeto (fora do Git).
   - ./sandstorm:/srv/sandstorm:ro
   ```

2. Em `deploy/Caddyfile`, ao lado do bloco do midia:

   ```caddyfile
   # Sandstorm Server Assist (outro projeto do mesmo dono): página de download, instalador e latest.json da
   # atualização automática em /sandstorm/. Os arquivos ficam em deploy/sandstorm (fora do Git) e são enviados
   # pelo "npm run site:publish" daquele projeto.
   redir /sandstorm /sandstorm/ 308
   handle_path /sandstorm/downloads/* {
   	root * /srv/sandstorm/downloads
   	header Content-Disposition attachment
   	header Cache-Control "no-cache"
   	file_server
   }
   handle_path /sandstorm/* {
   	root * /srv/sandstorm/site
   	import site_headers
   	file_server
   }
   ```

3. No servidor: criar `game-in/deploy/sandstorm/` e recriar o Caddy (`docker compose ... up -d caddy`).

## Publicar uma versão

```bash
npm run build:installer
npm run package
npm run site:preview
```

Conferir a página em `desktop/build/dist/site`. Depois, com o endereço SSH do servidor do game-in:

```bash
SSA_SSH_HOST=usuario@endereço npm run site:publish
```

O script envia o instalador e o zip, troca a página de uma vez e guarda as 3 últimas versões. Por último, envia o
`latest.json`: só então os apps instalados veem a versão nova. No fim, confere pela internet a página, o download
e o manifesto.

- A chave de assinatura precisa estar em `~/.sandstorm-server-assist/update-signing-key.pem`. O script recusa uma
  chave diferente da do app.
- `SSA_SSH_KEY` muda a chave SSH (padrão `~/.ssh/game-in.pem`). `SSA_SITE_URL` muda o endereço conferido.
- O endereço do servidor não fica no código porque o projeto será público.
