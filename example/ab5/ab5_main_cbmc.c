#include "F_ab5.h"
#include <stdio.h>
#include <stdint.h>

void print_data(uint32_t count, inputs_ab5_t inputs, outputs_ab5_t outputs) {
    printf("Tick %7d: A:%d, B:%d\r\n", count, inputs.A, outputs.B);
}

int main() {
    enforcervars_ab5_t enf;
    inputs_ab5_t inputs;
    outputs_ab5_t outputs;
    int count = 0;
    
    //set values to known state
    ab5_init_all_vars(&enf, &inputs, &outputs);
   
    //introduce nondeterminism
    //a nondet_xxxxx function name tells cbmc that it could be anything, but must be unique
    inputs.A = nondet_1();
    enf.v = nondet_2();

    //sanity check: if state can be anything, it could be a violation state
    //uncomment this to cause a verification failure
    //enf._policy_AB5_state = nondet_uint64sd_t(); 

    //run the enforcer (i.e. tell CBMC to check this out)
    ab5_run_via_enforcer(&enf, &inputs, &outputs);

    print_data(count, inputs, outputs);
    
}

void ab5_run(inputs_ab5_t *inputs, outputs_ab5_t *outputs) {
    //do nothing
}

