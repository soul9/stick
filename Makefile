include $(GOROOT)/src/Make.inc

TARG=stick
EXAMPLES=examples/
GOFILES=stick.go conf.go ircextras.go actions.go

include $(GOROOT)/src/Make.pkg 
