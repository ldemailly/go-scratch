// Demonstrate nil dereference without panic/segv, works on Linux after
// echo 0 > /proc/sys/vm/mmap_min_addr
// ported from https://github.com/ldemailly/experimental/blob/master/misc-c/happynull.c
package main

/*
#include <sys/mman.h>
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>

__attribute__((constructor))
static void pagingDrZero(void) {
  void *start = (void *)0;
  void *p = mmap(start, 4096, PROT_READ | PROT_WRITE, MAP_FIXED | MAP_ANON | MAP_PRIVATE, -1, 0);
  if (errno != 0) {
      perror("Failed page 0 alloc you need to run:\n\necho 0 > /proc/sys/vm/mmap_min_addr\n\nso normal segv/panic will happen");
  }
}
*/
import "C"
import "fmt"

func main() {
	var p *int
	fmt.Printf("Before nil deref: p=%p\n", p)
	// Dereferencing a nil pointer, won't crash if the above worked.
	*p = 42
	fmt.Printf("after nil deref: p=%p, *p=%d\n", p, *p)
}
