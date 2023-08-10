#include <algorithm>
#include <cstdio>
#include <iostream>
#include <string>
#include <vector>
#include <cstdlib>
#include <fcntl.h>
#include <unistd.h>
#include <sys/mman.h>

int
main() {
    int x = rand();
    std::string filename = "test.bin" + std::to_string(x);
    int fd =
        open(filename.data(), O_CREAT | O_TRUNC | O_RDWR, S_IRUSR | S_IWUSR);
    std::vector<char> data;
    data.resize(1024 * 1024);
    for (int i = 0; i < 1024; i++) {
        write(fd, data.data(), data.size());
    }
    fsync(fd);

    printf("will mmap...");
    getchar();
    auto ptr = (char*)mmap(nullptr,
                           1 << 30,
                           PROT_READ,
                           MAP_SHARED,
                           fd,
                           0);
    close(fd);
    printf("maped...");
    for (int i = 0; i < (1 << 30); i += (4 << 10)) {
        // ptr[i] = i;
        data[i % 1024] = ptr[i];
    }
    printf("populated...");
    getchar();

    munmap(ptr, 1 << 30);
    printf("unmap...");
    getchar();
    remove(filename.data());
    return 0;
}
