.SILENT: test test-unit test-bench
test-bench:
	@echo 'Run benchmark tests'

	go test  -bench=. -count=1 -benchmem -parallel=4

	rm ./cover.out &> /dev/null

test-unit:
	@echo 'Run unit tests'

	go test -count=1 -parallel=2 -coverprofile=./cover.out ./...
	go tool cover -func ./cover.out
	go tool cover -func ./cover.out|grep -Po  '(?<=\s)\d{1,2}'|tail --lines 1|xargs echo "Total coverage: "

test:
	$(MAKE) test-unit
	$(MAKE) test-bench