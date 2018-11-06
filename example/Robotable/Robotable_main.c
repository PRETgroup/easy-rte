#include "F_Robotable.h"
#include <stdio.h>

#define XMIN -10
#define XMAX 10
#define YMIN -10
#define YMAX 10

#define SPEED 2 //even though we set the speed to be >1, it will be clamped

void print_data(uint32_t count, inputs_Robotable_t* inputs, outputs_Robotable_t* outputs) {
    //print top border
    printf("#");
    for(int x = XMIN; x <= XMAX; x++) {
        printf("# ");
    }
    printf("#\r\n");

    //print side borders and contents
    for(int y = YMIN; y <= YMAX; y++) {
        printf("#");
        for(int x = XMIN; x <= XMAX; x++) {
            if(y == inputs->curLocY && x == inputs->curLocX) {
                printf("@");
            } else if (y == inputs->reqLocY && x == inputs->reqLocX) {
                printf("*");
            } else {
                printf(" ");
            }

            printf(" ");
        }
        printf("#\r\n");
    }

    //print bottom border
    printf("#");
    for(int x = XMIN; x <= XMAX; x++) {
        printf("# ");
    }
    printf("#\r\n");
    printf("Current: %i,%i\r\n", inputs->curLocX, inputs->curLocY);
}

int main() {
    enforcervars_Robotable_t enf;
    inputs_Robotable_t inputs;
    outputs_Robotable_t outputs;
    
    Robotable_init_all_vars(&enf, &inputs, &outputs);

    inputs.curLocX = -8;
    inputs.curLocY = 0;
    inputs.reqLocX = 12;
    inputs.reqLocY = -12;

    uint32_t count = 0;
    while(count++ < 20) {

        Robotable_run_via_enforcer(&enf, &inputs, &outputs);

        print_data(count, &inputs, &outputs);

        inputs.curLocX += outputs.driveX;
        inputs.curLocY += outputs.driveY;
    }
}

void Robotable_run(inputs_Robotable_t* inputs, outputs_Robotable_t *outputs) {
    if(inputs->reqLocX < inputs->curLocX) {
        outputs->driveX = -SPEED; 
    } else if(inputs->reqLocX > inputs->curLocX) {
        outputs->driveX = SPEED;
    } else {
        outputs->driveX = 0;
    }

    if(inputs->reqLocY < inputs->curLocY) {
        outputs->driveY = -SPEED;
    } else if(inputs->reqLocY > inputs->curLocY) {
        outputs->driveY = SPEED;
    } else {
        outputs->driveY = 0;
    }
}

