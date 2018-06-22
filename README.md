# easy-rte

WIP

This is a re-implementation of the enforcer semantics from goFB, such that they work with normal C functions.


Example:

`
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
		s0 {														//first state is initial, and represents "We're waiting for an A"
			-> s0 on (!A and !B): v := 0;							//if we receive neither A nor B, do nothing
			-> s1 on (A and !B): v := 0;							//if we receive an A only, head to state s1
			-> violation on ((!A and B) or (A and B));				//if we receive a B, or an A and a B (i.e. if we receive a B) then VIOLATION
		}

		s1 {														//s1 is "we're waiting for a B, and it needs to get here within 5 ticks"
			-> s1 on (!A and !B and v < 5);							//if we receive nothing, and we aren't over-time, then we do nothing
			-> s0 on (!A and B);									//if we receive a B only, head to state s0
			-> violation on ((v >= 5) or (A and B) or (A and !B));	//if we go overtime, or we receive another A, then VIOLATION
		}
	}
}
`