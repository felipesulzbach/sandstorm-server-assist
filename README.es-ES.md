# Sandstorm Server Assist

[English](README.md) · [Português](README.pt-PT.md) · **Español** · [Русский](README.ru.md) · [Українська](README.uk.md)

[![Última versión](https://img.shields.io/github/v/release/felipesulzbach/sandstorm-server-assist)](https://github.com/felipesulzbach/sandstorm-server-assist/releases/latest)
[![Licencia: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)

> Traducción del [README en inglés](README.md). Si hay diferencias, vale la versión en inglés.

Una app gratuita para Windows que te ayuda a **instalar, configurar y poner en marcha un servidor dedicado de
Insurgency: Sandstorm con mods de mod.io**, sin editar archivos `.bat` e `.ini` a mano.

> Herramienta no oficial de la comunidad. Sin relación con New World Interactive, Focus Entertainment, Valve ni mod.io.
> Consulta [DISCLAIMER.md](DISCLAIMER.md).

![Panel con el servidor en marcha](docs/images/dashboard.png)

## Descarga

Descarga la última versión en la **[página de versiones (Releases)](https://github.com/felipesulzbach/sandstorm-server-assist/releases/latest)**:

- **Instalador** (`…-windows-setup.exe`, recomendado): se instala para tu usuario de Windows, sin pedir
  administrador, y se actualiza solo.
- **Portátil** (`…-windows-portable.zip`): descomprime y ejecuta; solo avisa cuando hay una versión nueva.

Necesitas Windows 10 u 11. Windows puede mostrar un aviso de SmartScreen porque la app no tiene firma de código:
pulsa **Más información** › **Ejecutar de todas formas**. Cada versión muestra el SHA-256 de sus archivos.

## Funciones

- **Configuración paso a paso:** una lista en el Panel muestra qué está listo y qué falta.
- **Instalación con progreso:** descarga SteamCMD (se comprueba la firma de Valve) e instala o actualiza el servidor,
  con velocidad y tiempo restante.
- **Usa el servidor que ya tienes:** elige su carpeta y la app lee su configuración, preguntando antes de cambiar lo
  que escribiste.
- **Todos los ajustes en un sitio:** nombre, puertos, jugadores, mapas y rotación, mapas de mods, día o noche,
  mutators (también los de tus mods), 60 reglas del juego, administradores, contraseñas y tokens, con guías
  "¿Cómo conseguirlo?".
- **Mods de mod.io:** autoriza el servidor una vez, ve la descarga de cada mod, y la app reinicia una vez cuando los
  mods nuevos lo necesitan.
- **Ponerlo en marcha:** iniciar y parar con las fases en directo, cambiar de mapa con el servidor encendido, panel de
  jugadores (expulsar, banear, mensaje), mensaje del día y mensajes de entrada y salida.
- **Automatización:** mensajes programados, reinicio diario con aviso y actualización, reinicio tras un cierre
  inesperado, avisos en Discord, abrir con Windows.
- **Diagnóstico de puertos:** comprueba el firewall, tu router (UPnP), tu IP pública, CGNAT y Steam, y explica cómo
  confirmarlo con un amigo.
- **Errores explicados y registros limpios:** cada error conocido dice qué pasó, por qué y cómo arreglarlo; las
  contraseñas y los tokens siempre se ocultan.
- **Secretos protegidos y copia de seguridad portátil:** Windows cifra los secretos (DPAPI); un único archivo con
  contraseña lleva todo a otro PC.
- **Comentarios:** informa de un problema o sugiere una idea desde la app, con el registro y una captura opcional.
- **5 idiomas y 5 temas:** English, Português, Español, Русский, Українська.

## Inicio rápido

1. Instala la app y acepta los términos de uso.
2. **Instalación:** elige la carpeta del servidor (unos 25 GB libres), descarga SteamCMD e instala el servidor. Si ya
   tienes un servidor, elige su carpeta.
3. **Servidor:** define el nombre, los mapas y las contraseñas. Consigue un token de servidor de juego de Steam (GSLT)
   con la guía.
4. **Mods:** autoriza el servidor en mod.io y suscríbete a mods en la web de mod.io.
5. **Servidor › Red:** redirige los puertos de juego y de consulta (UDP) en tu router, crea la regla del firewall y
   usa **Diagnosticar puertos**.
6. **Panel:** pulsa **Iniciar servidor** y comparte el comando para entrar con tus jugadores.

El tutorial completo, con imágenes, está en la **[guía de uso](docs/es-ES/GUIDE.md)**.

## Capturas

| Instalación                                           | Ajustes del servidor                                     |
| ----------------------------------------------------- | -------------------------------------------------------- |
| ![Instalación](docs/images/install.png)               | ![Ajustes del servidor](docs/images/server-settings.png) |
| **Diagnóstico de puertos**                            | **Comentarios**                                          |
| ![Diagnóstico de puertos](docs/images/port-check.png) | ![Comentarios](docs/images/feedback.png)                 |

## Comentarios y ayuda

- **¿Encontraste un problema o tienes una idea?** En la app, abre **Comentarios**: prepara un archivo con los
  detalles y abre el formulario de GitHub ya rellenado. O
  [abre una incidencia (issue)](https://github.com/felipesulzbach/sandstorm-server-assist/issues/new/choose).
- **¿Un problema de seguridad?** No abras una incidencia pública: consulta [SECURITY.md](SECURITY.md).

## Compilar desde el código

Requisitos: Windows 10/11 (WebView2), Go 1.27, Node.js 24; NSIS para el instalador.

```bash
npm install
npm run dev               # app en modo desarrollo
npm test                  # pruebas de Go, del front-end y de la documentación
npm run build             # desktop/build/bin/sandstorm-server-assist.exe
npm run build:installer   # instalador en desktop/build/dist
```

Wails 2 (Go, Clean/Hexagonal) + React 19 + Vite. Las reglas están en Go (`desktop/internal`); las pantallas solo
envían ids y opciones. La documentación técnica está en inglés: [arquitectura](docs/ARCHITECTURE.md),
[decisiones](docs/DECISIONS.md), [publicación](docs/RELEASING.md) y [cómo contribuir](CONTRIBUTING.md).

## Traducciones

La app y esta documentación están en 5 idiomas. Los textos de las pantallas están en
`desktop/frontend/src/i18n/locales/` (`en-US.ts` es la referencia). Las traducciones del README y de la guía siguen
los archivos en inglés; si hay diferencias, vale la versión en inglés. Se agradece ayuda con las traducciones.

## Créditos

Creado por **Felipe Sulzbach** ([@felipesulzbach](https://github.com/felipesulzbach)), con la ayuda de
[Claude Code](https://claude.com/claude-code) de Anthropic. Hecho para la comunidad de Insurgency: Sandstorm.

## Licencia

[Apache License 2.0](LICENSE). Consulta también [NOTICE](NOTICE), [DISCLAIMER.md](DISCLAIMER.md),
[TERMS.md](TERMS.md) y [PRIVACY.md](PRIVACY.md).
