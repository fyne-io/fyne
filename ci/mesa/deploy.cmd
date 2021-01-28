@echo off
@cd /d "%~dp0"

:--------------------------------------
@echo Mesa3D system-wide deployment utility
@echo -------------------------------------

@set mesaloc=%~dp0
@IF "%mesaloc:~-1%"=="\" set mesaloc=%mesaloc:~0,-1%

rem desktopgl
@IF /I %PROCESSOR_ARCHITECTURE%==X86 copy "%mesaloc%\x86\opengl32.dll" "%windir%\System32\mesadrv.dll"
@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 copy "%mesaloc%\x86\opengl32.dll" "%windir%\SysWOW64\mesadrv.dll"
@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 copy "%mesaloc%\x64\opengl32.dll" "%windir%\System32\mesadrv.dll"
@IF /I %PROCESSOR_ARCHITECTURE%==X86 IF EXIST "%mesaloc%\x86\libglapi.dll" copy "%mesaloc%\x86\libglapi.dll" "%windir%\System32"
@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 IF EXIST "%mesaloc%\x86\libglapi.dll" copy "%mesaloc%\x86\libglapi.dll" "%windir%\SysWOW64"
@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 IF EXIST "%mesaloc%\x64\libglapi.dll" copy "%mesaloc%\x64\libglapi.dll" "%windir%\System32"

@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 REG ADD "HKEY_LOCAL_MACHINE\SOFTWARE\Wow6432Node\Microsoft\Windows NT\CurrentVersion\OpenGLDrivers\MSOGL" /v "DLL" /t REG_SZ /d "mesadrv.dll" /f
@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 REG ADD "HKEY_LOCAL_MACHINE\SOFTWARE\Wow6432Node\Microsoft\Windows NT\CurrentVersion\OpenGLDrivers\MSOGL" /v "DriverVersion" /t REG_DWORD /d "1" /f
@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 REG ADD "HKEY_LOCAL_MACHINE\SOFTWARE\Wow6432Node\Microsoft\Windows NT\CurrentVersion\OpenGLDrivers\MSOGL" /v "Flags" /t REG_DWORD /d "1" /f
@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 REG ADD "HKEY_LOCAL_MACHINE\SOFTWARE\Wow6432Node\Microsoft\Windows NT\CurrentVersion\OpenGLDrivers\MSOGL" /v "Version" /t REG_DWORD /d "2" /f
@REG ADD "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\OpenGLDrivers\MSOGL" /v "DLL" /t REG_SZ /d "mesadrv.dll" /f
@REG ADD "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\OpenGLDrivers\MSOGL" /v "DriverVersion" /t REG_DWORD /d "1" /f
@REG ADD "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\OpenGLDrivers\MSOGL" /v "Flags" /t REG_DWORD /d "1" /f
@REG ADD "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\OpenGLDrivers\MSOGL" /v "Version" /t REG_DWORD /d "2" /f
@echo.
@echo Desktop OpenGL drivers deploy complete.

echo "PROCESSOR_ARCHITECTURE=%PROCESSOR_ARCHITECTURE%"
rem osmesa 
@IF /I %PROCESSOR_ARCHITECTURE%==X86 IF EXIST "%mesaloc%\x86\osmesa.dll" copy "%mesaloc%\x86\osmesa.dll" "%windir%\System32"
@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 IF EXIST "%mesaloc%\x86\osmesa.dll" copy "%mesaloc%\x86\osmesa.dll" "%windir%\SysWOW64"
@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 IF EXIST "%mesaloc%\x64\osmesa.dll" copy "%mesaloc%\x64\osmesa.dll" "%windir%\System32"

@IF /I %PROCESSOR_ARCHITECTURE%==X86 IF EXIST "%mesaloc%\x86\libglapi.dll" copy "%mesaloc%\x86\libglapi.dll" "%windir%\System32"
@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 IF EXIST "%mesaloc%\x86\libglapi.dll" copy "%mesaloc%\x86\libglapi.dll" "%windir%\SysWOW64"
@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 IF EXIST "%mesaloc%\x64\libglapi.dll" copy "%mesaloc%\x64\libglapi.dll" "%windir%\System32"
rem@set osmesatype=gallium

rem@IF /I %PROCESSOR_ARCHITECTURE%==X86 copy "%mesaloc%\x86\osmesa-%osmesatype%\osmesa.dll" "%windir%\System32"
rem@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 copy "%mesaloc%\x86\osmesa-%osmesatype%\osmesa.dll" "%windir%\SysWOW64"
rem@IF /I %PROCESSOR_ARCHITECTURE%==AMD64 copy "%mesaloc%\x64\osmesa-%osmesatype%\osmesa.dll" "%windir%\System32"
@echo.
@echo Off-screen render driver deploy complete.

echo All done