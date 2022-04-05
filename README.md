Please paste this readme into a markdown renderer if you are reading this as plain text (https://dillinger.io/). 
# Elevator project group 78
Peculiar design choises listed below
### Use of timestamps
- Hall orders, and when they were served, are timestamped. This makes it easy to merge incomming state (both in normal operation and during reconnects): just keep the newest ones.
- Elevator state updates are also timestamped.

### No explicit error handling or detection
- Network disconnects, packet losses and motor disconnects are not explicitly detected nor handled. Door obstruction is only used to ensure correct door operation.

### Peer-to-peer continuous broadcasting, delayed lights and persistant storage - no acks
- The central state (containing orders, timestamps, and elevator state (basically the input for hall_request_assigner)) is broadcast and stored persistantly regularly.
- Because of persistant storage, no orders are lost. This means that the light-button contract is fullfilled. However, 
- we are also delaying the button light updates by 20 times the broadcast interval, so that sufficient attemps have been made to broadcast the new order.

### Continuous order assigments with elevator timeouts
- Becuase elevator state updates are timestamped, it is easy to detect an elevator which should have moved or broadcast something. Whether this is due to door obstruction, network disconnect, packet loss, software crash or motor disconnect, we do not care.
- Orders are continuously (re)assigned by each elevator. Other assumed faulty elevators (we do not care whether we ourselves are faulty) are removed from the assignments. The assigments are not broadcast, because all elevators should agree on which takes which in normal operation.

## Interesting modules
|Name|Input|Ouput|Side effects|
|--|--|--|--|
|elevator|assigned orders|served order, state (behaviour, direction, floor)|door (submodule), motor, floor indicator|
|assigner|central state|assigned orders|none|
|lights|central state|none|button lights|
|door|open|closed|door|

### Other modules
|Name|Comment|
|--|--|
|central|definition of the central state type|
|config|global constant parameters|
|[skv](https://github.com/rapidloop/skv)|persistant storage (simple key-value store)|
|[driver-go-group-78](https://github.com/TTK4145-2022-students/driver-go-group-78)||
|[Network-go-group-78](https://github.com/TTK4145-2022-students/Network-go-group-78)||

## Considerations regarding timestamps
While nodes are connected together and to the internet, their clocks are synced up (NTP). There would be no problem syncing clocks without a internet connection, but this has not been implemented because nodes are also connected to the internet when they are connected together in this course. However, if nodes are disconnected, their clocks may fall out of sync. This can be a problem in some very specific scenarios, but as this requires very long disconnect periods, this is deemed unproblematic in this context.
