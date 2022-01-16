# What is Babylon?

The Babylon project is a re-creation of many of my original PBX network
integration tools from the early 1990's, in go.  Many of these were originally
written in C for QNX, and then ported to Linux. This will eventually include
the F9600 mml server, a Panasonic PAPI DBS server, an smdi server, the SPO256
speaker, cdr logging servers, a cdr collection server, and other odd things.

One reason for this project is simply to help me decide best practices for
developing enterprise services in golang going forward.  This project will
eventually directly support building for traditional os packaging, offer
ansible deployment, or docker creation, using a single top level Makefile.

## Installation

Make "install" is sufficient to install these tools and daemons on a generic
posix system. This installs to /usr/local by default, and can be overridden
with a PREFIX setting, such as ''make PREFIX=/usr install''. The Makefile also
makes it easy to cross-compile, as well as managing separate debug and release
builds. It also should be easy to integrate with traditional OS packaging.

In git checkouts I manage a vendor directory outside of git.  This is because
it may generate different content when you update the go.mod file. Generating
a vendor branch means it also can get into the stand-alone dist tarball, and
that can then be used in network isolated build systems.  Since the builds are
cached anyway without a vendor directory, this has no impact on performance.
The vendor directory is only refreshed if the go.sum file changes.

## Participation

Babylon is written in Go. I use go modules support, so this project can be
cloned into any stand-alone directory and built there. I also use a front-end
Makefile to simplify builds and offer basic project automation.

My Makefile includes a "dist" target, which produces a stand-alone tarball that
can then be used to build install the tools detached from the repo. This is
particularly useful for OS packaging. I include a "lint", "test", and
"coverage" target to test and verify code. The default make builds binaries in
target/{build-type} in a manner like rust cargo does, where they can then be
tested. Cross-compiling goes into target/${GOOS}-${GOARCH}.  Many of these
special targets are standardized in the .make directory.

## Support

Support is offered thru https://git.gnutelephony.org/babylon/issues. When
entering a new support issue, please mark it part of the support project. I
also have dyfet@jabber.org. Babylon packaging build support for some GNU/Linux
distributions may be found on https://pkg.gnutelephony.org/babylon. I also
have my own build infrastructure for Alpine using ProduceIt, and I publish apk
binary packages thru https://public.tychosoft.com/alpine. In the future maybe
other means of support will become possible.

