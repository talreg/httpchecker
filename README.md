HTTP checker / loader
============
This is a simple http client with a web interface.
I have written it in order to load test/switch test servers.

Building
-------------
You'll need `go` installed.
There is a `Makefile` for it. If you do not have make, you can run the commands by yourself.
<br>The `make` command will build and set the file in the `bin` folder.

Running
----------
The runner has one parameter - `address`. It must be provided.

example: to runa check on google:`bin/checker --address=https://google.com`

The server will then run in the default port, 3000. You should be able to access it with http://localhost:3000.

You will see connection errors with timestamp if these happen. You'll also see bad and good http responses (based on the http code).

Changing the load
-------------------
You can adjust the `serverLoadFactor` variable as you see fit.
