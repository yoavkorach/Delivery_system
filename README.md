This is a basic Delivery system project using go, it supports placing a delivery and checking its status. 
Every delivery time is calculated by distance (in km) generated randomly and assuming a delivery guy travels one km per second. (completely arbitrary can be changed by changing the waiting time)
The number of delivery guys is constant, the program uses goncurrency in order to update time until delivery. And support new assign new orders if all delivery guys are occupied. 
Complexity:
Since number of delivery guys is constant most functions operate in O(1), exculding the assignOrderToPerson function which is dependant on number of unassigned orders n, so O(n).

