# Plano de desenvolvimento

> Documento vivo: atualizado a cada etapa. Última atualização: 2026-11-09, por volta de 21:00 (página publicada).
> Base: [RESUMO_PROJETO.md](RESUMO_PROJETO.md) · Decisões: [DECISOES.md](DECISOES.md)

## Situação

Versão **0.1.0** concluída: back-end e telas prontos, 5 idiomas, temas e skins, ícone próprio, `.exe` e zip
portátil gerados. Pendente só a validação humana na janela nativa (diálogos do Windows, importar o servidor
real, iniciar/parar).

| Fase                             | Entrega                                                                                                      | Situação       |
| -------------------------------- | ------------------------------------------------------------------------------------------------------------ | -------------- |
| 0. Fundação                      | Docs, `CLAUDE.md`, Wails 2 + React, scripts, lint, testes, skill `security-audit`                            | ✅ Concluída   |
| 1. Perfil, cofre e configuração  | Perfil, cofre DPAPI, tela Servidor (5 abas) com guias "Como conseguir?", importar servidor existente         | ✅ Concluída   |
| 2. Instalação com progresso      | Baixar SteamCMD (assinatura Valve conferida), instalar/atualizar/verificar com barra, velocidade e tempo     | ✅ Concluída\* |
| 3. Iniciar e acompanhar          | Linha de comando no Go, iniciar/parar/reiniciar, fases pelo log, servidor aberto fora do app, conexão        | ✅ Concluída   |
| 4. mod.io                        | Pedir código, autorizar (volta sozinho para none), detectar chave de API no log, mods com progresso          | ✅ Concluída\* |
| 5. Erros e logs                  | Catálogo de erros (o que/causa/solução + links), gaveta de logs fechada por padrão, filtros, máscara         | ✅ Concluída   |
| 6. Backup e restauração          | `.ssbackup` (Argon2id + AES-256-GCM), criar, conferir e restaurar em outra pasta/máquina                     | ✅ Concluída   |
| 7. Idiomas                       | en-US (padrão), pt-PT, es-ES, ru, uk; teste confere as chaves de todos                                       | ✅ Concluída   |
| 8. Projeto público               | Apache 2.0, NOTICE, DISCLAIMER, README em inglês, CONTRIBUTING, SECURITY, dados pessoais fora do código      | ✅ Concluída   |
| 9. Aparência                     | Temas Areia/Preto/Oliva/Azul/Claro e skins Lateral/Compacta/Cartões/Topo (ADR-011)                           | ✅ Concluída   |
| 10. Distribuição                 | Ícone próprio (`npm run icon`), zip portátil com LICENSE/NOTICE/DISCLAIMER (`npm run package`)               | ✅ Concluída   |
| 11. Servidor existente (ADR-013) | Completar os campos com o servidor já configurado (log + ini/txt), perguntar manter ou sobrescrever          | ✅ Concluída   |
| 12. Instalador e atualização     | Instalador por usuário, atualização automática assinada, aceite dos termos e do aviso de dados (ADR-014/015) | ✅ Concluída   |

\* Verificado por testes e com o app aberto; o download real do SteamCMD/servidor e o pedido de código ao
mod.io não foram executados de ponta a ponta (dependem da internet e da conta do usuário).

## Verificações (09/11/2026)

- `go test ./...` e `go vet ./...`: passam (domínio, runner com log real, backup completo com import e
  restauração, cofre DPAPI real, assinatura do steamcmd.exe real da Valve, cliente mod.io bloqueando outros hosts).
- `npm test` (12 testes: chaves de todos os idiomas, códigos de erro do Go com texto, ErrorCard, Guide, formatos).
- `npm run lint`, `npm run typecheck`, `npm run audit` (só a exceção de dev herdada, até 2026-12-31),
  `npm run audit:go` (nada alcançável), `npm run audit:code` (gosec limpo, supressões com motivo).
- `npm run build`: `.exe` de 13 MB.
- App aberto em modo dev: todas as telas abrem, tema e skin trocam, sessão real do mod.io lida.

