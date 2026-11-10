# Decisões de arquitetura (ADRs)

Registro curto: contexto, decisão e consequência. Decisão nova entra no fim.

## ADR-001 · App desktop com Wails 2 (Go + React), mesma base do midia-multi-tools

- **Contexto:** o app precisa rodar processos do Windows (SteamCMD, servidor), ler e gravar arquivos do
  servidor, guardar segredos com a criptografia do Windows e mostrar telas com progresso. Havia dois projetos
  de referência: `game-in` (Electron + React, com API Go e banco) e `midia-multi-tools` (Wails 2 + React, app
  local que roda processos externos).
- **Decisão:** Wails **2.17.1** + Go 1.27.3 + React 19 + Vite, como o midia-multi-tools. A CLI do Wails roda
  pela mesma versão do `go.mod` (`scripts/desktop.mjs`), sem instalação global.
- **Reaproveitado do midia-multi-tools:** `scripts/desktop.mjs`, `gosec.mjs`, `govulncheck.mjs`, configs de
  ESLint/Prettier/audit-ci, a skill `security-audit` (adaptada), a barra de título própria, o tema claro/escuro,
  o armazenamento atômico de preferências e a organização dos docs.
- **Consequência:** app pequeno (usa o WebView2 do Windows 11), sem servidor nem banco. Regras em Go testáveis
  sem janela. O Electron do game-in ficou de fora: traria Chromium e Node inteiros sem ganho aqui.

## ADR-002 · Segredos no cofre do Windows (DPAPI), fora dos `.ini` e `.bat`

- **Contexto:** hoje a senha de RCON fica em texto puro no `Game.ini`, e os tokens GSLT/GameStats em texto puro
  num `readme.txt`. Qualquer cópia da pasta leva os segredos junto.
- **Decisão:** senha de entrada, senha de RCON, GSLT, GameStats e e-mail do mod.io ficam em
  `%APPDATA%\sandstorm-server-assist\vault.bin`, criptografados com DPAPI (`CryptProtectData`, escopo do usuário
  do Windows, com entropia própria do app). Só esse usuário, nessa máquina, consegue abrir.
  Na hora de iniciar, o app passa os segredos pela linha de comando (`?Password=`, `-Rcon -RconPassword=`,
  `-GSLTToken=`, `-GameStatsToken=`). O `Game.ini` deixa de ter a senha de RCON (o app oferece mover o valor
  existente para o cofre).
- **Na tela:** campos secretos vêm só como "definido / não definido". Ver o valor exige clicar em "Mostrar".
- **No log:** toda linha mostrada passa por uma máscara (`api_key=`, `Password=`, `Token=`, `-SecurityCode=`).
- **Limite conhecido:** a linha de comando de um processo pode ser lida por outros programas do mesmo usuário.
  É o único jeito que o servidor aceita esses valores sem deixá-los em disco.

## ADR-003 · O app inicia o `InsurgencyServer.exe` direto, sem `.bat`

- **Decisão:** o Go monta os argumentos a partir do perfil (domínio puro e testado) e inicia o executável com
  `exec.Command` (lista, sem shell), na pasta do servidor. Os `.bat` antigos continuam funcionando, mas o app
  não depende deles.
- O servidor **não morre se o app fechar**: ao abrir, o app procura processos `InsurgencyServer*.exe` da pasta
  configurada e volta a acompanhá-los. Parar encerra os processos dessa pasta.
- **Regras do mod.io que o app aplica** (comprovadas em 09/11/2026, ver `resumo_sessao.md`):
  `-SecurityCode=none` no uso normal; código novo só uma vez para autorizar, e depois volta para `none`;
  mapa inicial sempre vanilla e o mapa de mod pelo `-ModDownloadTravelTo`.

## ADR-004 · Estado do servidor lido do log

- **Decisão:** o app acompanha o `Insurgency.log` (lendo o que foi acrescentado) e traduz linhas conhecidas em
  eventos: mod.io pronto, porta aberta, mod ativado, instalando/desinstalando mod, viagem para o mapa, mudança
  de estado da partida, erro fatal. As regras de leitura ficam em `domain/logparse`, com testes feitos a partir
  de linhas reais.
- **Consequência:** se uma atualização do jogo mudar o texto do log, só a fase mostrada fica menos precisa;
  iniciar e parar continuam funcionando.

