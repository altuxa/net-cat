## net-cat

### Objectives

This project consists on recreating the **NetCat in a Server-Client Architecture** that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

### Usage

```console
$ go run ./cmd 
Listening on the port :8989
$ go run ./cmd 2525
Listening on the port :2525
$ go run . 2525 localhost
[USAGE]: ./TCPChat $port
$
```

Server:

```console
$ go run ./cmd 2525
Listening on the port :2525
```

The client :

```console
$ nc $IP $port
```

### Allowed Packages

- io
- log
- os
- fmt
- net
- sync
- time
- bufio
- errors
- strings
- reflect