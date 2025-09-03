// ported from https://github.com/ldemailly/experimental/blob/master/misc-c/happynull.c
package main

/*
// Works on Linux after
// echo 0 > /proc/sys/vm/mmap_min_addr
#include <sys/mman.h>
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>

__attribute__((constructor))
static void pagingDrZero(void) {
  void *start = (void *)0;
  void *p = mmap(start, 4096, PROT_READ | PROT_WRITE, MAP_FIXED | MAP_ANON | MAP_PRIVATE, -1, 0);
  if (errno != 0) {
      perror("failed alloc did you echo 0 > /proc/sys/vm/mmap_min_addr ?");
      exit(1);
  }
  if (p) {
    printf("Allocated page at %p\n", p);
  }
}
*/
import "C"

func main() {
	var p *int
	// Dereferencing a nil pointer, won't crash if the above worked.
	*p = 42
}
