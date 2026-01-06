#include <unistd.h>
#include <stdio.h>

enum { bufSize = 128 }; // small on purpose to get lots of syscalls

int main(void) {
    char buf[bufSize];
    ssize_t bytesRead, /*bytesWritten,*/ total = 0, numCall =0;
    for (;;) {
        numCall++;
        bytesRead = read(STDIN_FILENO, buf, bufSize);
        if (bytesRead <= 0) {
            break;
        }
        /*
        numCall++;
        bytesWritten = write(STDOUT_FILENO, buf, bytesRead);
        if (bytesWritten < 0) {
            perror("write");
            return bytesWritten;
        }
        total += bytesWritten;
        */
       total += bytesRead;
    }
    if (bytesRead < 0) {
        perror("read");
    }
    fprintf(stderr, "%zd (%zd)\n", total, numCall);
    return bytesRead;
}
