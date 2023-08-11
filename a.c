#include <algorithm>
#include <cstdio>
#include <iostream>
#include <string>
#include <vector>
#include <cstdlib>
#include <fcntl.h>
#include <unistd.h>
#include <sys/mman.h>
#include <sys/stat.h>

#define handle_error(msg) \
   do { perror(msg); exit(EXIT_FAILURE); } while (0)

int
main() {
    int x = rand();
    std::string filename = "/tmp/data";
//    int fd =
//        open(filename.data(), O_CREAT | O_TRUNC | O_RDWR, S_IRUSR | S_IWUSR);
    int fd =
        open(filename.data(), O_RDONLY);
    printf("fd = %d\n", fd);
    if (fd == -1)
        handle_error("open");
    std::vector<char> data;
    data.resize(1024 * 1024);
//    for (int i = 0; i < 1024; i++) {
//        write(fd, data.data(), data.size());
//    }
//    fsync(fd);

    printf("will mmap...");
    getchar();
    auto filesize = 524288000;
    printf("filesize: %d\n", filesize);
    auto ptr = (char*)mmap(nullptr,
                           filesize,
                           PROT_READ,
                           MAP_PRIVATE,
                           fd,
                           0);
    close(fd);
    printf("maped...");
    for (int i = 0; i < filesize; i += (4 << 10)) {
        // ptr[i] = i;
        data[i % 1024] = ptr[i];
    }
    printf("populated...");
    getchar();

    munmap(ptr, filesize);
    printf("unmap...");
    getchar();
//    remove(filename.data());
    return 0;
}
