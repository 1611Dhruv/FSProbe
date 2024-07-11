# Makefile for compiling files from vsf/src and saving them in vsf/bin

# Compiler
CC = gcc

# Compiler flags
CFLAGS = -Wall -Werror -pedantic -std=gnu18 -g
FUSE_CFLAGS = `pkg-config fuse --cflags --libs`

# Source and binary directories
SRC_DIR = vsf/src
BIN_DIR = vsf/bin

# Source files
SRC_FILES = $(SRC_DIR)/mk_vsf.c $(SRC_DIR)/vsf.c

# Executables
MK_VSF_EXEC = $(BIN_DIR)/mk_vsf
VSF_EXEC = $(BIN_DIR)/vsf

# Default target
all: build_system run_fuse

# Create bin directory if it doesn't exist
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Rule to compile mk_vsf
$(MK_VSF_EXEC): $(SRC_DIR)/mk_vsf.c | $(BIN_DIR)
	$(CC) $(CFLAGS) -o $(MK_VSF_EXEC) $(SRC_DIR)/mk_vsf.c

# Rule to compile vsf
$(VSF_EXEC): $(SRC_DIR)/vsf.c | $(BIN_DIR)
	$(CC) $(CFLAGS) $(FUSE_CFLAGS) -o $(VSF_EXEC) $(SRC_DIR)/vsf.c


# Target to build the system
build_system: $(MK_VSF_EXEC) $(VSF_EXEC)
	@echo "File system built successfully."

# Target to run the FUSE file system
run_fuse: $(VSF_EXEC)
	@echo "Running FUSE file system..."

# Clean up
clean:
	rm -rf $(BIN_DIR)

.PHONY: all clean build_system run_fuse
