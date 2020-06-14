VERSION  				:= 			$(shell cat ./VERSION)

install:
	go install -v

test:
	go test ./... -v

fmt:
	go fmt ./... -v

.PHONY: install test fmt
