package inbound

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/listener"
	"github.com/sipt/shuttle/pkg/socks"
	"github.com/sirupsen/logrus"

	connpkg "github.com/sipt/shuttle/conn"
)

func init() {
	Register("socks", newSocksInbound)
}

func newSocksInbound(addr string, params map[string]string) (listen func(context.Context, listener.HandleFunc), err error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	addrPtr := &socks.Addr{}
	if host == "" {
		addrPtr.IP = net.IPv4(127, 0, 0, 1)
	}
	addrPtr.Port, err = strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	server, err := socks.NewServer(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "init socks server on [%s] failed", addr)
	}
	logrus.WithField("addr", fmt.Sprintf("socks5://%s", addr)).Info("socks listen starting")
	authType, ok := params[ParamsKeyAuthType]
	authFunc := socks.NoAuthRequired
	if ok && authType == AuthTypeBasic {
		authFunc, err = newAuthFunc(params)
		if err != nil {
			return
		}
	}
	return func(ctx context.Context, handleFunc listener.HandleFunc) {
		cmdFunc := NewCmdFunc(ctx, addrPtr, handleFunc)
		server.Serve(authFunc, cmdFunc)
		<-ctx.Done()
		_ = server.Close()
	}, nil
}

func newAuthFunc(params map[string]string) (func(net.Conn, []byte) error, error) {
	username := params[ParamsKeyUser]
	if len(username) == 0 {
		return nil, errors.New("[user] is empty")
	}
	password := params[ParamsKeyPassword]
	if len(password) == 0 {
		return nil, errors.New("[password] is empty")
	}
	authTyp := socks.AuthMethodUsernamePassword
	return func(conn net.Conn, b []byte) error {
		authRequest, err := socks.ParseAuthRequest(b)
		if err != nil {
			return err
		}
		if authRequest.Version != socks.Version5 {
			return errors.Errorf("not support version: %d", authRequest.Version)
		}
		replyAuthType := socks.AuthMethodNoAcceptableMethods
		for _, v := range authRequest.Methods {
			if v == authTyp {
				replyAuthType = authTyp
			}
		}
		reply, err := socks.MarshalAuthReply(socks.Version5, replyAuthType)
		if err != nil {
			return err
		}
		_, err = conn.Write(reply)
		if err != nil {
			return err
		}
		b = make([]byte, 512)
		n, err := conn.Read(b)
		if err != nil {
			return err
		}
		pass := socks.BasicAuthorization(username, password, b[:n])
		b, err = socks.MarshalAuthResponse(0x01, pass)
		if err != nil {
			return err
		}
		_, err = conn.Write(b)
		if err == nil && !pass {
			return errors.Errorf("unauthorized")
		}
		return err
	}, nil
}

func NewCmdFunc(ctx context.Context, addr *socks.Addr, handle listener.HandleFunc) func(net.Conn, []byte) error {
	return func(conn net.Conn, b []byte) error {
		cmdReq, err := socks.ParseCmdRequest(b)
		if err != nil {
			return err
		}
		if cmdReq.Version != socks.Version5 {
			return errors.Errorf("not support version: %d", cmdReq.Version)
		}
		req := &request{
			domain: cmdReq.Addr.Name,
			ip:     cmdReq.Addr.IP,
			port:   cmdReq.Addr.Port,
		}
		switch cmdReq.Cmd {
		case socks.CmdConnect:
			req.network = "tcp"
			b, err = socks.MarshalCmdReply(socks.Version5, socks.StatusSucceeded, addr)
			if err != nil {
				return err
			}
			_, err = conn.Write(b)
			if err != nil {
				return err
			}
			ctx = context.WithValue(ctx, constant.KeyRequestInfo, req)
			handle(connpkg.NewConn(conn, ctx))
		case socks.CmdUDPAssociate:
			req.network = "udp"
			err = udpAssociate(ctx, conn, req, handle)
			if err != nil {
				return err
			}
		default:
			return errors.Errorf("not support %s", cmdReq.Cmd)
		}
		return nil
	}
}

func udpAssociate(ctx context.Context, conn net.Conn, req *request, handle listener.HandleFunc) error {
	dst := &net.UDPAddr{
		IP:   req.ip,
		Port: req.port,
	}
	u, err := socks.NewUDPServer("", dst)
	if err != nil {
		return errors.Wrap(err, "create udp listen failed")
	}
	go u.Serve(func(pc net.PacketConn, remote net.Addr, dst net.Addr, b []byte) error {
		ctx = context.WithValue(ctx, constant.KeyRequestInfo, req)
		handle(connpkg.NewUDPConn(pc, ctx, remote, b))
		return nil
	})
	addr := u.LocalAddr().(*net.UDPAddr)
	fmt.Println(addr, u.LocalAddr())
	b, err := socks.MarshalCmdReply(socks.Version5, socks.StatusSucceeded, &socks.Addr{
		IP:   addr.IP,
		Port: addr.Port,
	})
	if err != nil {
		return err
	}
	_, err = conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}
