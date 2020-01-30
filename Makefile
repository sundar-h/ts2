
clean:
	rm log/*
	rm -rf out

install:
	go build -o ts2 ./main.go; mv ts2 ~/go/root/bin
