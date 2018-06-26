#include "F_AB5Function.h"
#include <stdio.h>

void print_data(uint32_t count, inputs_AB5Function_t inputs, outputs_AB5Function_t outputs) {
    printf("Tick %7d: A:%d, B:%d\r\n", count, inputs.A, outputs.B);
}

int main() {
    enforcervars_AB5Function_t enf;
    inputs_AB5Function_t inputs;
    outputs_AB5Function_t outputs;
    
    AB5Function_init_all_vars(&enf, &inputs, &outputs);

    uint32_t count = 0;
    while(count++ < 30) {
        if(count % 10 == 0) {
            inputs.A = true;
        } else {
            inputs.A = false;
        }

        AB5Function_run_via_enforcer(&enf, inputs, &outputs);

        print_data(count, inputs, outputs);
    }
}

void AB5Function_run(inputs_AB5Function_t inputs, outputs_AB5Function_t *outputs) {
    //do nothing
}

