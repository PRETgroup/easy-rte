default: easy-rte-c easy-rte-parser

easy-rte-c: rtec/* rtedef/*
	go get github.com/PRETgroup/goFB/goFB
	go build -o easy-rte-c -i ./rtec/main

easy-rte-parser: rteparser/* rtedef/*
	go get github.com/PRETgroup/goFB/goFB
	go build -o easy-rte-parser -i ./rteparser/main

example_ab5: example/ab5/ab5_main.c example/ab5/F_ab5Function.c
	gcc example/ab5/ab5_main.c example/ab5/F_ab5Function.c -o example_ab5

#example/%/F_%.c: easy-rte-c easy-rte-parser example/%/%.erte

example/ab5/F_ab5Function.c: easy-rte-c easy-rte-parser example/ab5/ab5.erte
	./easy-rte-parser -i example/ab5 -o example/ab5
	./easy-rte-c -i example/ab5 -o example/ab5

example_ab5_verilog: example/ab5/enforcer_ab5.v

example/ab5/enforcer_ab5.v: easy-rte-c easy-rte-parser example/ab5/ab5.erte
	./easy-rte-parser -i example/ab5 -o example/ab5
	./easy-rte-c -i example/ab5 -o example/ab5 -l=verilog

example_ab5seconds: example/ab5seconds/ab5_main.c example/ab5seconds/F_ab5Function.c
	gcc example/ab5seconds/ab5_main.c example/ab5seconds/F_ab5Function.c -o example_ab5

example/ab5seconds/F_ab5Function.c: easy-rte-c easy-rte-parser example/ab5seconds/ab5seconds.erte
	./easy-rte-parser -i example/ab5seconds -o example/ab5seconds
	./easy-rte-c -i example/ab5seconds -o example/ab5seconds

example_ab5seconds_verilog: example/ab5seconds/enforcer_ab5.v

example/ab5seconds/enforcer_ab5.v: easy-rte-c easy-rte-parser example/ab5seconds/ab5seconds.erte
	./easy-rte-parser -i example/ab5seconds -o example/ab5seconds
	./easy-rte-c -i example/ab5seconds -o example/ab5seconds -l=verilog

example_robotable: example/robotable/robotable_main.c example/robotable/F_Robotable.c
	gcc example/robotable/robotable_main.c example/robotable/F_Robotable.c -o example_robotable

example/robotable/F_Robotable.c: easy-rte-c easy-rte-parser example/robotable/robotable.erte
	./easy-rte-parser -i example/robotable -o example/robotable
	./easy-rte-c -i example/robotable -o example/robotable

example_robotable_verilog: example/robotable/enforcer_robotable.v

example/robotable/enforcer_robotable.v: easy-rte-c easy-rte-parser example/robotable/robotable.erte
	./easy-rte-parser -i example/robotable -o example/robotable
	./easy-rte-c -i example/robotable -o example/robotable -l=verilog

clean:
	rm -f easy-rte-c
	rm -f easy-rte-parser
	rm -f example_ab5
	rm -f example_robotable
	go get -u github.com/PRETgroup/goFB/goFB