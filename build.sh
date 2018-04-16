export GOOS=linux
export GOARCH=386

go get
go build lofai.go

docker build -t walked/lofai .

rm ./lofai