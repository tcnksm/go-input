SHELL := /bin/bash

test: tools
	rm -rf test-reports
	mkdir test-reports
	go clean -testcache
	GO111MODULE="off" go test -v 2>&1 ./... | go-junit-report -iocopy -set-exit-code -out test-reports/unit-test-report.xml

tools:
	go install github.com/jstemmer/go-junit-report/v2@v2.0.0
