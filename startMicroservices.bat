@echo off

echo Starting Authentication Microservice...
cd authenticationMicroservice
start cmd /k "go run main.go"
cd ..

echo Starting Vehicle Microservice...
cd vehicleMicroservice
start cmd /k "go run main.go"
cd ..

echo Starting User Microservice...
cd userMicroservice
start cmd /k "go run main.go"
cd ..

echo Starting Frontend Live Server...
cd frontend
start cmd /k "npx live-server"
cd ..
exit
