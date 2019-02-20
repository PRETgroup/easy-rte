# pacemaker examples

Actual times:

| timer | ms  |
| ----- | ---:|
| AVI   | 300 |
| AEI   | 800 |
| PVARP |  50 |
| VRP   | 150 |
| LRI   | 950 |
| URI   | 900 |

//P1: AP and VP cannot happen simultaneously.
//P2: VS or VP must be true within AVI after an atrial event AS or AP.
//P3: AS or AP must be true within AEI after a ventricular event VS or VP.
//P4: After a ventricular event VS or VP, another ventricular event can happen only after URI.
//P5: After a ventricular event VS or VP, another ventricular event should happen within LRI.

AS/AP --AVI--> VS/VP --AEI--> AS/AP --AVI--> VS/VP

AEI + AVI is maximum time between VS/VP and VS/VP
LRI is also maximum time between VS/VP and VS/VP

According to expert, 
URI and LRI should be less than AVI+AEI, i.e.
AVI + AEI < URI
AVI + AEI < LRI

