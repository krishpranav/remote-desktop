
ifndef encoders
encoders = h264
endif

tags = 
ifneq (,$(findstring h264,$(encoders)))
tags = h264enc
endif

ifneq (,$(findstring vp8,$(encoders)))
tags := $(tags) vp8enc
endif

tags := $(strip $(tags))

main.tar.gz: clean main
	@tar zcf main.tar.gz frontend main

main.zip: clean main
	@zip -r main.zip frontend main

main:
	go build -tags "$(tags)" cmd/main.go

.PHONY: clean
clean:
	@if [ -f main ]; then rm main; fi
	@if [ -f main.tar.gz ]; then rm main.tar.gz ; fi
	@if [ -f main.zip ]; then rm main.zip ; fi