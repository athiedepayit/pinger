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

