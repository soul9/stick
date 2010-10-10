include $(GOROOT)/src/Make.inc

TARG=stick
EXAMPLES=examples/mark.go
GOFILES=stick.go\
                 conf.go

include $(GOROOT)/src/Make.pkg 
