# Resumo do projeto: Sandstorm Server Assist

> Documento-base: o que é, stack, arquitetura e regras.
> Andamento: [PLANO_DESENVOLVIMENTO.md](PLANO_DESENVOLVIMENTO.md) · Decisões: [DECISOES.md](DECISOES.md)
> Como o servidor foi configurado à mão (histórico de sucesso): [../resumo_sessao.md](../resumo_sessao.md)

## 1. O que é

Um app desktop para Windows que ajuda a **instalar, configurar e iniciar** um servidor dedicado de
Insurgency: Sandstorm com mods do mod.io, sem precisar editar `.bat` e `.ini` à mão.

O que ele faz:

- **Primeiros passos:** uma lista de etapas mostra o que já está pronto e o que falta.
- **Instalação:** baixa o SteamCMD e instala ou atualiza o servidor (app 581330), com barra de progresso,
  velocidade e tempo restante.
- **Configuração:** nome do servidor, portas, jogadores, mapa inicial, mapa de mod (`-ModDownloadTravelTo`),
  administradores (SteamID64), rotação de mapas, senha de entrada, RCON, GSLT e GameStats. Cada informação
  que o usuário precisa buscar fora do app tem um guia "Como conseguir?" com passos e links.
- **mod.io:** pede o código de autorização por e-mail, autoriza o servidor, lista os mods inscritos e mostra o
  download de cada um.
- **Iniciar e parar:** o app monta a linha de comando sozinho e acompanha as fases pelo log: iniciando,
  mod.io pronto, mods ativados, viajando para o mapa, pronto para jogar.
- **Erros explicados:** cada erro conhecido mostra o que aconteceu, a causa provável e como resolver.
- **Logs:** escondidos por padrão, numa gaveta que abre com um clique, com filtro e busca. Segredos aparecem
  mascarados.
- **Segurança:** senhas e tokens ficam num cofre criptografado pelo Windows (DPAPI), não em `.ini` nem `.bat`.
- **Backup:** um arquivo único, protegido por senha, com tudo que é preciso para subir o servidor em outra
  máquina.

- **Idiomas:** inglês americano (padrão), português de Portugal, espanhol da Espanha, russo e ucraniano.
- **Aparência:** 5 temas de cor e 4 skins de layout, como no game-in.
- **Projeto público (futuro):** licença Apache 2.0, com aviso de isenção em `DISCLAIMER.md`.

Textos em **linguagem simples**.

## 2. Stack

| Camada         | Tecnologia                     | Papel                                                         |
| -------------- | ------------------------------ | ------------------------------------------------------------- |
| Janela / ponte | Wails 2.17 (WebView2)          | Janela nativa, ponte Go ↔ JavaScript, diálogos de arquivo     |
| Back-end       | Go 1.27 (Clean/Hexagonal)      | Regras, processos (SteamCMD, servidor), cofre, backup, mod.io |
| Front          | React 19 + TypeScript + Vite   | Telas                                                         |
| Externos       | SteamCMD, InsurgencyServer.exe | Instalar/atualizar e rodar o servidor (processos separados)   |

Mesma base do projeto `C:\workspace\midia-multi-tools` (scripts, lint, auditoria de segurança, barra de título,
tema). Motivos em [DECISOES.md](DECISOES.md) (ADR-001).

## 3. Arquitetura

```
[React (WebView2)] ──bindings Wails──> [adapters/desktop] ──> [app/* (casos de uso)] ──> [domain/* (regras puras)]
                                               │
                                               ├──> adapters/vault     (DPAPI do Windows)
                                               ├──> adapters/proc      (iniciar, achar e parar processos)
                                               ├──> adapters/steamcmd  (rodar o SteamCMD e ler o progresso)
                                               ├──> adapters/modio     (API do mod.io, só HTTPS e endereço fixo)
                                               └──> adapters/download  (baixar o SteamCMD com progresso)
```

```
desktop/internal/
  domain/   regras puras, sem I/O: perfil e linha de comando, eventos do log, catálogo de erros,
            progresso do SteamCMD, SteamID, INI, formato do backup, mascarar segredos
  app/      casos de uso: instalação, execução do servidor, mods, backup
  adapters/ Windows, processos, rede, arquivos e a ponte com a tela (desktop)
```

## 4. Regras do projeto

- O front manda **ids e opções**, nunca comandos nem caminhos livres. O Go monta os argumentos (lista, sem shell).
- Segredos nunca vão para o front em texto puro, exceto quando o usuário pede "Mostrar" num campo.
  Nunca vão para log.
- Links externos só por chave de uma lista fixa no Go.
- Internet só para: baixar o SteamCMD (endereço fixo da Valve) e a API do mod.io (endereço fixo).
- Nunca apagar arquivo do usuário sem confirmação. Arquivos de configuração são gravados de forma atômica e
  com cópia `.bak` antes da primeira alteração.
