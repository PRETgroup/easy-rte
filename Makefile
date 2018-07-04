default: easy-rte-c easy-rte-parser

easy-rte-c: rtec/* rtedef/*
	go get github.com/PRETgroup/goFB/goFB
	go build -o easy-rte-c -i ./rtec/main

easy-rte-parser: rteparser/* rtedef/*
	go get github.com/PRETgroup/goFB/goFB
	go build -o easy-rte-parser -i ./rteparser/main

example_ab5: example/ab5/ab5_main.c example/ab5/F_ab5Function.c
	gcc example/ab5/ab5_main.c example/ab5/F_ab5Function.c -o example_ab5

example/ab5/F_ab5Function.c: easy-rte-c easy-rte-parser example/ab5/ab5.erte
	./easy-rte-parser -i example/ab5 -o example/ab5
	./easy-rte-c -i example/ab5 -o example/ab5

example_ab5_vhdl: example/ab5/enforcer_ab5.vhdl

example/ab5/enforcer_ab5.vhdl: easy-rte-c easy-rte-parser example/ab5/ab5.erte
	./easy-rte-parser -i example/ab5 -o example/ab5
	./easy-rte-c -i example/ab5 -o example/ab5 -l=vhdl

example_robotable: example/robotable/robotable_main.c example/robotable/F_Robotable.c
	gcc example/robotable/robotable_main.c example/robotable/F_Robotable.c -o example_robotable

example/robotable/F_Robotable.c: easy-rte-c easy-rte-parser example/robotable/robotable.erte
	./easy-rte-parser -i example/robotable -o example/robotable
	./easy-rte-c -i example/robotable -o example/robotable

clean:
	rm -f easy-rte-c
	rm -f easy-rte-parser
	rm -f example_ab5
	rm -f example_robotable
	go get -u github.com/PRETgroup/goFB/goFB