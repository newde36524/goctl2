go build
./goctl2 api go --dir ./testApi --api ./testApi.api
rm -rf ./testApi
rm -rf ./goctl2
