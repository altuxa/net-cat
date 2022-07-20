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