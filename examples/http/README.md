# Domain sockets over SSH HTTP demo

This shows example code on how you could run and HTTP server over SSH using domain sockets. While this wouldn't be useful for user traffic, this can be quite useful for admin traffic to your device such as gathering metrics or other sensitive information.

## Running

You can run this demo over you local SSH server if you don't want to do a two box setup. You also can use the SSH agent, a key file or password file.

Running the server like so:

```bash
cd service
go run service.go --socket=./socket.sock
```

Now to run the client, which will do 100 queries against the server. For the example this will use the password file "./pass" 

```bash
cd app
go run app.go --endpoint=127.0.0.1:22 --socket=[path to this directory]/examples/http/service/socket.sock --pass=pass
```

That's it, you should see output:
```
2021/11/01 18:58:34 app.go:76: server returned:  Hello!
2021/11/01 18:58:34 app.go:76: server returned:  Hello!
2021/11/01 18:58:34 app.go:76: server returned:  Hello!
2021/11/01 18:58:34 app.go:78: attempt(70) was successful
...
```
