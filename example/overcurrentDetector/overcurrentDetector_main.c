#include "F_overcurrentDetector.h"
#include <stdio.h>

void print_data(uint64_t t, inputs_overcurrentDetector_t inputs, outputs_overcurrentDetector_t outputs) {
    printf("Time: %llu ms\tCurrent:%d mA\tRelay: %s\r\n", t/1000, inputs.I_mA, outputs.relay_en ? "enabled":"DISABLED");
}

int main() {
    enforcervars_overcurrentDetector_t enf;
    inputs_overcurrentDetector_t inputs;
    outputs_overcurrentDetector_t outputs;
    
    overcurrentDetector_init_all_vars(&enf, &inputs, &outputs);
    inputs.I_mA = 10000;
    uint64_t t = 0;
    while(t+=1000 < 300000) {
        if(t == 10000) {
            inputs.I_mA = 10000;
        } 

        overcurrentDetector_run_via_enforcer(&enf, &inputs, &outputs);

        print_data(t, inputs, outputs);
    }
}

void overcurrentDetector_run(inputs_overcurrentDetector_t *inputs, outputs_overcurrentDetector_t *outputs) {
    //do nothing
}

