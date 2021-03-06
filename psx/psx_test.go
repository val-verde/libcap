package psx

import (
	"runtime"
	"syscall"
	"testing"
)

func TestSyscall3(t *testing.T) {
	want := syscall.Getpid()
	if got, _, err := Syscall3(syscall.SYS_GETPID, 0, 0, 0); err != 0 {
		t.Errorf("failed to get PID via libpsx: %v", err)
	} else if int(got) != want {
		t.Errorf("pid mismatch: got=%d want=%d", got, want)
	}
	if got, _, err := Syscall3(syscall.SYS_CAPGET, 0, 0, 0); err != 14 {
		t.Errorf("malformed capget returned %d: %v (want 14: %v)", err, err, syscall.Errno(14))
	} else if ^got != 0 {
		t.Errorf("malformed capget did not return -1, got=%d", got)
	}
}

func TestSyscall6(t *testing.T) {
	want := syscall.Getpid()
	if got, _, err := Syscall6(syscall.SYS_GETPID, 0, 0, 0, 0, 0, 0); err != 0 {
		t.Errorf("failed to get PID via libpsx: %v", err)
	} else if int(got) != want {
		t.Errorf("pid mismatch: got=%d want=%d", got, want)
	}
	if got, _, err := Syscall6(syscall.SYS_CAPGET, 0, 0, 0, 0, 0, 0); err != 14 {
		t.Errorf("malformed capget errno %d: %v (want 14: %v)", err, err, syscall.Errno(14))
	} else if ^got != 0 {
		t.Errorf("malformed capget did not return -1, got=%d", got)
	}
}

// The man page for errno indicates that it is never set to zero, so
// validate that it retains its value over a successful Syscall[36]()
// and is overwritten on a failing syscall.
func TestErrno(t *testing.T) {
	// This testing is much easier if we don't have to guess which
	// thread is running this Go code.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Start from a known bad state and clean up afterwards.
	setErrno(int(syscall.EPERM))
	defer setErrno(0)

	v3, _, errno := Syscall3(syscall.SYS_GETUID, 0, 0, 0)
	if errno != 0 {
		t.Fatalf("psx getuid failed: %v", errno)
	}
	v6, _, errno := Syscall6(syscall.SYS_GETUID, 0, 0, 0, 0, 0, 0)
	if errno != 0 {
		t.Fatalf("psx getuid failed: %v", errno)
	}

	if v3 != v6 {
		t.Errorf("psx getuid failed to match v3=%d, v6=%d", v3, v6)
	}

	if v := setErrno(-1); v != int(syscall.EPERM) {
		t.Errorf("psx changes prevailing errno got=%v(%d) want=%v", syscall.Errno(v), v, syscall.EPERM)
	}
}
