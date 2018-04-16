export GOOS=linux
export GOARCH=386

go build lofai.go

docker build -t walked/lofai .