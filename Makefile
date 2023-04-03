GO = go build
GO_T = go test
GO_FLAGS =
GO_LFLAGS = 
BUILD_DIR = ./build/
TEST_DIR = ./test/
CMD_DIR = ./cmd/
SRC_NEW_MAP = ./cmd/new_map_cmd/new_map_cmd.go
SRC_SAMPLE_MAP = ./cmd/sample_map_cmd/sample_map_cmd.go

new_main:
	@$(GO) -o $(BUILD_DIR)$@ $(SRC_NEW_MAP)

sample_main:
	@$(GO) -o $(BUILD_DIR)$@ $(SRC_SAMPLE_MAP)

%:
	@$(GO) $(GO_FLAGS) -o $(BUILD_DIR)$@ $(CMD_DIR)$@/$@.go $(GO_LFLAGS)

run:
	@./build/main

clean:
	@-rm ./build/*_main
	
.PHONY: clean new_main sample_main %