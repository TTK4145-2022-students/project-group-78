# Elevator project report
# Delayed light, persistant storage, timesync
## Documentation
Module communication diagram shwoing what information is exchanged between modules.
![Communication diagram](communication_diagram.drawio.svg)

Peculiar design choises listed below
### Use of timestamps
- Hall orders, and when they were served, are timestamped. This makes it easy to merge incomming state (both in normal operation and during reconnects): just keep the newest ones.
- Elevator state updates are also timestamped.

### No explicit error detection
- Network disconnects, packet losses and motor disconnects are not explicitly detected.

### Peer-to-peer continuous broadcasting, delayed lights and persistant storage - no acks
- The central state (containing orders, timestamps, and elevator state (basically the input for hall_request_assigner)) is broadcast and stored persistantly regularly.
- Because of persistant storage, no orders are lost. This means that the light-button contract is fullfilled. However, 
- we are also delaying the button light updates by 20 times the broadcast interval, so that sufficient attemps have been made to broadcast the new order.

### Continuous order assigments with elevator timeouts
- Becuase elevator state updates are timestamped, it is easy to detect an elevator which should have moved or broadcast something. Whether this is due to door obstruction, network disconnect, packet loss, software crash or motor disconnect, we do not care.
- Orders are continuously (re)assigned by each elevator. Other assumed faulty elevators (we do not care whether we ourselves are faulty) are removed from the assignments. The assigments are not broadcast, because all elevators should agree on which takes which in normal operation.

### Considerations regarding timestamps
While nodes are connected together and to the internet, their clocks are synced up (NTP). There would be no problem syncing clocks without a internet connection, but this has not been implemented because nodes are also connected to the internet when they are connected together in this course. However, if nodes are disconnected, their clocks may fall out of sync. This can be a problem in some very specific scenarios, but as this requires very long disconnect periods, this is deemed unproblematic in this context.

## Case studies of important descicions

### Ulrik
A system that can be interacted with at multiple processes may suffer from the problem of knowing what happened first. In the case of the elevator project 
it occurs, for instance, when a order comes in in one process at merely the same time as it is removed by another process. This may result in a situation where the processes disagree upon the system state. Facing this problem we have come up with two reasonable solutions. One, in which sequence numbers are utilized and another that uses timestamp. Both approaches come with their set of new problems. 

When considering the sequence number solution, the problem of two processes claiming the same sequence number for an event, occurs. This is however not at problem when timestamping every event. That is, the instance of similarily timestamped events with nanosecond precision, is highly unlikely. And the resending of central state would in the event of this unlikelyhood resolve thtis issue at the next broadcast ~15 ms later. 

### The button light contract
The spec specifies that when an elevator button is pushed it should always try to send out the event to the other elevators before the actual light is turned on. Since this system does not operate in a way that events specifically is pushed out to the rest of the network whenever they happen, this is an issue with our design concerning this formulation. Since our design always sends out the current state of the elevator with a fixed time-interval we chose to solve this by delaying the action of turnng on the light just enough so we know that the light has been broadcasted. The design also stores this in a persistant storage, ensuring that the call never will be lost anyway. 

### Error detection and handling
As hall_request_assigner is used, error handling is as simple as excluding erroneous elevators when assiging. However, the task of detecting that an elevator is erroneous was an interesting design descicion. We found two alternatives: explicit and implicit error detection.

#### Explicit error detection
Detecting obstruction is trivial. Detecting motor stop requires a timer. Detecting crash/disconnect reqiures heartbeats and a timer. This is an explicit error detection.

#### Implicit error detection
Elevator state changes are timestamped. If we do not have a recent state change for an elevator which should be moving, we know that it is erroneous.
(An obstructed elevator will not change its state, neither will an elevator with motor stop. A crashed or disconnected node while not manage to send the change. Therefore we can rely on this mecanism)

#### Descicion
We went with the implicit one, beacuse it is simpler, and therefor more difficult to get wrong. It is also more robust. However, there are downsides to this descicion. If the spec changed to include some different behaviour for different errors, we would maybe have to restructure completly. Also this solution is slower than the explicit one, because we wait for a timeout on the obstruction, in stead of immideatly deem that the obstructed elevator is erroneous.
