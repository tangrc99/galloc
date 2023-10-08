package galloc

import (
	"fmt"
	"golang.org/x/sys/unix"
	"syscall"
)

func mmap(sz int) (error, Page) {
	b, err := unix.Mmap(-1, 0, sz, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|unix.MAP_ANONYMOUS)
	if err != nil {
		return err, Page{}
	}
	// Advise the kernel that the mmap is accessed randomly.
	err = unix.Madvise(b, syscall.MADV_RANDOM)
	if err != nil && err != syscall.ENOSYS {
		// Ignore not implemented error in kernel because it still works.
		return fmt.Errorf("madvise: %s", err), Page{}
	}

	return nil, Page{
		dataRef: b,
		size:    sz,
	}
}

func munmap(b []byte) error {
	return unix.Munmap(b)
}
