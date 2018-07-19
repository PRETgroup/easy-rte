default: easy-rte-c easy-rte-parser

#convert C build instruction to C target
c_enf: $(PROJECT)

#convert verilog build instruction to verilog target
verilog_enf: $(PROJECT)_V

easy-rte-c: rtec/* rtedef/*
	go get github.com/PRETgroup/goFB/goFB
	go build -o easy-rte-c -i ./rtec/main

easy-rte-parser: rteparser/* rtedef/*
	go get github.com/PRETgroup/goFB/goFB
	go build -o easy-rte-parser -i ./rteparser/main

#convert $(PROJECT) into the C binary name
$(PROJECT): example_$(PROJECT)

#generate the C binary from the C sources
example_$(PROJECT): example/$(PROJECT)/*.c example/$(PROJECT)/*.c
	gcc example/$(PROJECT)/*.c -o example_$(PROJECT)

#generate the C sources from the erte files
example/$(PROJECT)/*.c: default example/$(PROJECT)/*.erte
	./easy-rte-parser -i example/$(PROJECT) -o example/$(PROJECT)
	./easy-rte-c -i example/$(PROJECT) -o example/$(PROJECT)

#convert $(PROJECT)_V into the verilog names
$(PROJECT)_V: example/$(PROJECT)/*.v

example/$(PROJECT)/*.v: default example/$(PROJECT)/*.erte
	./easy-rte-parser -i example/$(PROJECT) -o example/$(PROJECT)
	./easy-rte-c -i example/$(PROJECT) -o example/$(PROJECT) -l=verilog

clean:
	rm -f easy-rte-c
	rm -f easy-rte-parser
	rm -f example_*
	go get -u github.com/PRETgroup/goFB/goFB