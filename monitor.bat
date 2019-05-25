@echo off
:start
tasklist |findstr /bc:"client.exe"&&(goto A) ||(goto B)
if "%errorlevel%"=="0" (goto B) else (goto A)
:A
ping -4 www.baidu.com
:B
E:\client.exe E:\ss.json

ping -4 www.baidu.com
goto start