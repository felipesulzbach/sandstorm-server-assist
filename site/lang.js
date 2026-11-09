// Idiomas da página de download (os mesmos 5 do app). O inglês está no próprio HTML; aqui ficam os outros.
// **texto** vira <strong> e `texto` vira código. Sem innerHTML: os nós são criados um a um.
'use strict';

const T = {
  pt: {
    language: 'Idioma',
    lead: 'Instale, configure e inicie o seu próprio servidor dedicado de Insurgency: Sandstorm com mods do mod.io. Passo a passo, com erros explicados e as suas senhas protegidas.',
    download: '⬇ Baixar para Windows',
    portable: 'Versão portátil (.zip)',
    version: 'Versão',
    requires: 'Windows 10 ou 11 (64 bits)',
    installer: 'instalador',
    portableShort: 'portátil',
    f1t: '⬇ Instalação com progresso',
    f1: 'Baixa o SteamCMD e o servidor do jogo das fontes oficiais, com velocidade e tempo restante.',
    f2t: '⚙ Configuração do servidor',
    f2: 'Nome, portas, jogadores, mapas, mapa de mod, mutators, administradores e rotação de mapas, em palavras simples.',
    f3t: '↺ Use o servidor que já tem',
    f3: 'Lê um servidor existente e preenche os campos. Se mudou algo, pergunta o que manter.',
    f4t: '🧩 Mods do mod.io',
    f4: 'Autoriza o servidor uma vez e mostra cada mod a ser baixado e montado.',
    f5t: '⚠ Erros explicados',
    f5: 'Cada erro conhecido diz o que aconteceu, a causa provável e como resolver. Os logs ficam a um clique.',
    f6t: '🔒 Segredos protegidos',
    f6: 'Senhas e tokens são cifrados pelo Windows e nunca gravados em arquivos .ini ou .bat.',
    f7t: '💾 Backup e restauração',
    f7: 'Um arquivo protegido por senha com tudo para subir o servidor noutro PC.',
    f8t: '⟳ Atualização automática',
    f8: 'O app instalado atualiza-se sozinho com versões assinadas. O servidor do jogo continua a funcionar.',
    f9t: '🌐 5 idiomas',
    howTitle: 'Como instalar',
    how1: 'Baixe e abra o **sandstorm-server-assist-setup.exe**.',
    how2: 'O Windows pode avisar que o app não é reconhecido (SmartScreen). Clique em **Mais informações** e depois em **Executar assim mesmo**.',
    how3: 'Siga o instalador (não precisa de administrador) e abra o **Sandstorm Server Assist** pelo menu Iniciar.',
    how4: 'Aceite os termos de uso e o aviso de privacidade e siga os primeiros passos no Painel.',
    howPortable:
      'Prefere não instalar? Baixe a **versão portátil**, extraia e abra o **sandstorm-server-assist.exe**. Ela não se atualiza sozinha: volte aqui para as versões novas.',
    verifyTitle: 'Conferir o download',
    verifyText:
      'Para ter certeza de que o arquivo chegou inteiro, compare o SHA-256 (no PowerShell: `Get-FileHash arquivo`).',
    published: 'publicado em',
    terms: 'Termos de uso',
    privacy: 'Privacidade',
    disclaimer:
      'Projeto não oficial, sem ligação nem aprovação da New World Interactive, Focus Entertainment, Valve ou mod.io. Insurgency, Steam e mod.io são marcas dos seus donos. O app não distribui o jogo: baixa-o das fontes oficiais.',
  },
  es: {
    language: 'Idioma',
    lead: 'Instala, configura e inicia tu propio servidor dedicado de Insurgency: Sandstorm con mods de mod.io. Paso a paso, con errores explicados y tus contraseñas protegidas.',
    download: '⬇ Descargar para Windows',
    portable: 'Versión portátil (.zip)',
    version: 'Versión',
    requires: 'Windows 10 u 11 (64 bits)',
    installer: 'instalador',
    portableShort: 'portátil',
    f1t: '⬇ Instalación con progreso',
    f1: 'Descarga SteamCMD y el servidor del juego de las fuentes oficiales, con velocidad y tiempo restante.',
    f2t: '⚙ Configuración del servidor',
    f2: 'Nombre, puertos, jugadores, mapas, mapa de mod, mutators, administradores y rotación de mapas, con palabras sencillas.',
    f3t: '↺ Usa el servidor que ya tienes',
    f3: 'Lee un servidor existente y rellena los campos. Si cambiaste algo, pregunta qué mantener.',
    f4t: '🧩 Mods de mod.io',
    f4: 'Autoriza el servidor una vez y muestra cada mod mientras se descarga y se monta.',
    f5t: '⚠ Errores explicados',
    f5: 'Cada error conocido dice qué pasó, la causa probable y cómo arreglarlo. Los registros están a un clic.',
    f6t: '🔒 Secretos protegidos',
    f6: 'Contraseñas y tokens se cifran con Windows y nunca se guardan en archivos .ini o .bat.',
    f7t: '💾 Copia de seguridad',
    f7: 'Un archivo protegido con contraseña con todo para montar el servidor en otro PC.',
    f8t: '⟳ Actualización automática',
    f8: 'La app instalada se actualiza sola con versiones firmadas. El servidor del juego sigue funcionando.',
    f9t: '🌐 5 idiomas',
    howTitle: 'Cómo instalar',
    how1: 'Descarga y abre **sandstorm-server-assist-setup.exe**.',
    how2: 'Windows puede avisar de que la app no se reconoce (SmartScreen). Haz clic en **Más información** y luego en **Ejecutar de todas formas**.',
    how3: 'Sigue el instalador (no necesita administrador) y abre **Sandstorm Server Assist** desde el menú Inicio.',
    how4: 'Acepta los términos de uso y el aviso de privacidad y sigue los primeros pasos del Panel.',
    howPortable:
      '¿Prefieres no instalar? Descarga la **versión portátil**, extráela y abre **sandstorm-server-assist.exe**. No se actualiza sola: vuelve aquí para las versiones nuevas.',
    verifyTitle: 'Comprobar la descarga',
    verifyText:
      'Para asegurarte de que el archivo llegó completo, compara el SHA-256 (en PowerShell: `Get-FileHash archivo`).',
    published: 'publicado el',
    terms: 'Términos de uso',
    privacy: 'Privacidad',
    disclaimer:
      'Proyecto no oficial, sin relación ni aprobación de New World Interactive, Focus Entertainment, Valve o mod.io. Insurgency, Steam y mod.io son marcas de sus dueños. La app no distribuye el juego: lo descarga de las fuentes oficiales.',
  },
  ru: {
    language: 'Язык',
    lead: 'Установите, настройте и запустите собственный выделенный сервер Insurgency: Sandstorm с модами mod.io. Шаг за шагом, с понятными ошибками и защищёнными паролями.',
    download: '⬇ Скачать для Windows',
    portable: 'Портативная версия (.zip)',
    version: 'Версия',
    requires: 'Windows 10 или 11 (64 бита)',
    installer: 'установщик',
    portableShort: 'портативная',
    f1t: '⬇ Установка с прогрессом',
    f1: 'Скачивает SteamCMD и игровой сервер из официальных источников, со скоростью и оставшимся временем.',
    f2t: '⚙ Настройка сервера',
    f2: 'Название, порты, игроки, карты, карта мода, мутаторы, администраторы и ротация карт — простыми словами.',
    f3t: '↺ Используйте свой сервер',
    f3: 'Читает существующий сервер и заполняет поля. Если вы что-то изменили, спрашивает, что оставить.',
    f4t: '🧩 Моды mod.io',
    f4: 'Один раз авторизует сервер и показывает загрузку и подключение каждого мода.',
    f5t: '⚠ Понятные ошибки',
    f5: 'Для каждой известной ошибки — что случилось, вероятная причина и как исправить. Логи открываются одним кликом.',
    f6t: '🔒 Защищённые секреты',
    f6: 'Пароли и токены шифруются Windows и никогда не записываются в файлы .ini или .bat.',
    f7t: '💾 Резервная копия',
    f7: 'Один защищённый паролем файл со всем, что нужно для сервера на другом ПК.',
    f8t: '⟳ Автообновление',
    f8: 'Установленное приложение само обновляется подписанными версиями. Игровой сервер продолжает работать.',
    f9t: '🌐 5 языков',
    howTitle: 'Как установить',
    how1: 'Скачайте и откройте **sandstorm-server-assist-setup.exe**.',
    how2: 'Windows может сообщить, что приложение не распознано (SmartScreen). Нажмите **Подробнее**, затем **Выполнить в любом случае**.',
    how3: 'Пройдите установку (права администратора не нужны) и откройте **Sandstorm Server Assist** из меню «Пуск».',
    how4: 'Примите условия использования и уведомление о конфиденциальности и выполните первые шаги на панели.',
    howPortable:
      'Не хотите устанавливать? Скачайте **портативную версию**, распакуйте и откройте **sandstorm-server-assist.exe**. Она не обновляется сама: за новыми версиями возвращайтесь сюда.',
    verifyTitle: 'Проверка загрузки',
    verifyText:
      'Чтобы убедиться, что файл скачался целиком, сравните SHA-256 (в PowerShell: `Get-FileHash файл`).',
    published: 'опубликовано',
    terms: 'Условия использования',
    privacy: 'Конфиденциальность',
    disclaimer:
      'Неофициальный проект, не связанный с New World Interactive, Focus Entertainment, Valve или mod.io и не одобренный ими. Insurgency, Steam и mod.io — товарные знаки их владельцев. Приложение не распространяет игру: оно скачивает её из официальных источников.',
  },
  uk: {
    language: 'Мова',
    lead: 'Встановіть, налаштуйте і запустіть власний виділений сервер Insurgency: Sandstorm з модами mod.io. Крок за кроком, зі зрозумілими помилками й захищеними паролями.',
    download: '⬇ Завантажити для Windows',
    portable: 'Портативна версія (.zip)',
    version: 'Версія',
    requires: 'Windows 10 або 11 (64 біти)',
    installer: 'інсталятор',
    portableShort: 'портативна',
    f1t: '⬇ Встановлення з прогресом',
    f1: 'Завантажує SteamCMD та ігровий сервер з офіційних джерел, зі швидкістю і часом, що залишився.',
    f2t: '⚙ Налаштування сервера',
    f2: 'Назва, порти, гравці, карти, карта моду, мутатори, адміністратори й ротація карт — простими словами.',
    f3t: '↺ Використайте свій сервер',
    f3: 'Читає наявний сервер і заповнює поля. Якщо ви щось змінили, питає, що залишити.',
    f4t: '🧩 Моди mod.io',
    f4: 'Один раз авторизує сервер і показує завантаження та підключення кожного моду.',
    f5t: '⚠ Зрозумілі помилки',
    f5: 'Для кожної відомої помилки — що сталося, ймовірна причина і як виправити. Логи відкриваються одним кліком.',
    f6t: '🔒 Захищені секрети',
    f6: 'Паролі й токени шифруються Windows і ніколи не записуються у файли .ini чи .bat.',
    f7t: '💾 Резервна копія',
    f7: 'Один захищений паролем файл з усім, що потрібно для сервера на іншому ПК.',
    f8t: '⟳ Автооновлення',
    f8: 'Встановлений застосунок сам оновлюється підписаними версіями. Ігровий сервер продовжує працювати.',
    f9t: '🌐 5 мов',
    howTitle: 'Як встановити',
    how1: 'Завантажте і відкрийте **sandstorm-server-assist-setup.exe**.',
    how2: 'Windows може повідомити, що застосунок не розпізнано (SmartScreen). Натисніть **Докладніше**, потім **Однаково запустити**.',
    how3: 'Пройдіть встановлення (права адміністратора не потрібні) і відкрийте **Sandstorm Server Assist** з меню «Пуск».',
    how4: 'Прийміть умови використання та повідомлення про конфіденційність і виконайте перші кроки на панелі.',
    howPortable:
      'Не хочете встановлювати? Завантажте **портативну версію**, розпакуйте і відкрийте **sandstorm-server-assist.exe**. Вона не оновлюється сама: за новими версіями повертайтеся сюди.',
    verifyTitle: 'Перевірка завантаження',
    verifyText:
      'Щоб переконатися, що файл завантажився повністю, порівняйте SHA-256 (у PowerShell: `Get-FileHash файл`).',
    published: 'опубліковано',
    terms: 'Умови використання',
    privacy: 'Конфіденційність',
    disclaimer:
      'Неофіційний проєкт, не пов’язаний з New World Interactive, Focus Entertainment, Valve чи mod.io і не схвалений ними. Insurgency, Steam і mod.io — торговельні марки їхніх власників. Застосунок не поширює гру: він завантажує її з офіційних джерел.',
  },
};