## Fase 11 · Servidor existente (09/11/2026, pedido do usuário)

- Domínio: `server/detect.go` (`ParseCommandLine`, `Compare`, `Merge`, `ValidAdmins`) com testes.
- `app/config`: `LogFiles` (por data, inclui `Insurgency_N.log`), `Detect`, `ScanExisting`, `ApplyExisting`,
  `ResolveServerDir` (pasta de cima → subpasta com o servidor; corrige config antiga ao abrir). Substitui o
  antigo `ImportExisting`.
- Ponte: `ConfigAPI.ScanExisting` / `ApplyExisting(keepMine)`. Front: `features/sync/ExistingSync.tsx` (janela com
  conflitos, preenchimentos e iguais), usado na Instalação (sozinho ao escolher a pasta + botão) e no Servidor.
- Textos nos 5 idiomas (`sync.*`). Estilos `.modal*` e `.sync*` com variáveis (`--backdrop`, `--shadow-modal`).
- Verificado: go test/vet, npm test, lint, typecheck, audit:code, build e zip; leitura real (só leitura) do
  servidor em `D:\_server`; janela aberta no modo dev com os dados reais e cancelada (nada aplicado).
- Achado no teste real: a configuração apontava para a pasta de cima do servidor; agora o app corrige sozinho.
- Correção junto: "Em execução desde" de um servidor aberto fora do app mostrava a hora em que o app o viu.
  Agora `proc.FindInDir` devolve a hora real de início (`GetProcessTimes`) e o runner usa a mais antiga
  (`startedAt`, com teste). Conferido no servidor real (01:39:21). `.exe` e zip regerados.

## Fase 12 · Instalador, atualização automática e termos (09/11/2026, pedido do usuário)

- Atualização automática (ADR-014, modelo do midia-multi-tools): `domain/update`, `app/updater`,
  `adapters/updater` (via `netx`, host `game-in.duckdns.org` na lista), chave Ed25519 criada fora do git
  (`~/.sandstorm-server-assist/update-signing-key.pem`; a pública em `pubkey.go`), scripts
  `update:keygen`/`update:pubkey`/`update:manifest`. Botão "Atualizar" na barra de título e cartão "Atualizações"
  no Sobre. Código `update-failed` no catálogo.
- Instalador NSIS por usuário (`build/windows/installer/project.nsi`, `npm run build:installer`), 5 idiomas,
  `/S /relaunch` para a atualização. Gerado: `desktop/build/dist/sandstorm-server-assist-0.1.0-windows-setup.exe`
  (7 MB); zip portátil regerado com TERMS e PRIVACY.
- Termos de uso e aviso de dados (ADR-015): `TERMS.md`, `PRIVACY.md`, textos `legal.*` nos 5 idiomas, tela de
  aceite no primeiro uso, versão em `domain/legal`, aceite em `config.json`.
- Verificado: go test/vet, npm test (12), lint, typecheck, audit:code (gosec limpo), audit:go; `latest.json`
  assinado pelo script e conferido pelo Go com a chave do app (SHA-256 e tamanho batem); tela de aceite vista no
  modo dev (abas, idioma, botão só libera com a caixa marcada), sem aceitar.

## Fase 13 · Página de download (09/11/2026, pedido do usuário)

- `site/` (página em 5 idiomas, CSP sem inline), `scripts/publish-site.mjs` (`site:preview`, `site:publish`),
  termos e privacidade em HTML gerados dos `.md`, guia [PAGINA_DOWNLOAD.md](PAGINA_DOWNLOAD.md) com o trecho do
  Caddy e do compose do game-in (ADR-016).
- Verificado: prévia servida com o mesmo CSP do Caddy, sem erros no console; troca de idioma (pt automático, uk,
  en), termos/privacidade, celular sem rolagem lateral; lint e audit:code limpos.

## Fase 14 · Correções do primeiro uso real (09/11/2026, versão 0.1.1)

