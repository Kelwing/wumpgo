# Shard

A Shard represents a single shard, or connection to the gateway.  It maintains the connection to Discord, handles reconnects, and forwards dispatch events on to the desired destination.  It doesn't do anything fancy to route events to handlers, doesn't do any sort of caching.  What it does do is accept an implementation of the Dispatcher interface and calls its Dispatch function for each dispatch or Op 0 event.

wumpgo provides the following dispatcher implementation with a matching receiver implementation for acting upon those events.
- NOOP[^noop]
- Local
- NATS
- Redis

The Local, NATS, and Redis dispatchers each have a matching receiver that allow you to assign functions to handle Discord events coming from each of those dispatchers.

[^noop]: The NOOP dispatcher simply logs the event payload and throws it away.