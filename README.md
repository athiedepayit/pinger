# pinger

design goals:

- Multi-directional service that checks if a configured "other side" is alive, and provides status updates based on success or failure
- "A" side:
    - Listen on a port
    - Attempt to request a remote resource ("B" side)
    - Provide a status page (/health) for this remote request
- "B" side:
    - Listen on a port
    - Attempt to request a remote resource ("A" side)
    - Provide a status page (/health) for this remote request
- Both sides:
    - run command (i.e. curl a webhook) on failure and recovery

Going to attempt to write this in go since it provides a good standard http library


How you run two copies on the same machine:

"A" side: `./pinger -interval 10 -port 8881 -remote "http://localhost:8882"`

"B" side: `./pinger -interval 10 -port 8882 -remote "http://localhost:8881"`

Using commands:

Optionally you can specify `-recoverycmd` and `-errorcmd` commands. These are strings that are just run as commands in a shell. Probably quote them when you run it on the command line. So if you run "curl" you'd better have curl installed. Making this just use external commands greatly simplifies things; on a desktop you can use `-errorcmd "notify-send -u critical SERVER DOWN"` for example.


## Run example

Example where I ran the "A" side first, then turned on the "B" side a few runs in.


### "A" side

```
$ go run . -interval 3 -remote "http://localhost:8761" -port 8701 -errorcmd "notify-send dead" -recoverycmd "notify-send alive"
---
- Interval: 3
- ListenPort: 8701
- RemoteHost: http://localhost:8761
- ErrCommand: notify-send dead
- RecCommand: notify-send alive
---
checking http://localhost:8761
Error running check: Get "http://localhost:8761": dial tcp [::1]:8761: connect: connection refused
checking...
result is false
checking http://localhost:8761
Error running check: Get "http://localhost:8761": dial tcp [::1]:8761: connect: connection refused
checking...
result is false
checking http://localhost:8761
Error running check: Get "http://localhost:8761": dial tcp [::1]:8761: connect: connection refused
checking...
result is false
checking http://localhost:8761
checking...
result is true
sending command 'notify-send alive'
checking http://localhost:8761
checking...
result is true
checking http://localhost:8761
checking...
result is true
```

### "B" side

```
$ go run . -interval 10 -remote "http://localhost:8701"
---
- Interval: 10
- ListenPort: 8761
- RemoteHost: http://localhost:8701
- ErrCommand:
- RecCommand:
---
checking http://localhost:8701
Error running check: Get "http://localhost:8701": dial tcp [::1]:8701: connect: connection refused
checking...
result is false
checking http://localhost:8701
Error running check: Get "http://localhost:8701": dial tcp [::1]:8701: connect: connection refused
checking...
result is false
```

## Curl output

### "Bad"

```
$ curl http://localhost:8761/health -v
* Host localhost:8761 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8761...
* Connected to localhost (::1) port 8761
* using HTTP/1.x
> GET /health HTTP/1.1
> Host: localhost:8761
> User-Agent: curl/8.10.1
> Accept: */*
>
* Request completely sent off
< HTTP/1.1 500 Internal Server Error
< Date: Sat, 05 Oct 2024 20:47:05 GMT
< Content-Length: 6
< Content-Type: text/plain; charset=utf-8
<
Error
* Connection #0 to host localhost left intact
```

### "Good"

```
$ curl http://localhost:8761/health -v
* Host localhost:8761 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8761...
* Connected to localhost (::1) port 8761
* using HTTP/1.x
> GET /health HTTP/1.1
> Host: localhost:8761
> User-Agent: curl/8.10.1
> Accept: */*
>
* Request completely sent off
< HTTP/1.1 200 OK
< Date: Sat, 05 Oct 2024 20:48:12 GMT
< Content-Length: 3
< Content-Type: text/plain; charset=utf-8
<
OK
* Connection #0 to host localhost left intact
```

