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
	"net/url"
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

var ( // these are modified by Makefile
	Version            = "0.3.x-unknown"
	DefaultDir  string = "" // directory if no -d flag
	DefaultBind        = "127.0.0.1"
)

func init() {
	flag.IntVar(&port, "p", 0, "Port to listen on (default: 0, look for free port)")
	flag.StringVar(&in, "i", DefaultBind, "Interface to listen on, default "+DefaultBind)
	flag.StringVar(&dir1, "d", DefaultDir, "Directory to serve (if empty: public-html, or current working dir)")
	flag.IntVar(&portMin, "minport", 8000, "Minimum port to try binding to")
	flag.IntVar(&portMax, "maxport", 8999, "Maximum port to try binding to")
	flag.BoolVar(&logging, "log", false, "Enable http request logging")
	flag.BoolVar(&singlePage, "single", false, "Single page mode (see below)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "servest v%s\n", Version)
		fmt.Fprintf(os.Stderr, "Source: %s\n", repolink)
		fmt.Fprintf(os.Stderr, "Usage: servest [flags]\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n\tservest -log -i 127.0.0.1 -p 8080 -d /var/www/html\n")
		fmt.Fprintf(os.Stderr, "Run without installing:\n\tgo run -v github.com/aerth/servest@latest -d /some/dir\n")
		fmt.Fprintf(os.Stderr, "Serve current directory:\n\tservest\n")

		fmt.Fprintf(os.Stderr, "\nSingle page mode: will serve index.html for all requests that do not match an existing file. The subdirectory may contain an index.html, and that will be served if nonzero size. \n")
	}
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
	var err error
	// User defined a directory to serve
	if dir1 != "" {
		servepath = dir1
		servepath, err = filepath.Abs(servepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[servest] fatal: abs %v", err)
			os.Exit(222)
		}
		fmt.Fprintf(os.Stderr, "[servest] serving %s directory\n", servepath)
	} else if _, err = os.Stat("public-html"); err == nil {
		// Else we serve public-html directory
		servepath = "public-html"
		servepath, err = filepath.Abs(servepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[servest] fatal: abs %v", err)
			os.Exit(222)
		}
		fmt.Fprintf(os.Stderr, "[servest] serving %s directory\n", servepath)
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
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// get file 	extension
	urlPath, err := url.PathUnescape(r.URL.Path)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if strings.Contains(urlPath, "..") {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	if urlPath == "/favicon.ico" {
		http.ServeFile(w, r, filepath.Join(www.dir, "favicon.ico"))
		return
	}
	if urlPath == "/" {
		urlPath = "/index.html"
	}
	ext := filepath.Ext(urlPath)
	if len(ext) > 10 { // indicates some/.hidden/file (like a git dir)
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	p := filepath.Join(www.dir, filepath.Clean(urlPath))
	if logging {
		fmt.Fprintf(os.Stderr, "[servest] %s %s %s %s %s - %s %s (file=%s)\n",
			time.Now().Format(time.ANSIC), r.Method, r.Host, urlPath, r.URL.RawPath, trimPort(r.RemoteAddr), r.UserAgent(), p)
	}
	if !strings.HasPrefix(p, www.dir) {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}

	// SPA chunk modified from github.com/roberthodgen/spa-server

	// p is the absolute path of the requested HTTP path
	println("p: ", p)
	stat, err := os.Lstat(p)
	doesntexist := os.IsNotExist(err)
	isdir := err == nil && stat.IsDir()
	if singlePage && isdir {
		println("isdir", p)
		p = filepath.Join(p, "index.html")
		println("trying", p)
		stat, err = os.Lstat(p) // recheck /foo/index.html
		if err == nil && stat.Size() > 8 {
			println("index.html exists", p)
			http.ServeFile(w, r, p)
			return
		}
	}
	if singlePage && err != nil {
		p = filepath.Join(www.dir, "index.html")
	}
	if doesntexist {
		http.ServeFile(w, r, filepath.Join(www.dir, "index.html"))
		return
	} else if isdir {
		http.ServeFile(w, r, filepath.Join(www.dir, "index.html"))
		return
	} else if !stat.Mode().IsRegular() {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	switch strings.ToLower(ext) {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "text/javascript")
	case ".html":
		w.Header().Set("Content-Type", "text/html")
	case ".json":
		w.Header().Set("Content-Type", "application/json")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	case ".webp":
		w.Header().Set("Content-Type", "image/webp")
	}
	fmt.Printf("serving %s\n", p)
	http.ServeFile(w, r, p)

	// unused now
	//www.underlyingHandler.ServeHTTP(w, r)
}
