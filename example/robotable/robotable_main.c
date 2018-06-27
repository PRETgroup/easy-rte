#include "F_Robotable.h"
#include <stdio.h>

void print_data(uint32_t count, inputs_Robotable_t inputs, outputs_Robotable_t outputs) {
    
}

int main() {
    enforcervars_Robotable_t enf;
    inputs_Robotable_t inputs;
    outputs_Robotable_t outputs;
    
    Robotable_init_all_vars(&enf, &inputs, &outputs);

    uint32_t count = 0;
    while(count++ < 30) {
        // if(count % 10 == 0) {
        //     inputs.A = true;
        // } else {
        //     inputs.A = false;
        // }

        Robotable_run_via_enforcer(&enf, inputs, &outputs);

        print_data(count, inputs, outputs);
    }
}

void Robotable_run(inputs_Robotable_t inputs, outputs_Robotable_t *outputs) {
    //do nothing
}

