main: main.go handler/handler.go parameter/parameter.go secret/secret.go snsevent/snsevent.go
	go build -o main
main.zip: main
	zip main.zip main
clean:
	rm -f main main.zip
