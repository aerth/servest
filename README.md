# servest

Lightweight and easy-to-configure http fileserver

When Makefile is used, default bind interface is 0.0.0.0

When `go run -v github.com/aerth/servest@latest` is used, default bind interface is 127.0.0.1

Without `-p` flag, `servest` finds an available port from 8000 to 8999, and serves current directory (or `public-html` if it exists)

With `-single` flag, `servest` will serve /index.html, or /{path}/index.html if it path has no extension

See Makefile for build customization, such as default directory and default bind address

```
servest v0.3.x-unknown
Source: https://github.com/aerth/servest
Usage: servest [flags]
  -d string
        Directory to serve (if empty: public-html, or current working dir)
  -i string
        Interface to listen on, default 127.0.0.1 (default "127.0.0.1")
  -log
        Enable http request logging
  -maxport int
        Maximum port to try binding to (default 8999)
  -minport int
        Minimum port to try binding to (default 8000)
  -p int
        Port to listen on (default: 0, look for free port)
  -single
        Single page mode (see below)
  -version
        Show version information and exit

Example:
        servest -log -i 127.0.0.1 -p 8080 -d /var/www/html
Run without installing:
        go run -v github.com/aerth/servest@latest -d /some/dir
Serve current directory:
        servest

Single page mode: will serve index.html for all requests that do not match an existing file. The subdirectory may contain an index.html, and that will be served if nonzero size.
```

TODO:

  * pidfile option?
  * logfile option?
  * daemonize option?
