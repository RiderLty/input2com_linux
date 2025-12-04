package main

import (
	"net"
	"os"
)

type UDSPacketConn struct {
	conn *net.UnixConn
	ch   chan []byte
}

func (u *UDSPacketConn) ReadChan() <-chan []byte { return u.ch }

func (u *UDSPacketConn) Close() error {
	return u.conn.Close()
}

func (u *UDSPacketConn) serve() {
	buf := make([]byte, 65536)
	for {
		n, _, err := u.conn.ReadFromUnix(buf)
		if err != nil {
			close(u.ch)
			return
		}
		pkt := make([]byte, n)
		copy(pkt, buf[:n])
		u.ch <- pkt
	}
}
func CreateUDSReader(address string) (*UDSPacketConn, error) {
	addr, err := net.ResolveUnixAddr("unixgram", address)
	if err != nil {
		return nil, err
	}
	if address[0] != '@' {
		_ = os.Remove(address)
	}
	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		return nil, err
	}
	u := &UDSPacketConn{
		conn: conn,
		ch:   make(chan []byte, 16),
	}
	go u.serve()
	return u, nil
}

type UDSWriter struct {
	conn *net.UnixConn
}

func (w *UDSWriter) Write(p []byte) (int, error) {
	return w.conn.Write(p)
}

func (w *UDSWriter) Close() error { return w.conn.Close() }

func CreateUDSWriter(address string) (*UDSWriter, error) {
	addr, err := net.ResolveUnixAddr("unixgram", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUnix("unixgram", nil, addr)
	if err != nil {
		return nil, err
	}
	return &UDSWriter{conn: conn}, nil
}

type UDPPacketConn struct {
	conn *net.UDPConn
	ch   chan []byte
}

func (u *UDPPacketConn) ReadChan() <-chan []byte { return u.ch }

func (u *UDPPacketConn) Close() error { return u.conn.Close() }

func (u *UDPPacketConn) serve() {
	buf := make([]byte, 65536)
	for {
		n, _, err := u.conn.ReadFromUDP(buf)
		if err != nil {
			close(u.ch)
			return
		}
		pkt := make([]byte, n)
		copy(pkt, buf[:n])
		u.ch <- pkt
	}
}

func CreateUDPReader(address string) (*UDPPacketConn, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	u := &UDPPacketConn{
		conn: conn,
		ch:   make(chan []byte, 16),
	}
	go u.serve()
	return u, nil
}

// ========== UDP Writer ==========

type UDPWriter struct {
	conn *net.UDPConn
}

func (w *UDPWriter) Write(p []byte) (int, error) {
	return w.conn.Write(p)
}

func (w *UDPWriter) Close() error { return w.conn.Close() }

func CreateUDPWriter(remoteAddr string) (*UDPWriter, error) {
	raddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}
	return &UDPWriter{conn: conn}, nil
}
