# Serve On SSH

## Introduction

There is often a need to offer services for administrative purposes on servers or even for microservices that are running on a device. And while there are many way to secure them, one of the most convenient ways is to do it behind SSH.

This package provides a `Forwarder` type and a `Dialer` that can be used with HTTP or gRPC libraries to serve connections behind SSH. In addition, we use Unix sockets for serving traffic, so there are no exposed ports on the server side. 

## Examples

We offer two examples of this running:
```
example/
	http/
	grpc/
```
Inside these directories are README.md that will explain how you can run the demo and see the code behind each.

## As always

Have fun!
