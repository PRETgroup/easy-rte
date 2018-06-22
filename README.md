# easy-rte

## About
This project provides an easy-to-use implementation of _bi-directional Runtime Enforcement_, based on the semantics originally presented in [Runtime Enforcement of Cyber-Physical Systems](https://dl-acm-org.ezproxy.auckland.ac.nz/citation.cfm?id=3126500) (ACM Transactions on Embedded Computing Systems (TECS) 2017).

While the original implementation was restricted to simple _boolean_ arguments only, and was implemented in Python for use 
with SCCharts, this project presents a more generalised any-type enforcement system, which can be used with any C project. 
_easy-rte_ was ported from [goFB](https://github.com/PRETgroup/goFB), which implemented the semantics in this way, but restricted them for use with IEC 61499 function blocks.

## What is Runtime Enforcement?

TODO

## Example of Use

Imagine a function which inputs boolean `A` and outputs boolean `B`. 
Let's give it the following properties:
1. A and B cannot happen simultaneously.
2. A and B alternate starting with an A. 
3. B should be true within 5 ticks after an occurance of A.

We can present this as the following _easy-rte_ policy format:
```
function AB5Function;
interface of AB5Function {
	in bool A;  //in here means that they're going from PLANT to CONTROLLER
	out bool B; //out here means that they're going from CONTROLLER to PLANT
}

policy AB5 of AB5Function {
	internals {
		dtimer v;
	}

	states {
		s0 {	//first state is initial, and represents "We're waiting for an A"
			-> s0 on (!A and !B): v := 0;			//if we receive neither A nor B, do nothing
			-> s1 on (A and !B): v := 0;			//if we receive an A only, head to state s1
			-> violation on ((!A and B) or (A and B));	//if we receive a B, or an A and a B (i.e. if we receive a B) then VIOLATION
		}

		s1 {	//s1 is "we're waiting for a B, and it needs to get here within 5 ticks"
			-> s1 on (!A and !B and v < 5);				//if we receive nothing, and we aren't over-time, then we do nothing
			-> s0 on (!A and B);					//if we receive a B only, head to state s0
			-> violation on ((v >= 5) or (A and B) or (A and !B));	//if we go overtime, or we receive another A, then VIOLATION
		}
	}
}
```

As can be seen, this can be thought of as a simple mealy finite state machine, which provides the rules for operation and when a transition to violation might occur.
However, we don't want these violations to occur. 
So, we shall convert the _policy_ to an _enforcer_, which is capable of modifying the I/O such that a violation cannot occur.

In _easy-rte_, this is done in two steps. Firstly, we convert the _erte_ file into an equivalent policy XML file (which makes it easier to understand, and allows portability between tools).
* `./easy-rte-parser -i example -o example`

Then, we convert this policy XML file into an _enforcer_, which is written in C. The process for doing this is provided in the linked paper, in the About section of this README.
* `./easy-rte-c -i example -o example`

## Build instructions

Download and install the latest version of [Go](https://golang.org/doc/install).

Then, download this repository, and run `make`.
