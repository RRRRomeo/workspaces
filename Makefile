GO = go build
SRC = sample_map.go new_map.go

all: clean main run

main:
	@$(GO) -o $@ $(SRC)

run:
	@./main

clean:
	@-rm ./main
.PHONY:all main