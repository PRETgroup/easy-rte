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
    //randomise inputs
    inputs.A = nondet_1();

    //randomise enforcer state, i.e. clock values and position (excepting violation state)
    enf.v = nondet_2();
    enf._policy_AB5_state = nondet_3() % 2; //Here, "% 2" ensures that it is in any state _other_ than the violation state

    //run the enforcer (i.e. tell CBMC to check this out)
    ab5_run_via_enforcer(&enf, &inputs, &outputs);

    print_data(count, inputs, outputs);
    
}

void ab5_run(inputs_ab5_t *inputs, outputs_ab5_t *outputs) {
    //randomise controller

    outputs->B = nondet_4(); 
}

