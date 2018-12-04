#!/bin/sh
cd ./8000  
nohup go run ./main.go > /dev/null 2&>1 & 
cd ../8001 
nohup go run ./main.go > /dev/null 2&>1 & 
cd ../8002 
nohup go run ./main.go > /dev/null 2&>1 & 
cd ../8003 
nohup go run ./main.go > /dev/null 2&>1 & 
cd ../
