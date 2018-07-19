# IDMT example

From 
Towards Threat of Implementation Attacks on
Substation Security: Case Study on Fault Detection
and Isolation
by
Anupam Chattopadhyay, Senior Member, IEEE, Abhisek Ukil, Senior Member, IEEE,
Dirmanto Jap, Member, IEEE, and Shivam Bhasin, Member, IEEE.

The standard fault detection for overcurrent faults is implemented via the
'Inverse Definite Minimum Time (IDMT)' curve, which is covered by the
IEC 60255 standard. 

The equation for this curve is as follows:

`t = (K*B) / ((I / Iset)^a - 1)`
Here, 
* `t` = time
* `K` = time multiplier
* `I` = measured current
* `Iset` = nominal current
* `a` = calibration parameter
* `B` = Calibration parameter

Utilities typically change the values of `a`, `B` to change the slope of the curve,
to configure the needs of their protection operations as necessary.

For instance, here are some common values:

| Type of Curve     |  `a`  |  `B`  |
|-------------------|-------|-------|
| Normal Inverse    | 0.02  | 0.14  |
| Very Inverse      | 1.0   | 13.5  |
| Extremely Inverse | 2.0   | 80.0  |
| Long-time Inverse | 1.0   | 120.0 |

For this case study, we will use "Very Inverse" type curve, with `a`=1.0 and `B`=13.5.

Using the IDMT curve equation, and the time multiplier `K`=0.1, for an overcurrent magnitude of 10, i.e. 10 times fault, we get the operating time as 0.15 seconds. However, for the same setting, if the overcurrent magnitude is lower, e.g. 5 times, the operating time is longer (0.34 seconds).

Attackers might and try cause a controller to fail to meet the overcurrent deadline, thus causing damaging errors.

We ensure:
* If an overcurrent occurs, the relay is disconnected before the time becomes unsafe.