# Domain sockets over SSH gRPC demo

This shows example code on how you could run a gRPC server over SSH using domain sockets. While this wouldn't be useful for user traffic, this can be quite useful for admin traffic to your device such as gathering metrics or other sensitive information.

## Running

You can run this demo over you local SSH server if you don't want to do a two box setup. You also can use the SSH agent, a key file or password file.

Running the server like so:

```bash
cd service
go run service.go --socket=./socket.sock
```

Now to run the client, which will do 100 queries against the server. For the example this will use the password file "./pass", which is a plain text file holding the user password. 

```bash
cd app
go run app.go --endpoint=127.0.0.1:22 --socket=[path to this directory]/examples/http/service/socket.sock --pass=pass
```

That's it, you should see output:
```
2021/11/01 22:51:28 app.go:83: attempt(82) was successful
2021/11/01 22:51:28 app.go:83: attempt(32) was successful
2021/11/01 22:51:28 app.go:83: attempt(22) was successful
2021/11/01 22:51:28 app.go:83: attempt(13) was successful
...
```

## Notes

### Might need to increase your ulimit
Especially if running on OSX, ulimit starts at like 256. A ulimit -n 5000 should easily be enough to run this demo.
