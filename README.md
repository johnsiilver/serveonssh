# Serve On SSH

## Introduction

There is often a need to offer services for administrative purposes on servers or even for microservices that are running on a device. And while there are many way to secure them, one of the most convenient ways is to do it behind SSH.

This package provides a `Proxy` type and a `Dialer` that can be used with HTTP/gRPC/... packages to serve connections behind SSH. In addition, we use Unix sockets for serving traffic, so there are no exposed ports on the server side. 

This has several advantages for real world admin traffic:
- SSH is much easier to setup that standard AAA for web services
- SSH is probably already running
- You can easily block SSH to servers from the public without complex filters

So for your admin traffic, block SSH externally and now only your VPN clients and internal service can reach any administrative endpoints. These endpoints are now secured using the same technology trusted for your logins. SSH is already logging and can be scraped for bad actors.

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
