
run:
	go run ./cmd/testing/

drive-add:
	go run ./cmd/cli drive add -n main

file-ls:
	go run ./cmd/cli ls -a main

upload:
	go run ./cmd/cli upload -a main -i ./README.md -o /dev/tmp_files/README.md

generate:
	wire ./...
