Unicode true

####
## Instalador do Sandstorm Server Assist (ADR-014), a partir do modelo do Wails e do midia-multi-tools.
## Gerado por "npm run build:installer" (wails build -nsis -installscope user). O wails_tools.nsh é recriado
## pelo Wails a cada build com os dados do wails.json.
##
## Instalação só para o usuário (em %LOCALAPPDATA%\Programs, sem pedir administrador): é o que permite a
## atualização automática instalar sozinha. A configuração e o cofre ficam em %APPDATA%\sandstorm-server-assist
## e não são apagados ao desinstalar (o aviso de dados explica como apagar).
####
!ifndef REQUEST_EXECUTION_LEVEL
  !define REQUEST_EXECUTION_LEVEL "user"
!endif

!include "wails_tools.nsh"
!include "LogicLib.nsh"
!include "FileFunc.nsh"

# A versão destes dois precisa ter 4 partes.
VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

ManifestDPIAware true

!include "MUI2.nsh"

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
!define MUI_ABORTWARNING
!define MUI_FINISHPAGE_RUN "$INSTDIR\${PRODUCT_EXECUTABLE}"

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

# Os mesmos 5 idiomas do app (ADR-009). O primeiro é o padrão.
!insertmacro MUI_LANGUAGE "English"
!insertmacro MUI_LANGUAGE "Portuguese"
!insertmacro MUI_LANGUAGE "Spanish"
!insertmacro MUI_LANGUAGE "Russian"
!insertmacro MUI_LANGUAGE "Ukrainian"

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\bin\${INFO_PROJECTNAME}-${ARCH}-installer.exe"
!ifdef WAILS_INSTALL_SCOPE
  !if "${WAILS_INSTALL_SCOPE}" == "user"
    InstallDir "$LOCALAPPDATA\Programs\${INFO_PRODUCTNAME}"
  !else
    InstallDir "$PROGRAMFILES64\${INFO_PRODUCTNAME}"
  !endif
!else
  InstallDir "$LOCALAPPDATA\Programs\${INFO_PRODUCTNAME}"
!endif
ShowInstDetails show

Function .onInit
   !insertmacro wails.checkArchitecture

   ${IfNot} ${Silent}
       !insertmacro MUI_LANGDLL_DISPLAY
   ${EndIf}

   # Atualização automática: o app roda "instalador /S /relaunch" e fecha. Espera o executável ser liberado
   # (até ~30 s) antes de copiar os arquivos novos.
   ${GetParameters} $R0
   ClearErrors
   ${GetOptions} $R0 "/relaunch" $R1
   ${IfNot} ${Errors}
       StrCpy $R2 0
       wait_app:
           ClearErrors
           FileOpen $R3 "$INSTDIR\${PRODUCT_EXECUTABLE}" a
           ${IfNot} ${Errors}
               FileClose $R3
               Goto app_closed
           ${EndIf}
           IfFileExists "$INSTDIR\${PRODUCT_EXECUTABLE}" 0 app_closed
           IntOp $R2 $R2 + 1
           IntCmp $R2 60 app_closed
           Sleep 500
           Goto wait_app
       app_closed:
   ${EndIf}
FunctionEnd

# Depois de uma atualização automática, abre o app de novo (já na versão nova).
Function .onInstSuccess
   ${GetParameters} $R0
   ClearErrors
   ${GetOptions} $R0 "/relaunch" $R1
   ${IfNot} ${Errors}
       Exec '"$INSTDIR\${PRODUCT_EXECUTABLE}"'
   ${EndIf}
FunctionEnd

Section
    !insertmacro wails.setShellContext

    !insertmacro wails.webview2runtime

    SetOutPath $INSTDIR

    !insertmacro wails.files

    # A licença Apache 2.0 pede LICENSE e NOTICE junto da distribuição; vão também o aviso de isenção, os
    # termos de uso e o aviso de dados.
    File "..\..\..\..\LICENSE"
    File "..\..\..\..\NOTICE"
    File "..\..\..\..\DISCLAIMER.md"
    File "..\..\..\..\TERMS.md"
    File "..\..\..\..\PRIVACY.md"

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
    CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"

    !insertmacro wails.associateFiles
    !insertmacro wails.associateCustomProtocols

    !insertmacro wails.writeUninstaller
SectionEnd

Section "uninstall"
    !insertmacro wails.setShellContext

    RMDir /r "$AppData\${PRODUCT_EXECUTABLE}" # dados do WebView2

    RMDir /r $INSTDIR

    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.unassociateFiles
    !insertmacro wails.unassociateCustomProtocols

    !insertmacro wails.deleteUninstaller
SectionEnd