## ADR-005 · Catálogo de erros com causa e solução

- **Decisão:** `domain/problems` tem um catálogo: código, título, o que aconteceu, causa provável, como resolver
  e links (por chave). Erros do app (`AppError`) e linhas do log apontam para um item do catálogo. A tela mostra
  sempre os três blocos ("O que aconteceu", "Possível causa", "Como resolver").
- Erro não catalogado vira um item genérico com a mensagem original (mascarada) e a sugestão de abrir os logs.

## ADR-006 · Backup portátil criptografado com senha

- **Contexto:** o cofre DPAPI só abre na mesma máquina e no mesmo usuário. O backup precisa abrir em outra.
- **Decisão:** arquivo `.ssbackup` = cabeçalho + dados cifrados com **AES-256-GCM**, chave derivada da senha do
  backup com **Argon2id** (64 MiB, 3 passadas, sal aleatório). Dentro, um zip com `manifest.json`, o perfil, os
  segredos, os arquivos de configuração do servidor (`Game.ini`, `Engine.ini`, `Admins.txt`, `MapCycle.txt`,
  `Bans.json`, `Motd.txt`) e, se o usuário marcar, a sessão do mod.io (`user.json`).
- Ao restaurar: pede a senha, mostra o conteúdo, pede a pasta do servidor na máquina nova e grava tudo.
  Os segredos voltam para o cofre DPAPI da máquina nova.
- Senha do backup: mínimo de 10 caracteres. Sem a senha, o backup não abre (não há recuperação).

## ADR-007 · Downloads com progresso visível

- **SteamCMD:** baixado de `https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip` (endereço fixo), com
  progresso por bytes. Do zip, só o `steamcmd.exe` é extraído, e a assinatura digital dele é conferida
  (`WinVerifyTrust`) antes de rodar.
- **Servidor:** o SteamCMD escreve `Update state (0x61) downloading, progress: 12.34 (x / y)`. O domínio lê
  essas linhas e a tela mostra fase, porcentagem, bytes, velocidade e tempo restante.
- **Mods:** o tamanho de cada mod vem da API do mod.io (dados públicos), e o progresso é o tamanho da pasta do
  mod em disco comparado ao tamanho final. Os eventos do log (instalando, instalado) marcam início e fim.

## ADR-008 · Licença Apache 2.0 e aviso de isenção (projeto público)

- **Contexto:** o usuário quer publicar no GitHub para outras pessoas usarem, contribuírem e fazerem fork, sem
  responder por mau uso ou uso ilegal.
- **Decisão:** **Apache License 2.0** (`LICENSE`) mais um `DISCLAIMER.md`.
  - Permite uso, cópia, mudança, fork e uso comercial. Quem redistribui mantém a licença e o aviso de autoria,
    e marca os arquivos que mudou.
  - As seções 7 e 8 dizem que o software vem "AS IS", sem garantia, e que o autor não responde por danos.
  - A seção 6 não dá direito de uso de marcas: ajuda a separar o projeto de Insurgency, Steam e mod.io.
  - Traz uma licença de patentes explícita, que a MIT não tem.
- **DISCLAIMER:** projeto não oficial, sem ligação com New World Interactive, Focus Entertainment, Valve ou
  mod.io. Cada usuário responde pelo uso do servidor, pelos termos do jogo, da Steam e do mod.io, e pelas leis
  do seu país. O app não distribui o jogo nem o SteamCMD: baixa das fontes oficiais.
- **Consequência:** dependências usadas (Wails MIT, React MIT, x/crypto BSD) são compatíveis. Isto não é parecer
  jurídico: para uso comercial ou dúvidas, consultar um advogado.

## ADR-009 · Interface em vários idiomas (en-US padrão, pt-PT, es-ES, ru, uk)

- **Decisão:** os textos da tela ficam no front, em `desktop/frontend/src/i18n/locales/<idioma>.ts`. O inglês
  (`en-US.ts`) é a referência: os outros idiomas têm o mesmo tipo, e o TypeScript acusa chave faltando.
  Sem biblioteca externa (um `t(chave, parâmetros)` pequeno, menos dependências para auditar).
- O Go **não tem texto para o usuário**: devolve códigos (`problems.*`), chaves de validação (`validation.*`)
  com parâmetros e eventos. O catálogo de erros (ADR-005) passa a ter os textos no front, em
  `problems.<código>.{title,what,cause,solution}`.
