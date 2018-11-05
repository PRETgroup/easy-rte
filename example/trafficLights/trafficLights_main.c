#include "F_trafficLights.h"
#include <stdio.h>

void print_data(uint32_t count, inputs_trafficLights_t* inputs, outputs_trafficLights_t* outputs) {
    char NS = 'R';
    char EW = 'R';
    if(outputs->NS_G) {
        NS = 'G';
    } else if (outputs->NS_Y) {
        NS = 'Y';
    }
    if(outputs->EW_G) {
        EW = 'G';
    } else if (outputs->EW_Y) {
        EW = 'Y';
    }

    printf("\r\n%i:\r\n", count);
    printf("  |   | \r\n");
    printf("  | %c | \r\n", NS);
    printf("--+   +--\r\n");
    printf(" %c     %c \r\n", EW, EW);
    printf("--+   +--\r\n");
    printf("  | %c | \r\n", NS);
    printf("  |   | \r\n");
    
}

int main() {
	enforcervars_trafficLights_t enf_trafficLights;
    inputs_trafficLights_t inputs_trafficLights;
    outputs_trafficLights_t outputs_trafficLights;
    
    trafficLights_init_all_vars(&enf_trafficLights, &inputs_trafficLights, &outputs_trafficLights);

    int count;
    while(count++ < 100) {

        trafficLights_run_via_enforcer(&enf_trafficLights, &inputs_trafficLights, &outputs_trafficLights);

        print_data(count, &inputs_trafficLights, &outputs_trafficLights);
    }
}

void trafficLights_run(inputs_trafficLights_t *inputs, outputs_trafficLights_t *outputs) {
    //do nothing

    outputs->NS_G = 0;
	outputs->NS_Y = 1;
	outputs->NS_R = 0;
	outputs->EW_G = 0;
	outputs->EW_Y = 1;
	outputs->EW_R = 0;
	 
}


