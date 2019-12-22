# logyeet

`logyeet` is a logging agent that is designed to get log data off of a server (or container) and
into your analysis software as quickly and efficiently as possible.  It communicates over plain
HTTPS (actually, gRPC), makes no attempt to understand the log data it is sending, and only uses a
tiny amount of RAM.

It gives your centralized analysis/storage server the bytes.  That's all it does.

gRPC is used so that you don't need to go through any special gymnastics to get log data from
isolated servers.  It's just another route in your edge proxy.
