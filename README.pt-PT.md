# Sandstorm Server Assist

[English](README.md) · **Português** · [Español](README.es-ES.md) · [Русский](README.ru.md) · [Українська](README.uk.md)

[![Última versão](https://img.shields.io/github/v/release/felipesulzbach/sandstorm-server-assist)](https://github.com/felipesulzbach/sandstorm-server-assist/releases/latest)
[![Licença: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)

> Tradução do [README em inglês](README.md). Se houver diferenças, vale a versão em inglês.

Uma app gratuita para Windows que ajuda a **instalar, configurar e pôr a correr um servidor dedicado de
Insurgency: Sandstorm com mods do mod.io**, sem editar ficheiros `.bat` e `.ini` à mão.

> Ferramenta não oficial da comunidade. Sem ligação à New World Interactive, Focus Entertainment, Valve ou mod.io.
> Veja o [DISCLAIMER.md](DISCLAIMER.md).

![Painel com o servidor a correr](docs/images/dashboard.png)

## Transferir

Descarregue a versão mais recente na **[página de versões (Releases)](https://github.com/felipesulzbach/sandstorm-server-assist/releases/latest)**:

- **Instalador** (`…-windows-setup.exe`, recomendado): instala para o seu utilizador do Windows, sem pedir
  administrador, e atualiza-se sozinho.
- **Portátil** (`…-windows-portable.zip`): descompacte e execute; só avisa quando há uma versão nova.

Precisa do Windows 10 ou 11. O Windows pode mostrar um aviso do SmartScreen porque a app não tem assinatura de
código: clique em **Mais informações** › **Executar mesmo assim**. Cada versão lista o SHA-256 dos ficheiros.

## Funcionalidades

- **Configuração passo a passo:** uma lista no Painel mostra o que está pronto e o que falta.
- **Instalação com progresso:** descarrega o SteamCMD (a assinatura da Valve é verificada) e instala ou atualiza o
  servidor, com velocidade e tempo restante.
- **Use o servidor que já tem:** escolha a pasta dele e a app lê as definições, perguntando antes de mudar algo que
  escreveu.
- **Todas as definições num só sítio:** nome, portas, jogadores, mapas e rotação, mapas de mods, dia ou noite,
  mutators (incluindo os dos seus mods), 60 regras do jogo, administradores, palavras-passe e tokens, com guias
  "Como obter?".
- **Mods do mod.io:** autorize o servidor uma vez, veja o download de cada mod, e a app reinicia uma vez quando mods
  novos o exigem.
- **Pôr a correr:** iniciar e parar com as fases em tempo real, trocar o mapa com o servidor ligado, painel de
  jogadores (expulsar, banir, mensagem), mensagem do dia e mensagens de entrada e saída.
- **Automação:** mensagens programadas, reinício diário com aviso e atualização, reinício depois de uma falha, avisos
  no Discord, abrir com o Windows.
- **Diagnóstico de portas:** verifica a firewall, o router (UPnP), o seu IP público, CGNAT e a Steam, e explica como
  confirmar com um amigo.
- **Erros explicados e registos limpos:** cada erro conhecido diz o que aconteceu, porquê e como resolver; senhas e
  tokens ficam sempre mascarados.
- **Segredos protegidos e cópia de segurança portátil:** os segredos são cifrados pelo Windows (DPAPI); um único
  ficheiro com palavra-passe leva tudo para outro PC.
- **Feedback:** relate um problema ou sugira uma ideia a partir da app, com o registo e uma captura opcional.
- **5 idiomas e 5 temas:** English, Português, Español, Русский, Українська.

## Início rápido

1. Instale a app e aceite os termos de utilização.
2. **Instalação:** escolha a pasta do servidor (cerca de 25 GB livres), transfira o SteamCMD e instale o servidor. Se
   já tem um servidor, escolha a pasta dele.
3. **Servidor:** defina o nome, os mapas e as palavras-passe. Obtenha um token de servidor de jogo da Steam (GSLT)
   com o guia.
4. **Mods:** autorize o servidor no mod.io e subscreva mods no site do mod.io.
5. **Servidor › Rede:** reencaminhe as portas de jogo e de consulta (UDP) no router, crie a regra da firewall e use
   **Diagnosticar portas**.
6. **Painel:** clique em **Iniciar servidor** e partilhe o comando de entrada com os jogadores.

O tutorial completo, com imagens, está no **[guia de utilização](docs/pt-PT/GUIDE.md)**.

## Imagens

| Instalação                                           | Definições do servidor                                     |
| ---------------------------------------------------- | ---------------------------------------------------------- |
| ![Instalação](docs/images/install.png)               | ![Definições do servidor](docs/images/server-settings.png) |
| **Diagnóstico de portas**                            | **Feedback**                                               |
| ![Diagnóstico de portas](docs/images/port-check.png) | ![Feedback](docs/images/feedback.png)                      |

## Feedback e ajuda

- **Encontrou um problema ou tem uma ideia?** Na app, abra **Feedback**: ela prepara um ficheiro com os detalhes e
  abre o formulário do GitHub já preenchido. Ou
  [abra um pedido (issue)](https://github.com/felipesulzbach/sandstorm-server-assist/issues/new/choose).
- **Problema de segurança?** Não abra um pedido público: veja o [SECURITY.md](SECURITY.md).

## Compilar a partir do código

Requisitos: Windows 10/11 (WebView2), Go 1.27, Node.js 24; NSIS para o instalador.

```bash
npm install
npm run dev               # app em modo de desenvolvimento
npm test                  # testes do Go, do front-end e da documentação
npm run build             # desktop/build/bin/sandstorm-server-assist.exe
npm run build:installer   # instalador em desktop/build/dist
```

Wails 2 (Go, Clean/Hexagonal) + React 19 + Vite. As regras ficam no Go (`desktop/internal`); os ecrãs só enviam ids e
opções. A documentação técnica está em inglês: [arquitetura](docs/ARCHITECTURE.md), [decisões](docs/DECISIONS.md),
[publicação](docs/RELEASING.md) e [como contribuir](CONTRIBUTING.md).

## Traduções

A app e esta documentação existem em 5 idiomas. Os textos dos ecrãs estão em `desktop/frontend/src/i18n/locales/`
(`en-US.ts` é a referência). As traduções do README e do guia seguem os ficheiros em inglês; se houver diferenças,
vale a versão em inglês. Ajuda com traduções é bem-vinda.

## Créditos

Criado por **Felipe Sulzbach** ([@felipesulzbach](https://github.com/felipesulzbach)), com a ajuda do
[Claude Code](https://claude.com/claude-code) da Anthropic. Feito para a comunidade de Insurgency: Sandstorm.

## Licença

[Apache License 2.0](LICENSE). Veja também [NOTICE](NOTICE), [DISCLAIMER.md](DISCLAIMER.md), [TERMS.md](TERMS.md) e
[PRIVACY.md](PRIVACY.md).
