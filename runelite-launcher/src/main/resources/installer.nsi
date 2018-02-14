;--------------------------------
;Include Modern UI

  !include "MUI2.nsh"

;--------------------------------
;Defines
  !define PROJECT_FILE "${project.artifactId}-windows-4.0-386.exe"
  !define PROJECT_NAME "${project.name}"

;--------------------------------
;General

  ;Name and file
  Name "%${PROJECT_NAME}"
  OutFile "%${PROJECT_NAME}.exe"

  ;Default installation folder
  InstallDir "$LOCALAPPDATA\%${PROJECT_NAME}"

  ;Get installation folder from registry if available
  InstallDirRegKey HKCU "Software\%${PROJECT_NAME}" ""

  ;Request application privileges for Windows Vista
  RequestExecutionLevel user

;--------------------------------
;Interface Settings

  !define MUI_ABORTWARNING
  !define MUI_ICON "runelite.ico"
  !define MUI_UNICON "runelite.ico"

;--------------------------------
;Pages

  !insertmacro MUI_PAGE_LICENSE "license.txt"
  !insertmacro MUI_PAGE_COMPONENTS
  !insertmacro MUI_PAGE_DIRECTORY
  !insertmacro MUI_PAGE_INSTFILES

  !insertmacro MUI_UNPAGE_CONFIRM
  !insertmacro MUI_UNPAGE_INSTFILES

;--------------------------------
;Languages

  !insertmacro MUI_LANGUAGE "English"

;--------------------------------
;Installer Sections

Section "Install"

  SetOutPath "$INSTDIR"

  File "%${PROJECT_FILE}"

  ; Create desktop shortcut
  CreateShortCut "$DESKTOP\%${PROJECT_NAME}.lnk" "$INSTDIR\%${PROJECT_FILE}" ""

  ; Create start-menu items
  CreateDirectory "$SMPROGRAMS\%${PROJECT_NAME}"
  CreateShortCut "$SMPROGRAMS\%${PROJECT_NAME}\Uninstall.lnk" "$INSTDIR\Uninstall.exe" "" "$INSTDIR\Uninstall.exe" 0
  CreateShortCut "$SMPROGRAMS\%${PROJECT_NAME}\%${PROJECT_NAME}.lnk" "$INSTDIR\%${PROJECT_FILE}" "" "$INSTDIR\%${PROJECT_FILE}" 0

  ;Store installation folder
  WriteRegStr HKCU "Software\%${PROJECT_NAME}" "" $INSTDIR

  ;Create uninstaller
  WriteUninstaller "$INSTDIR\Uninstall.exe"


SectionEnd

;--------------------------------
;Uninstaller Section

Section "Uninstall"

  ; Delete Files
  RMDir /r "$INSTDIR\*.*"

  ; Remove the installation directory
  RMDir "$INSTDIR"

  ; Delete shortcuts
  Delete "$DESKTOP\%${PROJECT_NAME}.lnk"
  Delete "$SMPROGRAMS\%${PROJECT_NAME}\*.*"
  RmDir  "$SMPROGRAMS\%${PROJECT_NAME}"

  DeleteRegKey /ifempty HKCU "Software\%${PROJECT_NAME}"

SectionEnd
