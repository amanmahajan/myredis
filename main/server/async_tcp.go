package server

import (
	"errors"
	"log"
	"myredis/config"
	"myredis/core"
	"net"
	"syscall"
	"time"
)

var clients int = 0

var lastCronExecution = time.Now()
var deleteInterval = 1 * time.Second

// RunAsyncTcpServer MACOS compatible Server
func RunAsyncTcpServer() error {

	log.Println("Starting AsyncTcpServer")

	max_clients := 10000

	kq, err := syscall.Kqueue()
	if err != nil {
		return err
	}

	defer syscall.Close(kq)

	events := make([]syscall.Kevent_t, max_clients)

	// Lets create a socket. It returns server file descriptor
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}

	defer syscall.Close(serverFD)

	// Non blocking mode for the server
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	/*
		ip4[12], ip4[13], ip4[14], ip4[15] are the last four bytes of the net.IP slice, which represent the IPv4 address.
		In Go, net.IP is a slice where the first 12 bytes are used for IPv6 addresses or padding, and the last 4 bytes are used for IPv4 addresses.
		By extracting these bytes, the code is converting the net.IP slice into the [4]byte array required by syscall.SockaddrInet4.
	*/

	ip4 := net.ParseIP(config.Host)
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[12], ip4[13], ip4[14], ip4[15]},
	}); err != nil {
		return err
	}

	if err = syscall.Listen(serverFD, max_clients); err != nil {
		return err
	}

	// Async io start here

	// Specify the events we want to get hints about
	// and set the socket on which
	// Prepare an event to monitor the server socket
	event := syscall.Kevent_t{
		Ident:  uint64(serverFD),                   // The file descriptor to monitor
		Filter: syscall.EVFILT_READ,                // Equivalent to EPOLLIN
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE, // Add and enable the event
	}

	// Register the server socket with kqueue
	if _, err := syscall.Kevent(kq, []syscall.Kevent_t{event}, nil, nil); err != nil {
		return err
	}

	for {

		// Run the expiration key algorithm after every 1 sec
		if time.Now().After(lastCronExecution.Add(deleteInterval)) {
			core.DeleteExpiredKey()
			lastCronExecution = time.Now()
		}

		// See if any file descriptor ready for io
		nevents, err := syscall.Kevent(kq, nil, events, nil)
		if err != nil {
			if errors.Is(err, syscall.EINTR) {
				continue // Retry the syscall if it was interrupted
			}

			return err
		}
		for i := 0; i < nevents; i++ {
			// If the server socket itself is ready for an IO
			if events[i].Ident == uint64(serverFD) {
				// Accept the incoming connection from a client
				fd, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("err", err)
					continue
				}

				// Increase the number of concurrent clients count
				// Current fd will be the File descriptor of the new client connection
				clients++
				syscall.SetNonblock(fd, true)

				// Add this new TCP connection to be monitored
				clientEvent := syscall.Kevent_t{
					Ident:  uint64(fd),
					Filter: syscall.EVFILT_READ,
					Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
				}
				if _, err := syscall.Kevent(kq, []syscall.Kevent_t{clientEvent}, nil, nil); err != nil {
					log.Fatal(err)
				}
			} else {
				comm := core.FileDescriptorComm{FileDescriptor: int(events[i].Ident)}
				cmds, err := readCommands(comm)
				if err != nil {
					syscall.Close(int(events[i].Ident))
					clients--
					continue
				}
				respond(cmds, comm)
			}
		}
	}

}
