package main

import (
	seccomp "github.com/elastic/go-seccomp-bpf"
)

func Seccomp() error {
	// Create a filter.
	filter := seccomp.Filter{
		NoNewPrivs: true,
		Flag:       seccomp.FilterFlagTSync,
		Policy: seccomp.Policy{
			DefaultAction: seccomp.ActionErrno,
			Syscalls: []seccomp.SyscallGroup{
				{
					Action: seccomp.ActionAllow,
					Names: []string{
						"close",
						"fstat",
						"poll",
						"lseek",
						"mmap",
						"mprotect",
						"munmap",
						"rt_sigaction",
						"rt_sigprocmask",
						"rt_sigreturn",
						"ioctl",
						"pread64",
						"readv",
						"access",
						"sched_yield",
						"dup",
						"dup2",
						"socketpair",
						"clone",
						"execve",
						"exit",
						"wait4",
						"uname",
						"fcntl",
						"flock",
						"getcwd",
						"readlink",
						"setuid",
						"setgid",
						"setgroups",
						"sigaltstack",
						"sched_getparam",
						"sched_getscheduler",
						"arch_prctl",
						"futex",
						"sched_getaffinity",
						"epoll_create",
						"getdents64",
						"set_tid_address",
						"clock_gettime",
						"exit_group",
						"epoll_wait",
						"epoll_ctl",
						"waitid",
						"openat",
						"newfstatat",
						"ppoll",
						"epoll_pwait",
						"eventfd2",
						"epoll_create1",
						"pipe2",
						"prlimit64",
						"sched_setattr",
						"userfaultfd",
						"pwritev2",
						"pkey_free",
						"io_uring_enter",
						"io_uring_register",
						"fsconfig",
						"clone3",
						"pidfd_getfd",
						"faccessat2",
						"epoll_pwait2",
						"quotactl_fd",
						"read",
						"write",
					},
				},
			},
		},
	}

	// Load it. This will set no_new_privs before loading.
	return seccomp.LoadFilter(filter)
}
