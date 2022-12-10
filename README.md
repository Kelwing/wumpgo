# wumpgo

<p align="center">
<img src="assets/wumpgo_white.png" width=50% height=50%>
</p>

A work in progress Golang Discord library.

## Installation

```
go get -u wumpgo.dev/wumpgo@v0.0.2
```

---

## Components

Below is a brief overview of each component of the library.  Most components are designed to work as independently of other components as possible.  For example: if you simply need a REST library to use in an API you're writing, you can use the rest component by itself.

### Objects

The objects package is the base of the library.  Pretty much all other components depend on it.  It contains the Discord API model definitions.

### Rest

The rest package contains a simple Discord REST API client library.  It is designed to be very tunable so it will fit many different usecases.  We provide defaults for the ratelimiter, cache, and proxy, all of which are omitted by default.  The default implementations follow an interface contract, so feel free to bring you own as well.

### Gateway

In the Gateway package you will find a simple and pluggable gateway implementation.  It uses a dispatcher and receiver model to control how events are delivered in your application.  Dispatchers cause the gateway to dispatch the events to your chosen delivery system for delivery to a receiver.  At the simplest you can use a local dispatcher and receiver for behavior similar to every other library you've used.  On the complex side you can use a message queue or pubsub system for a large, scalable implementation.  We provide a few implementations for common use-cases, but feel free to bring your own. 

### Interactions