# Elevator project report
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

### System state agreement
A system that can be interacted with at multiple processes may suffer from the problem of not knowing what happened first. In the case of the elevator project it occurs, for instance, when an order comes in in one process at the same time as it is removed by another process. This may result in a situation where the processes disagree upon the system state. Facing this problem we have come up with two reasonable solutions. One, in which sequence numbers are utilized and another that uses timestamp.

When considering the sequence number solution, the problem of two processes claiming the same sequence number for an event, occurs. This is however not a problem when timestamping every event. That is, the instance of identically timestamped events with nanosecond precision, is highly unlikely. And the resending of central state would in the event of this unlikelihood resolve thtis issue at the next broadcast ~15 ms later. 
As the timestamp solution is clearly easier to implement, we went for that solution. See the paragraph "considerations regarding timestamps" above.

### The button light contract
The spec specifies that when an elevator button is pushed it should always try to send out the event to the other elevators before the actual light is turned on. 

#### The ack way
This issue can be solved by sending the events out on the network and wait for an acknowlagement from the other elevators before turning on the light. The downside of using acks is that the lights will be dependent on the acks. Mixing these two modules together is complicating the system more than necessary.

#### Persistant storage
To satisfy the spec, one can resolve the issue without acknowlagements, but with a persistant storage and a small delay. In this case the light is turned on when the button is pushed and the action is stored in a persistant storage. This means that even if the system shuts down, the order is not lost because when booted up again it will load in the events stored in the persistant storage.

#### Desicion
Since our system does not handle any acks, the issue is handled by using persistant storage. To satisfy the spec completly the system is also adding a small delay before turning on the light, so that it is certain that the button event also has been sent out on the network. After the event has been stored in the persistant storage and sent out, the light is turned on, and it cannot be lost.

### Error detection and handling
As hall_request_assigner is used, error handling is as simple as excluding erroneous elevators when assiging. However, the task of detecting that an elevator is erroneous was an interesting design descicion. We found two alternatives: explicit and implicit error detection.

#### Explicit error detection
Detecting obstruction is trivial. Detecting motor stop requires a timer. Detecting crash/disconnect reqiures heartbeats and a timer. This is an explicit error detection.

#### Implicit error detection
Elevator state changes are timestamped. If we do not have a recent state change for an elevator which should be moving, we know that it is erroneous.
(An obstructed elevator will not change its state, neither will an elevator with motor stop. A crashed or disconnected node while not manage to send the change. Therefore we can rely on this mecanism)

#### Descicion
We went with the implicit one, beacuse it is simpler, and therefor more difficult to get wrong. It is also more robust. However, there are downsides to this descicion. If the spec changed to include some different behaviour for different errors, we would maybe have to restructure completly. Also this solution is slower than the explicit one, because we wait for a timeout on the obstruction, in stead of immediately deem that the obstructed elevator is erroneous.

## Lesson learned - pure functions and immutability
Thinking in terms of pure functions and immutability has proven itself in this project. Writing pure functions forces in many ways better code because it forces seperation of concern, and avoiding functions that does it all. It also makes code testable. Doing this from the start would have made things go alot faster and smoother.

Although, the code would be better from the start if design was done with a larger emphasis on pure functions, an outright mistake done during development was not thinking in terms of immutability, especally when writing concurrent code. Before ```CentralState``` was made immutable, weird stuff would happen when it was mutated, because at least two threads accesses it at a time.
