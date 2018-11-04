#include "F_ab5seconds.h"
#include <stdio.h>

void print_data(uint32_t count, inputs_ab5seconds_t inputs, outputs_ab5seconds_t outputs) {
    printf("Tick %7d: A:%d, B:%d\r\n", count, inputs.A, outputs.B);
}

int main() {
    enforcervars_ab5seconds_t enf;
    inputs_ab5seconds_t inputs;
    outputs_ab5seconds_t outputs;
    
    ab5seconds_init_all_vars(&enf, &inputs, &outputs);

    uint32_t count = 0;
    while(count++ < 50) {
        if(count % 10 == 0) {
            inputs.A = true;
        } else {
            inputs.A = false;
        }

        ab5seconds_run_via_enforcer(&enf, &inputs, &outputs);

        print_data(count, inputs, outputs);
    }
}

void ab5seconds_run(inputs_ab5seconds_t *inputs, outputs_ab5seconds_t *outputs) {
    //do nothing
}

