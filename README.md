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
