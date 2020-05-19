# servest

Lightweight and easy-to-configure http fileserver

With no arguments, binds to 0.0.0.0, finds a port from 8000 to 8999, and serves current directory

See Makefile for build customization, such as default directory and default bind address

```
Usage of ./servest:
  -d string
    	Directory to serve (current working dir if empty)
  -i string
    	Interface to listen on (default "0.0.0.0")
  -log
    	Enable http request logging
  -maxport int
    	Maximum port to try binding to (default 8999)
  -minport int
    	Minimum port to try binding to (default 8000)
  -p int
    	Port to listen on (default: 0, look for free port)
  -version
    	Show version information and exit
```

TODO:

  * pidfile option?
  * logfile option?
  * daemonize option?
