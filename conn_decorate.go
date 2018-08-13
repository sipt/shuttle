package shuttle

import (
	"net"
	"time"
	"bytes"
	"github.com/sipt/shuttle/pool"
	"github.com/sipt/shuttle/util"
	"net/http"
	"bufio"
)

var defaultTimeOut = 30 * time.Second

//
func DefaultDecorate(c net.Conn, network string) (IConn, error) {
	return &DefaultConn{
		Conn:    c,
		ID:      util.NextID(),
		Network: network,
	}, nil
}

func DefaultDecorateForTls(c net.Conn, network string, id int64) (IConn, error) {
	return &DefaultConn{
		Conn:    c,
		ID:      id,
		Network: network,
	}, nil
}

type DefaultConn struct {
	net.Conn
	ID      int64
	Network string
}

func (c *DefaultConn) GetID() int64 {
	return c.ID
}

func (c *DefaultConn) Flush() (int, error) {
	return 0, nil
}

func (c *DefaultConn) GetNetwork() string {
	return c.Network
}

func (c *DefaultConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	Logger.Tracef("[ID:%d] Read Data: %v", c.GetID(), b[:n])
	return
}

func (c *DefaultConn) Write(b []byte) (n int, err error) {
	Logger.Tracef("[ID:%d] Write Data: %v", c.GetID(), b)
	return c.Conn.Write(b)
}
func (c *DefaultConn) Close() error {
	Logger.Tracef("[ID:%d] close connection", c.GetID())
	go storage.Put(c.GetID(), RecordStatus, RecordStatusCompleted)
	return c.Conn.Close()
}

//超时装饰
func TimerDecorate(c IConn, rto, wto time.Duration) (IConn, error) {
	if rto == 0 {
		rto = defaultTimeOut
	}
	if wto == 0 {
		wto = defaultTimeOut
	}
	return &TimerConn{
		IConn:        c,
		ReadTimeOut:  rto,
		WriteTimeOut: wto,
	}, nil
}

type TimerConn struct {
	IConn
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
}

func (c *TimerConn) Read(b []byte) (n int, err error) {
	if c.ReadTimeOut > -1 {
		c.SetReadDeadline(time.Now().Add(c.ReadTimeOut))
	}
	n, err = c.IConn.Read(b)
	return
}

func (c *TimerConn) Write(b []byte) (n int, err error) {
	if c.WriteTimeOut > -1 {
		c.SetWriteDeadline(time.Now().Add(c.WriteTimeOut))
	}
	n, err = c.IConn.Write(b)
	return
}

//缓冲装饰
func BufferDecorate(c IConn) (IConn, error) {
	return &BufferConn{
		IConn:  c,
		buffer: bytes.NewBuffer(pool.GetBuf()[:0]),
	}, nil
}

type BufferConn struct {
	IConn
	buffer *bytes.Buffer
}

func (c *BufferConn) Write(b []byte) (n int, err error) {
	return c.buffer.Write(b)
}

func (c *BufferConn) Flush() (n int, err error) {
	n, err = c.IConn.Write(c.buffer.Bytes())
	if err != nil {
		return
	}
	c.buffer.Reset()
	n, err = c.IConn.Flush()
	return
}

//实时写出
func RealTimeDecorate(c IConn) (IConn, error) {
	return &RealTimeFlush{
		IConn: c,
	}, nil
}

type RealTimeFlush struct {
	IConn
}

func (r *RealTimeFlush) Write(b []byte) (n int, err error) {
	n, err = r.IConn.Write(b)
	if err != nil {
		return
	}
	_, err = r.IConn.Flush()
	return
}

//导出装饰器
func DumperDecorate(c IConn, allowDump bool, template *Record) (IConn, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 4096))
	dumper := &HttpDumper{
		IConn:      c,
		allowDump:  allowDump,
		template:   template,
		readWriter: bufio.NewReadWriter(bufio.NewReader(buffer), bufio.NewWriter(buffer)),
	}
	return dumper, nil
}

type HttpDumper struct {
	id, oldID  int64
	req        *http.Request
	allowDump  bool
	template   *Record
	readWriter *bufio.ReadWriter
	IConn
}

func (d *HttpDumper) Read(b []byte) (n int, err error) {
	n, err = d.IConn.Read(b)
	if d.allowDump {
		go func(id, oldID int64) {
			dump.WriteResponse(id, b[:n])
			if oldID != 0 && oldID != id {
				dump.Complete(oldID)
			}
		}(d.id, d.oldID)
	}
	d.oldID = d.id
	return
}

func (d *HttpDumper) Write(b []byte) (n int, err error) {
	return d.readWriter.Write(b)
}

func (d *HttpDumper) BufferWrite() (n int, err error) {
	d.req, err = http.ReadRequest(d.readWriter.Reader)
	if err != nil {
		return 0, err
	}
	d.id = util.NextID()
	err = d.req.Write(d.IConn)
	if d.id == 0 {
		d.id = d.GetID()
	} else {
		d.id = util.NextID()
	}
	if d.allowDump {
		go func(id int64, req *http.Request) {
			record := *d.template
			record.ID = d.id
			record.URL = d.req.URL.String()
			record.Status = RecordStatusActive
			record.Created = time.Now()
			recordChan <- &record
			dump.InitDump(d.id)
			writer := bytes.NewBuffer(pool.GetBuf()[:0])
			req.Write(writer)
			dump.WriteRequest(d.id, writer.Bytes())
		}(d.id, d.req)
	}
	return
}

func (d *HttpDumper) Close() (err error) {
	err = d.IConn.Close()
	if d.allowDump {
		go dump.Complete(d.id)
	}
	return
}
