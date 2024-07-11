#include "vsf.h"
#include <fuse.h>
#include <stdio.h>

int vsf_getattr(const char *path, struct stat *stbuf) {
  printf("Running vsf_getattr\n");
  return 0;
}

int vsf_readdir(const char *path, void *buf, fuse_fill_dir_t filler,
                off_t offset, struct fuse_file_info *fi) {
  printf("Running vsf_readdir\n");
  return 0;
}

int vsf_mkdir(const char *path, mode_t mode) {
  printf("Running vsf_mkdir\n");
  return 0;
}

int vsf_rmdir(const char *path) {
  printf("Running vsf_rmdir\n");
  return 0;
}

int vsf_create(const char *path, mode_t mode, struct fuse_file_info *fi) {
  printf("Running vsf_create\n");
  return 0;
}

int vsf_open(const char *path, struct fuse_file_info *fi) {
  printf("Running vsf_open\n");
  return 0;
}

int vsf_read(const char *path, char *buf, size_t size, off_t offset,
             struct fuse_file_info *fi) {
  printf("Running vsf_read\n");
  return 0;
}

int vsf_write(const char *path, const char *buf, size_t size, off_t offset,
              struct fuse_file_info *fi) {
  printf("Running vsf_write\n");
  return 0;
}

int vsf_unlink(const char *path) {
  printf("Running vsf_unlink\n");
  return 0;
}

int vsf_truncate(const char *path, off_t size) {
  printf("Running vsf_truncate\n");
  return 0;
}

int vsf_statfs(const char *path, struct statvfs *stbuf) {
  printf("Running vsf_statfs\n");
  return 0;
}

struct fuse_operations ops = {
    .getattr = vsf_getattr,
    .readdir = vsf_readdir,
    .mkdir = vsf_mkdir,
    .rmdir = vsf_rmdir,
    .create = vsf_create,
    .open = vsf_open,
    .read = vsf_read,
    .write = vsf_write,
    .unlink = vsf_unlink,
    .truncate = vsf_truncate,
};

void usage(char *name) {
  printf("Usage: %s disk_path [FUSE options] mount_point\n", name);
  printf("\tdisk_path the path to the file system image\n");
  printf("\t[FUSE options] the various fuse options\n");
  printf("\tmount_point the path of the directory where FUSE should mount the "
         "new file system\n");
}

int main(int argc, char *argv[]) {
  printf("Running VSF\n");
  return fuse_main(argc - 1, argv + 1, &ops, NULL);
}
