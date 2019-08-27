// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sockstest provides utilities for SOCKS testing.
package socks

import (
	"errors"
	"io"
	"net"
)

// An AuthRequest represents an authentication request.
type AuthRequest struct {
	Version int
	Methods []AuthMethod
}

// ParseAuthRequest parses an authentication request.
func ParseAuthRequest(b []byte) (*AuthRequest, error) {
	if len(b) < 2 {
		return nil, errors.New("short auth request")
	}
	if b[0] != Version5 {
		return nil, errors.New("unexpected protocol version")
	}
	if len(b)-2 < int(b[1]) {
		return nil, errors.New("short auth request")
	}
	req := &AuthRequest{Version: int(b[0])}
	if b[1] > 0 {
		req.Methods = make([]AuthMethod, b[1])
		for i, m := range b[2 : 2+b[1]] {
			req.Methods[i] = AuthMethod(m)
		}
	}
	return req, nil
}

// MarshalAuthReply returns an authentication reply in wire format.
func MarshalAuthReply(ver int, m AuthMethod) ([]byte, error) {
	return []byte{byte(ver), byte(m)}, nil
}

// BasicAuthorization check username and password.
func BasicAuthorization(username, password string, b []byte) bool {
	offset := byte(0)
	if b[offset] != 0x01 {
		return false
	}
	offset++
	uLen := b[offset]
	offset++
	if string(b[offset:offset+uLen]) != username {
		return false
	}
	offset += uLen
	pLen := b[offset]
	offset++
	if string(b[offset:offset+pLen]) != password {
		return false
	}
	return true
}

// MarshalAuthResponse returns an authentication response in wire format.
func MarshalAuthResponse(ver int, ok bool) ([]byte, error) {
	if ok {
		return []byte{byte(ver), byte(0x00)}, nil
	} else {
		return []byte{byte(ver), byte(0x01)}, nil
	}
}

// A CmdRequest repesents a command request.
type CmdRequest struct {
	Version int
	Cmd     Command
	Addr    Addr
	Offset  int
}

// ParseCmdRequest parses a command request.
func ParseCmdRequest(b []byte) (*CmdRequest, error) {
	if len(b) < 7 {
		return nil, errors.New("short cmd request")
	}
	if b[0] != Version5 {
		return nil, errors.New("unexpected protocol version")
	}
	if Command(b[1]) != CmdConnect && Command(b[1]) != CmdUDP {
		return nil, errors.New("unexpected command")
	}
	if b[2] != 0 {
		return nil, errors.New("non-zero reserved field")
	}
	req := &CmdRequest{Version: int(b[0]), Cmd: Command(b[1])}
	l := 2
	off := 4
	switch b[3] {
	case AddrTypeIPv4:
		l += net.IPv4len
		req.Addr.IP = make(net.IP, net.IPv4len)
	case AddrTypeIPv6:
		l += net.IPv6len
		req.Addr.IP = make(net.IP, net.IPv6len)
	case AddrTypeFQDN:
		l += int(b[4])
		off = 5
	default:
		return nil, errors.New("unknown address type")
	}
	if len(b[off:]) < l {
		return nil, errors.New("short cmd request")
	}
	if req.Addr.IP != nil {
		copy(req.Addr.IP, b[off:])
	} else {
		req.Addr.Name = string(b[off : off+l-2])
	}
	req.Addr.Port = int(b[off+l-2])<<8 | int(b[off+l-1])
	req.Offset = off + l
	return req, nil
}

// MarshalCmdReply returns a command reply in wire format.
func MarshalCmdReply(ver int, reply Reply, a *Addr) ([]byte, error) {
	b := make([]byte, 4)
	b[0] = byte(ver)
	b[1] = byte(reply)
	if a.Name != "" {
		if len(a.Name) > 255 {
			return nil, errors.New("fqdn too long")
		}
		b[3] = AddrTypeFQDN
		b = append(b, byte(len(a.Name)))
		b = append(b, a.Name...)
	} else if ip4 := a.IP.To4(); ip4 != nil {
		b[3] = AddrTypeIPv4
		b = append(b, ip4...)
	} else if ip6 := a.IP.To16(); ip6 != nil {
		b[3] = AddrTypeIPv6
		b = append(b, ip6...)
	} else {
		return nil, errors.New("unknown address type")
	}
	b = append(b, byte(a.Port>>8), byte(a.Port))
	return b, nil
}

// A Server repesents a server for handshake testing.
type Server struct {
	ln net.Listener
}

// Addr rerurns a server address.
func (s *Server) Addr() net.Addr {
	return s.ln.Addr()
}

// TargetAddr returns a fake final destination address.
//
// The returned address is only valid for testing with Server.
func (s *Server) TargetAddr() net.Addr {
	a := s.ln.Addr()
	switch a := a.(type) {
	case *net.TCPAddr:
		if a.IP.To4() != nil {
			return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 5963}
		}
		if a.IP.To16() != nil && a.IP.To4() == nil {
			return &net.TCPAddr{IP: net.IPv6loopback, Port: 5963}
		}
	}
	return nil
}

// Close closes the server.
func (s *Server) Close() error {
	return s.ln.Close()
}

func (s *Server) Serve(authFunc, cmdFunc func(net.Conn, []byte) error) {
	c, err := s.ln.Accept()
	if err != nil {
		return
	}
	defer c.Close()
	go s.Serve(authFunc, cmdFunc)
	b := make([]byte, 512)
	n, err := c.Read(b)
	if err != nil {
		return
	}
	if err := authFunc(c, b[:n]); err != nil {
		return
	}
	n, err = c.Read(b)
	if err != nil {
		return
	}
	if err := cmdFunc(c, b[:n]); err != nil {
		return
	}
}

// NewServer returns a new server.
//
// The provided authFunc and cmdFunc must parse requests and return
// appropriate replies to clients.
func NewServer(addr string) (*Server, error) {
	var err error
	s := new(Server)
	s.ln, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// NoAuthRequired handles a no-authentication-required signaling.
func NoAuthRequired(rw net.Conn, b []byte) error {
	req, err := ParseAuthRequest(b)
	if err != nil {
		return err
	}
	b, err = MarshalAuthReply(req.Version, AuthMethodNotRequired)
	if err != nil {
		return err
	}
	n, err := rw.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return errors.New("short write")
	}
	return nil
}

// NoProxyRequired handles a command signaling without constructing a
// proxy connection to the final destination.
func NoProxyRequired(rw io.ReadWriter, b []byte) error {
	req, err := ParseCmdRequest(b)
	if err != nil {
		return err
	}
	req.Addr.Port += 1
	if req.Addr.Name != "" {
		req.Addr.Name = "boundaddr.doesnotexist"
	} else if req.Addr.IP.To4() != nil {
		req.Addr.IP = net.IPv4(127, 0, 0, 1)
	} else {
		req.Addr.IP = net.IPv6loopback
	}
	b, err = MarshalCmdReply(Version5, StatusSucceeded, &req.Addr)
	if err != nil {
		return err
	}
	n, err := rw.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return errors.New("short write")
	}
	return nil
}
