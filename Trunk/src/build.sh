#!/bin/bash
cd ./charge
go build -o ../../bin/charge/charge .
cp -rf ./conf/ ../../bin/charge/

cd ../usercenter
go build -o ../../bin/usercenter/usercenter .
cp -rf ./conf ../../bin/usercenter/
cp -rf ./static ../../bin/usercenter/
cp -rf ./static_source ../../bin/usercenter/
cp -rf ./views ../../bin/usercenter/

cd ../controlcenter
go build -o ../../bin/controlcenter/controlcenter .
cp -rf ./conf ../../bin/controlcenter/
cp -rf ./static ../../bin/controlcenter/
cp -rf ./views ../../bin/controlcenter/
