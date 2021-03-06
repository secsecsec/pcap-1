include $(GOROOT)/src/Make.inc

TARG=github.com/davecheney/pcap
GOFILES=\
	packet.go\
	ethernet.go\
	frame.go\
	ip.go\
	mac.go\
	pcap.go\
	udp.go\
	file.go\

GOFILES_darwin=\
	bpf.go\

GOFILES_amd64=\
	ztypes_amd64.go\

GOFILES_386=\
	ztypes_386.go\

GOFILES+=$(GOFILES_$(GOOS))

GOFILES+=$(GOFILES_$(GOARCH))

CLEANFILES+=ztypes_*.go

include $(GOROOT)/src/Make.pkg

ztypes_386.go: types.c
	godefs -gpcap -f -m32 $^ | gofmt > $@

ztypes_amd64.go: types.c
	godefs -gpcap -f -m64 $^ | gofmt > $@