- Idioma inicial: o do Windows, se for um dos cinco; senão, inglês. O usuário troca em Configurações, e a escolha
  fica nas preferências do app.
- Os docs de desenvolvimento e os comentários do código continuam em português por enquanto. O `README.md`
  público fica em inglês.

## ADR-010 · Chave de API do mod.io fora do código

- **Contexto:** o `.bat` antigo usa a chave de API do mod.io escrita no arquivo. Num repositório público ela
  ficaria exposta para qualquer um usar em nome de outra pessoa.
- **Decisão:** a chave é um segredo do cofre (`modioApiKey`). O app detecta a chave no `Insurgency.log` do
  próprio servidor (as requisições do jogo trazem `api_key=`) ou o usuário cria uma em
  https://mod.io/me/access, com um guia na tela.
- Endereço da API fixo no código: `https://g-254.modapi.io/v1` (jogo 254). Só HTTPS.

## ADR-011 · Temas de cor e skins de layout (como no game-in)

- **Contexto:** o usuário mostrou o game-in, que troca a cor e a forma do layout, e quis o mesmo aqui.
- **Decisão:** mesmo modelo do game-in (`lib/appearance.ts` + `themes.css`): `<html data-theme>` troca só as
  variáveis de cor e `<html data-skin>` muda a forma do layout. Cores: Areia (padrão, acento âmbar), Preto,
  Oliva, Azul escuro e Claro. Skins: Barra lateral (padrão, igual ao midia-multi-tools), Compacta (só ícones),
  Cartões (painéis flutuantes) e Barra no topo. Escolha salva no `localStorage` (vale por computador), com
  try/catch: sem armazenamento, fica o padrão. Aplicada antes do primeiro desenho (sem piscar).
- **Consequência:** tema novo = um bloco de variáveis em `themes.css` + item em `THEMES` + tradução.

## ADR-012 · Notas do servidor do usuário fora do repositório

- `resumo_sessao.md` (raiz) tem dados da rede do usuário (Radmin, e-mail): fica no `.gitignore`. Os docs do
  projeto (`docs/`) não citam IPs, e-mails nem tokens.

## ADR-013 · Completar os campos com os dados de um servidor existente

- **Contexto:** quem já tem um servidor configurado à mão (rodando ou parado) não quer digitar tudo de novo. Se
  já mudou algum campo no app, o app não pode apagar isso sem perguntar.
- **Decisão:** a fonte principal é a linha `LogInit: Command Line:` que o próprio servidor grava no começo de cada
  log. Vale para o servidor rodando (log atual) e parado (último log), sem depender de onde estão os `.bat`. Os
  logs são lidos do mais recente para o mais velho pela data (`Insurgency.log`, `Insurgency_2.log`,
  `Insurgency-backup-*.log`). Completam: `Admins.txt`, `MapCycle.txt`, `[Rcon]` do `Game.ini` e a chave do mod.io
  achada nos logs (ADR-010).
- **Regra de comparação** (domínio, `server.Compare`/`server.Merge`): igual → nada; no app vazio ou com o valor
  padrão → preenche sem perguntar; diferente do padrão e do servidor → **conflito**, e a tela pergunta campo a
  campo "Meu (app)" ou "Servidor" (com "manter todos os meus" e "usar todos do servidor"). O "Aplicar" só libera
  com todos os conflitos escolhidos. Na tela Servidor, mudanças não salvas são salvas antes, para contar como
  "meu".
- **Segurança:** segredos nunca vão para a tela (só "guardada no app" / "diferente no servidor"); a tela manda só as
  chaves que quer manter (`keepMine`, conferidas contra a lista); o Go lê de novo os arquivos ao aplicar; valor
  do servidor inválido (ex.: portas repetidas) volta ao do app e aparece como "não usado".
- **Pasta de cima:** se a pasta escolhida não tem o `InsurgencyServer.exe` mas uma única subpasta tem, o app usa a
  subpasta (também corrige uma configuração antiga ao abrir). Evita instalar um segundo servidor por engano.
- **Quando roda:** sozinho logo depois de escolher a pasta de um servidor instalado (só abre a janela se houver
  algo a trazer) e pelos botões "Trazer dados do servidor" (Instalação) e "Trazer do servidor" (Servidor).

