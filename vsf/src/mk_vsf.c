// This file initializes a very simple filesystem.
// It initializes a given disk file's super block and
// clears out the remaining data.

#include "vsf.h"
#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>
#include <unistd.h>

// Global variables
// Path to the disk image
char *diskImage = NULL;
// Number of inodes and data blocks
int numInodes = 0;
int numBlocks = 0;

/**
 * Prints the usage for mk_vsf. It should print the options available like
 * -i to initialize number of inodes (divisible by 2) -b to set the number of
 * data blocks and disk image to check out the disk
 */
void usage(void) {
  printf("Usage: mk_vsf -i <number of inodes> -b <number of data blocks> "
         "-d <disk image>\n");
  exit(0);
}

int isPowerOfTwo(int n) { return n > 0 && (n & (n - 1)) == 0; }

int nextPowerOfTwo(int n) {
  int power = 1;
  while (power < n) {
    power *= 2;
  }
  return power;
}

/**
 * Parses the command line arguments and sets the global variables
 */
void parseArgs(int argc, char **argv) {
  int opt;
  while ((opt = getopt(argc, argv, "i:b:d:")) != -1) {
    switch (opt) {
    case 'i':
      numInodes = atoi(optarg);
      if (numInodes % 2 != 0) {
        fprintf(stderr, "Error: Number of inodes must be divisible by 2.\n");
        usage();
      }
      if (!isPowerOfTwo(numInodes)) {
        int newNumInodes = nextPowerOfTwo(numInodes);
        fprintf(stderr,
                "Warning: Number of inodes is not a power of two. "
                "Updating from %d to %d.\n",
                numInodes, newNumInodes);
        numInodes = newNumInodes;
      }
      break;
    case 'b':
      numBlocks = atoi(optarg);
      if (!isPowerOfTwo(numBlocks)) {
        int newNumBlocks = nextPowerOfTwo(numBlocks);
        fprintf(stderr,
                "Warning: Number of data blocks is not a power of two. "
                "Updating from %d to %d.\n",
                numBlocks, newNumBlocks);
        numBlocks = newNumBlocks;
      }
      break;
    case 'd':
      diskImage = optarg;
      break;
    default:
      usage();
      break;
    }
  }

  if (numInodes == 0 || numBlocks == 0 || diskImage == NULL) {
    usage();
  }
}

/**
 * Main function for mk_vsf
 */
int main(int argc, char **argv) {
  printf("Running mk_vsf\n");

  parseArgs(argc, argv);

  int fd = open(diskImage, O_RDWR,
                S_IRUSR | S_IWUSR | S_IRGRP | S_IWGRP | S_IROTH | S_IWOTH);
  if (fd < 0) {
    perror("Error opening disk image");
    return 1;
  }

  // Map the disk image into memory
  // 1 Super Block, Inode and Data Bitmaps, and remaining blocks
  int mapSize = BLOCK_SIZE * (1 + 2 + numInodes + numBlocks);
  char *map = mmap(NULL, mapSize, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0);
  if (map == MAP_FAILED) {
    perror("Error mapping disk image");
    return 1;
  }

  // Clear out the disk image
  memset(map, 0, mapSize);

  // Initialize the super block
  superblock_t *sb = (superblock_t *)map;
  sb->magic_number = VSF_MAGIC;
  sb->block_size = BLOCK_SIZE;
  sb->total_blocks = numBlocks;
  sb->blocks_bitmap = 1;
  sb->total_inodes = numInodes;
  sb->inodes_bitmap = 2;

  return 0;
}
