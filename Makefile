all:
	go build -o lilscan main.go
install:
	cp lilscan $$GOPATH/bin
