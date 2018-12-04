cd ../bin
rd /s /q .
md charge
md usercenter
md controlcenter
cd ../src/charge
go build -o ../../bin/charge/charge.exe .
xcopy .\conf\*.* ..\..\bin\charge\conf\ /S/E

cd ..\usercenter
go build -o ..\..\bin\usercenter\usercenter.exe .
xcopy .\conf\*.* ..\..\bin\usercenter\conf\ /S/E
xcopy .\static\*.* ..\..\bin\usercenter\static\ /S/E
xcopy .\static_source\*.* ..\..\bin\usercenter\static_source\ /S/E
xcopy .\views\*.* ..\..\bin\usercenter\views\ /S/E

cd ..\controlcenter
go build -o ..\..\bin\controlcenter\controlcenter.exe .
xcopy .\conf\*.* ..\..\bin\controlcenter\conf\ /S/E
xcopy .\static\*.* ..\..\bin\controlcenter\static\ /S/E
xcopy .\views\*.* ..\..\bin\controlcenter\views\ /S/E
pause

