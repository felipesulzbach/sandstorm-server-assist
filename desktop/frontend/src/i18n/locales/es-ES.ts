import type { Messages } from './en-US';

export const esES: Messages = {
  app: {
    name: 'Sandstorm Server Assist',
    beta: 'BETA',
    version: 'Versión {version}',
    loading: 'Cargando…',
  },
  titlebar: {
    help: 'Ayuda',
    minimize: 'Minimizar',
    maximize: 'Maximizar',
    restore: 'Restaurar',
    close: 'Cerrar',
  },
  nav: {
    label: 'Secciones',
    dashboard: 'Panel',
    install: 'Instalación',
    server: 'Servidor',
    mods: 'Mods',
    backup: 'Copia de seguridad',
    settings: 'Ajustes',
    feedback: 'Comentarios',
    about: 'Acerca de',
  },
  common: {
    save: 'Guardar',
    saved: 'Guardado',
    cancel: 'Cancelar',
    close: 'Cerrar',
    choose: 'Elegir…',
    notChosen: 'Sin elegir',
    copy: 'Copiar',
    copied: '¡Copiado!',
    show: 'Mostrar',
    hide: 'Ocultar',
    change: 'Cambiar',
    set: 'Definir',
    remove: 'Quitar',
    add: 'Añadir',
    open: 'Abrir',
    refresh: 'Actualizar',
    optional: 'opcional',
    yes: 'Sí',
    no: 'No',
    go: 'Ir',
    howTo: '¿Cómo lo consigo?',
    unsaved: 'Cambios sin guardar',
    appliesNextStart: 'Los cambios se aplican la próxima vez que se inicie el servidor.',
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
    left: 'quedan {time}',
  },
  phase: {
    'not-installed': 'No instalado',
    stopped: 'Detenido',
    starting: 'Iniciando',
    modio: 'Conectando con mod.io',
    mods: 'Cargando mods',
    traveling: 'Cargando el mapa',
    ready: 'Listo para jugar',
    stopping: 'Deteniendo',
    crashed: 'Se detuvo inesperadamente',
  },
  phaseHint: {
    'not-installed': 'Primero instala el servidor en Instalación.',
    stopped: 'El servidor no está en marcha.',
    starting: 'El servidor se está abriendo. Tarda unos segundos.',
    modio: 'Iniciando sesión en mod.io con la cuenta del servidor.',
    mods: 'Sincronizando suscripciones y montando los mods. Los mods nuevos se descargan ahora.',
    traveling: 'Cambiando al mapa del mod. Los mapas grandes tardan hasta un minuto.',
    ready: 'Los jugadores ya pueden entrar.',
    stopping: 'Cerrando el servidor.',
    crashed: 'El servidor se cerró solo. Mira el error de abajo y los registros.',
  },
  dashboard: {
    title: 'Panel',
    subtitle: 'Inicia tu servidor y sigue lo que está haciendo.',
    start: 'Iniciar servidor',
    stop: 'Detener servidor',
    restart: 'Reiniciar',
    timeline: 'Arranque',
    step: {
      start: 'Proceso iniciado',
      modio: 'mod.io conectado',
      mods: 'Mods montados',
      map: 'Mapa cargado',
      ready: 'Listo',
    },
    map: 'Mapa',
    scenario: 'Escenario',
    port: 'Puerto',
    since: 'En marcha desde',
    modsActive: '{n} mods activos',
    external:
      'Este servidor ya estaba en marcha cuando se abrió la aplicación, y no fue ella quien lo inició. La aplicación lo sigue y lo controla con normalidad; solo no lo reinicia sola para cargar mods nuevos.',
    mutatorsRefused:
      'El servidor rechazó estos mutators en este mapa: {names}. El juego no dice por qué: un mutator puede no funcionar en este modo de juego o junto con otro que elegiste (visto con Ultralethal en Checkpoint, con o sin Bullet Sponge). Quítalo en Servidor › Mapas.',
    authorizing:
      'Autorizando en mod.io con el código. El servidor se reiniciará solo en modo normal.',
    remount:
      'Se descargaron mods nuevos y el servidor no cargó los que ya estaban instalados (mapa del mod y mutators). Reiniciando una vez para cargarlos todos.',
    gameStatsMissing:
      'Sin token de GameStats: el servidor funciona con normalidad, pero los jugadores no ganan estadísticas/XP. Puedes añadir uno en Servidor › Contraseñas y tokens.',
    cmdLine: 'Línea de comandos (secretos ocultos)',
    mapChange: {
      title: 'Cambiar mapa',
      mapLabel: 'Mapa',
      modeLabel: 'Modo',
      text: 'Cambia el mapa con el servidor encendido. Los jugadores siguen conectados y van al nuevo mapa. Los mapas de mods vienen de los mods cargados en el servidor.',
      needsRcon:
        'Para cambiar el mapa con el servidor encendido, activa RCON (con contraseña) en Servidor › Contraseñas y tokens y reinicia el servidor.',
      openServer: 'Abrir ajustes del servidor',
      load: 'Ver mapas disponibles',
      reload: 'Pedir la lista al servidor otra vez',
      empty: 'El servidor no listó ningún mapa. Inténtalo de nuevo en unos segundos.',
      fromConfig: 'mapa de mod',
      fromConfigHint:
        'Los mapas de mods vienen de los mods cargados en esta partida. Un mod al que te acabas de suscribir aparece después de reiniciar el servidor.',
      notLoaded: 'mod no cargado',
      notLoadedHint:
        'Los mapas marcados con "mod no cargado" están en la rotación de mapas, pero su mod no está cargado. Quítalos de la rotación en Servidor › Mapas o, si te acabas de suscribir al mod, reinicia el servidor.',
      current: 'ahora',
      change: 'Cambiar mapa',
      confirm: '¿Cambiar a {map}? La partida actual termina y todos van al nuevo mapa.',
      confirmYes: 'Cambiar ahora',
      sent: 'Cambiando a {map}… El estado de arriba muestra cuándo está listo.',
      waitReady: 'Disponible cuando el servidor esté listo.',
      restartNote:
        'Cuando el servidor se reinicia, vuelve al mapa configurado en Servidor › Mapas.',
    },
    connect: {
      title: 'Cómo entran los jugadores',
      text: 'Los jugadores entran por internet (con los puertos del router abiertos) o por un programa de VPN que compartes con ellos.',
      console: 'En el juego, abre la consola (tecla ~) y escribe:',
      noIp: 'Indica la IP que usan los jugadores en Servidor › Red.',
    },
    checklist: {
      title: 'Primeros pasos',
      allDone: '¡Todo listo! Ya puedes iniciar el servidor.',
      serverDir: 'Elegir la carpeta del servidor',
      steamcmd: 'Obtener SteamCMD',
      serverInstalled: 'Instalar el servidor dedicado',
      hostname: 'Poner nombre al servidor',
      modioAuth: 'Autorizar el servidor en mod.io',
      admins: 'Añadir administradores',
      connectIp: 'Indicar la IP que usan los jugadores',
      secretsSafe: 'Sacar las contraseñas de los archivos de configuración',
      gslt: 'Añadir un token de servidor de Steam (GSLT)',
    },
  },
  install: {
    title: 'Instalación',
    subtitle: 'Descarga SteamCMD e instala o actualiza el servidor dedicado.',
    folders: 'Carpetas',
    serverDir: 'Carpeta del servidor',
    serverDirHint:
      'Necesita unos 25 GB libres, más los mods (los mods de mapa pueden ocupar 5 GB cada uno). Evita las carpetas del sistema.',
    steamcmdDir: 'Carpeta de SteamCMD',
    steamcmdDirHint:
      'Si está vacía, la aplicación usa una carpeta "steamcmd" junto a la carpeta del servidor.',
    steamcmd: {
      title: 'SteamCMD',
      text: 'Herramienta de línea de comandos de Valve que descarga el servidor. La aplicación comprueba la firma digital de Valve antes de ejecutarla.',
      fetch: 'Descargar SteamCMD',
      check: 'Comprobar SteamCMD',
      selfUpdates:
        'SteamCMD se actualiza solo cada vez que la aplicación lo usa. Este botón solo comprueba el archivo y lo vuelve a descargar si está dañado.',
      ready: 'SteamCMD listo',
      missing: 'SteamCMD no encontrado',
    },
    server: {
      title: 'Servidor dedicado',
      installed: 'Instalado',
      notInstalled: 'No instalado',
      build: 'Build {id}',
      size: '{size} en disco',
      updated: 'Actualizado el {date}',
      install: 'Instalar servidor',
      update: 'Buscar actualizaciones',
      validate: 'Verificar archivos',
      validateHint:
        'Comprueba todos los archivos y repara lo que esté dañado. Es más lento; úsalo si el servidor deja de abrirse.',
      runningWarn: 'Detén el servidor antes de instalar o actualizar.',
    },
    progress: {
      title: 'Progreso',
      'self-update': 'Actualizando SteamCMD',
      login: 'Conectando con Steam',
      preallocating: 'Reservando espacio en disco',
      downloading: 'Descargando',
      verifying: 'Verificando archivos',
      committing: 'Escribiendo archivos',
      reconfiguring: 'Configurando',
      done: 'Hecho',
      failed: 'Falló',
      waiting: 'Empezando…',
      cancel: 'Cancelar',
      finished: {
        steamcmd: 'SteamCMD descargado y verificado.',
        steamcmdPresent: 'SteamCMD ya está aquí y la firma de Valve es válida. Nada que descargar.',
        server: 'El servidor está instalado y actualizado.',
        upToDate: 'El servidor ya tiene la versión más reciente. Nada que descargar.',
        validated: 'Se comprobaron todos los archivos del servidor; lo dañado se reparó.',
      },
    },
    import: {
      title: '¿Ya tienes un servidor en esta carpeta?',
      text: 'Trae los datos del servidor ya configurado aquí, en marcha o parado: nombre, puertos, mapas, mods, administradores, rotación de mapas y contraseñas. La app lee el último registro del servidor y los archivos de configuración.',
      button: 'Traer datos del servidor',
    },
  },
  server: {
    travelNoMods:
      'No se encontró ningún mapa en los mods instalados. Suscríbete a un mod de mapa e inicia el servidor una vez.',
    travelNotInstalled: '{name} (no está en los mods instalados)',
    mutatorsHint: 'La app encuentra los mutators de los mods instalados y los marca como MOD.',
    mutatorsMods: 'De tus mods',
    mutatorsGame: 'Del juego',
    mutatorsBadge: 'MOD',
    mutatorsRefused:
      'En el mapa actual el servidor rechazó: {names}. El juego no dice por qué: puede no funcionar en este modo de juego o con otro mutator que elegiste (visto con Ultralethal en Checkpoint, con o sin Bullet Sponge).',
    mutatorsRefusedBadge: 'rechazado por el servidor',
    mutatorsFromMod: 'Mod: {name}',
    mutatorsMissing: 'Ningún mod instalado tiene este mutator: el servidor no lo cargará.',
    mutatorsNoMods:
      'No hay mutators en los mods instalados. Suscríbete a un mod de mutator en mod.io e inicia el servidor una vez para descargarlo.',
    modList: 'Enviar mi lista de mods al servidor (-CmdModList)',
    modListHint:
      'Avanzado. Normalmente el servidor carga los mods suscritos por su cuenta de mod.io. Actívalo solo si quieres definir la lista aquí; las versiones recientes del juego pueden ignorarla.',
    modIds: 'Ids de los mods',
    modIdsHint: 'Los números de los mods en mod.io, separados por comas.',
    title: 'Servidor',
    subtitle:
      'Nombre, mapas, administradores y contraseñas. Se usan la próxima vez que se inicie el servidor.',
    tabs: {
      rules: 'Reglas del juego',
      automation: 'Automatización',
      general: 'General',
      maps: 'Mapas',
      admins: 'Administradores',
      secrets: 'Contraseñas y tokens',
      network: 'Red',
      messages: 'Mensajes',
    },
    reset: 'Restaurar valores',
    hostname: 'Nombre del servidor',
    hostnameHint: 'Aparece en la lista de servidores. Hasta 64 caracteres, sin comillas.',
    maxPlayers: 'Máximo de jugadores',
    port: 'Puerto del juego',
    queryPort: 'Puerto de consulta',
    portsHint:
      'Los jugadores se conectan al puerto del juego. El puerto de consulta responde a la lista de servidores.',
    useMods: 'Usar mods de mod.io',
    useModsHint: 'El servidor descarga los mods a los que está suscrita su cuenta de mod.io.',
    startMap: 'Mapa inicial',
    startMapHint:
      'Debe ser un mapa del juego base. Los mapas de mods van en "Mapa del mod", abajo.',
    mode: 'Escenario',
    customScenario: 'Nombre del escenario',
    travel: 'Cambiar a un mapa de mod tras cargar los mods',
    travelHint:
      'Usa -ModDownloadTravelTo. Ejemplo: mapa Containment, escenario Scenario_Containment_Checkpoint_Security.',
    travelMap: 'Nombre del mapa del mod',
    travelScenario: 'Escenario del mapa del mod',
    mutators: 'Mutators',
    mapCycle: 'Rotación de mapas (MapCycle.txt)',
    mapCycleHint:
      'Una entrada por línea. Déjalo vacío para mantener el archivo que ya tiene el servidor.',
    admins: {
      title: 'Administradores',
      text: 'Los administradores pueden expulsar, banear y cambiar de mapa en el juego. Se identifican por su SteamID64.',
      steamId: 'SteamID64 o enlace del perfil de Steam',
      name: 'Nombre (opcional)',
      empty: 'Aún no hay administradores.',
      found: 'SteamID64 encontrado: {id}',
    },
    secrets: {
      vaultNote:
        'Guardado y cifrado por Windows (DPAPI). Solo tu usuario de Windows, en este PC, puede abrirlo. Los valores nunca van a archivos de configuración.',
      set: 'Guardado en la caja fuerte',
      notSet: 'Sin definir',
      serverPassword: 'Contraseña de entrada',
      serverPasswordHint: 'Déjala vacía para un servidor al que pueda entrar cualquiera de la red.',
      rcon: 'Activar RCON (administración remota)',
      rconPassword: 'Contraseña RCON',
      rconPort: 'Puerto RCON',
      gsltToken: 'Token de servidor de juego de Steam (GSLT)',
      gsltTokenHint:
        'Opcional en una red privada (VPN). Necesario para aparecer en la lista pública de servidores.',
      gameStatsToken: 'Token de GameStats',
      gameStatsTokenHint:
        'Opcional. Permite a los jugadores ganar estadísticas y XP en tu servidor.',
      rconInIni: 'La contraseña RCON está guardada en texto plano en Game.ini.',
      moveToVault: 'Mover a la caja fuerte',
      moved: 'Hecho: la contraseña ahora solo está en la caja fuerte.',
    },
    network: {
      bindIp: 'Tarjeta de red del servidor',
      bindIpHint:
        'Deja "Todas" en la mayoría de los casos. Elige una tarjeta solo si el servidor debe responder en una sola red, como la de un programa de VPN (-MultiHome).',
      bindIpMissing:
        'Esta dirección no está en este equipo ahora. Si es de un programa de VPN, conecta la VPN antes de iniciar el servidor.',
      bindIpAll: 'Todas las tarjetas de red',
      connectIp: 'IP que usan los jugadores',
      connectIpHint:
        'Por internet, es tu IP pública. Con un programa de VPN, es la IP que el programa muestra para este PC.',
      firewall:
        'Si tus amigos no pueden entrar, permite InsurgencyServer.exe en el Firewall de Windows, en el tipo de red que usas (Privada o Pública).',
      ports:
        'Puertos que abrir en el router (UDP): {port} (juego) y {queryPort} (consulta). No abras el puerto de RCON ({rconPort}).',
    },
    messages: {
      announcements: 'Mensajes programados en el chat',
      announcementsHint:
        'Un mensaje por línea, hasta 10. La app envía uno cada vez, en orden, mientras el servidor está listo.',
      announcementsExample: 'Únete a nuestro Discord: discord.gg/ejemplo',
      every: 'Cada (minutos)',
      motd: 'Mensaje del día',
      motdHint:
        'Se muestra a los jugadores en la pantalla de carga al entrar. Hasta 1000 caracteres; puede tener varias líneas. Se aplica en el próximo inicio.',
      join: 'Mensaje cuando entra un jugador',
      leave: 'Mensaje cuando sale un jugador',
      chatHint:
        'Se envía al chat del juego. {name} se cambia por el nombre del jugador. Hasta 150 caracteres. Vacío = desactivado.',
      joinExample: '¡Bienvenido, {name}! Buena partida.',
      leaveExample: '{name} salió del servidor.',
      howChat:
        'El juego no tiene ajuste para estos mensajes: la app lee en el registro del servidor quién entró o salió y envía el mensaje por RCON, mientras la app esté abierta.',
      needsRcon:
        'Los mensajes de entrada y salida se envían por RCON. Activa RCON (con contraseña) en Contraseñas y tokens para usarlos.',
    },
    portCheck: {
      button: 'Diagnosticar puertos',
      running: 'Comprobando… (hasta 15 s)',
      hint: 'Comprueba desde este PC: el servidor, el Firewall de Windows, el router (UPnP), tu IP pública y Steam. Solo una prueba desde fuera de tu red demuestra que los puertos están abiertos: mira los pasos de abajo.',
      gamePort:
        'Puerto de juego UDP {port}: redirígelo en el router junto con el puerto de consulta.',
      publicIp: 'Tu IP pública: {ip}',
      useIp: 'Usar como "IP que usan los jugadores"',
      verdict: {
        open: 'Los puertos parecen abiertos: el servidor respondió a través de tu IP pública.',
        likely:
          'Todo lo que este PC puede comprobar parece correcto. Confírmalo con una prueba desde fuera de tu red.',
        problem: 'Se encontró un problema. Corrige los elementos en rojo y vuelve a comprobar.',
        unsure: 'No se pudo confirmar desde este PC. Haz la prueba desde fuera de tu red (abajo).',
      },
      items: {
        serverOff:
          'El servidor no está en marcha. Inícialo para comprobar el puerto de consulta, el camino por la IP pública y Steam.',
        localOk: 'El servidor responde en el puerto de consulta UDP {port} en este PC.',
        localNo:
          'El servidor no responde en el puerto de consulta UDP {port}, ni siquiera en este PC. Espera a que termine de cargar o revisa el puerto de consulta.',
        firewallOk: 'El Firewall de Windows tiene las reglas del servidor.',
        firewallMissing:
          'No se encontraron las reglas del Firewall de Windows que crea la app. Créalas en Servidor › Red (ignóralo si hiciste tu propia regla).',
        routerUnknown:
          'El router no respondió por UPnP, así que la app no puede leer sus redirecciones. Revísalas en la página del router.',
        routerMapped: 'Router: el puerto UDP {port} está redirigido a este PC.',
        routerNoMap:
          'Router: no hay redirección UPnP para el puerto UDP {port}. Las redirecciones hechas a mano en la página del router a menudo no aparecen aquí; comprueba que existen.',
        routerOtherPc:
          'Router: el puerto UDP {port} está redirigido a otro dispositivo ({client}), no a este PC.',
        routerWrongPort:
          'Router: el puerto UDP {port} está redirigido al puerto interno {internal}. Los dos números deben ser iguales.',
        routerDisabled:
          'Router: la redirección del puerto UDP {port} existe, pero está desactivada.',
        rconForwarded:
          'Router: el puerto RCON TCP {port} está abierto a internet. Cualquiera puede probar tu contraseña RCON. Quita esta redirección.',
        cgnat:
          'La dirección de tu router en internet ({ip}) no es una IP pública: tu proveedor usa CGNAT. La redirección de puertos no funcionará; pide una IP pública a tu proveedor.',
        ipMismatch:
          'El router dice que su IP en internet es {router}, pero internet ve {public}. Tu proveedor puede usar CGNAT, o hay un programa de VPN activo.',
        publicOk: 'El router tiene tu IP pública ({ip}), así que no hay CGNAT.',
        publicIp: 'IP pública que ve internet: {ip}.',
        publicUnknown: 'No se pudo averiguar tu IP pública (sin internet, o STUN bloqueado).',
        loopbackOk:
          'El servidor respondió a través de tu IP pública {ip}:{port}: el router redirige el puerto de consulta.',
        loopbackNo:
          'Sin respuesta a través de tu IP pública desde este PC. Muchos routers no lo permiten, así que no significa que el puerto esté cerrado.',
        steamListed: 'Steam muestra tu servidor en tu IP pública.',
        steamNotListed:
          'Steam todavía no muestra tu servidor en tu IP pública. Justo después de iniciarlo puede tardar unos minutos; si sigue así, revisa el token de Steam (GSLT).',
      },
      manual: {
        title: 'Prueba desde fuera de tu red (la única prueba segura)',
        intro:
          'Un PC no puede probar sus propios puertos desde internet. Pide a un amigo que esté en otra red que haga esto con tu servidor en marcha:',
        yourIp: 'tu IP pública',
        steps: [
          'Abrir Steam e ir a Ver › Servidores de juego.',
          'Abrir la pestaña Favoritos y pulsar "Añadir un servidor" (el botón +).',
          'Escribir {ip}:{port} y pulsar "Buscar partidas en esta dirección".',
          'Si tu servidor aparece con su nombre y su mapa, los puertos están abiertos. Tu amigo ya puede entrar desde la lista.',
        ],
        tipsTitle: 'Consejos:',
        tips: [
          'No pruebes desde un PC de tu propia red: la mayoría de los routers bloquean el camino de vuelta a su propia IP pública.',
          'Los datos móviles (un móvil como punto de acceso) cuentan como otra red.',
          'Las webs que "comprueban puertos" suelen probar solo TCP. El juego usa UDP, así que pueden decir "cerrado" aunque esté abierto.',
          'Los jugadores necesitan redirigidos el puerto de juego y el de consulta (los dos UDP). Nunca redirijas el puerto RCON.',
        ],
      },
    },
  },
  sync: {
    button: 'Traer del servidor',
    buttonHint:
      'Rellenar los campos con los datos del servidor ya configurado en la carpeta elegida',
    needDir: 'Elige primero la carpeta del servidor en Instalación',
    title: 'Datos del servidor existente',
    source: 'Leído de {file}',
    sourceFiles: 'Leído de los archivos de configuración',
    running: 'servidor en marcha ahora',
    stopped: 'servidor parado',
    conflictsTitle: 'Cambiaste {n} campo(s) y el servidor tiene otro valor',
    conflictsText:
      'Elige qué valor se queda: el que escribiste en la app o el que usa el servidor.',
    keepAllMine: 'Mantener todos los míos',
    useAllServer: 'Usar todos los del servidor',
    mine: 'Mío (app)',
    server: 'Servidor',
    fillTitle: '{n} campo(s) se rellenarán con los datos del servidor',
    fillText: 'Estos campos estaban vacíos o aún tenían el valor por defecto en la app.',
    sameCount: '{n} campo(s) ya son iguales en la app y en el servidor.',
    chooseEach: 'Elige un valor para cada campo de arriba.',
    apply: 'Aplicar',
    secretMine: 'Guardada en la app',
    secretServer: 'Distinta en el servidor',
    secretFound: 'Encontrada en el servidor (va a la caja fuerte)',
    empty: 'vacío',
    items: '{n} elementos',
    upToDate: 'La app ya tiene los mismos datos que el servidor. Nada que cambiar.',
    notFound:
      'No se encontraron datos de un servidor configurado en esta carpeta. ¿El servidor se ha iniciado al menos una vez?',
    done: 'Listo: {applied} campo(s) actualizados con los datos del servidor y {kept} mantenidos como los escribiste.',
    skipped: 'No se usaron porque el valor del servidor no es válido aquí: {fields}.',
    rconInIni:
      'La contraseña RCON sigue escrita en Game.ini. Ve a Servidor › Contraseñas y tokens y usa "Mover a la caja fuerte".',
    fields: {
      startLighting: 'Iluminación del mapa inicial',
      bindIp: 'Tarjeta de red del servidor',
      modIds: 'Lista de mods',
      rules: 'Reglas del juego',
      motd: 'Mensaje del día',
      hostname: 'Nombre del servidor',
      port: 'Puerto del juego',
      queryPort: 'Puerto de consulta',
      maxPlayers: 'Máximo de jugadores',
      startMap: 'Mapa inicial',
      travelMap: 'Mapa de mod tras cargar los mods',
      useMods: 'Usar mods',
      mutators: 'Mutators',
      rcon: 'RCON',
      rconPort: 'Puerto RCON',
      admins: 'Administradores',
      mapCycle: 'Rotación de mapas',
      serverPassword: 'Contraseña del servidor',
      rconPassword: 'Contraseña RCON',
      gsltToken: 'Token GSLT',
      gameStatsToken: 'Token GameStats',
      modioApiKey: 'Clave de API de mod.io',
    },
  },
  maps: {
    mode: {
      Checkpoint_Security: 'Checkpoint (Seguridad)',
      Checkpoint_Insurgents: 'Checkpoint (Insurgentes)',
      custom: 'Otro (escribe el nombre)',
    },
  },
  mods: {
    title: 'Mods (mod.io)',
    subtitle: 'El servidor descarga los mods a los que está suscrita su cuenta de mod.io.',
    session: {
      title: 'Autorización del servidor',
      valid: 'Autorizado como {user} hasta el {date}.',
      expired: 'La autorización caducó el {date}. Vuelve a autorizar.',
      none: 'Este PC aún no está autorizado en mod.io.',
    },
    auth: {
      title: 'Autorizar el servidor',
      text: 'Usa una cuenta de mod.io solo para el servidor (no la personal). Se hace una vez y dura alrededor de un año.',
      email: 'Correo de la cuenta del servidor',
      apiKey: 'Clave de la API de mod.io',
      detect: 'Buscar en los registros del servidor',
      detected: 'Clave encontrada en los registros del servidor y guardada en la caja fuerte.',
      notFound: 'Aún no hay clave en los registros. Pega una creada en mod.io (enlace en la guía).',
      requestCode: 'Enviar código por correo',
      codeSent: 'Código enviado. Revisa la bandeja de entrada y la carpeta de spam.',
      code: 'Código de 5 caracteres',
      authorize: 'Autorizar',
      authorizeHint:
        'El servidor se inicia una vez con el código y después se reinicia solo en modo normal (SecurityCode=none).',
    },
    health: {
      title: 'Mods con problemas',
      text: 'La app compara los archivos de tus mods instalados. Dos mods que traen el mismo archivo del juego en versiones distintas pueden hacer que el juego o el servidor se cierren con "Serial size mismatch" al cargar un mapa.',
      recheck: 'Comprobar de nuevo',
      none: 'No hay conflictos entre tus {n} mods instalados.',
      found: '{n} par(es) de mods traen los mismos archivos del juego en versiones distintas:',
      and: 'y',
      pair: '{n} archivo(s) en versiones distintas',
      example: 'p. ej. {file}',
      advice:
        'Si el juego o el servidor se cierra al cargar un mapa, deja solo un mod del par: cancela la suscripción del otro en mod.io (cuenta del servidor). Las copias idénticas de un archivo no aparecen: no causan problemas.',
      unreadable:
        'No se pudieron leer los archivos de: {names}. Los conflictos con estos mods pueden no aparecer aquí.',
    },
    crash: {
      title: '¿Se cerró el juego de un jugador (o el servidor)? Encuentra el mod',
      text: 'Pega el mensaje del error (por ejemplo el que tiene "Serial size mismatch"). La app encuentra qué mod instalado trae el archivo que se nombra. El texto no se guarda ni se envía.',
      placeholder: 'StaticMesh /Game/…: Serial size mismatch: Got …, Expected …',
      find: 'Encontrar el mod',
      noObject:
        'No se encontró ningún archivo del juego en este mensaje. Pega la línea que nombra el archivo (empieza por /Game/ o por el nombre del mod).',
      inMods: 'está en {names}',
      inNoMod: 'ningún mod instalado tiene este archivo',
      adviceMismatch:
        '"Serial size mismatch" significa que este mod tiene una copia del archivo hecha para otra versión del juego. Cancela la suscripción del mod en mod.io (cuenta del servidor) hasta que su autor lo actualice. Si crees que la descarga está dañada, usa el botón de borrar en la fila del mod para descargarlo de nuevo.',
      advice:
        'Prueba a quitar este mod: cancela la suscripción en mod.io (cuenta del servidor), o bórralo en su fila para descargarlo de nuevo.',
      adviceNone:
        'El archivo no está en los mods del servidor. Puede venir del propio juego o de un mod que el jugador tiene y el servidor no: el jugador puede verificar los archivos del juego en Steam (Propiedades › Archivos instalados) y borrar la carpeta del mod en C:\\Users\\Public\\mod.io\\254\\mods para descargarlo de nuevo.',
    },
    list: {
      title: 'Mods',
      refresh: 'Leer suscripciones de mod.io',
      refreshed: 'Actualizado el {date}',
      empty: 'Aún no hay mods. Suscríbete a mods en mod.io con la cuenta del servidor.',
      name: 'Mod',
      size: 'Tamaño',
      status: 'Estado',
      openSite: 'Abrir mod.io',
      conflict:
        'Los mods de mapa que traen los mismos archivos del juego entran en conflicto (ejemplo: Containment con Old World o Containment Storm). Si el servidor se cierra con "Serial size mismatch", cancela la suscripción de uno de ellos.',
    },
    pending: {
      text: 'Aún no descargados: {names}. El servidor descarga los mods suscritos al arrancar.',
      stoppedHint:
        '"Descargar ahora" arranca el servidor; tras la descarga se reinicia una vez solo para cargar todos los mods.',
      runningHint:
        'Un servidor encendido no carga mods nuevos: "Reiniciar y descargar ahora" lo reinicia.',
      start: 'Descargar ahora',
      restart: 'Reiniciar y descargar ahora',
      confirm: 'La partida actual termina y los jugadores salen un momento. ¿Reiniciar ahora?',
      confirmYes: 'Reiniciar ahora',
      started:
        'El servidor está arrancando y descarga los mods nuevos. Sigue el progreso en la lista; se reinicia una vez solo para cargar todos los mods.',
    },
    delete: {
      column: 'Acciones',
      button: 'Borrar los archivos descargados de {name}',
      confirm:
        '¿Borrar los archivos descargados de {name}? El servidor vuelve a descargar el mod al arrancar, si sigues suscrito. Úsalo para arreglar un mod dañado o forzar una descarga nueva.',
      redownload: 'Borrar y descargar de nuevo',
      only: 'Solo borrar',
      working: 'Borrando…',
      runningHint:
        '"Borrar y descargar de nuevo" detiene el servidor, borra el mod y arranca el servidor, que lo vuelve a descargar. Para solo borrar, detén antes el servidor.',
      stoppedHint:
        '"Borrar y descargar de nuevo" arranca el servidor, que descarga el mod. "Solo borrar" deja la descarga para el próximo arranque.',
      doneRedownload:
        '{name} borrado. El servidor está arrancando y lo vuelve a descargar; sigue el progreso aquí.',
      done: '{name} borrado. El servidor lo vuelve a descargar en el próximo arranque, si sigues suscrito.',
    },
    state: {
      installed: 'Instalado',
      subscribed: 'Se descargará',
      orphan: 'Sin suscripción (se eliminará)',
      active: 'Activo',
      installing: 'Descargando',
      updating: 'Actualizando',
      uninstalling: 'Eliminando',
      removed: 'Eliminado',
      pending: 'En espera',
    },
  },
  backup: {
    title: 'Copia de seguridad',
    subtitle:
      'Un único archivo cifrado con todo lo necesario para montar este servidor en otro PC.',
    create: {
      title: 'Crear una copia',
      contents:
        'Incluye: ajustes del servidor, contraseñas y tokens, Game.ini, Engine.ini, Admins.txt, MapCycle.txt, Bans.json y Motd.txt.',
      password: 'Contraseña de la copia',
      confirm: 'Repite la contraseña',
      passwordHint:
        'Al menos 10 caracteres. Sin ella la copia no se puede abrir: no hay recuperación.',
      mismatch: 'Las contraseñas no coinciden.',
      includeModio: 'Incluir la sesión de mod.io (el nuevo PC queda autorizado)',
      includeModioWarn:
        'La sesión funciona como una contraseña. Guarda el archivo en un lugar seguro.',
      button: 'Guardar copia…',
      done: 'Copia guardada: {file}',
    },
    restore: {
      title: 'Restaurar en este PC',
      choose: 'Elegir archivo de copia…',
      file: 'Archivo: {file}',
      password: 'Contraseña de la copia',
      open: 'Abrir copia',
      createdAt: 'Creada',
      machine: 'PC',
      hostname: 'Servidor',
      files: 'Archivos',
      secrets: 'Contraseñas y tokens',
      modio: 'Sesión de mod.io',
      included: 'incluida',
      notIncluded: 'no incluida',
      serverDir: 'Carpeta del servidor en este PC',
      restoreModio: 'Restaurar la sesión de mod.io',
      button: 'Restaurar',
      done: 'Restaurado. Siguiente: instala el servidor en Instalación (si este PC no lo tiene) e inícialo.',
      runningWarn: 'Detén el servidor antes de restaurar.',
    },
  },
  settings: {
    startWithWindows: 'Abrir la app al iniciar Windows',
    startWithWindowsHint: 'La app se abre minimizada al entrar en Windows.',
    autoStartServer: 'Iniciar el servidor cuando la app se abre con Windows',
    autoStartServerHint: 'Útil en un equipo que queda encendido como servidor.',
    title: 'Ajustes',
    subtitle: 'Idioma, apariencia y preferencias de la aplicación.',
    language: 'Idioma',
    appearance: 'Apariencia',
    color: 'Color',
    skin: 'Diseño',
    savedHere: 'Guardado en este equipo.',
    themes: {
      sand: 'Arena',
      black: 'Negro',
      olive: 'Oliva',
      blue: 'Azul oscuro',
      light: 'Claro',
    },
    skins: {
      sidebar: 'Barra lateral',
      sidebarDesc: 'Menú con nombres a la izquierda.',
      compact: 'Compacto',
      compactDesc: 'Solo iconos a la izquierda, más espacio para el contenido.',
      cards: 'Tarjetas',
      cardsDesc: 'Paneles flotantes con esquinas redondeadas.',
      top: 'Barra superior',
      topDesc: 'Menú en una barra horizontal arriba.',
    },
    console: 'Mostrar la ventana de consola del servidor',
    consoleHint:
      'Desactivado: la consola queda oculta y los registros aparecen en la aplicación. Se aplica en el próximo inicio.',
    logs: {
      title: 'Registros del servidor',
      text: 'El servidor guarda un archivo de registro por inicio. Ahora hay {n} inicios guardados.',
      auto: 'Borrar registros antiguos automáticamente',
      autoOff: 'Desactivado: la app no borra ningún registro.',
      autoOn:
        'Activado: cada vez que el servidor se inicia, se borran los inicios más antiguos y se mantienen los más recientes.',
      keep: 'Inicios a mantener',
      keepHint:
        'De {min} a {max}. El mínimo es el inicio actual y el anterior: la app los lee para "Traer del servidor" y para encontrar la clave de mod.io.',
      cleanNow: 'Borrar registros antiguos ahora',
      confirm: '¿Borrar {n} inicios antiguos ({size})? Se mantienen los {keep} más recientes.',
      confirmYes: 'Borrar {n}',
      wouldDelete: 'Manteniendo {keep} inicios, se borrarían {n} inicios antiguos ({size}).',
      nothing: 'Nada que borrar: hay {keep} inicios o menos.',
      done: '{n} archivos borrados ({size} liberados).',
    },
  },
  about: {
    title: 'Acerca de',
    text: 'Sandstorm Server Assist te ayuda a instalar, configurar y ejecutar un servidor dedicado de Insurgency: Sandstorm con mods de mod.io.',
    license: 'Licencia',
    licenseText:
      'Apache License 2.0. Libre para usar, modificar y compartir, se ofrece "tal cual", sin garantía.',
    disclaimer: 'Aviso legal',
    disclaimerText:
      'Proyecto no oficial, sin relación ni respaldo de New World Interactive, Focus Entertainment, Valve ni mod.io. Insurgency, Steam y mod.io son marcas de sus propietarios. Eres responsable del uso que hagas de tu servidor y de cumplir los términos del juego, de Steam y de mod.io y las leyes de tu país. La aplicación no distribuye el juego ni SteamCMD: los descarga de las fuentes oficiales.',
    links: 'Enlaces útiles',
    serverGuide: 'Guía de servidores de la comunidad',
    modioGuide: 'Guía de administración de servidores de mod.io',
    credits: 'Créditos',
    creditsText:
      'Creado por Felipe Sulzbach, con la ayuda de Claude Code (Anthropic). Gratuito y de código abierto, hecho para la comunidad de Insurgency: Sandstorm.',
  },
  feedback: {
    title: 'Comentarios',
    subtitle: 'Informar de un problema o sugerir una mejora',
    kindLabel: 'Tipo',
    kind: {
      bug: 'Informar de un problema',
      idea: 'Sugerir una mejora',
    },
    intro: {
      bug: 'Cuenta qué salió mal. La app prepara un archivo con los detalles y abre el formulario de GitHub ya rellenado.',
      idea: 'Comparte una idea para mejorar la app. La app abre el formulario de GitHub ya rellenado.',
    },
    titleField: 'Título (opcional)',
    titleHint: 'Una línea. Si lo dejas vacío, la app escribe uno.',
    description: 'Descripción (opcional)',
    descriptionHint: {
      bug: 'Qué hiciste, qué esperabas y qué pasó.',
      idea: 'Qué te gustaría que hiciera la app y por qué ayudaría.',
    },
    includeLog: 'Incluir el registro de la app',
    includeLogHint:
      'Las últimas líneas del panel Registros. Las contraseñas, tokens, direcciones IP y SteamIDs se ocultan.',
    image: 'Captura de pantalla (opcional)',
    imageHint: 'PNG, JPG, GIF o WebP, hasta 5 MB.',
    chooseImage: 'Elegir imagen',
    removeImage: 'Quitar',
    create: 'Preparar y abrir GitHub',
    privacy:
      'La app no envía nada. Guarda un archivo .zip en este PC y abre el formulario en tu navegador; tú decides qué publicar. Necesitas una cuenta gratuita de GitHub.',
    done: 'Guardado {file} ({size}).',
    next: [
      'Inicia sesión en GitHub en la ventana del navegador que se abrió.',
      'Arrastra el archivo .zip (mostrado en la ventana del Explorador) a la última casilla del formulario.',
      'Revisa el texto y pulsa "Create".',
    ],
    showFile: 'Mostrar el archivo',
  },
  update: {
    title: 'Actualizaciones',
    check: 'Buscar actualizaciones',
    install: 'Actualizar',
    installFailed: 'La actualización no empezó',
    readyHint:
      'La versión {version} ya está descargada. Haz clic para reiniciar la app en ella (o se instala al cerrar la app).',
    manualHint:
      'La versión {version} está disponible. Esta es la versión portátil: descárgala en la página de descarga.',
    manualButton: 'Versión {version}',
    errorCheck: 'No se pudo buscar actualizaciones ahora. La app lo intenta más tarde.',
    errorDownload: 'No se pudo descargar la actualización. La app lo intenta más tarde.',
    hint: 'La app busca versiones nuevas sola y solo instala actualizaciones firmadas por el proyecto. El servidor del juego sigue funcionando durante la actualización.',
    status: {
      idle: 'La app busca versiones nuevas unos segundos después de abrirse.',
      checking: 'Buscando una versión nueva…',
      uptodate: 'Tienes la versión más reciente.',
      downloading: 'Descargando la versión {version}… {percent}%',
      ready:
        'La versión {version} está lista. Haz clic en Actualizar (o se instala al cerrar la app).',
      manual:
        'La versión {version} está disponible. Esta es la versión portátil: descárgala en la página de descarga.',
      dev: 'Versión de desarrollo o portátil sin actualizaciones: la actualización automática está desactivada.',
    },
  },
  legal: {
    title: 'Términos de uso y privacidad',
    intro: 'Antes de usar la app, lee y acepta los términos de uso y el aviso de privacidad.',
    tabs: {
      terms: 'Términos de uso',
      privacy: 'Privacidad y datos',
    },
    agree: 'He leído y acepto los términos de uso y el aviso de privacidad.',
    accept: 'Aceptar y continuar',
    exit: 'Salir',
    version: 'Versión de los términos: {version}',
    acceptedOn: 'Versión {version} aceptada el {date}.',
    terms: {
      about: {
        title: '1. Qué es esta app',
        text: [
          'Sandstorm Server Assist es una herramienta gratuita, de código abierto y no oficial para instalar, configurar e iniciar un servidor dedicado de Insurgency: Sandstorm.',
          'No está afiliada, aprobada ni patrocinada por New World Interactive, Focus Entertainment, Valve ni mod.io. Sus nombres y marcas solo se usan para decir con qué funciona la app.',
        ],
      },
      asIs: {
        title: '2. Sin garantía',
        text: [
          'La app se ofrece "tal cual", sin garantía de ningún tipo, bajo la Licencia Apache 2.0.',
          'Los autores y colaboradores no responden por daños, pérdida de datos, acciones sobre cuentas, baneos, costes ni consecuencias legales derivadas de su uso.',
        ],
      },
      responsibility: {
        title: '3. Tu responsabilidad',
        text: [
          'Eres responsable de cómo usas la app y de los servidores que inicias con ella.',
          'Esto incluye cumplir la licencia y los términos del juego, el Acuerdo de Suscriptor de Steam y las reglas de servidores de Steam, los términos de mod.io y las leyes de tu país y de tus jugadores.',
          'También eres responsable del contenido, los mods, los nombres y los mensajes de tu servidor, y de a quién dejas entrar y administrar.',
        ],
      },
      prohibited: {
        title: '4. Lo que no está permitido',
        text: [
          'No uses la app con fines ilegales, para dañar a otros, para saltarte protecciones del juego, de Steam o de mod.io, ni para compartir contenido que no tienes derecho a compartir.',
        ],
      },
      thirdParty: {
        title: '5. Servicios de terceros',
        text: [
          'La app no distribuye el servidor del juego, SteamCMD ni mods. Los descarga de las fuentes oficiales (Valve y mod.io), con sus propios términos.',
          'Los tokens y las cuentas de Steam, GameStats y mod.io son tuyos y siguen las reglas de cada servicio.',
        ],
      },
      updates: {
        title: '6. Actualizaciones automáticas',
        text: [
          'La app instalada busca versiones nuevas en las publicaciones del proyecto en GitHub, las descarga en segundo plano y las instala cuando haces clic en Actualizar o cierras la app.',
          'Solo se aceptan actualizaciones firmadas con la clave del proyecto. La versión portátil solo avisa de que hay una versión nueva.',
        ],
      },
      changes: {
        title: '7. Cambios en estos términos',
        text: [
          'Si estos términos o el aviso de privacidad cambian, la app los muestra de nuevo y pide tu aceptación antes de continuar.',
        ],
      },
    },
    privacy: {
      summary: {
        title: '1. En resumen',
        text: [
          'La app no tiene cuenta, anuncios ni telemetría. No envía tus datos a los autores del proyecto.',
          'Todo lo que escribes se queda en este PC, salvo lo que va a los servicios que eliges usar (abajo).',
        ],
      },
      stored: {
        title: '2. Qué se queda en este PC',
        text: [
          'Configuración del servidor y preferencias de la app: en la carpeta %APPDATA%\\sandstorm-server-assist.',
          'Contraseñas, tokens, el correo y la clave de API de mod.io: en un almacén protegido por Windows (DPAPI), que solo tu usuario de Windows en este PC puede abrir.',
          'Copias de seguridad: solo donde las guardes, cifradas con la contraseña que elijas.',
          'Registros del servidor: se leen de la carpeta del servidor para mostrarlos en la app. Contraseñas y tokens se ocultan antes de mostrarse.',
        ],
      },
      sent: {
        title: '3. Qué sale de este PC, y a quién',
        text: [
          'Valve (Steam): descarga de SteamCMD y del servidor del juego.',
          'mod.io: el correo de la cuenta del servidor (para recibir el código de acceso), el propio código y la clave de API del servidor, para autorizar y descargar mods.',
          'GitHub (página del proyecto): la búsqueda de actualizaciones pide el archivo de la última versión y descarga el instalador desde allí. Como cualquier web, GitHub ve tu dirección IP y la hora de la petición. No se envía ningún identificador tuyo.',
          'Diagnóstico de puertos (solo cuando pulsas "Diagnosticar puertos"): servidores STUN públicos de Cloudflare y Google responden con tu dirección IP pública; se pregunta a Steam (Valve) si muestra un servidor en esa dirección; se pregunta a tu router por sus redirecciones (UPnP, dentro de tu red). No se envía nada más.',
          'Discord (solo si defines un webhook en Servidor › Automatización): los textos de los avisos que escribiste, con el nombre del servidor, el mapa y los nombres de los jugadores que entran o salen, van a tu canal de Discord. Discord recibe tu IP, como cualquier sitio web.',
          'Comentarios (informar de un problema o sugerir una mejora): la app no envía nada. Guarda un archivo en este PC y abre el formulario de GitHub en tu navegador; tú eliges qué publicar.',
          'Los enlaces en los que haces clic se abren en tu navegador, con las reglas de cada sitio.',
        ],
      },
      server: {
        title: '4. Tu servidor del juego',
        text: [
          'El servidor del juego que inicias se conecta a Steam, GameStats, mod.io y los jugadores, con las reglas de esos servicios. Los SteamID de administradores que añadas se guardan en los archivos del servidor (Admins.txt).',
          'Si entran jugadores en tu servidor, tú eres responsable de sus datos en él (nombres, SteamID, chat, registros).',
        ],
      },
      control: {
        title: '5. Tu control',
        text: [
          'Puedes ver y cambiar todo en la app. Para borrarlo todo: desinstala la app y borra la carpeta %APPDATA%\\sandstorm-server-assist.',
          'Puedes revocar tokens y claves en cualquier momento en las webs de Steam y mod.io.',
        ],
      },
      changes: {
        title: '6. Cambios',
        text: ['Si este aviso cambia, la app lo muestra de nuevo y pide tu aceptación.'],
      },
    },
  },
  logs: {
    button: 'Registros',
    errors: '{n} errores',
    oneError: '1 error',
    title: 'Registros',
    filter: {
      all: 'Todo',
      important: 'Importantes',
      warning: 'Avisos',
      error: 'Errores',
    },
    search: 'Buscar…',
    follow: 'Seguir',
    copy: 'Copiar visibles',
    openFolder: 'Abrir carpeta de registros',
    clear: 'Limpiar',
    empty: 'Aún no hay líneas de registro.',
    source: {
      server: 'servidor',
      steamcmd: 'steamcmd',
      app: 'app',
    },
    run: 'Inicio del servidor',
    runLive: 'Inicio actual (desde {time})',
    runSession: 'Inicio actual',
    runAt: 'Inicio de {time} ({size})',
    runLast: 'Último inicio, {time} ({size})',
    runStarted: 'Servidor iniciado a las {time}',
    savedRun: 'Inicio guardado, leído del archivo de registro del servidor. {n} errores.',
    savedRunOneError: 'Inicio guardado, leído del archivo de registro del servidor. 1 error.',
    backToLive: 'Volver al inicio actual',
    loading: 'Cargando…',
    masked: 'Las contraseñas y los tokens aparecen como •••.',
  },
  error: {
    modFile: 'Archivo nombrado en el error: {object}. Ningún mod instalado lo tiene.',
    modOwners: 'Archivo nombrado en el error: {object}. Está en: {mods}.',
    screenTitle: 'Esta pantalla tuvo un error',
    screenText:
      'El resto de la app sigue funcionando. Vuelve a intentarlo o regresa al Panel. El detalle queda en Registros.',
    appTitle: 'La app tuvo un error',
    appText:
      'Recarga la app para seguir usándola. El servidor sigue en marcha. El detalle queda en Registros.',
    reload: 'Recargar la app',
    retry: 'Reintentar',
    home: 'Volver al Panel',
    what: 'Qué ha pasado',
    cause: 'Causa probable',
    solution: 'Cómo solucionarlo',
    details: 'Detalles técnicos',
    openLogs: 'Abrir registros',
    help: 'Ayuda',
  },
  links: {
    'steam-gslt': 'Steam: cuentas de servidor de juego',
    gamestats: 'Sandstorm GameStats',
    'steamid-io': 'steamid.io',
    'steam-profile': 'Tu perfil de Steam',
    modio: 'mod.io: Insurgency: Sandstorm',
    'modio-subscriptions': 'mod.io: suscripciones',
    'modio-api-access': 'mod.io: acceso a la API',
    'modio-status': 'Estado de mod.io',
    'modio-admin-guide': 'Guía de administración de mod.io',
    'steamcmd-docs': 'Documentación de SteamCMD',
    'server-guide': 'Guía de servidores de la comunidad',
    'apache-license': 'Apache License 2.0',
    wails: 'Wails',
    'download-page': 'Versiones en GitHub',
    project: 'Proyecto en GitHub',
    'user-guide': 'Guía de uso',
    author: 'Felipe Sulzbach en GitHub',
  },
  guides: {
    steamId: {
      title: 'Encontrar un SteamID64',
      steps: [
        'Abre el perfil de Steam de la persona (en Steam: haz clic en el nombre › Perfil).',
        'Si la dirección es steamcommunity.com/profiles/7656119…, el número de 17 dígitos es el SteamID64. Pega aquí el enlace completo: la aplicación lo encuentra.',
        'Si la dirección tiene un nombre personalizado (steamcommunity.com/id/nombre), pégalo en steamid.io y copia el "steamID64".',
        'En la aplicación de Steam, activa Parámetros › Interfaz › "Mostrar barras de dirección web cuando estén disponibles" para ver la dirección.',
      ],
    },
    gslt: {
      title: 'Obtener un token GSLT',
      steps: [
        'Abre la página de cuentas de servidor de juego de Steam e inicia sesión.',
        'En "App ID of base game" escribe 581320 (Insurgency: Sandstorm) y añade una nota, por ejemplo el nombre del servidor.',
        'Haz clic en Create y copia el "Login Token".',
        'Un token por servidor. Tu cuenta de Steam no puede estar limitada (necesita al menos una compra).',
      ],
    },
    gameStats: {
      title: 'Obtener un token de GameStats',
      steps: [
        'Abre la página de Sandstorm GameStats e inicia sesión con Steam.',
        'Crea un token para tu servidor y cópialo.',
        'Pégalo aquí. Sin él el servidor funciona igual; los jugadores solo no ganan estadísticas.',
      ],
    },
    modioAccount: {
      title: 'Cuenta de mod.io para el servidor',
      steps: [
        'Crea una cuenta de mod.io con un correo solo para el servidor (no uses tu cuenta personal).',
        'Con esa cuenta, abre Insurgency: Sandstorm en mod.io y suscríbete a los mods que debe usar el servidor.',
        'Escribe ese correo aquí, envía el código y autoriza. El servidor descarga los mods suscritos en cada inicio.',
      ],
    },
    modioApiKey: {
      title: 'Obtener una clave de la API de mod.io',
      steps: [
        'Prueba primero "Buscar en los registros del servidor": el servidor escribe su clave en el registro tras iniciarse una vez con mods.',
        'Si no aparece: inicia sesión en mod.io y abre "API access".',
        'Crea una clave de API de solo lectura y cópiala (32 letras y números).',
        'La clave se queda en la caja fuerte. No forma parte del código de la aplicación porque el proyecto es público.',
      ],
    },
    ports: {
      title: 'Por internet: abrir los puertos del router',
      steps: [
        'Abre la página de administración del router (normalmente 192.168.0.1 o 192.168.1.1) y busca "Redirección de puertos" (puede llamarse "Servidor virtual" o "NAT").',
        'Redirige UDP {port} y UDP {queryPort} a la IP de este PC en tu red local (ejecuta ipconfig y mira "Dirección IPv4").',
        'Permite InsurgencyServer.exe en el Firewall de Windows.',
        'Los jugadores usan tu IP pública (busca "cuál es mi IP"). Escríbela en "IP que usan los jugadores".',
        'No redirijas el puerto de RCON ({rconPort}): controla el servidor.',
        'Si aun así no funciona, tu proveedor puede compartir una IP pública entre muchos clientes (CGNAT). Pide una IP pública al proveedor o usa un programa de VPN.',
      ],
    },
    vpn: {
      title: 'Con un programa de VPN (red virtual)',
      steps: [
        'Instala el mismo programa de red virtual en el PC del servidor y en el de cada jugador.',
        'Crea una red en el PC del servidor y comparte el nombre y la contraseña con los jugadores.',
        'La IP del servidor es la que el programa muestra para este PC. Escríbela en "IP que usan los jugadores".',
        'No hace falta abrir puertos en el router.',
      ],
    },
  },
  lighting: {
    label: 'Iluminación',
    default: 'Predeterminado del escenario',
    Day: 'Día',
    Night: 'Noche',
    start: 'Iluminación del mapa inicial',
    cycleTitle: 'Día o noche en la rotación',
    allDay: 'Todo de día',
    allNight: 'Todo de noche',
    allDefault: 'Todo predeterminado',
    unknownLine: 'línea que la app no lee',
    cycleMissing:
      'Algunos mapas de la rotación no están en el juego base ni en ningún mod instalado. Cuando la rotación llegue a ellos, el servidor no tendrá adónde ir.',
    cycleRemoveMissing: 'Quitar estos mapas',
    notInstalled: 'mod no instalado',
    remove: 'Quitar',
  },
  mutators: {
    AllYouCanEat: 'Empiezas con 100 puntos de suministro.',
    AntiMaterielRiflesOnly: 'Solo fusiles antimaterial, más el equipo normal y explosivos.',
    BoltActionsOnly: 'Solo fusiles de cerrojo, más el equipo normal y explosivos.',
    Broke: 'Empiezas con 0 puntos de suministro.',
    BulletSponge: 'Más salud.',
    Competitive: 'Equipo más caro, rondas más cortas y capturas más rápidas.',
    CompetitiveLoadouts: 'Clases cambiadas por las del Competitivo.',
    FastMovement: 'Movimiento más rápido.',
    Frenzy:
      'Enemigos de la IA que solo atacan cuerpo a cuerpo. Cuidado con los enemigos especiales.',
    Guerrillas: 'Empiezas con 5 puntos de suministro.',
    Hardcore: 'Movimiento más lento y capturas más largas.',
    HeadshotOnly: 'Los jugadores solo reciben daño por disparos a la cabeza.',
    HotPotato: 'Cae una granada activa cuando un jugador muere.',
    LockedAim: 'Las armas siempre apuntan al centro de la pantalla.',
    NoAim: 'Sin apuntar con la mira.',
    PistolsOnly: 'Solo pistolas, más el equipo normal y explosivos.',
    ShotgunsOnly: 'Solo escopetas, más el equipo normal y explosivos.',
    SlowCaptureTimes: 'Los objetivos tardan más en capturarse.',
    SlowMovement: 'Movimiento más lento.',
    SoldierOfFortune: 'Ganas puntos de suministro según sube tu puntuación.',
    SpecialOperations: 'Empiezas con 30 puntos de suministro.',
    Strapped: 'Empiezas con 1 punto de suministro.',
    Ultralethal: 'Todos mueren de un disparo.',
    Vampirism: 'Recuperas salud igual al daño que haces a los enemigos.',
    Warlords: 'Empiezas con 10 puntos de suministro.',
  },
  rules: {
    manage: 'Definir las reglas del juego aquí',
    manageHint:
      'La app escribe estas opciones en Game.ini y Engine.ini cuando el servidor se inicia. Las opciones vacías usan el valor del juego. Desactivado: la app no toca estos archivos.',
    fromServer:
      'El Game.ini del servidor ya define {n} reglas del juego. La aplicación aún no las gestiona, así que deja los archivos como están.',
    useServer: 'Usar las reglas del servidor',
    default: 'Predeterminado: {value}',
    gameDefault: 'Valor del juego',
    on: 'Activado',
    off: 'Desactivado',
    groups: {
      coop: 'Coop: bots',
      checkpoint: 'Checkpoint',
      objectives: 'Objetivos',
      friendlyFire: 'Fuego amigo y equipos',
      teamKills: 'Muertes de compañeros',
      chat: 'Chat, voz y registro de bajas',
      match: 'Partida',
      supply: 'Puntos de suministro',
      votes: 'Votación para expulsar',
      network: 'Rendimiento (Engine.ini)',
    },
    items: {
      aiDifficulty: {
        label: 'Dificultad de los bots',
        hint: '0 es lo más fácil, 1 lo más difícil',
      },
      minimumEnemies: {
        label: 'Mínimo de enemigos',
        hint: 'Bots enemigos con pocos jugadores',
      },
      maximumEnemies: {
        label: 'Máximo de enemigos',
        hint: 'Bots enemigos con el servidor lleno',
      },
      friendlyBotQuota: {
        label: 'Bots aliados',
        hint: 'Bots en el equipo de los jugadores',
      },
      vehicleInsertion: {
        label: 'Los enemigos llegan en vehículo',
        hint: 'Entrada con camiones donde el mapa lo tiene',
      },
      defendTimer: {
        label: 'Tiempo del contraataque (s)',
        hint: 'Tiempo para defender cada punto',
      },
      defendTimerFinal: {
        label: 'Tiempo del contraataque final (s)',
        hint: 'Tiempo para defender el último punto',
      },
      retreatTimer: {
        label: 'Tiempo de retirada de los bots (s)',
        hint: 'Después del contraataque',
      },
      respawnDpr: {
        label: 'Proporción para que vuelvan los bots',
        hint: 'Parte de bots muertos antes de volver',
      },
      respawnDelay: {
        label: 'Espera para que vuelvan los bots (s)',
        hint: 'Antes de que vuelvan los bots',
      },
      postCaptureRushTimer: {
        label: 'Carga tras el alijo (s)',
        hint: 'Los bots avanzan tras destruir un alijo',
      },
      counterAttackRespawnDpr: {
        label: 'Proporción en el contraataque',
        hint: 'Parte de bots muertos antes de volver',
      },
      counterAttackRespawnDelay: {
        label: 'Espera en el contraataque (s)',
        hint: 'Antes de que vuelvan los bots',
      },
      enemyRespawnMin: {
        label: 'Oleadas enemigas con pocos jugadores',
        hint: 'Multiplicador',
      },
      enemyRespawnMax: {
        label: 'Oleadas enemigas con el servidor lleno',
        hint: 'Multiplicador',
      },
      finalCacheBots: {
        label: 'Bots en el alijo final',
        hint: 'Multiplicador',
      },
      captureTime: {
        label: 'Tiempo de captura (s)',
        hint: 'Tiempo para tomar un punto',
      },
      resetTime: {
        label: 'Tiempo para reiniciar la captura (s)',
        hint: '-1 mantiene el progreso',
      },
      captureSpeedup: {
        label: 'Aceleración por jugador',
        hint: 'Cada jugador extra captura más rápido',
      },
      captureSpeedupPlayers: {
        label: 'Jugadores que aceleran',
        hint: 'Máximo contado',
      },
      friendlyFire: {
        label: 'Fuego amigo',
        hint: 'Daño a los compañeros',
      },
      friendlyFireModifier: {
        label: 'Daño del fuego amigo',
        hint: '0 es nada, 1 es el daño completo',
      },
      friendlyFireReflect: {
        label: 'Daño devuelto al atacante',
        hint: '0 es nada, 1 es el daño completo',
      },
      autoBalance: {
        label: 'Equilibrar equipos',
        hint: 'Mueve jugadores cuando los equipos están desiguales',
      },
      autoBalanceDelay: {
        label: 'Espera para equilibrar (s)',
        hint: 'Después de que los equipos queden desiguales',
      },
      teamSelect: {
        label: 'Los jugadores eligen equipo',
        hint: 'Permite cambiar de equipo',
      },
      teamKillLimit: {
        label: 'Muertes de compañeros antes de expulsar',
        hint: '0 desactiva la expulsión',
      },
      teamKillGrace: {
        label: 'Tiempo entre muertes de compañeros (s)',
        hint: 'Las muertes en este intervalo cuentan una vez',
      },
      teamKillReduceTime: {
        label: 'Perdonar una muerte de compañero tras (s)',
        hint: 'La cuenta baja en 1',
      },
      killFeed: {
        label: 'Registro de bajas',
        hint: 'Muestra quién mató a quién',
      },
      killFeedSpectator: {
        label: 'Registro de bajas para espectadores',
        hint: 'Y en las repeticiones',
      },
      killerInfo: {
        label: 'Mostrar quién te mató',
        hint: 'Atacante y arma en la pantalla de muerte',
      },
      killDistance: {
        label: 'Mostrar la distancia de la baja',
        hint: 'En la pantalla de muerte',
      },
      deadSay: {
        label: 'Los vivos leen a los muertos',
        hint: 'Chat de los jugadores muertos',
      },
      deadSayTeam: {
        label: 'Leer a compañeros muertos',
        hint: 'Chat de equipo de los muertos',
      },
      deadVoice: {
        label: 'Los vivos oyen a los muertos',
        hint: 'Voz de los jugadores muertos',
      },
      enemyHearsLocal: {
        label: 'Los enemigos oyen la voz local',
        hint: 'Voz de proximidad',
      },
      intermissionVoice: {
        label: 'Todos hablan entre rondas',
        hint: 'Los dos equipos',
      },
      mapVoting: {
        label: 'Votación de mapa',
        hint: 'Al final de la partida',
      },
      useMapCycle: {
        label: 'Usar la rotación de mapas',
        hint: 'Desactivado repite el mismo mapa',
      },
      minimumPlayers: {
        label: 'Jugadores para empezar',
        hint: 'Por equipo',
      },
      roundLimit: {
        label: 'Rondas por partida',
        hint: 'Varía según el modo',
      },
      winLimit: {
        label: 'Rondas para ganar',
        hint: 'Varía según el modo',
      },
      roundTime: {
        label: 'Tiempo de ronda (s)',
        hint: 'Varía según el modo',
      },
      preRoundTime: {
        label: 'Tiempo congelado al inicio (s)',
        hint: 'Al empezar la ronda',
      },
      postRoundTime: {
        label: 'Tiempo entre rondas (s)',
        hint: 'Pantalla de intermedio',
      },
      idleLimit: {
        label: 'Expulsar inactivos tras (s)',
        hint: '0 nunca expulsa',
      },
      idleLimitLow: {
        label: 'Límite de inactividad con pocos refuerzos (s)',
        hint: 'Cuando el equipo tiene pocas vidas',
      },
      bots: {
        label: 'Bots en versus',
        hint: 'Completa los equipos con bots',
      },
      botQuota: {
        label: 'Número de bots en versus',
        hint: 'Con los bots activados',
      },
      initialSupply: {
        label: 'Puntos de suministro iniciales',
        hint: 'Para comprar equipo',
      },
      maximumSupply: {
        label: 'Máximo de puntos de suministro',
        hint: 'Límite que se puede ganar',
      },
      supplyGain: {
        label: 'Ganar suministro durante la ronda',
        hint: 'Por puntuación',
      },
      voteKick: {
        label: 'Votación para expulsar',
        hint: 'Los jugadores abren la votación con Esc › Call Vote',
      },
      voteMinPlayers: {
        label: 'Jugadores para abrir una votación',
        hint: 'Mínimo en el servidor',
      },
      votePassRatio: {
        label: 'Votos a favor para aprobar',
        hint: 'Parte del equipo (0 a 1)',
      },
      voteBanDuration: {
        label: 'Bloqueo tras la votación (s)',
        hint: '-1 solo expulsa',
      },
      tickRate: {
        label: 'Tick rate',
        hint: 'Actualizaciones por segundo; más usa más CPU y ancho de banda',
      },
      maxClientRate: {
        label: 'Tasa máxima por cliente (red local)',
        hint: 'Bytes por segundo por jugador',
      },
      maxInternetClientRate: {
        label: 'Tasa máxima por cliente (internet)',
        hint: 'Bytes por segundo por jugador',
      },
    },
  },
  profiles: {
    label: 'Perfil',
    saveFirst: 'Guarda o deshaz los cambios antes de cambiar de perfil.',
    new: 'Perfil nuevo',
    rename: 'Renombrar',
    switchHint:
      'Cambiar de perfil requiere el servidor detenido. Las contraseñas y los tokens se comparten.',
    namePlaceholder: 'Ejemplo: Coop, Versus',
    copyCurrent: 'Copiar el perfil actual',
    create: 'Crear',
    delete: 'Borrar perfil',
    deleteYes: 'Borrar',
  },
  players: {
    title: 'Jugadores',
    load: 'Mostrar jugadores',
    reload: 'Actualizar',
    empty: 'No hay nadie en el servidor.',
    name: 'Nombre',
    score: 'Puntos',
    actions: 'Acciones',
    bot: 'bot',
    showBots: 'Mostrar los {n} bots',
    hideBots: 'Ocultar los {n} bots',
    noHumans: 'Ningún jugador, solo bots.',
    kick: 'Expulsar',
    ban: 'Bloquear',
    confirmKick: '¿Expulsar a {name}?',
    confirmBan: '¿Bloquear a {name}?',
    duration: 'Duración',
    reason: 'Motivo (opcional)',
    confirmYes: 'Confirmar',
    kicked: '{name} fue expulsado.',
    banned: '{name} fue bloqueado.',
    unbanned: 'Bloqueo retirado.',
    bans: 'Bloqueados',
    noBans: 'Ningún bloqueo.',
    unban: 'Quitar bloqueo',
    permanent: 'Permanente',
    days: '{n} día(s)',
    minutes: '{n} minutos',
    sayPlaceholder: 'Mensaje para todos en el servidor',
    say: 'Enviar',
    said: 'Mensaje enviado.',
  },
  automation: {
    restartTitle: 'Reinicio automático',
    daily: 'Reiniciar todos los días a las',
    dailyHint: 'Vacío = sin reinicio diario. Hora del equipo.',
    next: 'Próximo reinicio: {when}',
    update: 'Actualizar el servidor con SteamCMD en el reinicio diario',
    updateHint:
      'Las actualizaciones salen con los parches del juego; el reinicio tarda unos minutos más.',
    updateNoSteamCmd: 'Define primero la carpeta de SteamCMD en Instalación.',
    warn: 'Aviso en el chat antes del reinicio',
    warnHint: 'Se envía 5 y 1 minuto antes. {minutes} se cambia por el número.',
    warnExample: 'El servidor se reinicia en {minutes} min',
    onCrash: 'Iniciar el servidor otra vez si se cierra solo',
    onCrashHint: 'Espera 10 segundos. Se detiene tras 3 veces en 30 minutos.',
    appOpen: 'La automatización solo funciona con la app abierta (puede estar minimizada).',
    discordTitle: 'Avisos en Discord',
    discordText:
      'La app publica en un canal de Discord cuando el servidor arranca, se detiene, se cierra solo, o cuando alguien entra o sale. Vacío = ese aviso está desactivado.',
    webhook: 'Dirección del webhook',
    webhookHint:
      'En Discord: ajustes del canal › Integraciones › Webhooks › Nuevo webhook › Copiar URL.',
    discord: {
      ready: 'Servidor en línea',
      stopped: 'Servidor detenido',
      crashed: 'Servidor cerrado solo',
      join: 'Jugador entró',
      leave: 'Jugador salió',
    },
    discordExample: {
      ready: '🟢 {server} está en línea en {map}',
      stopped: '⚪ {server} se detuvo',
      crashed: '🔴 {server} se cerró solo',
      join: '➡️ {name} entró en {server}',
      leave: '⬅️ {name} salió de {server}',
    },
    fillExamples: 'Rellenar con ejemplos',
    test: 'Probar',
    testSent: 'Enviado. Revisa el canal.',
    tags: '{server} = nombre del servidor, {map} = mapa, {name} = jugador. La prueba usa los textos guardados.',
  },
  firewall: {
    title: 'Firewall de Windows',
    notInstalled: 'Instala primero el servidor.',
    ok: 'Reglas para el servidor encontradas en el Firewall de Windows.',
    text: 'Sin regla, Windows puede bloquear a los jugadores, sobre todo en redes marcadas como "Pública" (habitual con programas de VPN).',
    create: 'Crear la regla del firewall',
    uac: 'Windows pide la confirmación de administrador.',
  },
  validation: {
    markupPlain: 'El juego no muestra formato y mostraría {tag} tal cual. Quita las etiquetas.',
    rule: 'El valor de {rule} debe estar entre {min} y {max}.',
    ruleBool: 'Elige activado o desactivado en {rule}.',
    ruleMinMax: 'El mínimo de enemigos no puede ser mayor que el máximo.',
    lighting: 'Elige día, noche o el predeterminado del escenario.',
    lightingNoMap:
      'Este mapa no permite elegir día o noche desde la app. Déjalo en el predeterminado.',
    bindIp: 'Elige una dirección IPv4 válida de este equipo.',
    modIdsMax: 'Hasta {max} mods.',
    modListEmpty: 'Escribe al menos un id de mod, o desactiva la lista de mods.',
    webhook:
      'Pega la dirección del webhook copiada de Discord (empieza por https://discord.com/api/webhooks/).',
    announcementsMax: 'Hasta {max} mensajes programados.',
    announceEvery: 'Elige entre {min} y {max} minutos.',
    clock: 'Usa el formato HH:MM, por ejemplo 05:00.',
    discordText: 'Hasta {max} caracteres, en una línea.',
    reason: 'El motivo puede tener hasta {max} caracteres, en una línea y sin comillas.',
    netId: 'Este jugador no está en la lista. Actualiza la lista y vuelve a intentarlo.',
    banMinutes: 'Elige una duración de hasta {max} minutos.',
    profileName: 'Escribe un nombre de 1 a 40 caracteres.',
    profilesMax: 'Hasta {max} perfiles.',
    motd: 'El mensaje del día tiene hasta {max} caracteres y ningún carácter especial.',
    chatMessage: 'El mensaje tiene hasta {max} caracteres, en una línea.',
    chatNeedsRcon: 'Activa RCON en Contraseñas y tokens: la app envía este mensaje por RCON.',
    scenarioUnknown: 'Este mapa no está en la lista. Pide la lista al servidor otra vez.',
    scenarioNotLoaded: 'El mod de este mapa no está cargado en el servidor.',
    logKeep: 'Los inicios a mantener deben estar entre {min} y {max}.',
    unknownField: 'Campo desconocido.',
    password: 'Usa de 4 a 64 caracteres: letras, números y ! @ # $ % * . _ + - (sin espacios).',
    token: 'El token solo tiene letras y números (sin espacios). Cópialo de nuevo del sitio.',
    apiKey: 'La clave de la API tiene 32 letras y números.',
    email: 'Este correo no parece válido.',
    emailTypo: '¿Quisiste decir {suggestion}?',
    absolutePath: 'La carpeta debe ser una ruta completa.',
    hostnameLength: 'El nombre del servidor necesita de 1 a 64 caracteres.',
    hostnameQuotes: 'El nombre del servidor no puede llevar comillas (").',
    port: 'Los puertos van de 1024 a 65535.',
    portsDistinct: 'Los puertos de juego, consulta y RCON deben ser distintos.',
    maxPlayers: 'Los jugadores deben estar entre 1 y 64.',
    map: 'Nombre de mapa no válido (solo letras, números y _).',
    scenario:
      'El escenario empieza por Scenario_ (ejemplo: Scenario_Farmhouse_Checkpoint_Security).',
    travelNeedsMods: 'El mapa del mod solo funciona con los mods activados.',
    mutator: 'Mutator no válido: {value}.',
    ip: 'Escribe solo la IP, por ejemplo 203.0.113.10.',
    steamId: 'SteamID64 no válido: {value}. Tiene 17 dígitos y empieza por 7656119.',
    steamIdNotFound:
      'No se ha encontrado ningún SteamID64. Pega el número de 17 dígitos o un enlace steamcommunity.com/profiles/….',
    adminName: 'Nombre del administrador: hasta 40 caracteres, en una línea.',
    mapCycleLine: 'Cada línea de la rotación de mapas tiene hasta 200 caracteres.',
    modId: 'ID de mod no válido: {value}.',
    code: 'El código tiene 5 letras o números.',
    backupPassword: 'La contraseña de la copia necesita al menos {min} caracteres.',
  },
  problems: {
    unknown: {
      title: 'Algo ha salido mal',
      what: 'La aplicación encontró un error que no conoce.',
      cause: 'Puede ser algo nuevo, tras una actualización del juego o de Windows.',
      solution: [
        'Abre los registros y mira las últimas líneas en rojo.',
        'Inténtalo de nuevo. Si se repite, reinicia la aplicación.',
      ],
    },
    'server-not-installed': {
      title: 'Servidor no instalado',
      what: 'No se encontró InsurgencyServer.exe en la carpeta del servidor.',
      cause: 'Aún no se ha hecho la instalación, o la carpeta elegida no es la del servidor.',
      solution: [
        'Abre Instalación y haz clic en Instalar servidor.',
        'Si el servidor está instalado en otro sitio, elige esa carpeta en Instalación.',
      ],
    },
    'server-dir-invalid': {
      title: 'Carpeta no válida',
      what: 'La carpeta elegida no se puede usar.',
      cause:
        'No es una ruta completa, es la raíz de un disco o está dentro de una carpeta del sistema.',
      solution: [
        'Elige una carpeta en un disco con espacio (unos 25 GB más los mods), por ejemplo D:\\Servidores\\Sandstorm.',
        'Evita C:\\Windows, C:\\Program Files y raíces de disco como D:\\.',
      ],
    },
    'server-already-running': {
      title: 'El servidor ya está en marcha',
      what: 'Ya hay abierto un InsurgencyServer.exe de esta carpeta.',
      cause: 'Se inició antes, desde la aplicación o desde un archivo .bat.',
      solution: [
        'Usa Detener servidor y luego Iniciar, o Reiniciar.',
        'No ejecutes dos copias: se pelean por el mismo puerto.',
      ],
    },
    'server-not-running': {
      title: 'El servidor no está en marcha',
      what: 'Esta acción necesita el servidor en marcha.',
      cause: 'El servidor está detenido.',
      solution: ['Inicia el servidor y vuelve a intentarlo.'],
    },
    'server-start-failed': {
      title: 'El servidor no se abrió',
      what: 'Windows no pudo iniciar InsurgencyServer.exe.',
      cause: 'Puede que el archivo falte, esté dañado o lo bloquee el antivirus.',
      solution: [
        'En Instalación, haz clic en Verificar archivos para que SteamCMD repare la instalación.',
        'Comprueba que el antivirus no haya puesto InsurgencyServer.exe en cuarentena.',
      ],
    },
    'server-crashed': {
      title: 'El servidor se cerró solo',
      what: 'El proceso del servidor terminó sin que la aplicación lo pidiera.',
      cause:
        'Normalmente es un error fatal del juego: mods en conflicto, un archivo dañado o falta de memoria.',
      solution: [
        'Abre los registros y busca la última línea "Fatal error".',
        'Si el error nombra un archivo de mod, puede que dos mods traigan el mismo archivo (conflicto).',
        'En Instalación, usa Verificar archivos.',
      ],
    },
    'server-stop-failed': {
      title: 'No se pudo detener el servidor',
      what: 'Windows se negó a terminar el proceso del servidor.',
      cause: 'Puede que el servidor se iniciara como administrador, o ya se estaba cerrando.',
      solution: [
        'Espera unos segundos y vuelve a intentarlo.',
        'Si sigue igual, cierra InsurgencyServer.exe en el Administrador de tareas.',
      ],
    },
    'port-in-use': {
      title: 'Puerto en uso',
      what: 'El servidor no pudo abrir su puerto de juego o de consulta.',
      cause: 'Otro programa, u otra copia del servidor, está usando el puerto.',
      solution: [
        'Detén la otra copia del servidor, si la hay.',
        'O cambia el puerto en Servidor › General (por ejemplo 27103) y avisa a los jugadores.',
      ],
    },
    'mod-conflict': {
      title: 'Mods en conflicto',
      what: 'El servidor se cerró con "Serial size mismatch" o "Corrupt data found" al cargar el mapa.',
      cause:
        'Dos mods traen el mismo archivo del juego en versiones distintas, y uno sobrescribe al otro. Ya visto: Containment con Old World o Containment Storm.',
      solution: [
        'Mira arriba el archivo nombrado en el error y el mod que lo tiene, o usa Mods › Mods con problemas.',
        'Cancela la suscripción de ese mod en mod.io (cuenta del servidor) y vuelve a iniciar.',
      ],
    },
    'modio-code-invalid': {
      title: 'Código de mod.io no válido',
      what: 'mod.io respondió 401 "Invalid security code".',
      cause:
        'El código de 5 caracteres solo funciona una vez y por poco tiempo. Este ya se usó o caducó.',
      solution: ['En Mods, pide un código nuevo por correo y vuelve a autorizar.'],
    },
    'modio-not-authorized': {
      title: 'Servidor no autorizado en mod.io',
      what: 'El servidor no tiene una sesión válida de mod.io.',
      cause: 'Nunca se autorizó en este PC, o la autorización caducó.',
      solution: [
        'En Mods, escribe el correo de la cuenta del servidor y envía el código.',
        'Escribe el código recibido y haz clic en Autorizar.',
      ],
    },
    'modio-request-failed': {
      title: 'mod.io no respondió como se esperaba',
      what: 'La petición a la API de mod.io falló.',
      cause: 'Sin internet, mod.io caído o demasiadas peticiones.',
      solution: [
        'Comprueba internet y vuelve a intentarlo en unos minutos.',
        'Mira la página de estado de mod.io.',
      ],
    },
    'modio-email-invalid': {
      title: 'Correo no válido',
      what: 'El correo no parece válido.',
      cause: 'Una errata (por ejemplo "gamil.com" en lugar de "gmail.com").',
      solution: ['Revisa el correo de la cuenta de mod.io del servidor y vuelve a intentarlo.'],
    },
    'modio-no-api-key': {
      title: 'Falta la clave de la API de mod.io',
      what: 'Para pedir el código hace falta una clave de la API de mod.io.',
      cause: 'Aún no hay clave en la caja fuerte.',
      solution: [
        'En Mods, haz clic en "Buscar en los registros del servidor".',
        'O crea una clave en mod.io (acceso a la API) y pégala.',
      ],
    },
    'steamcmd-download': {
      title: 'No se pudo descargar SteamCMD',
      what: 'La descarga de SteamCMD no terminó.',
      cause: 'Sin internet, servidores de Valve inestables o el antivirus bloqueó el archivo.',
      solution: [
        'Comprueba internet y vuelve a intentarlo.',
        'Si ya tienes SteamCMD en otra carpeta, elígela en Instalación.',
      ],
    },
    'steamcmd-signature': {
      title: 'Firma de SteamCMD no válida',
      what: 'steamcmd.exe no tiene la firma digital esperada de Valve, así que no se ejecutó.',
      cause: 'Puede que el archivo se modificara por el camino o esté dañado.',
      solution: [
        'Haz clic en Descargar SteamCMD para descargarlo de nuevo.',
        'Si se repite, analiza el PC con el antivirus.',
      ],
    },
    'steamcmd-failed': {
      title: 'SteamCMD no pudo instalar',
      what: 'SteamCMD terminó con un error.',
      cause:
        'Habituales: "Missing configuration" (temporal), "Disk write failure" (carpeta sin permiso o antivirus), "No subscription" (aplicación equivocada).',
      solution: [
        'Vuelve a intentarlo: "Missing configuration" suele pasar a la segunda.',
        'Comprueba el espacio libre en disco y que la carpeta no esté protegida.',
      ],
    },
    'disk-full': {
      title: 'Disco lleno',
      what: 'No hay suficiente espacio libre en el disco.',
      cause:
        'El servidor necesita unos 25 GB, y los mods de mapa pueden superar los 5 GB cada uno.',
      solution: ['Libera espacio o elige una carpeta en otro disco.'],
    },
    'no-internet': {
      title: 'Sin conexión',
      what: 'La aplicación no pudo llegar a internet.',
      cause: 'Red desconectada, firewall bloqueando la aplicación o el sitio está caído.',
      solution: ['Comprueba la conexión y vuelve a intentarlo.'],
    },
    'invalid-field': {
      title: 'Información no válida',
      what: 'Un campo tiene un valor que el servidor no acepta.',
      cause: 'Mira el mensaje junto al campo.',
      solution: ['Corrige el campo y guarda de nuevo.'],
    },
    'vault-unavailable': {
      title: 'Caja fuerte no disponible',
      what: 'La aplicación no pudo abrir la caja fuerte donde guarda contraseñas y tokens.',
      cause:
        'La caja fuerte la creó otro usuario de Windows u otro PC (el cifrado de Windows solo se abre para quien la creó).',
      solution: [
        'En un PC nuevo, usa Copia de seguridad › Restaurar con el archivo de la copia y su contraseña.',
        'O vuelve a escribir las contraseñas y tokens en Servidor › Contraseñas y tokens.',
      ],
    },
    'backup-password': {
      title: 'Contraseña de la copia incorrecta',
      what: 'La copia no se abrió con esta contraseña.',
      cause: 'La contraseña es incorrecta. Sin la correcta la copia no se puede abrir.',
      solution: ['Revisa mayúsculas, minúsculas y el Bloq Mayús, y vuelve a intentarlo.'],
    },
    'backup-corrupt': {
      title: 'Archivo de copia no válido',
      what: 'El archivo no es una copia de esta aplicación, o está dañado.',
      cause: 'El archivo se modificó, se cortó al copiarlo o es de otro programa.',
      solution: ['Copia de nuevo el archivo desde el origen y vuelve a intentarlo.'],
    },
    'backup-write': {
      title: 'No se pudo guardar la copia',
      what: 'No se pudo escribir el archivo de la copia.',
      cause: 'Carpeta sin permiso, disco lleno o se quitó el pendrive.',
      solution: ['Elige otra carpeta y vuelve a intentarlo.'],
    },
    'config-write': {
      title: 'No se pudo guardar la configuración',
      what: 'No se pudo escribir un archivo de configuración del servidor.',
      cause: 'Puede que el servidor tenga el archivo abierto, o la carpeta está protegida.',
      solution: [
        'Detén el servidor y guarda de nuevo.',
        'Revisa los permisos de la carpeta del servidor.',
      ],
    },
    busy: {
      title: 'Espera a la tarea actual',
      what: 'Todavía hay otra tarea en curso.',
      cause:
        'La instalación, la copia de seguridad y el inicio del servidor no pueden ir a la vez.',
      solution: ['Espera a que termine o cancélala.'],
    },
    canceled: {
      title: 'Cancelado',
      what: 'La tarea se canceló.',
      cause: 'Hiciste clic en Cancelar.',
      solution: ['Vuelve a iniciarla cuando quieras.'],
    },
    'update-failed': {
      title: 'No se pudo instalar la actualización',
      what: 'No se pudo iniciar la nueva versión de la app.',
      cause: 'El instalador descargado se borró (por ejemplo, el antivirus) o Windows lo bloqueó.',
      solution: [
        'Cierra y abre la app: vuelve a descargar la nueva versión.',
        'O descarga el instalador en la página de descarga y ejecútalo.',
      ],
    },
    'mods-not-mounted': {
      title: 'Algunos mods no se cargaron',
      what: 'El servidor descargó mods nuevos y no cargó los que ya estaban instalados, así que pueden faltar el mapa del mod y los mutators.',
      cause:
        'Con suscripciones nuevas, el servidor solo carga lo que acaba de descargar. En el siguiente inicio carga todo.',
      solution: [
        'Haz clic en Reiniciar. La app lo hace sola cuando fue ella quien inició el servidor.',
      ],
    },
    'mod-delete': {
      title: 'No se pudo borrar el mod',
      what: 'La aplicación no pudo borrar los archivos del mod.',
      cause:
        'Un archivo del mod está en uso (por el servidor u otro programa) o Windows bloqueó la carpeta.',
      solution: [
        'Cierra los programas que puedan usar la carpeta del mod e inténtalo de nuevo.',
        'Si sigue fallando, abre Registros para ver qué archivo se bloqueó.',
      ],
    },
    'log-delete': {
      title: 'Algunos registros no se pudieron borrar',
      what: 'La app borró lo que pudo, pero quedaron algunos archivos de registro antiguos.',
      cause:
        'El archivo está abierto en otro programa (por ejemplo, un editor de texto) o Windows lo bloqueó.',
      solution: [
        'Cierra los programas que puedan tener el registro abierto e inténtalo de nuevo.',
        'O abre la carpeta de registros y bórralo a mano.',
      ],
    },
    'rcon-off': {
      title: 'RCON está desactivado',
      what: 'La app cambia el mapa por RCON, y no está activado.',
      cause: 'RCON está desactivado o no tiene contraseña en la app.',
      solution: [
        'Activa RCON y pon una contraseña en Servidor › Contraseñas y tokens.',
        'Reinicia el servidor para que abra RCON.',
      ],
    },
    'rcon-failed': {
      title: 'No se pudo hablar con el servidor',
      what: 'La app no pudo conectarse al servidor por RCON.',
      cause:
        'El servidor se inició antes de activar RCON, todavía está cargando, u otro programa usa el puerto de RCON.',
      solution: [
        'Espera a que el servidor esté listo e inténtalo de nuevo.',
        'Si acabas de activar RCON, reinicia el servidor.',
        'Revisa el puerto de RCON en Servidor › Contraseñas y tokens.',
      ],
    },
    'rcon-auth': {
      title: 'Contraseña de RCON rechazada',
      what: 'El servidor no aceptó la contraseña de RCON guardada en la app.',
      cause:
        'El servidor usa otra contraseña (por ejemplo la del Game.ini, o una cambiada después de iniciarlo).',
      solution: [
        'Reinicia el servidor desde la app para que use la contraseña guardada.',
        'O pon la misma contraseña en Servidor › Contraseñas y tokens.',
      ],
    },
    'map-change-failed': {
      title: 'El servidor no cambió el mapa',
      what: 'El servidor rechazó el cambio de mapa.',
      cause:
        'El mapa no está cargado en el servidor (por ejemplo, un mapa de mod cuyo mod no se cargó).',
      solution: [
        'Pide la lista al servidor otra vez y elige un mapa de ella.',
        'Para un mapa de mod, comprueba en la página Mods que el mod está activo.',
      ],
    },
    'discord-failed': {
      title: 'El mensaje de Discord no se envió',
      what: 'Discord no aceptó el mensaje del webhook.',
      cause: 'El webhook se borró, la dirección es incorrecta o no hay internet.',
      solution: [
        'En Discord, abre los ajustes del canal › Integraciones › Webhooks y copia la dirección otra vez.',
        'Pégala en Servidor › Automatización y pulsa Probar.',
      ],
    },
    'firewall-failed': {
      title: 'La regla del firewall no se creó',
      what: 'Windows no creó la regla para el servidor.',
      cause:
        'La confirmación de administrador falló u otro programa de seguridad bloqueó el cambio.',
      solution: [
        'Pulsa el botón otra vez y confirma la ventana de Windows.',
        'O crea la regla a mano en Firewall de Windows Defender › Configuración avanzada › Reglas de entrada.',
      ],
    },
    'autostart-failed': {
      title: 'No se pudo cambiar el inicio de Windows',
      what: 'La app no pudo añadirse ni quitarse de los programas que se abren con Windows.',
      cause: 'Un programa de seguridad o una directiva de Windows bloqueó el cambio.',
      solution: ['Vuelve a intentarlo.', 'O usa Configuración de Windows › Aplicaciones › Inicio.'],
    },
    'server-update-failed': {
      title: 'El servidor no se actualizó',
      what: 'El reinicio diario no pudo actualizar el servidor con SteamCMD.',
      cause:
        'La carpeta de SteamCMD no está definida, otra instalación estaba en curso o SteamCMD falló.',
      solution: [
        'Revisa Instalación: la carpeta de SteamCMD debe estar definida.',
        'Actualiza una vez a mano en Instalación para ver el error.',
      ],
    },
    'report-failed': {
      title: 'No se pudo preparar el informe',
      what: 'La app no pudo guardar el archivo del informe en este PC.',
      cause: 'El disco está lleno o no se puede escribir en la carpeta de la app.',
      solution: ['Libera espacio en el disco y vuelve a intentarlo.'],
    },
    'report-image': {
      title: 'Esta imagen no se puede usar',
      what: 'El archivo elegido no se añadió al informe.',
      cause: 'Pesa más de 5 MB, o no es una imagen PNG, JPG, GIF o WebP.',
      solution: ['Elige una imagen más pequeña, o guarda la captura como PNG o JPG.'],
    },
  },
};
