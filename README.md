One elevator or node consists of one ditributor and several components. The distributor receives events from producing components and distributes these to consuming components and to other distributors. This way events are exchanged between components internally in one node, and externally between nodes.

A producer has an Out chan and a consumer has an In chan. A componenet might be both a producer and a consumer. Each component initializes and then registers its channels with the distributor.

All events are defined in the events package. Because one does not care where an event comes from and it would be non trivial where to place the event interface