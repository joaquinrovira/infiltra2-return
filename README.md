To build and run

go run ./build
./bin/app

VSCode launch configuration is ready. press f5 and start debugging.


Or run delve:

go run ./build -debug
dlv --listen=:2345 --headless=true --api-version=2 exec ./bin/app