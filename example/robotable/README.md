# Robotable

Plant emits (and controller inputs)
* reqLocX
* reqLocY
* curLocX
* curLocY

Controller emits (and plant inputs)
* driveX
* driveY

The system is trying to keep a robot on a table with specified bounds. 
The controller is trying to get the robot to {reqLocX, reqLocY} from {curLocX, curLocY}.
It has outputs {driveX, driveY} to actuate the motor.

We ensure:
* OUPUT: ensure that the drive command does not exceed a safe value
* INPUT: ensure requested location is on the table
* INPUT: ensure sensed location is on the table
* OUTPUT: ensure that the current sense location augmented with the drive command does not push us off the table
					