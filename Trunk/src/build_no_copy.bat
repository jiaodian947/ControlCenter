cd ../bin

cd ../src/charge
go build -o ../../bin/charge/charge.exe .

cd ..\usercenter
go build -o ..\..\bin\usercenter\usercenter.exe .

cd ..\controlcenter
go build -o ..\..\bin\controlcenter\controlcenter.exe .

pause

