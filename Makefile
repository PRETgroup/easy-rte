default: easy-rte-c easy-rte-parser

easy-rte-c: rtec/* rtedef/*
	go get github.com/PRETgroup/goFB/goFB
	go build -o easy-rte-c -i ./rtec/main

easy-rte-parser: rteparser/* rtedef/*
	go get github.com/PRETgroup/goFB/goFB
	go build -o easy-rte-parser -i ./rteparser/main

example_ab5: example/ab5/example_main.c example/ab5/F_AB5Function.c
	gcc example/ab5/example_main.c example/ab5/F_AB5Function.c -o example_AB5

example/ab5/F_AB5Function.c: easy-rte-c easy-rte-parser example/ab5/ab5.erte
	./easy-rte-parser -i example/ab5 -o example/ab5
	./easy-rte-c -i example/ab5 -o example/ab5

clean:
	rm -f easy-rte-c
	rm -f easy-rte-parser
	rm -f example_AB5
	go get -u github.com/PRETgroup/goFB/goFB