- Relato do usuário: muitos erros, só 2 de 7 mods carregados, mapa do jogo base em vez do Containment, mutator sem
  efeito, mod "Removido" na lista. Causa no log: inscrições novas fazem o servidor montar só os mods novos (ADR-017).
- Runner reinicia uma vez sozinho (`checkRemount`, aviso `dashboard.remount`), problema `mods-not-mounted`;
  `server.FixMapCycleLine` (8 linhas com `Scenario=` repetido no MapCycle do usuário); lista de Mods sem os
  removidos; `logparse.Level` trata erros conhecidos como aviso (no log real: 233 → 57 erros, e os 57 eram da rotação
  e dos mods não montados).
- Verificado: testes com as linhas reais (reinicia uma vez; não reinicia sem necessidade; servidor de fora só
  avisa), go vet, npm test, lint, typecheck, audit:code. Instalador e zip 0.1.1 em `desktop/build/dist`.

## Fase 15 · Logs por partida e limpeza (09/11/2026, versão 0.1.2)

- `domain/serverlogs` (nomes, ordem, `ToDelete`), `config.SetLogRetention`/`ServerLogs`/`DeleteOldLogs`,
  `logbuf.StartRun`/`Format`, runner abre partida e limpa uma vez depois da troca do log, `LogAPI`
  (`Runs`, `ReadRun`, `Cleanup`, `SetCleanup`, `CleanOld`), problema `log-delete`, `validation.logKeep`.
- Tela: gaveta com seletor de partida e divisória; cartão "Registos do servidor" nas Definições (desligado por
  padrão, 2 a 100, apagar agora com confirmação). Textos nos 5 idiomas.
- Cópia única do app (`SingleInstanceLock`), achada ao ver que o aceite dos termos sumiu do `config.json`.
- Verificado: testes novos (regras de retenção, apagar só logs do servidor, partida e limpeza no runner), npm test,
  lint, typecheck, audit:code; no modo dev com os logs reais: 14 partidas no seletor, 6 erros na atual, cartão e
  validação (1 recusado). Nada foi apagado. Instalador e zip 0.1.2.

## Repositório público (09/11/2026)

- https://github.com/felipesulzbach/sandstorm-server-assist, branch `main`. Publicação na AWS a partir de outro
  computador (precisa da chave de assinatura das atualizações, copiada à mão).

## Página de download publicada (09/11/2026, ~21:00)

- **No ar:** https://game-in.duckdns.org/sandstorm/ com a versão 0.1.2 (instalador, zip portátil e `latest.json`
  assinado). A atualização automática dos apps instalados passa a achar o manifesto.
- game-in: rota `/sandstorm/` no `main` (`ff6f4f3`, mais `8292f5a`, que pôs a linha da tabela no lugar certo do
  `docs/PLANO.md`); no servidor, `git pull`, pasta `deploy/sandstorm/` e só o Caddy recriado (`up -d --no-deps
  caddy`). game-in (`/healthz`, `/`) e `/midia/` continuaram respondendo 200.
- Gerado no PC de publicação: NSIS 3.13 instalado pelo winget (`NSIS.NSIS`), `npm run build:installer`, `package`,
  `site:preview` e `site:publish`.
- Conferido de fora: página, termos, privacidade, CSP, `Content-Disposition: attachment` nos downloads, SHA-256 do
  instalador baixado igual ao do `latest.json`, página no navegador sem erros no console.

## Pendências

- Feito depois da primeira rodada: ícone, zip portátil, testes de installer/config/logbuf e teste real do
  download do SteamCMD (`SSA_STEAMCMD_E2E=1`, passou).
- Validação humana na janela nativa: escolher pasta, importar o servidor de `D:\_server`, mover a senha de
  RCON para o cofre, iniciar/parar, backup/restauração com os diálogos do Windows.

## Ideias de evolução

- Detectar conflito de mods antes de iniciar (procurar caminhos repetidos nos `.pak`, como feito à mão).
- RCON pelo app (jogadores, trocar mapa, mensagens).
- Criar a regra do Firewall do Windows para a rede do Radmin (pede administrador).
- Repositório público no GitHub com CI (quando o usuário pedir).
