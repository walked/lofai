export GOOS=linux
export GOARCH=386

go get
go build -o ./artifacts/lofai

docker build -t walked/lofai .

rm -rf ./artifacts