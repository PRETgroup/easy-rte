.PHONY: default c_enf verilog_enf $(PROJECT) c_build

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

run_cbmc: default c_enf
	cbmc example/$(PROJECT)/cbmc_main_$(PROJECT).c example/$(PROJECT)/F_$(PROJECT).c

run_ebmc: default verilog_enf
	ebmc example/$(PROJECT)/ebmc_F_$(PROJECT).sv

#convert $(PROJECT) into the C binary name
$(PROJECT): example/$(PROJECT)/F_$(PROJECT).c

#generate the C sources from the erte files
example/$(PROJECT)/F_$(PROJECT).c: easy-rte-c easy-rte-parser example/$(PROJECT)/$(PROJECT).erte
	./easy-rte-parser -i example/$(PROJECT) -o example/$(PROJECT)
	./easy-rte-c -i example/$(PROJECT) -o example/$(PROJECT)

#convert $(PROJECT)_V into the verilog names
$(PROJECT)_V: example/$(PROJECT)/*.v

#generate the Verilog sources from the erte files
example/$(PROJECT)/*.v: default example/$(PROJECT)/*.erte
	./easy-rte-parser -i example/$(PROJECT) -o example/$(PROJECT)
	./easy-rte-c -i example/$(PROJECT) -o example/$(PROJECT) -l=verilog

#Bonus: C compilation: convert $(PROJECT) into the C binary name
c_build: example_$(PROJECT)

#generate the C binary from the C sources
example_$(PROJECT): example/$(PROJECT)/$(PROJECT)_main.c example/$(PROJECT)/F_$(PROJECT).c
	gcc example/$(PROJECT)/$(PROJECT)_main.c example/$(PROJECT)/F_$(PROJECT).c -o example_$(PROJECT)

clean:
	rm -f easy-rte-c
	rm -f easy-rte-parser
	rm -f example_*
	go get -u github.com/PRETgroup/goFB/goFB