## ADR-014 · Instalador por usuário e atualização automática com manifesto assinado

- **Contexto:** o usuário pediu que o app se atualize sozinho, como o game-in e o midia-multi-tools, e vai ter uma
  página de download parecida. O Wails não tem atualizador pronto.
- **Decisão:** o mesmo modelo do midia-multi-tools.
  - **Instalador NSIS por usuário** (`npm run build:installer` → `wails build -nsis -installscope user`):
    instala em `%LOCALAPPDATA%\Programs\Sandstorm Server Assist` sem pedir administrador, em 5 idiomas, com
    LICENSE, NOTICE, DISCLAIMER, TERMS e PRIVACY. Sai em
    `desktop/build/dist/sandstorm-server-assist-<versão>-windows-setup.exe`. Com `/S /relaunch` espera o app
    fechar, instala em silêncio e reabre o app.
  - **Manifesto `latest.json`** na pasta de downloads da página (`https://game-in.duckdns.org/sandstorm/downloads/`,
    mesmo servidor dos outros dois projetos): versão, nome do instalador (com a mesma versão), SHA-256, tamanho e
    assinatura **Ed25519**. A chave privada fica só na máquina de quem publica
    (`~/.sandstorm-server-assist/update-signing-key.pem`, `npm run update:keygen`); a pública está no app
    (`adapters/updater/pubkey.go`). `npm run update:manifest` assina e recusa se a chave não for a do app.
  - **App** (`domain/update` regras puras, `app/updater` caso de uso, `adapters/updater` HTTP/disco/lançador):
    procura 10 s depois de abrir e a cada 4 h, só pelo `netx` (HTTPS e host na lista fixa). Confere a assinatura
    antes de baixar e o SHA-256/tamanho depois. Pronta, a barra de título mostra **Atualizar** (instala e reabre);
    se o usuário só fechar o app, instala ao fechar sem reabrir. Não instala com instalação do servidor em
    andamento. A versão portátil (sem `uninstall.exe` ao lado) só avisa e abre a página.
  - **Servidor do jogo:** é outro processo e continua no ar durante a atualização; o app reaberto o reconhece
    como "aberto fora do app".
- **Publicar (quando a página existir):** `npm run build:installer`, `npm run update:manifest`, enviar o
  instalador e só depois o `latest.json`.
