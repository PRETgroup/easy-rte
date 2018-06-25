default: easy-rte-c easy-rte-parser

easy-rte-c: rtec/* rtedef/*
	go get github.com/PRETgroup/goFB/goFB
	go build -o easy-rte-c -i ./rtec/main

easy-rte-parser: rteparser/* rtedef/*
	go get github.com/PRETgroup/goFB/goFB
	go build -o easy-rte-parser -i ./rteparser/main

example: example/example_main.c example/F_AB5Function.c
	gcc example/example_main.c example/F_AB5Function.c -o example_AB5

example/F_AB5Function.c: easy-rte-c easy-rte-parser example/ab5.erte
	./easy-rte-parser -i example -o example
	./easy-rte-c -i example -o example

clean:
	rm -f easy-rte-c
	rm -f easy-rte-parser
	rm -f example_AB5
	go get -u github.com/PRETgroup/goFB/goFB