#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <sched.h>
#include <sys/stat.h>

void nsenter() {
    char *mnt_path = getenv("MNT_PATH");
    if (mnt_path != NULL) {
        int fd = open(mnt_path, O_RDONLY);
        if (fd == -1) {
            perror("open");
            exit(EXIT_FAILURE);
        }

        int ns = setns(fd, 0);
        if (ns == -1) {
            perror("setns");
            exit(EXIT_FAILURE);
        }
    }

    unsetenv("LD_PRELOAD");
}

__attribute__((constructor)) void ctor() {
    nsenter();
}
