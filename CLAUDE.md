# Sandstorm Server Assist

App desktop (Windows) para instalar, configurar e iniciar um servidor dedicado de Insurgency: Sandstorm com
mods do mod.io. `desktop/` é o projeto Wails 2 (Go, Clean/Hexagonal) e `desktop/frontend/` o React + Vite
(workspace `@ssa/frontend`). Detalhes no README; decisões em `docs/DECISOES.md`; andamento em
`docs/PLANO_DESENVOLVIMENTO.md`; última sessão em `docs/RESUMO_SESSAO.md`. O histórico de como o servidor foi
configurado à mão está em `resumo_sessao.md` (raiz): é a fonte das regras do mod.io que o app aplica.

## Regras do projeto

- **Plano sempre atualizado:** ao terminar uma etapa ou decidir algo, atualize `docs/PLANO_DESENVOLVIMENTO.md`
  (e `docs/DECISOES.md` se for decisão de arquitetura). No fim da sessão, atualize `docs/RESUMO_SESSAO.md`.
- **Segurança a cada mudança de código:** rode a skill `security-audit` (no mínimo `npm run audit:code` e a
  revisão do diff) e corrija antes de concluir. Nunca `npm audit fix --force`.
- **Git:** commit e push só com autorização do usuário.
- **Regras só no Go** (`internal/domain` puro + `internal/app` casos de uso). O front só mostra e coleta escolhas.
- **O front nunca manda comando nem caminho livre:** manda ids e opções. Caminhos só vêm de diálogos do Windows
  abertos pelo Go. O Go monta os argumentos do servidor e do SteamCMD (lista, sem shell).
- **Segredos** (senhas, tokens, sessão do mod.io) só no cofre DPAPI (`adapters/vault`) ou no backup cifrado.
  Nunca em log, nunca no front sem o usuário clicar em "Mostrar". Toda linha de log mostrada passa por
  `domain/logparse.Redact`.
- **Erros** sempre como item do catálogo `domain/problems` (o que aconteceu, causa, solução).
- **Links externos** só por chave (`links` em `adapters/desktop/system_api.go`).
- **Wails:** depois de mudar métodos expostos em `internal/adapters/desktop`, rode `npm run gen:bindings`.
- **Textos da tela em 5 idiomas** (ADR-009): en-US é a referência (`desktop/frontend/src/i18n/locales/en-US.ts`);
  pt-PT, es-ES, ru e uk com as mesmas chaves (o `npm test` confere). O Go não tem texto para o usuário: devolve
  códigos (`problems.*`) e chaves (`validation.*`). Linguagem simples.
- **Projeto que será público** (ADR-008, ADR-012): Apache 2.0; nada de chave, token, IP, e-mail ou nome real
  no código, nos testes ou nos docs. `resumo_sessao.md` (notas do servidor do usuário) fica fora do git.
- **Termos e dados** (ADR-015): mudou `TERMS.md`/`PRIVACY.md` ou os textos `legal.*`, suba `legal.Version`.
- **Aparência** (ADR-011): cores só por variáveis CSS; tema novo em `styles/themes.css` + `shared/appearance.ts`.

## Comandos

- `npm install` → dependências.
- `npm run dev` → app em modo desenvolvimento.
- `npm run build` → `desktop/build/bin/sandstorm-server-assist.exe`.
- `npm run build:installer` → instalador NSIS em `desktop/build/dist` (ADR-014); `npm run package` → zip portátil.
- `npm run site:preview` / `npm run site:publish` → página de download (ADR-016, `docs/PAGINA_DOWNLOAD.md`;
  publicar envia ao servidor: só com autorização do usuário).
- `npm run update:manifest` → `latest.json` assinado (chave em `~/.sandstorm-server-assist/`, nunca no git).
- `npm test` (Go + front), `npm run lint`, `npm run typecheck`.
- `npm run audit`, `npm run audit:go`, `npm run audit:code`.
