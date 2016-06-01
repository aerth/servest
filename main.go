// The MIT License (MIT)
//
// Copyright (c) 2016 aerth
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

// Servest is a quick http server with a few options
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var port int
var in string
var dir string
var servepath string
var portMin = 8000
var portMax = 8999

func init() {
	flag.IntVar(&port, "p", 0, "Port to listen on (default: 0, look for free port)")
	flag.StringVar(&in, "i", "0.0.0.0", "Interface to listen on (default: 0.0.0.0)")
	flag.StringVar(&dir, "d", "", "Directory to serve (default: cwd)")
	flag.IntVar(&portMin, "minport", 8000, "Minimum port to try binding to")
	flag.IntVar(&portMax, "maxport", 8999, "Maximum port to try binding to")
}

func main() {
	flag.Parse()
	fmt.Println("[servest]")
	fmt.Println("https://github.com/aerth/servest")

	// User defined a directory to serve
	if dir != "" {
		servepath = dir
	} else {
		// Else we serve current working directory
		servepath, _ = os.Getwd()
	}

	// User defined a port for binding
	if port != 0 {
		fmt.Printf("\nServing %s on %s:%d\n", servepath, in, port)
		fmt.Println(http.ListenAndServe(fmt.Sprintf("%s:%d", in, port), http.FileServer(http.Dir(servepath))))
		os.Exit(1)
	}

	fmt.Printf("\nServing %s on %s\n", servepath, in)
	fmt.Printf("\nLooking for an available port between %d and %d \n", portMin, portMax)
	// Here we search for an open port within the boundries of portMin and portMax.
	for port := portMin; port <= portMax; port++ {
		// We print the port we are *trying* to bind to, if it isn't possible we keep trying different ports.
		fmt.Printf("Port: %d.\n", port)
		fmt.Println(http.ListenAndServe(fmt.Sprintf("%s:%d", in, port), http.FileServer(http.Dir(servepath))))
	}
	os.Exit(1)
}
