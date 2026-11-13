import type { Messages } from './en-US';

export const ptPT: Messages = {
  app: {
    name: 'Sandstorm Server Assist',
    beta: 'BETA',
    version: 'Versão {version}',
    loading: 'A carregar…',
  },
  titlebar: {
    help: 'Ajuda',
    minimize: 'Minimizar',
    maximize: 'Maximizar',
    restore: 'Restaurar',
    close: 'Fechar',
  },
  nav: {
    label: 'Secções',
    dashboard: 'Painel',
    install: 'Instalação',
    server: 'Servidor',
    mods: 'Mods',
    backup: 'Cópia de segurança',
    settings: 'Definições',
    feedback: 'Feedback',
    about: 'Sobre',
  },
  common: {
    save: 'Guardar',
    saved: 'Guardado',
    cancel: 'Cancelar',
    close: 'Fechar',
    choose: 'Escolher…',
    notChosen: 'Não escolhida',
    copy: 'Copiar',
    copied: 'Copiado!',
    show: 'Mostrar',
    hide: 'Ocultar',
    change: 'Alterar',
    set: 'Definir',
    remove: 'Remover',
    add: 'Adicionar',
    open: 'Abrir',
    refresh: 'Atualizar',
    optional: 'opcional',
    yes: 'Sim',
    no: 'Não',
    go: 'Ir',
    howTo: 'Como obter isto?',
    unsaved: 'Alterações por guardar',
    appliesNextStart: 'As alterações aplicam-se no próximo arranque do servidor.',
  },
  units: {
    b: 'B',
    kb: 'KB',
    mb: 'MB',
    gb: 'GB',
    perSecond: '{value}/s',
    seconds: '{n} s',
    minutes: '{n} min',
    hours: '{h} h {m} min',
    left: 'faltam {time}',
  },
  phase: {
    'not-installed': 'Não instalado',
    stopped: 'Parado',
    starting: 'A arrancar',
    modio: 'A ligar ao mod.io',
    mods: 'A carregar os mods',
    traveling: 'A carregar o mapa',
    ready: 'Pronto para jogar',
    stopping: 'A parar',
    crashed: 'Parou inesperadamente',
  },
  phaseHint: {
    'not-installed': 'Instale primeiro o servidor em Instalação.',
    stopped: 'O servidor não está em execução.',
    starting: 'O servidor está a abrir. Demora alguns segundos.',
    modio: 'A iniciar sessão no mod.io com a conta do servidor.',
    mods: 'A sincronizar as subscrições e a montar os mods. Os mods novos são transferidos agora.',
    traveling: 'A mudar para o mapa do mod. Mapas grandes demoram até um minuto.',
    ready: 'Os jogadores já podem entrar.',
    stopping: 'A fechar o servidor.',
    crashed: 'O servidor fechou sozinho. Veja o erro abaixo e os registos.',
  },
  dashboard: {
    title: 'Painel',
    subtitle: 'Arranque o servidor e acompanhe o que está a fazer.',
    start: 'Iniciar servidor',
    stop: 'Parar servidor',
    restart: 'Reiniciar',
    timeline: 'Arranque',
    step: {
      start: 'Processo iniciado',
      modio: 'mod.io ligado',
      mods: 'Mods montados',
      map: 'Mapa carregado',
      ready: 'Pronto',
    },
    map: 'Mapa',
    scenario: 'Cenário',
    port: 'Porta',
    since: 'Em execução desde',
    modsActive: '{n} mods ativos',
    external:
      'Este servidor já estava a correr quando a aplicação abriu, e não foi ela que o iniciou. A aplicação acompanha-o e controla-o normalmente; só não o reinicia sozinha para carregar mods novos.',
    mutatorsRefused:
      'O servidor recusou estes mutators neste mapa: {names}. O jogo não diz porquê: um mutator pode não funcionar neste modo de jogo ou junto com outro que escolheu (visto com o Ultralethal no Checkpoint, com ou sem o Bullet Sponge). Tire-o em Servidor › Mapas.',
    authorizing: 'A autorizar no mod.io com o código. O servidor reinicia sozinho em modo normal.',
    remount:
      'Foram baixados mods novos e o servidor não carregou os que já estavam instalados (mapa do mod e mutators). A reiniciar uma vez para carregar todos.',
    gameStatsMissing:
      'Sem token GameStats: o servidor funciona normalmente, mas os jogadores não ganham estatísticas/XP. Pode adicionar um em Servidor › Palavras-passe e tokens.',
    cmdLine: 'Linha de comandos (segredos ocultos)',
    mapChange: {
      title: 'Trocar mapa',
      mapLabel: 'Mapa',
      modeLabel: 'Modo',
      text: 'Troque o mapa com o servidor ligado. Os jogadores continuam ligados e vão para o novo mapa. Os mapas de mods vêm dos mods carregados no servidor.',
      needsRcon:
        'Para trocar o mapa com o servidor ligado, ligue o RCON (com palavra-passe) em Servidor › Palavras-passe e tokens e reinicie o servidor.',
      openServer: 'Abrir definições do servidor',
      load: 'Ver mapas disponíveis',
      reload: 'Pedir a lista ao servidor de novo',
      empty: 'O servidor não listou nenhum mapa. Tente de novo daqui a alguns segundos.',
      fromConfig: 'mapa de mod',
      fromConfigHint:
        'Os mapas de mods vêm dos mods carregados nesta partida. Um mod em que acabou de se inscrever aparece depois de o servidor reiniciar.',
      notLoaded: 'mod não carregado',
      notLoadedHint:
        'Os mapas marcados com "mod não carregado" estão na rotação de mapas, mas o mod deles não está carregado. Remova-os da rotação em Servidor › Mapas ou, se acabou de se inscrever no mod, reinicie o servidor.',
      current: 'agora',
      change: 'Trocar mapa',
      confirm: 'Trocar para {map}? A partida atual termina e todos vão para o novo mapa.',
      confirmYes: 'Trocar agora',
      sent: 'A trocar para {map}… O estado acima mostra quando estiver pronto.',
      waitReady: 'Disponível quando o servidor estiver pronto.',
      restartNote: 'Quando o servidor reinicia, volta ao mapa definido em Servidor › Mapas.',
    },
    connect: {
      title: 'Como os jogadores entram',
      text: 'Os jogadores entram pela internet (com as portas do router abertas) ou por um programa de VPN que você e os jogadores partilham.',
      console: 'No jogo, abra a consola (tecla ~) e escreva:',
      noIp: 'Defina o IP que os jogadores usam em Servidor › Rede.',
    },
    checklist: {
      title: 'Primeiros passos',
      allDone: 'Tudo pronto! Já pode iniciar o servidor.',
      serverDir: 'Escolher a pasta do servidor',
      steamcmd: 'Obter o SteamCMD',
      serverInstalled: 'Instalar o servidor dedicado',
      hostname: 'Dar um nome ao servidor',
      modioAuth: 'Autorizar o servidor no mod.io',
      admins: 'Adicionar administradores',
      connectIp: 'Definir o IP que os jogadores usam',
      secretsSafe: 'Tirar as palavras-passe dos ficheiros de configuração',
      gslt: 'Adicionar um token de servidor Steam (GSLT)',
    },
  },
  install: {
    title: 'Instalação',
    subtitle: 'Transfira o SteamCMD e instale ou atualize o servidor dedicado.',
    folders: 'Pastas',
    serverDir: 'Pasta do servidor',
    serverDirHint:
      'Precisa de cerca de 25 GB livres, mais os mods (mods de mapa podem ter 5 GB cada). Evite pastas do sistema.',
    steamcmdDir: 'Pasta do SteamCMD',
    steamcmdDirHint:
      'Se ficar vazia, a aplicação usa uma pasta "steamcmd" ao lado da pasta do servidor.',
    foldersNested:
      'A pasta do SteamCMD está dentro da pasta do servidor (ou o contrário). O SteamCMD não instala aí: poria o servidor na pasta dele. Escolha uma pasta do SteamCMD ao lado da pasta do servidor, não dentro dela.',
    steamcmd: {
      title: 'SteamCMD',
      text: 'Ferramenta de linha de comandos da Valve que transfere o servidor. A aplicação confirma a assinatura digital da Valve antes de o executar.',
      fetch: 'Transferir SteamCMD',
      check: 'Verificar SteamCMD',
      selfUpdates:
        'O SteamCMD atualiza-se sozinho sempre que a aplicação o usa. Este botão só confere o ficheiro e transfere-o de novo se estiver danificado.',
      ready: 'SteamCMD pronto',
      missing: 'SteamCMD não encontrado',
    },
    server: {
      title: 'Servidor dedicado',
      installed: 'Instalado',
      notInstalled: 'Não instalado',
      build: 'Build {id}',
      size: '{size} em disco',
      updated: 'Atualizado a {date}',
      install: 'Instalar servidor',
      update: 'Procurar atualizações',
      validate: 'Verificar ficheiros',
      validateHint:
        'Verifica todos os ficheiros e repara o que estiver danificado. É mais lento; use-o se o servidor deixar de abrir.',
      runningWarn: 'Pare o servidor antes de instalar ou atualizar.',
    },
    progress: {
      title: 'Progresso',
      'self-update': 'A atualizar o SteamCMD',
      login: 'A ligar à Steam',
      preallocating: 'A reservar espaço em disco',
      downloading: 'A transferir',
      verifying: 'A verificar ficheiros',
      committing: 'A gravar ficheiros',
      reconfiguring: 'A configurar',
      done: 'Concluído',
      failed: 'Falhou',
      waiting: 'A começar…',
      cancel: 'Cancelar',
      finished: {
        steamcmd: 'SteamCMD transferido e verificado.',
        steamcmdPresent:
          'O SteamCMD já está aqui e a assinatura da Valve é válida. Nada para transferir.',
        server: 'O servidor está instalado e atualizado.',
        upToDate: 'O servidor já tem a versão mais recente. Nada para transferir.',
        validated:
          'Todos os ficheiros do servidor foram verificados; o que estava danificado foi reparado.',
      },
    },
    import: {
      title: 'Já tem um servidor nesta pasta?',
      text: 'Traga os dados do servidor já configurado aqui, a correr ou parado: nome, portas, mapas, mods, administradores, rotação de mapas e palavras-passe. A app lê o último registo do servidor e os ficheiros de configuração.',
      button: 'Trazer dados do servidor',
    },
  },
  server: {
    travelNoMods:
      'Nenhum mapa encontrado nos mods instalados. Subscreva um mod de mapa e inicie o servidor uma vez.',
    travelNotInstalled: '{name} (não está nos mods instalados)',
    mutatorsHint: 'Os mutators dos mods instalados são encontrados pelo app e marcados como MOD.',
    mutatorsMods: 'Dos seus mods',
    mutatorsGame: 'Do jogo',
    mutatorsBadge: 'MOD',
    mutatorsRefused:
      'No mapa atual o servidor recusou: {names}. O jogo não diz porquê: pode não funcionar neste modo de jogo ou com outro mutator que escolheu (visto com o Ultralethal no Checkpoint, com ou sem o Bullet Sponge).',
    mutatorsRefusedBadge: 'recusado pelo servidor',
    mutatorsFromMod: 'Mod: {name}',
    mutatorsMissing: 'Nenhum mod instalado tem este mutator: o servidor não o vai carregar.',
    mutatorsNoMods:
      'Nenhum mutator nos mods instalados. Subscreva um mod de mutator no mod.io e inicie o servidor uma vez para o baixar.',
    modList: 'Enviar a minha lista de mods ao servidor (-CmdModList)',
    modListHint:
      'Avançado. Normalmente o servidor carrega os mods subscritos pela conta dele no mod.io. Ligue só se quiser definir a lista aqui; versões recentes do jogo podem ignorá-la.',
    modIds: 'Ids dos mods',
    modIdsHint: 'Os números dos mods no mod.io, separados por vírgulas.',
    title: 'Servidor',
    subtitle:
      'Nome, mapas, administradores e palavras-passe. Usados no próximo arranque do servidor.',
    tabs: {
      rules: 'Regras do jogo',
      automation: 'Automação',
      general: 'Geral',
      maps: 'Mapas',
      admins: 'Administradores',
      secrets: 'Palavras-passe e tokens',
      network: 'Rede',
      messages: 'Mensagens',
    },
    reset: 'Repor predefinições',
    hostname: 'Nome do servidor',
    hostnameHint: 'Aparece na lista de servidores. Até 64 caracteres, sem aspas.',
    maxPlayers: 'Máximo de jogadores',
    port: 'Porta do jogo',
    queryPort: 'Porta de consulta',
    portsHint:
      'Os jogadores ligam-se à porta do jogo. A porta de consulta responde à lista de servidores.',
    useMods: 'Usar mods do mod.io',
    useModsHint: 'O servidor transfere os mods subscritos pela sua conta do mod.io.',
    startMap: 'Mapa inicial',
    startMapHint: 'Tem de ser um mapa do jogo base. Os mapas de mods vão em "Mapa do mod", abaixo.',
    mode: 'Cenário',
    customScenario: 'Nome do cenário',
    travel: 'Mudar para um mapa de mod depois de carregar os mods',
    travelHint:
      'Usa -ModDownloadTravelTo. Exemplo: mapa Containment, cenário Scenario_Containment_Checkpoint_Security.',
    travelMap: 'Nome do mapa do mod',
    travelScenario: 'Cenário do mapa do mod',
    mutators: 'Mutators',
    mapCycle: 'Rotação de mapas (MapCycle.txt)',
    mapCycleHint:
      'Uma entrada por linha. Deixe vazio para manter o ficheiro que o servidor já tem.',
    admins: {
      title: 'Administradores',
      text: 'Os administradores podem expulsar, banir e mudar de mapa no jogo. São identificados pelo SteamID64.',
      steamId: 'SteamID64 ou ligação do perfil Steam',
      name: 'Nome (opcional)',
      empty: 'Ainda não há administradores.',
      found: 'SteamID64 encontrado: {id}',
    },
    secrets: {
      vaultNote:
        'Guardado e cifrado pelo Windows (DPAPI). Só o seu utilizador do Windows, neste PC, o consegue abrir. Os valores nunca vão para ficheiros de configuração.',
      set: 'Guardado no cofre',
      notSet: 'Não definido',
      serverPassword: 'Palavra-passe de entrada',
      serverPasswordHint: 'Deixe vazio para um servidor onde qualquer pessoa da rede pode entrar.',
      rcon: 'Ativar RCON (administração remota)',
      rconPassword: 'Palavra-passe RCON',
      rconPort: 'Porta RCON',
      gsltToken: 'Token de servidor de jogo Steam (GSLT)',
      gsltTokenHint:
        'Opcional numa rede privada (VPN). Necessário para aparecer na lista pública de servidores.',
      gameStatsToken: 'Token GameStats',
      gameStatsTokenHint:
        'Opcional. Permite que os jogadores ganhem estatísticas e XP no seu servidor.',
      rconInIni: 'A palavra-passe RCON está guardada em texto simples no Game.ini.',
      moveToVault: 'Mover para o cofre',
      moved: 'Feito: a palavra-passe está agora só no cofre.',
    },
    network: {
      bindIp: 'Placa de rede do servidor',
      bindIpHint:
        'Deixe em "Todas" na maioria dos casos. Escolha uma placa só se o servidor tiver de responder numa única rede, como a rede de um programa de VPN (-MultiHome).',
      bindIpMissing:
        'Este endereço não está neste computador agora. Se for de um programa de VPN, ligue a VPN antes de iniciar o servidor.',
      bindIpAll: 'Todas as placas de rede',
      connectIp: 'IP que os jogadores usam',
      connectIpHint:
        'Pela internet, é o seu IP público. Por um programa de VPN, é o IP que o programa mostra para este PC.',
      firewall:
        'Se os amigos não conseguirem entrar, permita o InsurgencyServer.exe na Firewall do Windows, no tipo de rede que usa (Privada ou Pública).',
      ports:
        'Portas a abrir no router (UDP): {port} (jogo) e {queryPort} (consulta). Não abra a porta do RCON ({rconPort}).',
    },
    messages: {
      announcements: 'Mensagens programadas no chat',
      announcementsHint:
        'Uma mensagem por linha, até 10. O app envia uma de cada vez, por ordem, enquanto o servidor está pronto.',
      announcementsExample: 'Entre no nosso Discord: discord.gg/exemplo',
      every: 'A cada (minutos)',
      motd: 'Mensagem do dia',
      motdHint:
        'Aparece aos jogadores no ecrã de carregamento quando entram. Até 1000 caracteres; pode ter várias linhas. Vale no próximo arranque.',
      join: 'Mensagem quando um jogador entra',
      leave: 'Mensagem quando um jogador sai',
      chatHint:
        'Vai para o chat do jogo. {name} vira o nome do jogador. Até 150 caracteres. Vazio = desligado.',
      joinExample: 'Bem-vindo, {name}! Bom jogo.',
      leaveExample: '{name} saiu do servidor.',
      howChat:
        'O jogo não tem definição para estas mensagens: o app lê no registo do servidor quem entrou ou saiu e envia a mensagem pelo RCON, enquanto o app estiver aberto.',
      needsRcon:
        'As mensagens de entrada e saída são enviadas pelo RCON. Ligue o RCON (com palavra-passe) em Palavras-passe e tokens para usá-las.',
    },
    portCheck: {
      button: 'Diagnosticar portas',
      running: 'A verificar… (até 15 s)',
      hint: 'Verifica a partir deste PC: o servidor, a Firewall do Windows, o router (UPnP), o seu IP público e a Steam. Só um teste de fora da sua rede prova que as portas estão abertas: veja os passos abaixo.',
      gamePort: 'Porta de jogo UDP {port}: reencaminhe-a no router junto com a porta de consulta.',
      publicIp: 'O seu IP público: {ip}',
      useIp: 'Usar como "IP que os jogadores usam"',
      verdict: {
        open: 'As portas parecem abertas: o servidor respondeu pelo seu IP público.',
        likely:
          'Tudo o que este PC consegue verificar parece certo. Confirme com um teste de fora da sua rede.',
        problem: 'Foi encontrado um problema. Corrija os itens a vermelho e verifique de novo.',
        unsure:
          'Não foi possível confirmar a partir deste PC. Faça o teste de fora da sua rede (abaixo).',
      },
      items: {
        serverOff:
          'O servidor não está ligado. Inicie-o para verificar a porta de consulta, o caminho pelo IP público e a Steam.',
        localOk: 'O servidor responde na porta de consulta UDP {port} neste PC.',
        localNo:
          'O servidor não responde na porta de consulta UDP {port}, nem neste PC. Espere que termine de carregar ou confira a porta de consulta.',
        firewallOk: 'A Firewall do Windows tem as regras do servidor.',
        firewallMissing:
          'As regras da Firewall do Windows criadas pela app não foram encontradas. Crie-as em Servidor › Rede (ignore se criou a sua própria regra).',
        routerUnknown:
          'O router não respondeu por UPnP, por isso a app não consegue ler os reencaminhamentos. Confira-os na página do router.',
        routerMapped: 'Router: a porta UDP {port} está reencaminhada para este PC.',
        routerNoMap:
          'Router: não há reencaminhamento UPnP para a porta UDP {port}. Reencaminhamentos feitos à mão na página do router muitas vezes não aparecem aqui; confirme que existem.',
        routerOtherPc:
          'Router: a porta UDP {port} está reencaminhada para outro aparelho ({client}), não para este PC.',
        routerWrongPort:
          'Router: a porta UDP {port} está reencaminhada para a porta interna {internal}. Os dois números têm de ser iguais.',
        routerDisabled:
          'Router: o reencaminhamento da porta UDP {port} existe, mas está desligado.',
        rconForwarded:
          'Router: a porta RCON TCP {port} está aberta para a internet. Qualquer pessoa pode tentar a sua senha RCON. Remova este reencaminhamento.',
        cgnat:
          'O endereço do seu router na internet ({ip}) não é um IP público: o seu fornecedor usa CGNAT. O reencaminhamento de portas não vai funcionar; peça um IP público ao fornecedor.',
        ipMismatch:
          'O router diz que o IP dele na internet é {router}, mas a internet vê {public}. O fornecedor pode usar CGNAT, ou há um programa de VPN ativo.',
        publicOk: 'O router tem o seu IP público ({ip}), por isso não há CGNAT.',
        publicIp: 'IP público visto pela internet: {ip}.',
        publicUnknown:
          'Não foi possível descobrir o seu IP público (sem internet, ou STUN bloqueado).',
        loopbackOk:
          'O servidor respondeu pelo seu IP público {ip}:{port}: o router reencaminha a porta de consulta.',
        loopbackNo:
          'Sem resposta pelo seu IP público a partir deste PC. Muitos routers não permitem isto, por isso não quer dizer que a porta esteja fechada.',
        steamListed: 'A Steam lista o seu servidor no seu IP público.',
        steamNotListed:
          'A Steam ainda não lista o seu servidor no seu IP público. Logo depois de iniciar pode levar alguns minutos; se continuar assim, confira o token da Steam (GSLT).',
      },
      manual: {
        title: 'Teste de fora da sua rede (o único teste seguro)',
        intro:
          'Um PC não consegue testar as próprias portas pela internet. Peça a um amigo noutra rede que faça isto com o seu servidor ligado:',
        yourIp: 'o seu IP público',
        steps: [
          'Abrir a Steam e ir a Ver › Servidores de jogo.',
          'Abrir o separador Favoritos e clicar em "Adicionar servidor" (o botão +).',
          'Escrever {ip}:{port} e clicar em "Procurar jogos neste endereço".',
          'Se o servidor aparecer com o nome e o mapa, as portas estão abertas. O amigo pode então entrar pela lista.',
        ],
        tipsTitle: 'Dicas:',
        tips: [
          'Não teste a partir de um PC da sua própria rede: a maioria dos routers bloqueia o caminho de volta para o próprio IP público.',
          'Os dados móveis (telemóvel como hotspot) contam como outra rede.',
          'Sites que "testam portas" normalmente testam só TCP. O jogo usa UDP, por isso podem dizer "fechada" mesmo com a porta aberta.',
          'Os jogadores precisam da porta de jogo e da porta de consulta reencaminhadas (as duas UDP). Nunca reencaminhe a porta RCON.',
        ],
      },
    },
  },
  sync: {
    button: 'Trazer do servidor',
    buttonHint: 'Preencher os campos com os dados do servidor já configurado na pasta escolhida',
    needDir: 'Escolha primeiro a pasta do servidor em Instalação',
    title: 'Dados do servidor existente',
    source: 'Lido de {file}',
    sourceFiles: 'Lido dos ficheiros de configuração',
    running: 'servidor a correr agora',
    stopped: 'servidor parado',
    conflictsTitle: 'Alterou {n} campo(s) e o servidor tem outro valor',
    conflictsText: 'Escolha que valor fica: o que escreveu na app ou o que o servidor está a usar.',
    keepAllMine: 'Manter todos os meus',
    useAllServer: 'Usar todos os do servidor',
    mine: 'Meu (app)',
    server: 'Servidor',
    fillTitle: '{n} campo(s) vão ser preenchidos com os dados do servidor',
    fillText: 'Estes campos estavam vazios ou ainda tinham o valor padrão na app.',
    sameCount: '{n} campo(s) já são iguais na app e no servidor.',
    chooseEach: 'Escolha um valor para cada campo acima.',
    apply: 'Aplicar',
    secretMine: 'Guardada na app',
    secretServer: 'Diferente no servidor',
    secretFound: 'Encontrada no servidor (vai para o cofre)',
    empty: 'vazio',
    items: '{n} itens',
    upToDate: 'A app já tem os mesmos dados do servidor. Nada a alterar.',
    notFound:
      'Não foram encontrados dados de um servidor configurado nesta pasta. O servidor já foi iniciado pelo menos uma vez?',
    done: 'Pronto: {applied} campo(s) atualizados com os dados do servidor e {kept} mantidos como escreveu.',
    skipped: 'Não usados porque o valor no servidor não é válido aqui: {fields}.',
    rconInIni:
      'A palavra-passe RCON ainda está escrita no Game.ini. Vá a Servidor › Palavras-passe e tokens e use "Mover para o cofre".',
    fields: {
      startLighting: 'Iluminação do mapa inicial',
      bindIp: 'Placa de rede do servidor',
      modIds: 'Lista de mods',
      rules: 'Regras do jogo',
      motd: 'Mensagem do dia',
      hostname: 'Nome do servidor',
      port: 'Porta do jogo',
      queryPort: 'Porta de consulta',
      maxPlayers: 'Máximo de jogadores',
      startMap: 'Mapa inicial',
      travelMap: 'Mapa de mod depois de carregar os mods',
      useMods: 'Usar mods',
      mutators: 'Mutators',
      rcon: 'RCON',
      rconPort: 'Porta RCON',
      admins: 'Administradores',
      mapCycle: 'Rotação de mapas',
      serverPassword: 'Palavra-passe do servidor',
      rconPassword: 'Palavra-passe RCON',
      gsltToken: 'Token GSLT',
      gameStatsToken: 'Token GameStats',
      modioApiKey: 'Chave de API do mod.io',
    },
  },
  maps: {
    mode: {
      Checkpoint_Security: 'Checkpoint (Segurança)',
      Checkpoint_Insurgents: 'Checkpoint (Insurgentes)',
      custom: 'Outro (escreva o nome)',
    },
  },
  mods: {
    title: 'Mods (mod.io)',
    subtitle: 'O servidor transfere os mods subscritos pela sua conta do mod.io.',
    session: {
      title: 'Autorização do servidor',
      valid: 'Autorizado como {user} até {date}.',
      expired: 'A autorização expirou a {date}. Autorize de novo.',
      none: 'Este PC ainda não está autorizado no mod.io.',
    },
    auth: {
      title: 'Autorizar o servidor',
      text: 'Use uma conta do mod.io só para o servidor (não a sua pessoal). Faz-se uma vez e dura cerca de um ano.',
      email: 'E-mail da conta do servidor',
      apiKey: 'Chave da API do mod.io',
      detect: 'Procurar nos registos do servidor',
      detected: 'Chave encontrada nos registos do servidor e guardada no cofre.',
      notFound: 'Ainda não há chave nos registos. Cole uma criada no mod.io (ligação no guia).',
      requestCode: 'Enviar código por e-mail',
      codeSent: 'Código enviado. Veja a caixa de entrada e a pasta de spam.',
      code: 'Código de 5 caracteres',
      authorize: 'Autorizar',
      authorizeHint:
        'O servidor arranca uma vez com o código e depois reinicia sozinho em modo normal (SecurityCode=none).',
    },
    health: {
      title: 'Mods com problemas',
      text: 'A app compara os ficheiros dos mods instalados. Dois mods que trazem o mesmo ficheiro do jogo em versões diferentes podem fazer o jogo ou o servidor fechar com "Serial size mismatch" ao carregar um mapa.',
      recheck: 'Verificar de novo',
      none: 'Nenhum conflito entre os seus {n} mods instalados.',
      found: '{n} par(es) de mods trazem os mesmos ficheiros do jogo em versões diferentes:',
      and: 'e',
      pair: '{n} ficheiro(s) em versões diferentes',
      example: 'ex.: {file}',
      advice:
        'Se o jogo ou o servidor fechar ao carregar um mapa, mantenha só um mod do par: cancele a subscrição do outro no mod.io (conta do servidor). Cópias idênticas de um ficheiro não aparecem: não causam problemas.',
      unreadable:
        'Não foi possível ler os ficheiros de: {names}. Conflitos com estes mods podem não aparecer aqui.',
    },
    crash: {
      title: 'O jogo de um jogador (ou o servidor) fechou? Encontre o mod',
      text: 'Cole a mensagem do erro (por exemplo a que tem "Serial size mismatch"). A app encontra o mod instalado que traz o ficheiro indicado. O texto não é guardado nem enviado.',
      placeholder: 'StaticMesh /Game/…: Serial size mismatch: Got …, Expected …',
      find: 'Encontrar o mod',
      noObject:
        'Não foi encontrado nenhum ficheiro do jogo nesta mensagem. Cole a linha que indica o ficheiro (começa por /Game/ ou pelo nome do mod).',
      inMods: 'está em {names}',
      inNoMod: 'nenhum mod instalado tem este ficheiro',
      adviceMismatch:
        '"Serial size mismatch" quer dizer que este mod tem uma cópia do ficheiro feita para outra versão do jogo. Cancele a subscrição do mod no mod.io (conta do servidor) até o autor o atualizar. Se acha que o download ficou estragado, use o botão de apagar na linha do mod para o transferir de novo.',
      advice:
        'Experimente tirar este mod: cancele a subscrição no mod.io (conta do servidor), ou apague-o na linha dele para o transferir de novo.',
      adviceNone:
        'O ficheiro não está nos mods do servidor. Pode vir do próprio jogo ou de um mod que o jogador tem e o servidor não: o jogador pode verificar os ficheiros do jogo na Steam (Propriedades › Ficheiros instalados) e apagar a pasta do mod em C:\\Users\\Public\\mod.io\\254\\mods para o transferir de novo.',
    },
    list: {
      title: 'Mods',
      refresh: 'Ler subscrições do mod.io',
      refreshed: 'Atualizado a {date}',
      empty: 'Ainda não há mods. Subscreva mods no mod.io com a conta do servidor.',
      name: 'Mod',
      size: 'Tamanho',
      status: 'Estado',
      openSite: 'Abrir mod.io',
      conflict:
        'Mods de mapa que trazem os mesmos ficheiros do jogo entram em conflito (exemplo: Containment com Old World ou Containment Storm). Se o servidor fechar com "Serial size mismatch", cancele a subscrição de um deles.',
    },
    pending: {
      text: 'Ainda não transferidos: {names}. O servidor transfere os mods inscritos quando arranca.',
      stoppedHint:
        '"Transferir agora" arranca o servidor; depois da transferência ele reinicia uma vez sozinho para carregar todos os mods.',
      runningHint:
        'Um servidor ligado não carrega mods novos: "Reiniciar e transferir agora" reinicia-o.',
      start: 'Transferir agora',
      restart: 'Reiniciar e transferir agora',
      confirm: 'A partida atual termina e os jogadores saem por um momento. Reiniciar agora?',
      confirmYes: 'Reiniciar agora',
      started:
        'O servidor está a arrancar e transfere os mods novos. Acompanhe o progresso na lista; ele reinicia uma vez sozinho para carregar todos os mods.',
    },
    delete: {
      column: 'Ações',
      button: 'Apagar os ficheiros transferidos de {name}',
      confirm:
        'Apagar os ficheiros transferidos de {name}? O servidor transfere o mod de novo quando arrancar, se continuar inscrito. Use para corrigir um mod danificado ou forçar uma transferência nova.',
      redownload: 'Apagar e transferir de novo',
      only: 'Só apagar',
      working: 'A apagar…',
      runningHint:
        '"Apagar e transferir de novo" para o servidor, apaga o mod e arranca o servidor, que o transfere de novo. Para só apagar, pare o servidor primeiro.',
      stoppedHint:
        '"Apagar e transferir de novo" arranca o servidor, que transfere o mod. "Só apagar" deixa a transferência para o próximo arranque.',
      doneRedownload:
        '{name} apagado. O servidor está a arrancar e transfere-o de novo; acompanhe o progresso aqui.',
      done: '{name} apagado. O servidor transfere-o de novo no próximo arranque, se continuar inscrito.',
    },
    state: {
      installed: 'Instalado',
      subscribed: 'Vai ser transferido',
      orphan: 'Sem subscrição (vai ser removido)',
      active: 'Ativo',
      installing: 'A transferir',
      updating: 'A atualizar',
      uninstalling: 'A remover',
      removed: 'Removido',
      pending: 'Em espera',
    },
  },
  backup: {
    title: 'Cópia de segurança',
    subtitle:
      'Um único ficheiro cifrado com tudo o que é preciso para montar este servidor noutro PC.',
    create: {
      title: 'Criar cópia de segurança',
      contents:
        'Inclui: definições do servidor, palavras-passe e tokens, Game.ini, Engine.ini, Admins.txt, MapCycle.txt, Bans.json e Motd.txt.',
      password: 'Palavra-passe da cópia',
      confirm: 'Repita a palavra-passe',
      passwordHint: 'Pelo menos 10 caracteres. Sem ela a cópia não abre: não há recuperação.',
      mismatch: 'As palavras-passe não coincidem.',
      includeModio: 'Incluir a sessão do mod.io (o novo PC fica autorizado)',
      includeModioWarn:
        'A sessão funciona como uma palavra-passe. Guarde o ficheiro num local seguro.',
      button: 'Guardar cópia…',
      done: 'Cópia guardada: {file}',
    },
    restore: {
      title: 'Restaurar neste PC',
      choose: 'Escolher ficheiro de cópia…',
      file: 'Ficheiro: {file}',
      password: 'Palavra-passe da cópia',
      open: 'Abrir cópia',
      createdAt: 'Criada',
      machine: 'PC',
      hostname: 'Servidor',
      files: 'Ficheiros',
      secrets: 'Palavras-passe e tokens',
      modio: 'Sessão do mod.io',
      included: 'incluída',
      notIncluded: 'não incluída',
      serverDir: 'Pasta do servidor neste PC',
      restoreModio: 'Restaurar a sessão do mod.io',
      button: 'Restaurar',
      done: 'Restaurado. A seguir: instale o servidor em Instalação (se este PC não o tiver) e inicie-o.',
      runningWarn: 'Pare o servidor antes de restaurar.',
    },
  },
  settings: {
    startWithWindows: 'Abrir o app quando o Windows inicia',
    startWithWindowsHint: 'O app abre minimizado quando entra no Windows.',
    autoStartServer: 'Iniciar o servidor quando o app abre com o Windows',
    autoStartServerHint: 'Útil num computador que fica ligado como servidor.',
    title: 'Definições',
    subtitle: 'Idioma, aspeto e preferências da aplicação.',
    language: 'Idioma',
    appearance: 'Aspeto',
    color: 'Cor',
    skin: 'Disposição',
    savedHere: 'Guardado neste computador.',
    themes: {
      sand: 'Areia',
      black: 'Preto',
      olive: 'Azeitona',
      blue: 'Azul escuro',
      light: 'Claro',
    },
    skins: {
      sidebar: 'Barra lateral',
      sidebarDesc: 'Menu com nomes à esquerda.',
      compact: 'Compacta',
      compactDesc: 'Só ícones à esquerda, mais espaço para o conteúdo.',
      cards: 'Cartões',
      cardsDesc: 'Painéis flutuantes com cantos arredondados.',
      top: 'Barra no topo',
      topDesc: 'Menu numa barra horizontal no topo.',
    },
    console: 'Mostrar a janela de consola do servidor',
    consoleHint:
      'Desligado: a consola fica oculta e os registos aparecem na aplicação. Aplica-se no próximo arranque.',
    logs: {
      title: 'Registos do servidor',
      text: 'O servidor guarda um ficheiro de registo por arranque. Agora há {n} arranques guardados.',
      auto: 'Apagar registos antigos automaticamente',
      autoOff: 'Desligado: o app não apaga nenhum registo.',
      autoOn:
        'Ligado: a cada arranque do servidor, os arranques mais antigos são apagados e os mais recentes ficam.',
      keep: 'Arranques a manter',
      keepHint:
        'De {min} a {max}. O mínimo é o arranque atual e o anterior: o app lê-os para "Trazer do servidor" e para achar a chave do mod.io.',
      cleanNow: 'Apagar registos antigos agora',
      confirm: 'Apagar {n} arranques antigos ({size})? Os {keep} mais recentes ficam.',
      confirmYes: 'Apagar {n}',
      wouldDelete: 'A manter {keep} arranques, seriam apagados {n} arranques antigos ({size}).',
      nothing: 'Nada para apagar: há {keep} arranques ou menos.',
      done: '{n} ficheiros apagados ({size} libertados).',
    },
  },
  about: {
    title: 'Sobre',
    text: 'O Sandstorm Server Assist ajuda a instalar, configurar e executar um servidor dedicado de Insurgency: Sandstorm com mods do mod.io.',
    license: 'Licença',
    licenseText:
      'Apache License 2.0. Livre para usar, alterar e partilhar, fornecido "tal como está", sem garantia.',
    disclaimer: 'Aviso legal',
    disclaimerText:
      'Projeto não oficial, sem ligação nem aprovação da New World Interactive, Focus Entertainment, Valve ou mod.io. Insurgency, Steam e mod.io são marcas dos respetivos donos. É responsável pela forma como usa o seu servidor e por cumprir os termos do jogo, da Steam e do mod.io e as leis do seu país. A aplicação não distribui o jogo nem o SteamCMD: transfere-os das fontes oficiais.',
    links: 'Ligações úteis',
    serverGuide: 'Guia de servidores da comunidade',
    modioGuide: 'Guia de administração de servidores do mod.io',
    credits: 'Créditos',
    creditsText:
      'Criado por Felipe Sulzbach, com a ajuda do Claude Code (Anthropic). Gratuito e de código aberto, feito para a comunidade de Insurgency: Sandstorm.',
  },
  feedback: {
    title: 'Feedback',
    subtitle: 'Relatar um problema ou sugerir uma melhoria',
    kindLabel: 'Tipo',
    kind: {
      bug: 'Relatar um problema',
      idea: 'Sugerir uma melhoria',
    },
    intro: {
      bug: 'Conte o que correu mal. A app prepara um ficheiro com os detalhes e abre o formulário do GitHub já preenchido.',
      idea: 'Partilhe uma ideia para melhorar a app. A app abre o formulário do GitHub já preenchido.',
    },
    titleField: 'Título (opcional)',
    titleHint: 'Uma linha. Se ficar vazio, a app escreve um.',
    description: 'Descrição (opcional)',
    descriptionHint: {
      bug: 'O que fez, o que esperava e o que aconteceu.',
      idea: 'O que gostaria que a app fizesse, e porque ajudaria.',
    },
    includeLog: 'Incluir o registo da app',
    includeLogHint:
      'As últimas linhas do painel Registos. Senhas, tokens, endereços IP e SteamIDs ficam mascarados.',
    image: 'Captura de ecrã (opcional)',
    imageHint: 'PNG, JPG, GIF ou WebP, até 5 MB.',
    chooseImage: 'Escolher imagem',
    removeImage: 'Remover',
    create: 'Preparar e abrir o GitHub',
    privacy:
      'A app não envia nada. Guarda um ficheiro .zip neste PC e abre o formulário no navegador; é você quem decide o que publicar. É preciso uma conta gratuita no GitHub.',
    done: 'Guardado {file} ({size}).',
    next: [
      'Entre no GitHub na janela do navegador que abriu.',
      'Arraste o ficheiro .zip (mostrado na janela do Explorador) para a última caixa do formulário.',
      'Reveja o texto e clique em "Create".',
    ],
    showFile: 'Mostrar o ficheiro',
  },
  update: {
    title: 'Atualizações',
    check: 'Procurar atualizações',
    install: 'Atualizar',
    installFailed: 'A atualização não começou',
    readyHint:
      'A versão {version} já foi baixada. Clique para reiniciar o app nela (ou ela é instalada quando fechar o app).',
    manualHint:
      'A versão {version} está disponível. Esta é a versão portátil: baixe a nova na página de download.',
    manualButton: 'Versão {version}',
    errorCheck: 'Não foi possível procurar atualizações agora. O app tenta de novo mais tarde.',
    errorDownload: 'Não foi possível baixar a atualização. O app tenta de novo mais tarde.',
    hint: 'O app procura versões novas sozinho e só instala atualizações assinadas pelo projeto. O servidor do jogo continua a funcionar durante a atualização.',
    status: {
      idle: 'O app procura versões novas alguns segundos depois de abrir.',
      checking: 'A procurar uma versão nova…',
      uptodate: 'Tem a versão mais recente.',
      downloading: 'A baixar a versão {version}… {percent}%',
      ready:
        'A versão {version} está pronta. Clique em Atualizar (ou ela é instalada quando fechar o app).',
      manual:
        'A versão {version} está disponível. Esta é a versão portátil: baixe a nova na página de download.',
      dev: 'Versão de desenvolvimento ou portátil sem atualizações: a atualização automática está desligada.',
    },
  },
  legal: {
    title: 'Termos de uso e privacidade',
    intro: 'Antes de usar o app, leia e aceite os termos de uso e o aviso de privacidade.',
    tabs: {
      terms: 'Termos de uso',
      privacy: 'Privacidade e dados',
    },
    agree: 'Li e aceito os termos de uso e o aviso de privacidade.',
    accept: 'Aceitar e continuar',
    exit: 'Sair',
    version: 'Versão dos termos: {version}',
    acceptedOn: 'Versão {version} aceite em {date}.',
    terms: {
      about: {
        title: '1. O que é este app',
        text: [
          'O Sandstorm Server Assist é uma ferramenta gratuita, de código aberto e não oficial para instalar, configurar e iniciar um servidor dedicado de Insurgency: Sandstorm.',
          'Não tem ligação, aprovação nem patrocínio da New World Interactive, da Focus Entertainment, da Valve ou do mod.io. Os nomes e marcas deles aparecem só para dizer com o que o app funciona.',
        ],
      },
      asIs: {
        title: '2. Sem garantia',
        text: [
          'O app é fornecido "tal como está", sem garantia de qualquer tipo, sob a Licença Apache 2.0.',
          'Os autores e colaboradores não respondem por danos, perda de dados, ações sobre contas, banimentos, custos ou consequências legais que venham do uso.',
        ],
      },
      responsibility: {
        title: '3. A sua responsabilidade',
        text: [
          'Você é responsável pela forma como usa o app e pelos servidores que inicia com ele.',
          'Isso inclui seguir a licença e os termos do jogo, o Acordo de Assinante Steam e as regras de servidores da Steam, os termos do mod.io e as leis do seu país e dos seus jogadores.',
          'Também é responsável pelo conteúdo, mods, nomes e mensagens do servidor, e por quem deixa entrar e administrar.',
        ],
      },
      prohibited: {
        title: '4. O que não é permitido',
        text: [
          'Não use o app para fins ilegais, para prejudicar outras pessoas, para contornar proteções do jogo, da Steam ou do mod.io, nem para partilhar conteúdo que não tem direito de partilhar.',
        ],
      },
      thirdParty: {
        title: '5. Serviços de terceiros',
        text: [
          'O app não distribui o servidor do jogo, o SteamCMD nem mods. Ele baixa-os das fontes oficiais (Valve e mod.io), com os termos delas.',
          'Tokens e contas da Steam, do GameStats e do mod.io são seus e seguem as regras de cada serviço.',
        ],
      },
      updates: {
        title: '6. Atualizações automáticas',
        text: [
          'A app instalada procura versões novas nas publicações do projeto no GitHub, descarrega-as em segundo plano e instala-as quando clica em Atualizar ou fecha a app.',
          'Só são aceites atualizações assinadas com a chave do projeto. A versão portátil só avisa que há uma versão nova.',
        ],
      },
      changes: {
        title: '7. Mudanças nestes termos',
        text: [
          'Se estes termos ou o aviso de privacidade mudarem, o app mostra-os de novo e pede o seu aceite antes de continuar.',
        ],
      },
    },
    privacy: {
      summary: {
        title: '1. Resumo',
        text: [
          'O app não tem conta, anúncios nem telemetria. Não envia os seus dados para os autores do projeto.',
          'Tudo o que escreve fica neste PC, exceto o que vai para os serviços que escolhe usar (abaixo).',
        ],
      },
      stored: {
        title: '2. O que fica neste PC',
        text: [
          'Configurações do servidor e preferências do app: na pasta %APPDATA%\\sandstorm-server-assist.',
          'Senhas, tokens, o e-mail e a chave de API do mod.io: num cofre protegido pelo Windows (DPAPI), que só o seu utilizador do Windows neste PC consegue abrir.',
          'Backups: só onde os guardar, cifrados com a senha que escolher.',
          'Logs do servidor: lidos da pasta do servidor para mostrar no app. Senhas e tokens são mascarados antes de aparecerem.',
        ],
      },
      sent: {
        title: '3. O que sai deste PC, e para quem',
        text: [
          'Valve (Steam): download do SteamCMD e do servidor do jogo.',
          'mod.io: o e-mail da conta do servidor (para receber o código de acesso), o próprio código e a chave de API do servidor, para autorizar e baixar mods.',
          'GitHub (página do projeto): a procura de atualizações pede o ficheiro da versão mais recente e descarrega o instalador de lá. Como qualquer site, o GitHub vê o seu endereço IP e a hora do pedido. Nenhum identificador seu é enviado.',
          'Diagnóstico de portas (só quando clica em "Diagnosticar portas"): servidores STUN públicos da Cloudflare e da Google respondem com o seu endereço IP público; a Steam (Valve) é consultada sobre se lista um servidor nesse endereço; o seu router é consultado sobre os reencaminhamentos (UPnP, dentro da sua rede). Nada mais é enviado.',
          'Discord (só se definir um webhook em Servidor › Automação): os textos dos avisos que escreveu, com o nome do servidor, o mapa e os nomes dos jogadores que entram ou saem, vão para o seu canal do Discord. O Discord recebe o seu IP, como qualquer site.',
          'Feedback (relatar um problema ou sugerir uma melhoria): a app não envia nada. Guarda um ficheiro neste PC e abre o formulário do GitHub no navegador; é você quem escolhe o que publicar.',
          'Os links em que clica abrem no seu navegador, com as regras de cada site.',
        ],
      },
      server: {
        title: '4. O seu servidor do jogo',
        text: [
          'O servidor do jogo que inicia liga-se à Steam, ao GameStats, ao mod.io e aos jogadores, com as regras desses serviços. Os SteamIDs de administradores que adicionar são gravados nos arquivos do servidor (Admins.txt).',
          'Se jogadores entrarem no seu servidor, é você o responsável pelos dados deles nele (nomes, SteamIDs, chat, logs).',
        ],
      },
      control: {
        title: '5. O seu controlo',
        text: [
          'Pode ver e mudar tudo no app. Para apagar tudo: desinstale o app e apague a pasta %APPDATA%\\sandstorm-server-assist.',
          'Tokens e chaves podem ser revogados a qualquer momento nos sites da Steam e do mod.io.',
        ],
      },
      changes: {
        title: '6. Mudanças',
        text: ['Se este aviso mudar, o app mostra-o de novo e pede o seu aceite.'],
      },
    },
  },
  logs: {
    button: 'Registos',
    errors: '{n} erros',
    oneError: '1 erro',
    title: 'Registos',
    filter: {
      all: 'Tudo',
      important: 'Importantes',
      warning: 'Avisos',
      error: 'Erros',
    },
    search: 'Pesquisar…',
    follow: 'Acompanhar',
    copy: 'Copiar visíveis',
    openFolder: 'Abrir pasta de registos',
    clear: 'Limpar',
    empty: 'Ainda não há linhas de registo.',
    source: {
      server: 'servidor',
      steamcmd: 'steamcmd',
      app: 'app',
    },
    run: 'Arranque do servidor',
    runLive: 'Arranque atual (iniciado {time})',
    runSession: 'Arranque atual',
    runAt: 'Arranque de {time} ({size})',
    runLast: 'Último arranque, {time} ({size})',
    runStarted: 'Servidor iniciado às {time}',
    savedRun: 'Arranque guardado, lido do ficheiro de registo do servidor. {n} erros.',
    savedRunOneError: 'Arranque guardado, lido do ficheiro de registo do servidor. 1 erro.',
    backToLive: 'Voltar ao arranque atual',
    loading: 'A carregar…',
    masked: 'Palavras-passe e tokens aparecem como •••.',
  },
  error: {
    modFile: 'Ficheiro indicado no erro: {object}. Nenhum mod instalado o tem.',
    modOwners: 'Ficheiro indicado no erro: {object}. Está em: {mods}.',
    screenTitle: 'Este ecrã teve um erro',
    screenText:
      'O resto da aplicação continua a funcionar. Tente de novo ou volte ao Painel. O detalhe fica nos Registos.',
    appTitle: 'A aplicação teve um erro',
    appText:
      'Recarregue a aplicação para continuar a usá-la. O servidor continua a correr. O detalhe fica nos Registos.',
    reload: 'Recarregar a aplicação',
    retry: 'Tentar de novo',
    home: 'Voltar ao Painel',
    what: 'O que aconteceu',
    cause: 'Causa provável',
    solution: 'Como resolver',
    details: 'Detalhes técnicos',
    openLogs: 'Abrir registos',
    help: 'Ajuda',
  },
  links: {
    'steam-gslt': 'Steam: contas de servidor de jogo',
    gamestats: 'Sandstorm GameStats',
    'steamid-io': 'steamid.io',
    'steam-profile': 'O seu perfil Steam',
    modio: 'mod.io: Insurgency: Sandstorm',
    'modio-subscriptions': 'mod.io: subscrições',
    'modio-api-access': 'mod.io: acesso à API',
    'modio-status': 'Estado do mod.io',
    'modio-admin-guide': 'Guia de administração do mod.io',
    'steamcmd-docs': 'Documentação do SteamCMD',
    'server-guide': 'Guia de servidores da comunidade',
    'apache-license': 'Apache License 2.0',
    wails: 'Wails',
    'download-page': 'Versões no GitHub',
    project: 'Projeto no GitHub',
    'user-guide': 'Guia de utilização',
    author: 'Felipe Sulzbach no GitHub',
  },
  guides: {
    steamId: {
      title: 'Encontrar um SteamID64',
      steps: [
        'Abra o perfil Steam da pessoa (na Steam: clique no nome › Perfil).',
        'Se o endereço for steamcommunity.com/profiles/7656119…, o número de 17 dígitos é o SteamID64. Cole aqui a ligação completa: a aplicação encontra-o.',
        'Se o endereço tiver um nome personalizado (steamcommunity.com/id/nome), cole-o em steamid.io e copie o "steamID64".',
        'Na aplicação Steam, ative Definições › Interface › "Mostrar barras de endereço web quando disponíveis" para ver o endereço.',
      ],
    },
    gslt: {
      title: 'Obter um token GSLT',
      steps: [
        'Abra a página de contas de servidor de jogo da Steam e inicie sessão.',
        'Em "App ID of base game" escreva 581320 (Insurgency: Sandstorm) e dê-lhe uma nota, por exemplo o nome do servidor.',
        'Clique em Create e copie o "Login Token".',
        'Um token por servidor. A sua conta Steam não pode ser limitada (precisa de pelo menos uma compra).',
      ],
    },
    gameStats: {
      title: 'Obter um token GameStats',
      steps: [
        'Abra a página Sandstorm GameStats e inicie sessão com a Steam.',
        'Crie um token para o seu servidor e copie-o.',
        'Cole-o aqui. Sem ele o servidor funciona na mesma; os jogadores só não ganham estatísticas.',
      ],
    },
    modioAccount: {
      title: 'Conta do mod.io para o servidor',
      steps: [
        'Crie uma conta no mod.io com um e-mail só para o servidor (não use a sua conta pessoal).',
        'Com essa conta, abra o Insurgency: Sandstorm no mod.io e subscreva os mods que o servidor deve usar.',
        'Escreva esse e-mail aqui, envie o código e autorize. O servidor transfere os mods subscritos em cada arranque.',
      ],
    },
    modioApiKey: {
      title: 'Obter uma chave da API do mod.io',
      steps: [
        'Experimente primeiro "Procurar nos registos do servidor": o servidor escreve a chave no registo depois de arrancar uma vez com mods.',
        'Se não encontrar: inicie sessão no mod.io e abra "API access".',
        'Crie uma chave de API só de leitura e copie-a (32 letras e números).',
        'A chave fica no cofre. Não faz parte do código da aplicação, porque o projeto é público.',
      ],
    },
    ports: {
      title: 'Pela internet: abrir as portas do router',
      steps: [
        'Abra a página de administração do router (normalmente 192.168.0.1 ou 192.168.1.1) e procure "Reencaminhamento de portas" (pode chamar-se "Servidor virtual" ou "NAT").',
        'Reencaminhe UDP {port} e UDP {queryPort} para o IP deste PC na rede local (execute ipconfig e veja "Endereço IPv4").',
        'Permita o InsurgencyServer.exe na Firewall do Windows.',
        'Os jogadores usam o seu IP público (pesquise "qual é o meu IP"). Escreva-o em "IP que os jogadores usam".',
        'Não reencaminhe a porta do RCON ({rconPort}): ela controla o servidor.',
        'Se mesmo assim não funcionar, o seu fornecedor pode partilhar um IP público entre muitos clientes (CGNAT). Peça um IP público ao fornecedor ou use um programa de VPN.',
      ],
    },
    vpn: {
      title: 'Por um programa de VPN (rede virtual)',
      steps: [
        'Instale o mesmo programa de rede virtual no PC do servidor e no PC de cada jogador.',
        'Crie uma rede no PC do servidor e partilhe o nome e a palavra-passe com os jogadores.',
        'O IP do servidor é o que o programa mostra para este PC. Escreva-o em "IP que os jogadores usam".',
        'Não é preciso abrir portas no router.',
      ],
    },
  },
  lighting: {
    label: 'Iluminação',
    default: 'Padrão do cenário',
    Day: 'Dia',
    Night: 'Noite',
    start: 'Iluminação do mapa inicial',
    cycleTitle: 'Dia ou noite na rotação',
    allDay: 'Tudo de dia',
    allNight: 'Tudo de noite',
    allDefault: 'Tudo no padrão',
    unknownLine: 'linha que o app não lê',
    cycleMissing:
      'Alguns mapas da rotação não estão no jogo base nem em nenhum mod instalado. Quando a rotação chegar neles, o servidor não tem para onde ir.',
    cycleRemoveMissing: 'Remover estes mapas',
    notInstalled: 'mod não instalado',
    remove: 'Remover',
  },
  mutators: {
    AllYouCanEat: 'Começa com 100 pontos de abastecimento.',
    AntiMaterielRiflesOnly: 'Só espingardas antimaterial, mais o equipamento normal e explosivos.',
    BoltActionsOnly: 'Só espingardas de ferrolho, mais o equipamento normal e explosivos.',
    Broke: 'Começa com 0 pontos de abastecimento.',
    BulletSponge: 'Mais vida.',
    Competitive: 'Equipamento mais caro, rondas mais curtas, capturas mais rápidas.',
    CompetitiveLoadouts: 'Classes trocadas pelas do Competitivo.',
    FastMovement: 'Movimento mais rápido.',
    Frenzy: 'Inimigos da IA que só atacam corpo a corpo. Cuidado com os inimigos especiais.',
    Guerrillas: 'Começa com 5 pontos de abastecimento.',
    Hardcore: 'Movimento mais lento e capturas mais longas.',
    HeadshotOnly: 'Os jogadores só levam dano com tiros na cabeça.',
    HotPotato: 'Uma granada ativa cai quando um jogador morre.',
    LockedAim: 'As armas apontam sempre para o centro do ecrã.',
    NoAim: 'Sem mira pela alça.',
    PistolsOnly: 'Só pistolas, mais o equipamento normal e explosivos.',
    ShotgunsOnly: 'Só caçadeiras, mais o equipamento normal e explosivos.',
    SlowCaptureTimes: 'Os objetivos demoram mais a capturar.',
    SlowMovement: 'Movimento mais lento.',
    SoldierOfFortune: 'Ganha pontos de abastecimento conforme a pontuação sobe.',
    SpecialOperations: 'Começa com 30 pontos de abastecimento.',
    Strapped: 'Começa com 1 ponto de abastecimento.',
    Ultralethal: 'Todos morrem com um tiro.',
    Vampirism: 'Ganha vida igual ao dano que causa aos inimigos.',
    Warlords: 'Começa com 10 pontos de abastecimento.',
  },
  rules: {
    manage: 'Definir as regras do jogo aqui',
    manageHint:
      'O app grava estas opções no Game.ini e no Engine.ini quando o servidor inicia. Opções vazias usam o padrão do jogo. Desligado: o app não mexe nestes ficheiros.',
    fromServer:
      'O Game.ini do servidor já define {n} regras do jogo. A aplicação ainda não as gere, por isso deixa os ficheiros como estão.',
    useServer: 'Usar as regras do servidor',
    default: 'Padrão: {value}',
    gameDefault: 'Padrão do jogo',
    on: 'Ligado',
    off: 'Desligado',
    groups: {
      coop: 'Coop: bots',
      checkpoint: 'Checkpoint',
      objectives: 'Objetivos',
      friendlyFire: 'Fogo amigo e equipas',
      teamKills: 'Mortes de colegas',
      chat: 'Chat, voz e lista de mortes',
      match: 'Partida',
      supply: 'Pontos de abastecimento',
      votes: 'Votação para expulsar',
      network: 'Desempenho (Engine.ini)',
    },
    items: {
      aiDifficulty: {
        label: 'Dificuldade dos bots',
        hint: '0 é o mais fácil, 1 o mais difícil',
      },
      minimumEnemies: {
        label: 'Mínimo de inimigos',
        hint: 'Bots inimigos com poucos jogadores',
      },
      maximumEnemies: {
        label: 'Máximo de inimigos',
        hint: 'Bots inimigos com o servidor cheio',
      },
      friendlyBotQuota: {
        label: 'Bots aliados',
        hint: 'Bots na equipa dos jogadores',
      },
      vehicleInsertion: {
        label: 'Inimigos chegam de veículo',
        hint: 'Entrada com camiões onde o mapa tem',
      },
      defendTimer: {
        label: 'Tempo do contra-ataque (s)',
        hint: 'Tempo para defender cada ponto',
      },
      defendTimerFinal: {
        label: 'Tempo do contra-ataque final (s)',
        hint: 'Tempo para defender o último ponto',
      },
      retreatTimer: {
        label: 'Tempo de recuo dos bots (s)',
        hint: 'Depois do contra-ataque',
      },
      respawnDpr: {
        label: 'Proporção para os bots voltarem',
        hint: 'Parte dos bots mortos antes de voltarem',
      },
      respawnDelay: {
        label: 'Espera para os bots voltarem (s)',
        hint: 'Antes de os bots voltarem',
      },
      postCaptureRushTimer: {
        label: 'Investida depois do esconderijo (s)',
        hint: 'Os bots avançam depois de um esconderijo destruído',
      },
      counterAttackRespawnDpr: {
        label: 'Proporção no contra-ataque',
        hint: 'Parte dos bots mortos antes de voltarem',
      },
      counterAttackRespawnDelay: {
        label: 'Espera no contra-ataque (s)',
        hint: 'Antes de os bots voltarem',
      },
      enemyRespawnMin: {
        label: 'Vagas inimigas com poucos jogadores',
        hint: 'Multiplicador',
      },
      enemyRespawnMax: {
        label: 'Vagas inimigas com o servidor cheio',
        hint: 'Multiplicador',
      },
      finalCacheBots: {
        label: 'Bots no esconderijo final',
        hint: 'Multiplicador',
      },
      captureTime: {
        label: 'Tempo de captura (s)',
        hint: 'Tempo para tomar um ponto',
      },
      resetTime: {
        label: 'Tempo para a captura voltar (s)',
        hint: '-1 mantém o progresso',
      },
      captureSpeedup: {
        label: 'Aceleração por jogador',
        hint: 'Cada jogador a mais captura mais depressa',
      },
      captureSpeedupPlayers: {
        label: 'Jogadores que aceleram',
        hint: 'Máximo contado',
      },
      friendlyFire: {
        label: 'Fogo amigo',
        hint: 'Dano aos colegas',
      },
      friendlyFireModifier: {
        label: 'Dano do fogo amigo',
        hint: '0 é nenhum, 1 é o dano inteiro',
      },
      friendlyFireReflect: {
        label: 'Dano devolvido ao atacante',
        hint: '0 é nenhum, 1 é o dano inteiro',
      },
      autoBalance: {
        label: 'Equilibrar as equipas',
        hint: 'Muda jogadores quando as equipas ficam desiguais',
      },
      autoBalanceDelay: {
        label: 'Espera para equilibrar (s)',
        hint: 'Depois de as equipas ficarem desiguais',
      },
      teamSelect: {
        label: 'Jogadores escolhem a equipa',
        hint: 'Permite trocar de equipa',
      },
      teamKillLimit: {
        label: 'Mortes de colegas até expulsar',
        hint: '0 desliga a expulsão',
      },
      teamKillGrace: {
        label: 'Tempo entre mortes de colegas (s)',
        hint: 'Mortes neste intervalo contam uma vez',
      },
      teamKillReduceTime: {
        label: 'Perdoar uma morte de colega após (s)',
        hint: 'A contagem desce 1',
      },
      killFeed: {
        label: 'Lista de mortes',
        hint: 'Mostra quem matou quem',
      },
      killFeedSpectator: {
        label: 'Lista de mortes para espectadores',
        hint: 'E nas repetições',
      },
      killerInfo: {
        label: 'Mostrar quem o matou',
        hint: 'Atacante e arma no ecrã de morte',
      },
      killDistance: {
        label: 'Mostrar a distância da morte',
        hint: 'No ecrã de morte',
      },
      deadSay: {
        label: 'Os vivos leem os mortos',
        hint: 'Chat dos jogadores mortos',
      },
      deadSayTeam: {
        label: 'Ler colegas mortos',
        hint: 'Chat de equipa dos mortos',
      },
      deadVoice: {
        label: 'Os vivos ouvem os mortos',
        hint: 'Voz dos jogadores mortos',
      },
      enemyHearsLocal: {
        label: 'Inimigos ouvem a voz local',
        hint: 'Voz de proximidade',
      },
      intermissionVoice: {
        label: 'Todos falam entre rondas',
        hint: 'As duas equipas',
      },
      mapVoting: {
        label: 'Votação de mapa',
        hint: 'No fim da partida',
      },
      useMapCycle: {
        label: 'Usar a rotação de mapas',
        hint: 'Desligado repete o mesmo mapa',
      },
      minimumPlayers: {
        label: 'Jogadores para começar',
        hint: 'Por equipa',
      },
      roundLimit: {
        label: 'Rondas por partida',
        hint: 'Varia com o modo',
      },
      winLimit: {
        label: 'Rondas para vencer',
        hint: 'Varia com o modo',
      },
      roundTime: {
        label: 'Tempo da ronda (s)',
        hint: 'Varia com o modo',
      },
      preRoundTime: {
        label: 'Tempo parado no início (s)',
        hint: 'No começo da ronda',
      },
      postRoundTime: {
        label: 'Tempo entre rondas (s)',
        hint: 'Ecrã de intervalo',
      },
      idleLimit: {
        label: 'Expulsar inativos após (s)',
        hint: '0 nunca expulsa',
      },
      idleLimitLow: {
        label: 'Limite de inatividade com poucos reforços (s)',
        hint: 'Quando a equipa tem poucas vidas',
      },
      bots: {
        label: 'Bots no versus',
        hint: 'Completa as equipas com bots',
      },
      botQuota: {
        label: 'Número de bots no versus',
        hint: 'Com os bots ligados',
      },
      initialSupply: {
        label: 'Pontos de abastecimento iniciais',
        hint: 'Para comprar equipamento',
      },
      maximumSupply: {
        label: 'Máximo de pontos de abastecimento',
        hint: 'Limite que se pode ganhar',
      },
      supplyGain: {
        label: 'Ganhar abastecimento na ronda',
        hint: 'Pela pontuação',
      },
      voteKick: {
        label: 'Votação para expulsar',
        hint: 'Os jogadores abrem a votação em Esc › Call Vote',
      },
      voteMinPlayers: {
        label: 'Jogadores para abrir votação',
        hint: 'Mínimo no servidor',
      },
      votePassRatio: {
        label: 'Votos sim para passar',
        hint: 'Parte da equipa (0 a 1)',
      },
      voteBanDuration: {
        label: 'Banimento após a votação (s)',
        hint: '-1 só expulsa',
      },
      tickRate: {
        label: 'Tick rate',
        hint: 'Atualizações por segundo; mais usa mais CPU e banda',
      },
      maxClientRate: {
        label: 'Taxa máxima por cliente (rede local)',
        hint: 'Bytes por segundo por jogador',
      },
      maxInternetClientRate: {
        label: 'Taxa máxima por cliente (internet)',
        hint: 'Bytes por segundo por jogador',
      },
    },
  },
  profiles: {
    label: 'Perfil',
    saveFirst: 'Guarde ou desfaça as mudanças antes de trocar de perfil.',
    new: 'Perfil novo',
    rename: 'Mudar o nome',
    switchHint:
      'Trocar de perfil exige o servidor parado. As palavras-passe e os tokens são partilhados.',
    namePlaceholder: 'Exemplo: Coop, Versus',
    copyCurrent: 'Copiar o perfil atual',
    create: 'Criar',
    delete: 'Apagar perfil',
    deleteYes: 'Apagar',
  },
  players: {
    title: 'Jogadores',
    load: 'Mostrar jogadores',
    reload: 'Atualizar',
    empty: 'Ninguém no servidor.',
    name: 'Nome',
    score: 'Pontos',
    actions: 'Ações',
    bot: 'bot',
    showBots: 'Mostrar os {n} bots',
    hideBots: 'Esconder os {n} bots',
    noHumans: 'Nenhum jogador, só bots.',
    kick: 'Expulsar',
    ban: 'Banir',
    confirmKick: 'Expulsar {name}?',
    confirmBan: 'Banir {name}?',
    duration: 'Duração',
    reason: 'Motivo (opcional)',
    confirmYes: 'Confirmar',
    kicked: '{name} foi expulso.',
    banned: '{name} foi banido.',
    unbanned: 'Banimento retirado.',
    bans: 'Banidos',
    noBans: 'Nenhum banimento.',
    unban: 'Retirar banimento',
    permanent: 'Permanente',
    days: '{n} dia(s)',
    minutes: '{n} minutos',
    sayPlaceholder: 'Mensagem para todos no servidor',
    say: 'Enviar',
    said: 'Mensagem enviada.',
  },
  automation: {
    restartTitle: 'Reinício automático',
    daily: 'Reiniciar todos os dias às',
    dailyHint: 'Vazio = sem reinício diário. Hora do computador.',
    next: 'Próximo reinício: {when}',
    update: 'Atualizar o servidor pelo SteamCMD no reinício diário',
    updateHint:
      'As atualizações saem com os patches do jogo; o reinício demora alguns minutos a mais.',
    updateNoSteamCmd: 'Defina primeiro a pasta do SteamCMD em Instalação.',
    warn: 'Aviso no chat antes do reinício',
    warnHint: 'Enviado 5 e 1 minuto antes. {minutes} vira o número.',
    warnExample: 'O servidor reinicia em {minutes} min',
    onCrash: 'Iniciar o servidor de novo se ele fechar sozinho',
    onCrashHint: 'Espera 10 segundos. Para após 3 vezes em 30 minutos.',
    appOpen: 'A automação só funciona com o app aberto (pode ficar minimizado).',
    discordTitle: 'Avisos no Discord',
    discordText:
      'O app publica num canal do Discord quando o servidor fica no ar, para, fecha sozinho, ou quando alguém entra ou sai. Vazio = esse aviso está desligado.',
    webhook: 'Endereço do webhook',
    webhookHint:
      'No Discord: definições do canal › Integrações › Webhooks › Novo webhook › Copiar URL.',
    discord: {
      ready: 'Servidor no ar',
      stopped: 'Servidor parado',
      crashed: 'Servidor fechou sozinho',
      join: 'Jogador entrou',
      leave: 'Jogador saiu',
    },
    discordExample: {
      ready: '🟢 {server} está no ar em {map}',
      stopped: '⚪ {server} parou',
      crashed: '🔴 {server} fechou sozinho',
      join: '➡️ {name} entrou em {server}',
      leave: '⬅️ {name} saiu de {server}',
    },
    fillExamples: 'Preencher com exemplos',
    test: 'Testar',
    testSent: 'Enviado. Veja o canal.',
    tags: '{server} = nome do servidor, {map} = mapa, {name} = jogador. O teste usa os textos guardados.',
  },
  firewall: {
    title: 'Firewall do Windows',
    notInstalled: 'Instale primeiro o servidor.',
    ok: 'Regras para o servidor encontradas na Firewall do Windows.',
    text: 'Sem regra, o Windows pode bloquear os jogadores, sobretudo em redes marcadas como "Pública" (comum com programas de VPN).',
    create: 'Criar a regra da firewall',
    uac: 'O Windows pede a confirmação de administrador.',
  },
  validation: {
    markupPlain: 'O jogo não mostra formatação e mostraria {tag} escrito. Tire as tags.',
    rule: 'O valor de {rule} tem de estar entre {min} e {max}.',
    ruleBool: 'Escolha ligado ou desligado em {rule}.',
    ruleMinMax: 'O mínimo de inimigos não pode ser maior que o máximo.',
    lighting: 'Escolha dia, noite ou o padrão do cenário.',
    lightingNoMap: 'Este mapa não permite escolher dia ou noite pelo app. Deixe no padrão.',
    bindIp: 'Escolha um endereço IPv4 válido deste computador.',
    modIdsMax: 'Até {max} mods.',
    modListEmpty: 'Escreva pelo menos um id de mod, ou desligue a lista de mods.',
    webhook:
      'Cole o endereço do webhook copiado do Discord (começa com https://discord.com/api/webhooks/).',
    announcementsMax: 'Até {max} mensagens programadas.',
    announceEvery: 'Escolha entre {min} e {max} minutos.',
    clock: 'Use o formato HH:MM, por exemplo 05:00.',
    discordText: 'Até {max} caracteres, numa linha.',
    reason: 'O motivo pode ter até {max} caracteres, numa linha e sem aspas.',
    netId: 'Este jogador não está na lista. Atualize a lista e tente de novo.',
    banMinutes: 'Escolha uma duração de até {max} minutos.',
    profileName: 'Escreva um nome com 1 a 40 caracteres.',
    profilesMax: 'Até {max} perfis.',
    motd: 'A mensagem do dia tem até {max} caracteres e nenhum caractere especial.',
    chatMessage: 'A mensagem tem até {max} caracteres, numa linha.',
    chatNeedsRcon: 'Ligue o RCON em Palavras-passe e tokens: o app envia esta mensagem pelo RCON.',
    scenarioUnknown: 'Este mapa não está na lista. Peça a lista ao servidor de novo.',
    scenarioNotLoaded: 'O mod deste mapa não está carregado no servidor.',
    logKeep: 'Os arranques a manter têm de estar entre {min} e {max}.',
    unknownField: 'Campo desconhecido.',
    password: 'Use de 4 a 64 caracteres: letras, números e ! @ # $ % * . _ + - (sem espaços).',
    token: 'O token só tem letras e números (sem espaços). Copie-o de novo do site.',
    apiKey: 'A chave da API tem 32 letras e números.',
    email: 'Este e-mail não parece válido.',
    emailTypo: 'Quis dizer {suggestion}?',
    absolutePath: 'A pasta tem de ser um caminho completo.',
    hostnameLength: 'O nome do servidor precisa de 1 a 64 caracteres.',
    hostnameQuotes: 'O nome do servidor não pode ter aspas (").',
    port: 'As portas vão de 1024 a 65535.',
    portsDistinct: 'As portas do jogo, de consulta e RCON têm de ser diferentes.',
    maxPlayers: 'Os jogadores têm de ser entre 1 e 64.',
    map: 'Nome de mapa inválido (só letras, números e _).',
    scenario: 'O cenário começa por Scenario_ (exemplo: Scenario_Farmhouse_Checkpoint_Security).',
    travelNeedsMods: 'O mapa do mod só funciona com os mods ligados.',
    mutator: 'Mutator inválido: {value}.',
    ip: 'Escreva só o IP, por exemplo 203.0.113.10.',
    steamId: 'SteamID64 inválido: {value}. Tem 17 dígitos e começa por 7656119.',
    steamIdNotFound:
      'Nenhum SteamID64 encontrado. Cole o número de 17 dígitos ou uma ligação steamcommunity.com/profiles/….',
    adminName: 'Nome do administrador: até 40 caracteres, numa linha.',
    mapCycleLine: 'Cada linha da rotação de mapas tem até 200 caracteres.',
    modId: 'ID de mod inválido: {value}.',
    code: 'O código tem 5 letras ou números.',
    backupPassword: 'A palavra-passe da cópia precisa de pelo menos {min} caracteres.',
  },
  problems: {
    unknown: {
      title: 'Algo correu mal',
      what: 'A aplicação encontrou um erro que não conhece.',
      cause: 'Pode ser uma situação nova, depois de uma atualização do jogo ou do Windows.',
      solution: [
        'Abra os registos e veja as últimas linhas a vermelho.',
        'Tente de novo. Se repetir, reinicie a aplicação.',
      ],
    },
    'server-not-installed': {
      title: 'Servidor não instalado',
      what: 'O InsurgencyServer.exe não foi encontrado na pasta do servidor.',
      cause: 'A instalação ainda não foi feita, ou a pasta escolhida não é a do servidor.',
      solution: [
        'Abra Instalação e clique em Instalar servidor.',
        'Se o servidor estiver instalado noutro sítio, escolha essa pasta em Instalação.',
      ],
    },
    'server-dir-invalid': {
      title: 'Pasta inválida',
      what: 'A pasta escolhida não pode ser usada.',
      cause:
        'Não é um caminho completo, é a raiz de um disco ou fica dentro de uma pasta do sistema.',
      solution: [
        'Escolha uma pasta num disco com espaço (cerca de 25 GB mais os mods), por exemplo D:\\Servidores\\Sandstorm.',
        'Evite C:\\Windows, C:\\Program Files e raízes de disco como D:\\.',
      ],
    },
    'server-already-running': {
      title: 'O servidor já está em execução',
      what: 'Já está aberto um InsurgencyServer.exe desta pasta.',
      cause: 'Foi iniciado antes, pela aplicação ou por um ficheiro .bat.',
      solution: [
        'Use Parar servidor e depois Iniciar, ou Reiniciar.',
        'Não execute duas cópias: disputam a mesma porta.',
      ],
    },
    'server-not-running': {
      title: 'O servidor não está em execução',
      what: 'Esta ação precisa do servidor em execução.',
      cause: 'O servidor está parado.',
      solution: ['Inicie o servidor e tente de novo.'],
    },
    'server-start-failed': {
      title: 'O servidor não abriu',
      what: 'O Windows não conseguiu iniciar o InsurgencyServer.exe.',
      cause: 'O ficheiro pode estar em falta, danificado ou bloqueado pelo antivírus.',
      solution: [
        'Em Instalação, clique em Verificar ficheiros para o SteamCMD reparar a instalação.',
        'Confirme que o antivírus não pôs o InsurgencyServer.exe em quarentena.',
      ],
    },
    'server-crashed': {
      title: 'O servidor fechou sozinho',
      what: 'O processo do servidor terminou sem a aplicação o pedir.',
      cause:
        'Normalmente é um erro fatal do jogo: mods em conflito, um ficheiro danificado ou falta de memória.',
      solution: [
        'Abra os registos e procure a linha "Fatal error" mais recente.',
        'Se o erro indicar um ficheiro de mod, dois mods podem trazer o mesmo ficheiro (conflito).',
        'Em Instalação, use Verificar ficheiros.',
      ],
    },
    'server-stop-failed': {
      title: 'Não foi possível parar o servidor',
      what: 'O Windows recusou terminar o processo do servidor.',
      cause: 'O servidor pode ter sido iniciado como administrador, ou já estava a fechar.',
      solution: [
        'Espere alguns segundos e tente de novo.',
        'Se continuar, feche o InsurgencyServer.exe no Gestor de Tarefas.',
      ],
    },
    'port-in-use': {
      title: 'Porta ocupada',
      what: 'O servidor não conseguiu abrir a porta do jogo ou de consulta.',
      cause: 'Outro programa, ou outra cópia do servidor, está a usar a porta.',
      solution: [
        'Pare a outra cópia do servidor, se existir.',
        'Ou altere a porta em Servidor › Geral (por exemplo 27103) e avise os jogadores.',
      ],
    },
    'mod-conflict': {
      title: 'Mods em conflito',
      what: 'O servidor fechou com "Serial size mismatch" ou "Corrupt data found" ao carregar o mapa.',
      cause:
        'Dois mods trazem o mesmo ficheiro do jogo em versões diferentes, e um sobrepõe o outro. Já visto: Containment com Old World ou Containment Storm.',
      solution: [
        'Veja acima o ficheiro indicado no erro e o mod que o tem, ou use Mods › Mods com problemas.',
        'Cancele a subscrição desse mod no mod.io (conta do servidor) e inicie de novo.',
      ],
    },
    'modio-code-invalid': {
      title: 'Código do mod.io inválido',
      what: 'O mod.io respondeu 401 "Invalid security code".',
      cause:
        'O código de 5 caracteres só funciona uma vez e por pouco tempo. Este já foi usado ou expirou.',
      solution: ['Em Mods, peça um código novo por e-mail e autorize de novo.'],
    },
    'modio-not-authorized': {
      title: 'Servidor não autorizado no mod.io',
      what: 'O servidor não tem uma sessão válida do mod.io.',
      cause: 'Nunca foi autorizado neste PC, ou a autorização expirou.',
      solution: [
        'Em Mods, escreva o e-mail da conta do servidor e envie o código.',
        'Escreva o código recebido e clique em Autorizar.',
      ],
    },
    'modio-request-failed': {
      title: 'O mod.io não respondeu como esperado',
      what: 'O pedido à API do mod.io falhou.',
      cause: 'Sem internet, mod.io em baixo, ou demasiados pedidos.',
      solution: [
        'Confirme a internet e tente de novo dentro de alguns minutos.',
        'Veja a página de estado do mod.io.',
      ],
    },
    'modio-email-invalid': {
      title: 'E-mail inválido',
      what: 'O e-mail não parece válido.',
      cause: 'Um erro de escrita (por exemplo "gamil.com" em vez de "gmail.com").',
      solution: ['Confirme o e-mail da conta do mod.io do servidor e tente de novo.'],
    },
    'modio-no-api-key': {
      title: 'Falta a chave da API do mod.io',
      what: 'Pedir o código precisa de uma chave da API do mod.io.',
      cause: 'Ainda não há chave no cofre.',
      solution: [
        'Em Mods, clique em "Procurar nos registos do servidor".',
        'Ou crie uma chave no mod.io (acesso à API) e cole-a.',
      ],
    },
    'steamcmd-download': {
      title: 'Não foi possível transferir o SteamCMD',
      what: 'A transferência do SteamCMD não terminou.',
      cause: 'Sem internet, servidores da Valve instáveis, ou o antivírus bloqueou o ficheiro.',
      solution: [
        'Confirme a internet e tente de novo.',
        'Se já tiver o SteamCMD noutra pasta, escolha-a em Instalação.',
      ],
    },
    'steamcmd-signature': {
      title: 'Assinatura do SteamCMD inválida',
      what: 'O steamcmd.exe não tem a assinatura digital esperada da Valve, por isso não foi executado.',
      cause: 'O ficheiro pode ter sido alterado pelo caminho ou está danificado.',
      solution: [
        'Clique em Transferir SteamCMD para o transferir de novo.',
        'Se repetir, analise o PC com o antivírus.',
      ],
    },
    'steamcmd-failed': {
      title: 'O SteamCMD não conseguiu instalar',
      what: 'O SteamCMD terminou com um erro.',
      cause:
        'Frequentes: "Missing configuration" (temporário), "Disk write failure" (pasta sem permissão ou antivírus), "No subscription" (aplicação errada).',
      solution: [
        'Tente de novo: "Missing configuration" costuma passar à segunda.',
        'Confirme o espaço livre no disco e que a pasta não está protegida.',
      ],
    },
    'folders-nested': {
      title: 'Pasta do SteamCMD dentro da pasta do servidor',
      what: 'O SteamCMD não instala o servidor numa pasta que contém o próprio SteamCMD. Em vez disso, instala-o na pasta dele (steamapps\\common), e a app não encontra o servidor.',
      cause:
        'A pasta do SteamCMD está dentro da pasta do servidor, ou a pasta do servidor está dentro da pasta do SteamCMD.',
      solution: [
        'Deixe as duas pastas lado a lado, por exemplo D:\\Servers\\Sandstorm para o servidor e D:\\Servers\\steamcmd para o SteamCMD.',
        'Em Instalação, escolha as pastas de novo e clique em Instalar servidor.',
        'Se o SteamCMD já transferiu o servidor para a pasta dele (steamcmd\\steamapps\\common\\sandstorm_server), pode apagar essa pasta para libertar cerca de 5 GB.',
      ],
    },
    'disk-full': {
      title: 'Disco cheio',
      what: 'Não há espaço livre suficiente no disco.',
      cause: 'O servidor precisa de cerca de 25 GB, e os mods de mapa podem ter mais de 5 GB cada.',
      solution: ['Liberte espaço ou escolha uma pasta noutro disco.'],
    },
    'no-internet': {
      title: 'Sem ligação',
      what: 'A aplicação não conseguiu chegar à internet.',
      cause: 'Rede desligada, firewall a bloquear a aplicação, ou o site está em baixo.',
      solution: ['Confirme a ligação e tente de novo.'],
    },
    'invalid-field': {
      title: 'Informação inválida',
      what: 'Um campo tem um valor que o servidor não aceita.',
      cause: 'Veja a mensagem junto ao campo.',
      solution: ['Corrija o campo e guarde de novo.'],
    },
    'vault-unavailable': {
      title: 'Cofre de palavras-passe indisponível',
      what: 'A aplicação não conseguiu abrir o cofre onde guarda palavras-passe e tokens.',
      cause:
        'O cofre foi criado por outro utilizador do Windows ou noutro PC (a cifra do Windows só abre para quem o criou).',
      solution: [
        'Num PC novo, use Cópia de segurança › Restaurar com o ficheiro da cópia e a respetiva palavra-passe.',
        'Ou escreva de novo as palavras-passe e tokens em Servidor › Palavras-passe e tokens.',
      ],
    },
    'backup-password': {
      title: 'Palavra-passe da cópia errada',
      what: 'A cópia não abriu com esta palavra-passe.',
      cause: 'A palavra-passe está errada. Sem a correta a cópia não pode ser aberta.',
      solution: ['Confirme maiúsculas, minúsculas e o Caps Lock e tente de novo.'],
    },
    'backup-corrupt': {
      title: 'Ficheiro de cópia inválido',
      what: 'O ficheiro não é uma cópia desta aplicação, ou está danificado.',
      cause: 'O ficheiro foi alterado, ficou cortado ao copiar ou pertence a outro programa.',
      solution: ['Copie o ficheiro de novo a partir da origem e tente outra vez.'],
    },
    'backup-write': {
      title: 'Não foi possível guardar a cópia',
      what: 'O ficheiro da cópia não pôde ser gravado.',
      cause: 'Pasta sem permissão, disco cheio, ou a pen USB foi removida.',
      solution: ['Escolha outra pasta e tente de novo.'],
    },
    'config-write': {
      title: 'Não foi possível guardar a configuração',
      what: 'Um ficheiro de configuração do servidor não pôde ser gravado.',
      cause: 'O servidor pode ter o ficheiro aberto, ou a pasta está protegida.',
      solution: [
        'Pare o servidor e guarde de novo.',
        'Confirme as permissões da pasta do servidor.',
      ],
    },
    busy: {
      title: 'Espere pela tarefa atual',
      what: 'Ainda há outra tarefa em curso.',
      cause:
        'A instalação, a cópia de segurança e o arranque do servidor não podem decorrer ao mesmo tempo.',
      solution: ['Espere que termine ou cancele-a.'],
    },
    canceled: {
      title: 'Cancelado',
      what: 'A tarefa foi cancelada.',
      cause: 'Clicou em Cancelar.',
      solution: ['Inicie-a de novo quando quiser.'],
    },
    'update-failed': {
      title: 'Não foi possível instalar a atualização',
      what: 'A versão nova do app não conseguiu ser iniciada.',
      cause:
        'O instalador baixado foi apagado (por exemplo, pelo antivírus) ou o Windows o bloqueou.',
      solution: [
        'Feche e abra o app: ele baixa a versão nova de novo.',
        'Ou baixe o instalador na página de download e execute-o.',
      ],
    },
    'mods-not-mounted': {
      title: 'Alguns mods não foram carregados',
      what: 'O servidor baixou mods novos e não carregou os que já estavam instalados, por isso o mapa do mod e os mutators podem faltar.',
      cause:
        'Com subscrições novas, o servidor só carrega o que acabou de baixar. No arranque seguinte carrega tudo.',
      solution: [
        'Clique em Reiniciar. O app faz isto sozinho quando foi ele que iniciou o servidor.',
      ],
    },
    'mod-delete': {
      title: 'Não foi possível apagar o mod',
      what: 'A aplicação não conseguiu apagar os ficheiros do mod.',
      cause:
        'Um ficheiro do mod está em uso (pelo servidor ou por outro programa) ou o Windows bloqueou a pasta.',
      solution: [
        'Feche programas que possam estar a usar a pasta do mod e tente de novo.',
        'Se continuar a falhar, abra os Registos para ver qual ficheiro foi bloqueado.',
      ],
    },
    'log-delete': {
      title: 'Alguns registos não foram apagados',
      what: 'O app apagou o que conseguiu, mas alguns ficheiros de registo antigos ficaram.',
      cause:
        'O ficheiro está aberto noutro programa (por exemplo, um editor de texto) ou o Windows bloqueou-o.',
      solution: [
        'Feche os programas que possam ter o registo aberto e tente de novo.',
        'Ou abra a pasta de registos e apague-o à mão.',
      ],
    },
    'rcon-off': {
      title: 'O RCON está desligado',
      what: 'O app troca o mapa pelo RCON, e ele não está ligado.',
      cause: 'O RCON está desligado ou sem palavra-passe no app.',
      solution: [
        'Ligue o RCON e defina uma palavra-passe em Servidor › Palavras-passe e tokens.',
        'Reinicie o servidor para ele abrir o RCON.',
      ],
    },
    'rcon-failed': {
      title: 'Não foi possível falar com o servidor',
      what: 'O app não conseguiu ligar-se ao servidor pelo RCON.',
      cause:
        'O servidor foi iniciado antes de o RCON ser ligado, ainda está a carregar, ou outro programa está a usar a porta do RCON.',
      solution: [
        'Espere o servidor ficar pronto e tente de novo.',
        'Se o RCON foi ligado agora, reinicie o servidor.',
        'Confira a porta do RCON em Servidor › Palavras-passe e tokens.',
      ],
    },
    'rcon-auth': {
      title: 'Palavra-passe do RCON recusada',
      what: 'O servidor não aceitou a palavra-passe do RCON guardada no app.',
      cause:
        'O servidor está a usar outra palavra-passe (por exemplo a do Game.ini, ou uma trocada depois de ele iniciar).',
      solution: [
        'Reinicie o servidor pelo app para ele usar a palavra-passe guardada.',
        'Ou defina a mesma palavra-passe em Servidor › Palavras-passe e tokens.',
      ],
    },
    'map-change-failed': {
      title: 'O servidor não trocou o mapa',
      what: 'O servidor recusou a troca de mapa.',
      cause:
        'O mapa não está carregado no servidor (por exemplo, um mapa de mod cujo mod não foi carregado).',
      solution: [
        'Peça a lista ao servidor de novo e escolha um mapa dela.',
        'Para um mapa de mod, confira na página Mods se o mod está ativo.',
      ],
    },
    'discord-failed': {
      title: 'A mensagem do Discord não foi enviada',
      what: 'O Discord não aceitou a mensagem do webhook.',
      cause: 'O webhook foi apagado, o endereço está errado ou não há internet.',
      solution: [
        'No Discord, abra as definições do canal › Integrações › Webhooks e copie o endereço de novo.',
        'Cole-o em Servidor › Automação e clique em Testar.',
      ],
    },
    'firewall-failed': {
      title: 'A regra da firewall não foi criada',
      what: 'O Windows não criou a regra para o servidor.',
      cause:
        'A confirmação de administrador falhou ou outro programa de segurança bloqueou a mudança.',
      solution: [
        'Clique no botão de novo e confirme a janela do Windows.',
        'Ou crie a regra à mão em Firewall do Windows Defender › Definições avançadas › Regras de entrada.',
      ],
    },
    'autostart-failed': {
      title: 'Não foi possível mudar o arranque do Windows',
      what: 'O app não conseguiu juntar-se ou sair dos programas que abrem com o Windows.',
      cause: 'Um programa de segurança ou uma política do Windows bloqueou a mudança.',
      solution: ['Tente de novo.', 'Ou use Definições do Windows › Aplicações › Arranque.'],
    },
    'server-update-failed': {
      title: 'O servidor não foi atualizado',
      what: 'O reinício diário não conseguiu atualizar o servidor pelo SteamCMD.',
      cause:
        'A pasta do SteamCMD não está definida, outra instalação estava a correr ou o SteamCMD falhou.',
      solution: [
        'Veja a Instalação: a pasta do SteamCMD tem de estar definida.',
        'Atualize uma vez à mão em Instalação para ver o erro.',
      ],
    },
    'report-failed': {
      title: 'Não foi possível preparar o relato',
      what: 'A app não conseguiu guardar o ficheiro do relato neste PC.',
      cause: 'O disco está cheio ou a pasta da app não pode ser escrita.',
      solution: ['Liberte espaço em disco e tente de novo.'],
    },
    'report-image': {
      title: 'Esta imagem não pode ser usada',
      what: 'O ficheiro escolhido não foi adicionado ao relato.',
      cause: 'Tem mais de 5 MB, ou não é uma imagem PNG, JPG, GIF ou WebP.',
      solution: ['Escolha uma imagem menor, ou guarde a captura como PNG ou JPG.'],
    },
  },
};
