# Plano de desenvolvimento

> Documento vivo: atualizado a cada etapa. Última atualização: 2026-11-11 (0.1.9 publicada).
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

## Fase 16 · Trocar mapa com o servidor ligado (10/11/2026)

- `domain/rcon` (pacotes), `adapters/rcon` (cliente em 127.0.0.1), `domain/server/scenarios.go` (`ParseScenarios`,
  `SplitScenario`, `MapChoices`, `TravelCommand`, `TravelRefused`), `app/mapchange` (lista e troca, uma conexão
  por ação), `ServerAPI.MapChange/Scenarios/ChangeMap`, cartão "Trocar mapa" no Painel, textos nos 5 idiomas.
- Verificado: testes novos (pacotes, servidor RCON falso com senha certa/errada e resposta em pedaços, regras da
  lista e da troca), npm test, lint, typecheck, audit:code. No servidor real, só leitura: RCON conectou, 278
  cenários do jogo base, mais 6 mapas de mods da rotação; cartão conferido no modo dev. A troca em si não foi feita.

## Fase 17 · Mensagens do servidor (10/11/2026, versão 0.1.3)

- `domain/server/messages.go` (validação, `MotdFile`, `ParseMotd`, `ChatCommand`), `logparse` com
  `PlayerLogin`/`PlayerJoined`/`PlayerLeft` e `Stamp`, `app/chat` (fila, conexão única, reconexão), runner
  `OnPlayer`, `config.Prepare`/`Detect` com `Motd.txt`, campo "motd" no "Trazer do servidor", aba Mensagens.
- Verificado: testes (linhas no formato do log real, nomes cortados, saída de quem não entrou, linhas antigas,
  reconexão, `Motd.txt`), aba conferida no modo dev. O envio no chat não foi testado com jogadores.

## Fase 18 · Rede neutra e teste de portas (10/11/2026, versão 0.1.3)

- Radmin fora do app (textos, link, guia, exemplo de IP). Guias "portas no roteador" e "programa de VPN".
- `checker/` (verificador, testes com servidor falso e um teste real opcional `SSA_A2S_TEST`), `domain/netcheck`,
  `adapters/portcheck`, `ServerAPI.CheckPorts`, bloco "Testar portas pela internet" em Servidor › Rede, problemas
  `port-check-*`. game-in: serviço `sandstorm-check` e rota `/sandstorm/api/` (falta subir no AWS).
- Termos/aviso de dados 2026-11-10. `npm test`, `vet` e `audit:code` passam a incluir o `checker`.
- Verificado: consulta A2S contra o servidor real desta máquina (só leitura: mapa e vagas lidos); no modo dev, o
  teste mostra o erro esperado enquanto o verificador não está no AWS.

## Fase 19 · Mais configurações do servidor (10/11/2026, versão 0.1.3)

- Regras do jogo (catálogo de 60 opções, aba Regras do jogo), mutators com descrição, dia ou noite (início, rotação,
  troca de mapa), painel de jogadores (lista, expulsar, banir, tirar banimento, mensagem a todos), mensagens
  programadas, reinício automático (diário com aviso e atualização; ao fechar sozinho), avisos no Discord, abrir com
  o Windows e iniciar o servidor junto, regra do Firewall (UAC), placa de rede (`-MultiHome`), tick rate e banda,
  vários perfis, lista de mods (`-CmdModList`). ADR-022.
- Verificado: testes (domínio, players, watchdog, notify, perfis, regras no Game.ini/Engine.ini, winsys só
  leitura), `listplayers`/`listbans`/`help` lidos do servidor real (formato e sintaxe), telas no modo dev
  (abas, catálogo, mutators, rotação com dia/noite, validação, rede, firewall, definições).
- Não testado ao vivo: ver ADR-022.

## Versão 0.1.3 e verificador de portas publicados (10/11/2026, ~13:30)

- **Verificador** (`sandstorm-check`, ADR-021) no servidor: `git pull` do game-in (`35e366a`), `build` e `up -d` do
  serviço (saudável) e Caddy recriado com `--force-recreate` (sem isso o contêiner continua com o Caddyfile antigo:
  o `git pull` troca o arquivo e a montagem segue no antigo). Conferido de fora: `/sandstorm/api/healthz` 200;
  `check` responde "fechada" sem servidor do jogo, recusa porta inválida e ignora `X-Forwarded-For` falso.
- **0.1.3** publicada no PC de publicação: testes Go (com o `checker`) e front, lint, typecheck, `audit`,
  `audit:go`, `audit:code` limpos; `build:installer`, `package`, `site:publish`. O `latest.json` aponta para a
  0.1.3 e o SHA-256 do instalador baixado confere; página, termos (2026-11-10) e privacidade no ar; game-in e
  `/midia/` seguiram respondendo 200.

## Fase 20 · Mutators dos mods (10/11/2026, versão 0.1.4)

- `domain/pak` (índice do .pak, teste com fuzz), `modio.Mutators`, `ConfigAPI.ModMutators`, lista com rolagem e
  marca MOD, sem o campo de texto. ADR-023.