- **Consequência:** perder a chave privada impede atualizar os apps instalados (seria preciso instalar à mão uma
  versão com outra chave). Até a página existir, o app instalado procura e não acha nada ("não foi possível
  procurar agora"), sem atrapalhar.

## ADR-015 · Termos de uso e aviso de dados (aceite no primeiro uso)

- **Contexto:** o usuário pediu que o usuário aceite termos de uso e de dados, se fizesse sentido ter os dois.
- **Decisão:** sim, os dois, num único aceite (clickwrap, como no game-in):
  - **Termos de uso:** o que o app é (não oficial), sem garantia (Apache 2.0), responsabilidade de quem usa
    (termos do jogo, Steam, mod.io, leis), o que não é permitido, serviços de terceiros, atualização automática.
  - **Aviso de dados (privacidade):** mesmo sem coletar nada, o app guarda segredos e manda dados a terceiros
    (e-mail e chave ao mod.io, pedido de atualização à página). O aviso diz o que fica no PC, o que sai e para
    quem, o papel de quem roda o servidor com os dados dos jogadores e como apagar tudo. Sem conta, anúncios nem
    telemetria.
- **Como:** textos nos 5 idiomas (`legal.*`) e em `TERMS.md`/`PRIVACY.md`. A versão fica em `domain/legal.Version`;
  o aceite (versão e hora) fica nas preferências (`config.json`). Sem aceitar a versão em vigor, a tela de aceite
  cobre o app (a barra de título continua: dá para fechar), com escolha de idioma, as duas abas, a caixa "Li e
  aceito" e os botões Sair / Aceitar. A procura de atualizações só começa depois do aceite. O Sobre mostra a
  versão aceita e reabre os textos.
- **Consequência:** mudou o texto, sobe `legal.Version` e todos aceitam de novo.

## ADR-016 · Página de download no servidor do game-in

- **Contexto:** o usuário pediu uma página de download como a dos outros dois projetos, que também serve o
  `latest.json` da atualização automática (ADR-014).
- **Decisão:** mesmo modelo do midia-multi-tools: página estática em `site/` servida pelo Caddy do game-in em
  `/sandstorm/`, downloads em `/sandstorm/downloads/` (como anexo, sem cache). Sem script nem estilo inline (CSP
  `site_headers`); o inglês fica no HTML e `site/lang.js` troca para os outros 4 idiomas do app (idioma do
  navegador ou escolha na página, sem `innerHTML`). `terms.html` e `privacy.html` são gerados de `TERMS.md` e
  `PRIVACY.md`.
- **Publicação:** `scripts/publish-site.mjs` usa os arquivos de `desktop/build/dist` (ainda não há repositório nem
  CI): `npm run site:preview` só monta; `npm run site:publish` envia por SSH, troca a página de uma vez, guarda 3
  versões e manda o `latest.json` por último. O endereço SSH vem de `SSA_SSH_HOST` (o projeto é público: nada de
  IP no código). Passo a passo e trecho do Caddy em [PAGINA_DOWNLOAD.md](PAGINA_DOWNLOAD.md).
- **Consequência:** o repositório do game-in precisa da rota `/sandstorm/` e do volume `./sandstorm` (uma vez).

## ADR-017 · Reiniciar uma vez quando mods novos deixam os instalados de fora

- **Contexto:** no primeiro uso do instalador (09/11/2026), o servidor subiu em mapa do jogo base, sem o
  Containment e sem o mutator NoRestrictedArea. O log mostrou a regra do servidor: quando há inscrições novas
  (`PendingChanges with state Added`), ele baixa e monta só os mods novos e viaja para o `-ModDownloadTravelTo` sem
  montar os que já estavam instalados. A viagem cai no mapa inicial e o mutator fica "invalid". Sem nada novo, a
  partida seguinte monta todos (como nos logs que funcionaram).
- **Decisão:** o runner guarda os mods "Added" da partida. Na viagem, se houve mod novo e algum mod instalado
  (`state.json`, `State 1`) não foi ativado (`OnModActivatedFromCloud`), o app reinicia o servidor **uma vez**
  (só servidor iniciado por ele; aviso no Painel). Servidor aberto fora do app só recebe o problema
  `mods-not-mounted` ("clique em Reiniciar"). Nunca reinicia duas vezes seguidas.
- **Junto:** a rotação corrige `Scenario="Scenario=...` (colado da página do mod; o jogo recusa a entrada) ao salvar,
  ao abrir e ao gravar o MapCycle.txt; mods desinstalados pelo servidor saem da lista de Mods; erros que o jogo grava
  sempre (cosméticos, replay, voz, moldura do GameStats, recursos repetidos entre mapas) contam como aviso.

## ADR-018 · Logs por partida, limpeza configurável e uma cópia do app

- **Contexto:** a gaveta juntava todas as partidas vistas desde que o app abriu (934 "erros" numa tarde) e nada
  apagava os logs antigos. O usuário pediu separar por partida e poder apagar os antigos, com a escolha dele e,
  por padrão, sem apagar nada.
- **Decisão:**
  - **Partidas:** o servidor já grava um arquivo por partida (`Insurgency.log` e `Insurgency-backup-<data>.log`). O
    `logbuf` numera as partidas (`StartRun` a cada início, também para servidor aberto fora do app) e a gaveta mostra
    só a atual, com uma linha "Servidor iniciado às …" e a contagem de erros dela. Um seletor abre as partidas
    guardadas pelo servidor (`LogAPI.Runs`/`ReadRun`: a tela manda só o nome, conferido contra a lista; últimas 6000
    linhas, mascaradas).
  - **Limpeza** (`domain/serverlogs`): só nomes que o servidor cria; nunca o log atual; mantém N partidas (mínimo 2, a
    atual e a anterior, que o app lê para "Trazer do servidor" e a chave do mod.io; máximo 100; sugestão 10).
    Definições › Registos do servidor: limpeza automática **desligada por padrão**; ligada, roda uma vez por partida
    iniciada pelo app, depois que o servidor troca o log (fase mod.io), para a partida atual contar entre as N.
    Botão "Apagar registos antigos agora" com confirmação (quantos e quanto espaço).
  - **Ruído:** `Could not find outer … to create` (mapas de mod sem objetos internos; centenas por partida na
    Containment) conta como aviso. Na partida real de 09/11 19:56: 934 → 6 erros.
  - **Uma cópia só** (`SingleInstanceLock`): duas cópias abertas salvavam o `config.json` uma por cima da outra (o
    aceite dos termos sumiu assim) e acompanhariam o mesmo servidor em dobro. Abrir de novo traz a janela para a
    frente. A versão de desenvolvimento usa outro id.

## ADR-019 · Trocar o mapa com o servidor ligado, pelo RCON

- **Contexto:** o usuário pediu para trocar o mapa pelo app, inclusive para mapas de mods. O perfil já tinha o mapa
  inicial e o mapa de mod, mas só valiam no próximo início.
- **Decisão:**
  - **RCON próprio** (protocolo Source RCON): `domain/rcon` monta e lê pacotes (tamanho limitado a 1 MB);
    `adapters/rcon` conecta só em `127.0.0.1` (o app roda na máquina do servidor), com a senha do cofre.
  - **Uma conexão por clique**, nunca em intervalos: o RCON do Sandstorm deixa uma thread presa no servidor a cada
    conexão. Uma ação de cada vez (`busy`).
  - **Lista:** comando `scenarios`. No servidor real (10/11/2026, com o Containment carregado) ele lista só os
    mapas do jogo base (278 cenários). Os mapas de mods entram pela rotação de mapas (MapCycle), pelo mapa de mod e
    pelo mapa inicial do perfil (`server.MapChoices`, origem "config"). A tela agrupa por mapa.
  - **Troca:** `travelscenario <cenário>`, que mantém senha, vagas e mutators (`travel <url>` trocaria a URL
    inteira). A tela manda só o id; o Go aceita só um id da última lista e recusa respostas de erro do servidor.
  - Problemas novos: `rcon-off`, `rcon-failed`, `rcon-auth`, `map-change-failed`.
  - A troca não muda o perfil: ao reiniciar, o servidor volta ao mapa de Servidor › Mapas.
- **Pendente:** a troca para mapa de mod com `travelscenario` não foi testada no servidor real (precisa trocar o
  mapa com jogadores possivelmente ligados). Se o servidor recusar, a tela mostra `map-change-failed`.

## ADR-020 · Mensagem do dia e mensagens de entrada e saída

- **Contexto:** o usuário pediu para configurar o texto que aparece quando um jogador entra ou sai.
- **Decisão:**
  - **Mensagem do dia:** o servidor lê `Insurgency/Config/Server/Motd.txt` ao iniciar e mostra na tela de
    carregamento (guias de hospedagem; sem parâmetro de linha de comando). `config.Prepare` grava o arquivo quando
    o texto não está vazio (até 1000 caracteres, várias linhas); "Trazer do servidor" lê um `Motd.txt` existente.
  - **Entrada e saída:** o jogo não tem texto próprio para isso. O app lê no log (`logparse`):
    `Login request: …?Name=<nome> userId: <id>`, depois `Join succeeded: <nome>` (às vezes cortado; casa pelo
    começo do nome) e `UNetConnection::Close … GameNetDriver … PC: INSPlayerController … UniqueId: <id>` (saída,
    só com o id). Manda `say <texto>` pelo RCON, com `{name}` trocado pelo nome (limpo de controles, até 40
    caracteres; a mensagem até 150).
  - **Uma conexão de RCON aberta** para as mensagens (`app/chat`), em fila, sem travar a leitura do log; cai →
    reabre uma vez. Fecha quando o servidor para. Linhas com mais de 2 minutos não geram mensagem (o app pode ler
    o log de uma partida já em andamento).
  - Exige RCON ligado (validação `chatNeedsRcon`); só funciona com o app aberto (dito na tela).
- **Consequência:** se uma atualização do jogo mudar essas linhas do log, as mensagens param e nada mais quebra.

## ADR-021 · Rede neutra (sem Radmin) e teste de portas pela internet

- **Contexto:** o app citava o Radmin VPN como o jeito de jogar. O usuário pediu para tirar: quem sobe o servidor
  pode usar qualquer programa de VPN ou abrir as portas no roteador, e quer saber se as portas estão abertas.
- **Decisão:**
  - **Textos neutros:** dois guias em Servidor › Rede (portas no roteador, com as portas do perfil; ou programa de
    VPN), a lista de portas UDP a abrir e o aviso para não abrir a do RCON. O link do Radmin saiu.
  - **Verificador próprio** (`checker/`, Go só com a biblioteca padrão), no servidor da página de download, atrás
    do Caddy em `/sandstorm/api/` (serviço `sandstorm-check` do compose do game-in, montado direto do GitHub).
    Testa **só o IP de origem do pedido** (último valor do `X-Forwarded-For` posto pelo Caddy): consulta do Steam
    (A2S_INFO, com desafio) na porta de consulta e conexão TCP na porta do RCON. Recusa IP privado e IPv6, campos a
    mais e portas fora de 1024–65535; 6 testes por minuto por IP, 120 no total, um por vez por IP. Não grava nem
    registra nada.
  - A porta do jogo não dá para testar de fora sem o protocolo do Unreal: a tela diz para abri-la junto com a de
    consulta. RCON aberto para a internet aparece como risco.
  - O resultado traz o IP público, com botão para usá-lo em "IP que os jogadores usam".
  - Aviso de dados atualizado (o verificador recebe o IP e as portas quando o usuário clica): `legal.Version` =
    2026-11-10, o app pede o aceite de novo.
- **Alternativas:** API do Steam (`GetServersAtAddress`) só diz que o Steam conhece o servidor, não que os
  jogadores alcançam a porta; sites de terceiros testam TCP, e o jogo usa UDP.

## ADR-022 · Mais configurações do servidor (regras, jogadores, automação, Windows, perfis)

- **Contexto:** o usuário pediu as 14 configurações sugeridas depois da 0.1.3. Fonte das opções: o guia oficial de
  administração do servidor (Server Admin Guide da New World Interactive) e o `help` do RCON do servidor real.
- **Decisão:**
  - **Regras do jogo** (`domain/server/rules.go`): catálogo de 60 opções do Game.ini (coop, checkpoint,
    objetivos, fogo amigo, team kills, chat e voz, partida, suprimentos, votação para expulsar) e do Engine.ini
    (tick rate e banda). `Profile.Rules` nil = o app não mexe nos arquivos (servidor configurado à mão continua
    igual); ligado, o app grava as opções com valor e apaga as vazias ao iniciar (`config.Prepare`, cópia .bak).
    "Trazer do servidor" lê as opções do catálogo. A tela monta os campos a partir do catálogo.
  - **Mutators:** lista oficial (25) com descrição traduzida, mais mutators de mods pelo nome.
  - **Dia ou noite:** `?Lighting=` na URL de início, `Lighting="Day|Night"` nas linhas da rotação e na troca de
    mapa (com iluminação vai `travel <mapa>?Scenario=…?Lighting=…?MaxPlayers=…[?Password=…]`, porque o
    travelscenario não aceita a iluminação; sem iluminação continua o travelscenario).
  - **Jogadores** (`app/players`, `domain/admin`): `listplayers` e `listbans` numa conexão por clique;
    expulsar (`kick`), banir (`banid` com minutos ou `permban`), tirar banimento (`unban`) e mensagem a todos
    (`say`). Só net IDs da última lista; motivo sem aspas nem controle, até 100 caracteres.
  - **Mensagens programadas:** até 10, uma por vez a cada N minutos (2–240), pela mesma conexão do chat (ADR-020).
  - **Reinício automático** (`app/watchdog`): diário num horário, com aviso no chat 5 e 1 minuto antes (texto do
    usuário com {minutes}) e atualização pelo SteamCMD; e ao fechar sozinho (espera 10 s, para depois de 3 vezes em
    30 minutos). Só com o app aberto.
  - **Discord** (`app/notify`, `adapters/discord`): webhook no cofre (só endereços oficiais do Discord, hosts no
    netx); textos do usuário para no ar, parado, fechou sozinho, entrou e saiu ({server}, {map}, {name}); nomes
    escapados e `allowed_mentions` vazio (ninguém é chamado). Fila de 32, um envio por segundo.
  - **Windows** (`adapters/winsys`): abrir com o Windows (registro `HKCU\…\Run`, abre minimizado com
    `--autostart`; opção de iniciar o servidor junto); regra do Firewall: o app abre uma cópia de si como
    administrador (UAC) com `--firewall-rule`, que lê a pasta do servidor da configuração e cria as regras de
    entrada para `InsurgencyServer.exe` e `InsurgencyServer-Win64-Shipping.exe` (netsh, nomes fixos); lista as
    placas de rede para o `-MultiHome`.
  - **Perfis** (`app/config/profiles.go`): o perfil em uso continua em `file.Profile`; os outros ficam guardados.
    Trocar só com o servidor parado. Segredos do cofre compartilhados. Um servidor por vez.
  - **Lista de mods** (`-CmdModList`): opcional e desligada; com o mod.io atual quem manda são as inscrições da
    conta do servidor, e a tela avisa que versões novas do jogo podem ignorar a lista.
  - Aviso de dados: linha do Discord. A versão 2026-11-10 ainda não foi publicada, então não muda.
- **Não confirmado ao vivo:** o efeito das opções no jogo, o travel com iluminação, `kick`/`banid` com jogadores, o
  `-CmdModList` com o mod.io atual, a regra do Firewall (pede administrador) e o webhook real.

## ADR-023 · Mutators dos mods instalados

- **Contexto:** a lista de mutators só mostrava os do jogo; os de mods eram digitados pelo nome. O usuário pediu
  que o app achasse os mutators dos mods sozinho.
- **Decisão:** um mutator de mod é um arquivo `<Plugin>/Content/Mutators/<Nome>.uasset` dentro do `.pak` do mod,
  e `<Nome>` é o que vai em `-mutators=` (conferido no mod NoRestrictedArea). `domain/pak` lê o rodapé do
  `.pak` (versão 10 e 11), o índice principal e o índice de pastas, com todo tamanho e posição conferidos (arquivo
  de terceiros; teste com fuzz). `adapters/modio` lê só esses pedaços de cada `.pak` dos mods instalados
  (`state.json`, pasta dentro de `mod.io/254/mods`). Se nenhum `.pak` do mod puder ser lido e o mod tem a tag
  "Mutator" no mod.io, usa o `PluginName` do `metadata_blob`.
- **Tela:** uma lista com rolagem: "Dos seus mods" (marca MOD e o nome do mod) e "Do jogo". Saiu o campo de texto;
  mutator do perfil que nenhum mod instalado traz continua marcado, com aviso.
- **Verificado:** nos 7 mods desta máquina (15 `.pak`, cerca de 11 GB) a leitura leva 0,01 s e acha só o
  NoRestrictedArea; nenhum mod de mapa deu falso positivo.

## ADR-024 · Formatação das mensagens e mapa de mod escolhido da lista

- **Contexto:** o usuário pediu formatação visual na mensagem do dia e nas mensagens programadas, com um editor, e
  que o mapa de mod depois de carregar fosse escolhido numa lista em vez de digitado.
- **O que o jogo aceita** (lido dos arquivos do servidor 1.22, índice dos .pak e o executável): a tela de carregamento
  usa texto rico do Unreal com a tabela `MotdTextStyleSet` (linhas Default, Bold, Header1, Header2) e o decorador
  `WebImageDecorator` com a tag `extimg` e o atributo `src`. O chat tem `ChatTextStyleSet` (Default, Bold); se o
  `say` passa pelo texto rico não foi confirmado com o jogo aberto. Não há tag de cor livre.
- **Decisão:**
  - `domain/server/markup.go`: `<Bold>`, `<Header1>`, `<Header2>` (fechar com `</>`) e
    `<extimg src="https://..."/>` na mensagem do dia; só `<Bold>` nas mensagens de chat. Tag desconhecida, sem
    fechar, aninhada ou imagem sem https vira erro de validação. `ConfigAPI.MarkupPreview` devolve os pedaços para a
    prévia (o app não baixa a imagem; quem baixa é o jogo de cada jogador).
  - Editor na tela: botões que envolvem a seleção com a tag, campo para o endereço da imagem e prévia.
  - Mapa de mod: `pak.Scan` também lista mapas (.umap) e cenários (Scenario_*.uasset); o mapa de um cenário é o
    .umap de nome mais curto tal que o cenário é `Scenario_<mapa>_<modo>` (os mods têm submapas com o nome do
    modo). `ConfigAPI.ModScenarios` alimenta os dois seletores (mapa e cenário). Valor do perfil que não está nos mods
    instalados continua na lista, marcado.
- **Verificado:** nos 6 mods de mapa desta máquina, 37 cenários, todos ligados ao mapa principal; editor e seletores no
  modo dev. Não testado com o jogo: a aparência real dos estilos e da imagem, e negrito no chat.
