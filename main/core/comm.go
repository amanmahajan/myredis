package core

import "syscall"

type FileDescriptorComm struct {
	FileDescriptor int
}

func (f FileDescriptorComm) Write(b []byte) (int, error) {
	return syscall.Write(f.FileDescriptor, b)
}

func (f FileDescriptorComm) Read(b []byte) (int, error) {
	return syscall.Read(f.FileDescriptor, b)
}
