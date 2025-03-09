// The MIT License (MIT)
//
// Copyright (c) 2016-2020 aerth
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Servest is a quick http fileserver with a few options
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	port       int
	in         string
	dir1       string
	servepath  string
	portMin    = 8000
	portMax    = 8999
	logging    = false // to stderr
	singlePage = false
)

const (
	repolink = "https://github.com/aerth/servest"
)

var (
	Version     = "0.2.x-unknown"
	DefaultDir  string // cwd if empty
	DefaultBind = "127.0.0.1"
)

func init() {
	flag.IntVar(&port, "p", 0, "Port to listen on (default: 0, look for free port)")
	flag.StringVar(&in, "i", DefaultBind, "Interface to listen on")
	flag.StringVar(&dir1, "d", DefaultDir, "Directory to serve (current working dir if empty)")
	flag.IntVar(&portMin, "minport", 8000, "Minimum port to try binding to")
	flag.IntVar(&portMax, "maxport", 8999, "Maximum port to try binding to")
	flag.BoolVar(&logging, "log", false, "Enable http request logging")
	flag.BoolVar(&singlePage, "single", false, "Single page mode (404s with no ext serve index.html)")
}

func main() {
	var showVersion = false
	flag.BoolVar(&showVersion, "version", false, "Show version information and exit")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "[servest] server require no arguments, only command line flags. see -h to print flags and defaults\n")
		os.Exit(111)
	}
	if showVersion {
		fmt.Printf("servest v%s\n", Version)
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, "[servest] starting...")
	fmt.Fprintln(os.Stderr, "[servest] source code: https://github.com/aerth/servest")
	t1 := time.Now()

	// User defined a directory to serve
	if dir1 != "" {
		servepath = dir1
	} else {
		// Else we serve current working directory (if possible)
		var err error
		servepath, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[servest] fatal: %v", err)
			os.Exit(222)
		}
		fmt.Fprintln(os.Stderr, "[servest] serving current working dir:", servepath)
	}

	webhandler := wrapHandler{servepath}
	// User defined a port for binding
	if port != 0 {
		fmt.Fprintf(os.Stderr, "[servest] serving %s on %s:%d\n", servepath, in, port)
		fmt.Fprintln(os.Stderr, http.ListenAndServe(fmt.Sprintf("%s:%d", in, port),
			webhandler))

		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "[servest] serving %s on %s (first available port)\n", servepath, in)
	fmt.Fprintf(os.Stderr, "[servest] looking for an available port between %d and %d \n", portMin, portMax)
	// Here we search for an open port within the boundries of portMin and portMax.
	for port := portMin; port <= portMax; port++ {
		// We print the port we are *trying* to bind to, if it isn't possible we keep trying different ports.
		go func() {
			<-time.After(time.Second)
			fmt.Fprintf(os.Stderr, "[servest] listening on port: %d.\n", port)
		}()
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", in, port),
			webhandler,
		)
		if !strings.Contains(err.Error(), "already in use") {
			fmt.Fprintln(os.Stderr, "[servest] error:", err)
			if time.Since(t1) < time.Second*3 {
				fmt.Fprintln(os.Stderr, "[servest] server encountered boot error. please report this as an issue at the source code repository:", repolink)
			}
			os.Exit(111)
		}
	}
	os.Exit(1)
}

type wrapHandler struct {
	dir string
}

func trimPort(s string) string {
	host, _, err := net.SplitHostPort(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[servest] error: %s\n", err.Error())
	}
	return host
}

func (www wrapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// insert any routers or www middleware handlers you need here?

	p := filepath.Join(www.dir, filepath.Clean(r.URL.Path))
	if logging {
		fmt.Fprintf(os.Stderr, "[servest] %s %s %s %s %s - %s %s (file=%s)\n", time.Now().Format(time.ANSIC), r.Method, r.Host, r.URL.Path, r.URL.RawPath, trimPort(r.RemoteAddr), r.UserAgent(), p)
	}

	if strings.HasSuffix(r.URL.Path, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(r.URL.Path, ".js") {
		w.Header().Set("Content-Type", "text/javascript")
	} else if strings.HasSuffix(r.URL.Path, ".html") {
		w.Header().Set("Content-Type", "text/html")
	} else if strings.HasSuffix(r.URL.Path, ".json") {
		w.Header().Set("Content-Type", "application/json")
	}

	// SPA chunk modified from github.com/roberthodgen/spa-server
	if info, err := os.Stat(p); err != nil {
		http.ServeFile(w, r, filepath.Join(www.dir, "index.html"))
		return
	} else if info.IsDir() {
		http.ServeFile(w, r, filepath.Join(www.dir, "index.html"))
		return
	}
	http.ServeFile(w, r, p)

	// unused now
	//www.underlyingHandler.ServeHTTP(w, r)
}
