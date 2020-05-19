#
# Easily Configurable Variables
#

# Consider changing DEFAULTBIND to 127.0.0.1
DEFAULTBIND ?= 0.0.0.0

# VERSION grabs version from VERSION file and git commit
VERSION ?= ${shell cat VERSION}-${shell git rev-parse --short HEAD} 

# DefaultDir is used when no command line arguments are given
# For security purposes, consider /var/www/html or similar
#DEFAULTDIR=

# DESTDIR to install to
DESTDIR ?= /usr/local/bin

# build servest
go != which go
ldflags="-w -s -X main.Version=${VERSION} -X main.DefaultDir=${DEFAULTDIR} -X main.DefaultBind=${DEFAULTBIND}"
buildflags=-v -ldflags ${ldflags} -tags osusergo,netgo

servest: *.go
	${go} build ${buildflags} -o $@

dev: *.go
	env DEFAULTBIND=127.0.0.1 ${MAKE} clean servest

install:
	install servest ${DESTDIR}/servest

clean:
	rm -f servest
