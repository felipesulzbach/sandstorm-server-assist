---
name: security-audit
description: Verifica e corrige vulnerabilidades no código e nas dependências do Sandstorm Server Assist (app desktop Wails com Go, front React em npm, cofre DPAPI, processos do SteamCMD e do servidor). Use depois de QUALQUER mudança de código, ao adicionar/atualizar dependências, quando pedirem auditoria de segurança ou quando o CI de segurança falhar.
---

# Auditoria e correção de vulnerabilidades

Adaptada da skill `security-audit` do projeto `midia-multi-tools`, com os riscos deste app: segredos do
servidor, processos externos (SteamCMD, InsurgencyServer.exe), rede (mod.io, Valve) e backups.

| Comando                | O que faz                                                                                   |
| ---------------------- | ------------------------------------------------------------------------------------------- |
| `npm run audit`        | `audit-ci` (falha a partir de _moderate_, config em `audit-ci.jsonc`) + `npm audit signatures` |
| `npm run audit:go`     | `govulncheck ./...` em `desktop/` (só o alcançável; exceções em `govulncheck-allowlist.json`) |
| `npm run audit:code`   | `gosec ./...` em `desktop/` + ESLint com regras de segurança em `desktop/frontend/src`       |

## Quando rodar o quê

| Mudança                               | Verificação mínima                                   |
| ------------------------------------- | ---------------------------------------------------- |
| Qualquer código                       | `npm run audit:code` + revisão do diff (abaixo)      |
| Dependência npm adicionada/atualizada | + `npm run audit`                                    |
| `go.mod`, versão do Go ou do Wails    | + `npm run audit:go`                                 |

## Revisão do diff (o que as ferramentas não pegam)

- **Segredos:** senha de entrada, RCON, GSLT, GameStats, chave de API e e-mail do mod.io só no cofre
  (`adapters/vault`, DPAPI com entropia) ou no backup cifrado. Nunca em log, `.ini`, `.bat`, mensagem de erro
  ou JSON para a tela. Todo texto de log passa por `domain/redact` (`logbuf.Add` já faz). Erros passam por
  `problems.New/Wrap` (mascaram o detalhe). A tela só recebe o valor com "Mostrar" (`RevealSecret`).
- **Processos:** `InsurgencyServer.exe` só com a linha montada por `server.BuildLaunch` (campos validados por
  regex: sem aspas, espaço, `?`, `&`). SteamCMD só com `steamcmd.Args` e só depois de `trust.VerifySigner(..., "Valve")`.
  Nada da tela vira argumento livre. Parar processos só os da pasta do servidor (`proc.FindInDir`).
- **Métodos expostos (Wails):** todo método exportado das `*API` em `adapters/desktop` é chamável pelo JS.
  Recebem ids, chaves e opções; caminhos só de diálogos abertos no Go; links só por chave da lista `links`.
- **Arquivos:** gravação atômica (`fsx.WriteAtomic`) com `.bak` antes da primeira troca de arquivo do usuário.
  Pastas escolhidas passam por `fsx.SafeDir`. Leituras com limite (`fsx.ReadLimited`).
- **Backup:** só nomes da lista `backup.ConfigFiles` + `modio/user.json`; limites de tamanho ao abrir;
  parâmetros do Argon2id limitados; cabeçalho autenticado.
- **Rede:** só HTTPS e hosts de `netx.AllowedHosts`; redirecionamento conferido.
- **Front:** nada de `dangerouslySetInnerHTML`, `eval` ou HTML vindo de log/nome de mod; CSP do build sem
  `unsafe-eval` e sem origens externas.
- **Projeto público:** nenhuma chave, token ou dado pessoal no repositório (a chave do mod.io não vai no
  código, ADR-010). `resumo_sessao.md` (notas do servidor do usuário) fica fora do git.

## Correção

1. npm: `npm audit fix` (nunca `--force`); Go: `go get <módulo>@<versão corrigida>` + `go mod tidy`.
2. Supressão do gosec só com `// #nosec G<regra> -- motivo` na mesma linha.
3. Sem correção possível: exceção com justificativa e data de expiração (`audit-ci.jsonc` ou
   `govulncheck-allowlist.json`).
4. Validar: `npm test`, `npm run lint`, `npm run typecheck`, `npm run build`.
5. Registrar abaixo as decisões que não são óbvias.

## Histórico de decisões

- 2026-11-09: projeto criado; supressões do gosec com motivo: G204 (processos do servidor, SteamCMD e
  explorer.exe com caminhos do perfil), G304 (leituras em caminhos montados pelo app), G103 (DPAPI e
  WinVerifyTrust exigem ponteiros), G115 (tamanhos limitados).
