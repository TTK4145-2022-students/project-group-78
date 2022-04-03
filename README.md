# Elevator project group 78
Peculiar design choises listed below
- Use of timestamps
    - Hall orders, and when they were served, are timestamped. This makes it easy to merge incomming state (both in normal operation and during reconnects): just keep the newest ones.
    - Elevator state updates are also timestamped.

- No explicit error handling or detection
    - Network disconnects, packet losses and motor disconnects are not explicitly detected nor handled. Door obstruction is only used to ensure correct door operation.

- Peer-to-peer continous broadcasting, delayed lights and persistant storage - no acks
    - The central state (containing orders, timestamps, and elevator state (basically the input for hall_request_assigner)) is broadcast and stored persistantly (LINK SKV) regularly.
    - Because of persistant storage, no orders are lost. This means that the light-button contract is fullfilled. However, 
    - we are also delaying the button light updates by 20 times the broadcast interval, so that sufficient attemps have been made to broadcast the new order.

- Continous order assigments with elevator timeouts
    - Becuase elevator state updates are timestamped, it is easy to detect that an elevator that has not moved when it should have. Whether this is due to door obstruction, network disconnect, packet loss, software crash or motor disconnect, we don't care.
    - Orders are continously (re)assigned by each elevator. Other assumed faulty elevators (we never check whether we ourselves are faulty) are removed from the assignments. The assigments are not broadcast, because all elevators should agree on which takes which in normal operation.

# Modules


Use ```timedatectl``` to check if NTP is active