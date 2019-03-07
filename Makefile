.PHONY: default c_enf verilog_enf $(PROJECT) c_build
#.PRECIOUS: %.xml 

# run this makefile with the following options
# make [c_enf] [c_build] [run_(c/e)bmc] PROJECT=XXXXX FILE=YYYYY
#   PROJECT = name of project directory
#   FILE    = name of file within project directory (default = PROJECT, e.g. example/ab5/ab5.whatever)
#
#   c_enf: make a C enforcer for the project
#   c_build: compile the C enforcer with a main file (this will need to be provided manually)
#   run_cbmc: check the compiled C enforcer to ensure correctness
#
# make [verilog_enf] [run_ebmc] PROJECT=XXXXX
#   verilog_enf: make a Verilog enforcer for the project
#   run_ebmc: check the compiled Verilog enforcer to ensure correctness

FILE ?= $(PROJECT)
PARSEARGS ?=

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

run_cbmc: default 
	cbmc example/$(PROJECT)/cbmc_main_$(PROJECT).c example/$(PROJECT)/F_$(PROJECT).c

run_ebmc: default 
	#$(foreach file,$(wildcard example/$(PROJECT)/*.sv), time --format="took %E" ebmc $(file) --k-induction --trace --top F_combinatorialVerilog_$(word 3,$(subst _, ,$(basename $(notdir $(file)))));)
	time --format="took %E" ebmc example/$(PROJECT)/test_F_$(FILE).sv --k-induction --trace --module F_combinatorialVerilog_$(FILE)
	#ebmc $^ --k-induction --trace

#convert $(PROJECT) into the C binary name
$(PROJECT): ./example/$(PROJECT)/$(FILE).c

#generate the C sources from the erte files
%.c: %.xml
	./easy-rte-c -i $^ -o example/$(PROJECT)

#convert $(PROJECT)_V into the verilog names
$(PROJECT)_V: ./example/$(PROJECT)/$(FILE).sv

#generate the xml from the erte files
%.xml: %.erte
	./easy-rte-parser $(PARSEARGS) -i $^ -o $@

#generate the Verilog sources from the xml files
%.sv: %.xml
	./easy-rte-c -i $^ -o example/$(PROJECT) -l=verilog

#Bonus: C compilation: convert $(PROJECT) into the C binary name
c_build: example_$(PROJECT)

#generate the C binary from the C sources
example_$(PROJECT): example/$(PROJECT)/$(PROJECT)_main.c example/$(PROJECT)/F_$(PROJECT).c
	gcc example/$(PROJECT)/$(PROJECT)_main.c example/$(PROJECT)/F_$(PROJECT).c -o example_$(PROJECT)

#Bonus: C assembly
c_asm: example/$(PROJECT)/F_$(PROJECT).c
	gcc -S example/$(PROJECT)/F_$(PROJECT).c -o example/$(PROJECT)/F_$(PROJECT).s

clean: clean_examples
	rm -f easy-rte-c
	rm -f easy-rte-parser
	go get -u github.com/PRETgroup/goFB/goFB

clean_examples:
	rm -f example_*
	rm -f ./example/*/F_*
	rm -f ./example/*/*.h
	rm -f ./example/*/*.v
	rm -f ./example/*/*.sv
	rm -f ./example/*/*.xml