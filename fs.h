
#include <stdint.h>
#include <time.h>

#define FUSE_USE_VERSION 26

#define VSF_MAGIC 0x9B

#define BLOCK_SIZE 4096

#define NUM_DIRECT_POINTERS 12
#define NUM_INDIRECT_POINTERS 1

// Common superblock structure for various file systems
typedef struct superblock {
  uint32_t magic_number;      // Magic number to identify the file system
  uint32_t block_size;        // Size of a block in bytes
  uint32_t total_blocks;      // Total number of blocks in the file system
  uint32_t blocks_bitmap;     // Block bitmap for VSF
  uint32_t total_inodes;      // Total number of inodes
  uint32_t inodes_bitmap;     // Inode bitmap for VSF
  uint32_t checkpoint_region; // Checkpoint region for LFS
  uint32_t reserved[8];       // Reserved for future use
} superblock_t;
/*
 * This defines the inode structure for our very
 * simple file system.
 */
typedef struct inode {
  uint32_t inode_number; // Unique inode number
  uint16_t file_type;    // File type (e.g., regular file, directory)
  uint16_t permissions;  // File access permissions
  uint32_t link_count;   // Number of hard links
  uint32_t uid;          // User ID of the file owner
  uint32_t gid;          // Group ID of the file owner
  uint64_t size;         // Size of the file in bytes
  time_t atime;          // Last access time
  time_t mtime;          // Last modification time
  time_t ctime;          // Creation time

  // Data block pointers
  uint32_t
      direct_pointers[NUM_DIRECT_POINTERS]; // Direct pointers to data blocks
  uint32_t indirect_pointers[NUM_INDIRECT_POINTERS]; // Single indirect pointers

  // Additional fields
  uint32_t flags;               // File attribute flags
  uint32_t extended_attributes; // Pointer to extended attributes (optional)
} inode_t;

/*
 * This defines the directory entry structure for our
 * very simple file system.
 */
typedef struct direntry {
  char name[256];
  uint32_t inode_number;
} direntry_t;
