#include <unistd.h>
#include <stdio.h>

enum { bufSize = 1024 };

int main(void) {
    char buf[bufSize];
    ssize_t bytesRead, bytesWritten, total = 0;
    while ((bytesRead = read(STDIN_FILENO, buf, bufSize)) > 0) {
        bytesWritten = write(STDOUT_FILENO, buf, bytesRead);
        if (bytesWritten < 0) {
            perror("write");
            return bytesWritten;
        }
        total += bytesWritten;
    }
    if (bytesRead < 0) {
        perror("read");
    }
    fprintf(stderr, "%zd\n", total);
    return bytesRead;
}