- **Publicada** no PC de publicação: testes, lint, typecheck, `audit`, `audit:go` e `audit:code` limpos;
  `build:installer`, `package`, `site:publish`. `latest.json` na 0.1.4, SHA-256 do instalador baixado confere,
  página mostra 0.1.4; verificador, game-in e `/midia/` seguiram respondendo 200. Sem mudança no servidor.

## Fase 21 · Formatação das mensagens e mapa de mod da lista (10/11/2026, versão 0.1.5)

- `domain/server/markup.go`, `ConfigAPI.MarkupPreview`, editor com barra e prévia (mensagem do dia, entrada, saída
  e programadas), `pak.Scan` com mapas e cenários, `modio.Scenarios`, seletores do mapa de mod. ADR-024.
- **Publicada** no PC de publicação: testes, lint, typecheck, `audit`, `audit:go` e `audit:code` limpos;
  `build:installer`, `package`, `site:publish`. `latest.json` na 0.1.5, SHA-256 do instalador baixado confere,
  página mostra 0.1.5; verificador, game-in e `/midia/` seguiram respondendo 200. Sem mudança no servidor.

## Fase 22 · Correção da tela preta e proteção contra erro de tela (10/11/2026, versão 0.1.6)

- Listas vazias em vez de nulas no Go; tela tolerante a `null`; `ErrorBoundary` por página e no app; erro de tela
  nos Registos. ADR-025.
- **Publicada** no PC de publicação: testes, lint, typecheck, `audit`, `audit:go` e `audit:code` limpos;
  `build:installer`, `package`, `site:publish`. `latest.json` na 0.1.6, SHA-256 do instalador baixado confere,
  página mostra 0.1.6; verificador, game-in e `/midia/` seguiram respondendo 200. Sem mudança no servidor.

## Fase 23 · O app reaberto continua com o próprio servidor (10/11/2026, versão 0.1.7)

- Registro da partida em `server-run.json`; ao reabrir, o servidor iniciado pelo app continua como dele (linha de
  comando, reinício para carregar mods novos). Aviso de servidor de fora reescrito nos 5 idiomas. ADR-026.
- **Publicada** no PC de publicação: testes, lint, typecheck, `audit`, `audit:go` e `audit:code` limpos;
  `build:installer`, `package`, `site:publish`. `latest.json` na 0.1.7, SHA-256 do instalador baixado confere,
  página mostra 0.1.7; verificador, game-in e `/midia/` seguiram respondendo 200. Sem mudança no servidor.

## Fase 24 · Painel preso em "a carregar o mapa" depois de trocar de mapa (10/11/2026, versão 0.1.8)

- A espera pelo mapa de mod da inicialização acaba na primeira viagem; servidor de fora não espera. ADR-027.
- `LogWorldSettings: Error: No "Max Raycast Distance" defined` (Border) conta como aviso, não como erro; o
  contador dos Registos usa o singular com 1 erro.
- **Publicada** em 11/11/2026 no PC de publicação: testes, lint, typecheck, `audit`, `audit:go` e `audit:code`
  limpos; `build:installer`, `package`, `site:publish`. `latest.json` na 0.1.8, SHA-256 do instalador baixado
  confere, página mostra 0.1.8; verificador, game-in e `/midia/` seguiram respondendo 200. Sem mudança no servidor.

## Fase 25 · Painel mostrava o mapa anterior depois de trocar para um mapa de mod (11/11/2026, versão 0.1.9)

- A viagem para mapa de mod vem com o caminho do mapa no log; o app passou a lê-la. ADR-028.
- **Publicada** em 11/11/2026 no PC de publicação: testes, lint, typecheck, `audit`, `audit:go` e `audit:code`
  limpos; `build:installer`, `package`, `site:publish`. `latest.json` na 0.1.9, SHA-256 do instalador baixado
  confere, página mostra 0.1.9; verificador, game-in e `/midia/` seguiram respondendo 200. Sem mudança no servidor.

## Fase 26 · Troca de mapa com os mods carregados e listas de mapa e modo (11/11/2026, versão 0.1.10)

- Mapas de mod da troca vêm dos mods ativos na partida (não mais da rotação); listas separadas de mapa e modo.
  ADR-029.
- Mapas da rotação sem mod: aviso com "Remover" em Servidor › Mapas (existência lida dos .pak do jogo base e dos
  mods instalados) e item desativado na troca de mapa. ADR-030.
- Mensagens sem formatação: o jogo mostrou `<Bold>` escrito na tela de carregamento e no chat. Editor de
  formatação removido; tags recusadas ao salvar e tiradas do `Motd.txt` e do chat. ADR-031.

## Pendências

- Feito depois da primeira rodada: ícone, zip portátil, testes de installer/config/logbuf e teste real do
  download do SteamCMD (`SSA_STEAMCMD_E2E=1`, passou).
- Validação humana na janela nativa: escolher pasta, importar o servidor de `D:\_server`, mover a senha de
  RCON para o cofre, iniciar/parar, backup/restauração com os diálogos do Windows.

## Ideias de evolução

- Detectar conflito de mods antes de iniciar (procurar caminhos repetidos nos `.pak`, como feito à mão).
- Repositório público no GitHub com CI (quando o usuário pedir).
