# Elevatpr project report
## Documentation
Delayed light, persistant storage, timesync


![Communication diagram](communication_diagram.drawio.svg)

## Case studies of important descicions

A system that can be interacted with at multiple processes may suffer from the problem of knowing what happened first. In the case of the elevator project 
it occurs, for instance, when a order comes in in one process at merely the same time as it is removed by another process. This may result in a situation where the processes disagree upon the system state. Facing this problem we have come up with two reasonable solutions. One, in which sequence numbers are utilized and another that uses timestamp. Both approaches come with their set of new problems. 

When considering the sequence number solution, the problem of two processes claiming the same sequence number for an event, occurs. This is however not a problem when timestamping every event. That is, the instance of similarily timestamped events with nanosecond precision, is highly unlikely. And the resending of central state would in the event of this unlikelyhood resolve thtis issue at the next broadcast ~15 ms later. 
