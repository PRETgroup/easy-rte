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

    //a nondet_xxxxx function name tells cbmc that it could be anything
    inputs.A = nondet_bool();
    enf.v = nondet_uint64_t();
    //enf._policy_AB5_state = nondet_uint64sd_t(); //sanity check: if state can be anything, it could be a violation state

    ab5seconds_run_via_enforcer(&enf, &inputs, &outputs);

}

void ab5seconds_run(inputs_ab5seconds_t *inputs, outputs_ab5seconds_t *outputs) {
    outputs->B = nondet_bool();

}