const HTML_LANG = { en: 'en', pt: 'pt-PT', es: 'es-ES', ru: 'ru', uk: 'uk' };

/** Monta os nós de um texto com **negrito** e `código`. */
function nodes(text) {
  const out = [];
  for (const part of text.split(/(\*\*[^*]+\*\*|`[^`]+`)/)) {
    if (!part) continue;
    if (part.startsWith('**')) {
      const el = document.createElement('strong');
      el.textContent = part.slice(2, -2);
      out.push(el);
    } else if (part.startsWith('`')) {
      const el = document.createElement('span');
      el.className = 'code';
      el.textContent = part.slice(1, -1);
      out.push(el);
    } else {
      out.push(document.createTextNode(part));
    }
  }
  return out;
}

const original = new Map();

function apply(lang) {
  const dict = T[lang];
  for (const el of document.querySelectorAll('[data-i18n]')) {
    if (!original.has(el))
      original.set(
        el,
        [...el.childNodes].map((n) => n.cloneNode(true)),
      );
    const text = dict && dict[el.dataset.i18n];
    el.replaceChildren(...(text ? nodes(text) : original.get(el).map((n) => n.cloneNode(true))));
  }
  document.documentElement.lang = HTML_LANG[lang] || 'en';
}

function detect() {
  try {
    const saved = localStorage.getItem('lang');
    if (saved && HTML_LANG[saved]) return saved;
  } catch {
    // sem armazenamento: segue pelo navegador
  }
  const nav = (navigator.language || '').slice(0, 2).toLowerCase();
  return HTML_LANG[nav] ? nav : 'en';
}

document.addEventListener('DOMContentLoaded', () => {
  const select = document.getElementById('lang');
  const lang = detect();
  if (select) {
    select.value = lang;
    select.addEventListener('change', () => {
      apply(select.value);
      try {
        localStorage.setItem('lang', select.value);
      } catch {
        // sem armazenamento: vale só nesta visita
      }
    });
  }
  apply(lang);
});
