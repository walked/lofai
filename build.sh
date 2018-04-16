export GOOS=linux
export GOARCH=386

go build lofai.go

sudo docker build -t walked/lofai .