./build-binary.sh

docker build --tag datasail/extract-weather:1.0 .

docker push datasail/extract-weather:1.0
