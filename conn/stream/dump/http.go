package dump

import (
	"context"
	"io"

	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/conn/stream"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sirupsen/logrus"
)

func init() {
	stream.RegisterStream("data-dump", newHTTPDump)
}

var keyEncode = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb2dJQkFBS0NBUUVBc3BSdGM2SmFtQkVQNU1VM1NVVzJoenRHNENjSmRjNGpMS2hPZnY4bGZvWCtSUVhhCjBhdlZjNGpWTXZ5QVlXNGNZTXhFYkZ6TUYxbDUyd0lSWE52bndsQVI4RDJoYnVJd29MMkM2UExQVjR1Y0NpdnMKUFF3SjdjRTRFYnZLdzRkOGFPWWJXWnlCUUtZY25Ic2FMV29GR3o3OVhHak0vRkRZcGZ4ZCtvaXYxZ1F2QXRrQQoyaW04NGYvWk9lRzRwaXRLT3JDYWNOMThDWFJCV1d4N3ZETGEwNTFmM205TjJKUnJyZGMxa0VKaGh3UUV5WnRCCmdXYWxIMDZDTXpBaFB5MVdYc3paWjJrSStLQmMyQ255VkdvR2pGUDYvSU15dmtyNU82QnlxMDRFOEk4V3RZcXgKSTJ0K2lVRmpqVGE3c0F1Y3RFVHFLQ0xoU0FjN3AwN2ZWa3FjdXdJREFRQUJBb0lCQUdBc25iR2I2MHhnUy8zNQoxR2VLdXQyanArMEtPUWNQNkZPaHBQeXlMcUF3UzVzaXB4RXFpTDg3SHc3aGU4WjlCWjJBQlEyVEFIdEd2ZUNjCkFYdlFGc1hJVjVEWnNEcEdhTWY0cUNzS3NXM0ZpMWpUQk54dndsMGdKVEV2d09pQzdCYVdibjVaVWliZUR5U2IKQzZNUHFRWmVheGE4ZmtFWXpVUy9ZR0dRQVpxeEtiWkxVSzFRMnBKYTZDUDE4d0RuZmROcDhRRk92aDE5eUY1ZgpRd1RGOTlDTXloR2VITFE0bmxaRnNNSE5hZ3JqN3l4d3dMeDVIdFVnd0wzOVdXMVBMQXNYZUs2N0dPeXFrRnJsCmE2NDBsdE5NVGlVVzlHMEt4NzdGWDJpaWc1MlN1VWh1RjlvcC9XNHJmOFRyVTBzQ281c1h1RTIrakU1T0ZIWCsKdndwM2paa0NnWUVBM0JDSHAxVXpwa0VKditpdXZxMWxLd09SbHFkMkxRcC94Zlo5bHV5VSs1VS9NTjRGYVdHOQo2VEVnVUdpY1A2eGZjMVc0K2pGdnYwV2kwTmEwWGI1eG93YTNZNUF4bG5ab0hoS1Z1d24yNTRYQmR5TU9LRHdlCkVCVFY2RDgwNEtQRlJ4TE9oelVMVFBIaHJPbmdrczhHZ25maGVqMjN5VUNPSHJ1MnNsTFF0TjBDZ1lFQXo3MncKVHVPUlNzcTZqN25SWkk3UzB2eEJtaDc2WW9xZGZScTFmeXdkd1Y4VFF1VGhIZzhmYmNtWEkvLzNPVDVIa1ZzTQp0VTRGaWkwRzR6eGRweCt1Z3dPd2pseTBoOGlzWEZxajN4RHgydlNtR2tNdmg4VThrSWtvVmxESWpOcC9lWDBuCldVeTBRcnlSak44UXlmdktzL0I3V1JKNHFDbzJnSEN6NldCYkVuY0NnWUFYVTcwOWRKK2orUEx5bjlTZUs3MDAKb1EwMnZndWQzS1lNc3dNL0UxYjdrQ2VCbzVkSlEyNGhJTzcrOXdmUkRCR1dKVGtWZEZZWEhXZVQ0WjUrN1dnWQpVdWJ2cStKRnc4bG5ucXEyaCtqZlErTnRJSThvbnE3Rkg2QkpIU0lheWVGb2xrckVORkE5V01xR1RNaGNaNHVXCkd1VVEweWYvTWxPZVdHR1daNGJ1RlFLQmdCUFhCdWFSMzBkb1V5YjAxU1dvYWtRU0tXWEJ2YUg1b1E3WXBTclAKR014bCt4M3hZL3FOOFM5NENFSTg2T3lEb1N3bHFQSUwwSVdneFQ2Z2ZrVSt4bGptMms3T1ZjTitDOUFLTEFwYQp3TzVyWFEyM0N1d0pqejR5aXpLckptd2xWZlZSV1pleXRxaUUvOVdYWERBZUp2N0dZZEZnN1RzS1JRaEJPejEzCm9Wc0RBb0dBQ2pPc0FPQndGa1htdHUxdEh1S09Cakxqa0xDVFRWK29IYU9hQ2VPUFpSRWl3dlNqMFNKYzlNb0MKQ0ZJTDNnVHRpRmN0eVF6SStOVWZtdTQ0RWlTTGxTVzJBTGhkK1FDYlkvTytkeFZIN0wvbmt4Y21ySjdhMEQ5TAp4cjJEUnFjTml1M3ZuQ1lGaUNRUUYrTlpvbXFuNGtmY0RJNkVUNERpYUNQY0ZXUnE1YVU9Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg"
var caEncode = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURIekNDQWdlZ0F3SUJBZ0lCQVRBTkJna3Foa2lHOXcwQkFRc0ZBREF4TVJBd0RnWURWUVFLRXdkVGFIVjAKZEd4bE1SMHdHd1lEVlFRREV4UlRhSFYwZEd4bElFZGxibVZ5WVhSbFpDQkRRVEFlRncweU1EQXlNRGt3T1RBMgpNVE5hRncweU5UQXlNRGt3T1RBMk1UTmFNREV4RURBT0JnTlZCQW9UQjFOb2RYUjBiR1V4SFRBYkJnTlZCQU1UCkZGTm9kWFIwYkdVZ1IyVnVaWEpoZEdWa0lFTkJNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUIKQ2dLQ0FRRUFzcFJ0YzZKYW1CRVA1TVUzU1VXMmh6dEc0Q2NKZGM0akxLaE9mdjhsZm9YK1JRWGEwYXZWYzRqVgpNdnlBWVc0Y1lNeEViRnpNRjFsNTJ3SVJYTnZud2xBUjhEMmhidUl3b0wyQzZQTFBWNHVjQ2l2c1BRd0o3Y0U0CkVidkt3NGQ4YU9ZYldaeUJRS1ljbkhzYUxXb0ZHejc5WEdqTS9GRFlwZnhkK29pdjFnUXZBdGtBMmltODRmL1oKT2VHNHBpdEtPckNhY04xOENYUkJXV3g3dkRMYTA1MWYzbTlOMkpScnJkYzFrRUpoaHdRRXladEJnV2FsSDA2QwpNekFoUHkxV1hzelpaMmtJK0tCYzJDbnlWR29HakZQNi9JTXl2a3I1TzZCeXEwNEU4SThXdFlxeEkydCtpVUZqCmpUYTdzQXVjdEVUcUtDTGhTQWM3cDA3ZlZrcWN1d0lEQVFBQm8wSXdRREFPQmdOVkhROEJBZjhFQkFNQ0FvUXcKSFFZRFZSMGxCQll3RkFZSUt3WUJCUVVIQXdFR0NDc0dBUVVGQndNQ01BOEdBMVVkRXdFQi93UUZNQU1CQWY4dwpEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBQXFnT1VrUldSWXRUYmlXZFBLMDFqTFBvWFJzenpMQXdBMzFUcXFKCjNDdms3WTJjU1RMVHpiNlFCWlZOVlJPKzhhSXJReTgvOEYwZ1l5TytlTHJPMk5Ec0dkODNoZGF5UTUyRVViYmkKZVJabFZ2MFZOam1zN2FITmpZL25pTFVWUlJsMHpXYldJSkNMME1KY2F1Qkcxd1VHR0xRNXk3dVZab3pYaHZxMQpOVVRQcDBaK1N6cjdLelRVdkZPTGhabzd5dkptajdoWmJ1OUQ4NWNuOUk0R0JhUFA2Ni9ick5tcFFhMWJsTk1RCkUwRmp2bFVwZ1hta3JtRElVUVZ5NmQxMVlTUWVRSnRUWmxTNlJaUnl0M2dXdEYwVGZkNHFNY1VXeXQrWXB6L1kKRUMveWJrUkxOdVowdmZ0Q3NNUmlrcy9neEtndnRrMG9xNExrVWpIVktGVU8vUHM9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"

func checkAllowDump(c context.Context, protocol ...string) bool {
	if len(protocol) > 0 && len(protocol[0]) > 0 {
		p := protocol[0]
		return allowDump && (p == constant.ProtocolHTTP || p == constant.ProtocolHTTPS && mitmEabled)
	} else {
		p, ok := c.Value(constant.KeyProtocol).(string)
		return ok && allowDump && (p == constant.ProtocolHTTP || p == constant.ProtocolHTTPS && mitmEabled)
	}
}

func newHTTPDump(ctx context.Context, params map[string]string) (stream.DecorateFunc, error) {
	allowDump = params["enabled"] == "true"
	err := InitDumpStorage(params["dump_path"])
	if err != nil {
		return nil, err
	}
	err = InitMITM(keyEncode, caEncode)
	if err != nil {
		return nil, err
	}
	go AutoSave(ctx)
	return func(c conn.ICtxConn) conn.ICtxConn {
		p, ok := c.Value(constant.KeyProtocol).(string)
		if !ok {
			return c
		}
		if !checkAllowDump(c, p) {
			return c
		}
		if p == constant.ProtocolHTTPS {
			lc, err := Mitm(c)
			if err != nil {
				logrus.WithError(err).Error("call mitm failed")
				_ = c.Close()
			}
			c = lc
		}
		rc := &httpDumpConn{
			ICtxConn: c,
		}
		var wt io.WriterTo
		var rf io.ReaderFrom
		if _, ok := c.(io.WriterTo); ok {
			wt = &recordTrafficConnWithWriteTo{
				ICtxConn: rc,
				WriterTo: &writeTo{
					ICtxConn: c,
				},
			}
		}
		if _, ok := c.(io.ReaderFrom); ok {
			rf = &recordTrafficConnWithReadFrom{
				ICtxConn: rc,
				ReaderFrom: &readFrom{
					ICtxConn: c,
				},
			}
		}
		switch {
		case wt == nil && rf == nil:
			return rc
		case wt != nil && rf == nil:
			return wt.(conn.ICtxConn)
		case wt == nil && rf != nil:
			return rf.(conn.ICtxConn)
		default:
			return &recordTrafficConnWithWriteToAndReadFrom{
				ICtxConn:   rc,
				WriterTo:   wt,
				ReaderFrom: rf,
			}
		}
	}, nil
}

type httpDumpConn struct {
	conn.ICtxConn
}

func (h *httpDumpConn) Read(b []byte) (n int, err error) {
	n, err = h.ICtxConn.Read(b)
	if n > 0 {
		if id, ok := recordID(h); ok {
			err := SaveRequest(id, b[:n])
			if err != nil {
				logrus.WithField("record_id", id).WithError(err).Error("[data_dump] save request failed")
			}
		}
	}
	return
}

func (h *httpDumpConn) Write(b []byte) (n int, err error) {
	n, err = h.ICtxConn.Write(b)
	if n > 0 {
		if id, ok := recordID(h); ok {
			err := SaveResponse(id, b[:n])
			if err != nil {
				logrus.WithField("record_id", id).WithError(err).Error("[data_dump] save response failed")
			}
		}
	}
	return
}

func (h *httpDumpConn) Close() (err error) {
	err = h.ICtxConn.Close()
	if id, ok := recordID(h); ok {
		err := CloseFiles(id)
		if err != nil {
			logrus.WithField("record_id", id).WithError(err).Error("[data_dump] close files failed")
		}
	}
	return err
}

type recordTrafficConnWithWriteTo struct {
	conn.ICtxConn
	io.WriterTo
}

type recordTrafficConnWithReadFrom struct {
	conn.ICtxConn
	io.ReaderFrom
}
type recordTrafficConnWithWriteToAndReadFrom struct {
	conn.ICtxConn
	io.WriterTo
	io.ReaderFrom
}

type writeTo struct {
	conn.ICtxConn
}

func (r *writeTo) WriteTo(w io.Writer) (n int64, err error) {
	wr := &writer{Writer: w, ICtxConn: r.ICtxConn}
	n, err = r.ICtxConn.(io.WriterTo).WriteTo(wr)
	return n, err
}

type readFrom struct {
	conn.ICtxConn
}

func (r *readFrom) ReadFrom(re io.Reader) (n int64, err error) {
	rr := &reader{Reader: re, ICtxConn: r.ICtxConn}
	n, err = r.ICtxConn.(io.ReaderFrom).ReadFrom(rr)
	return n, err
}

type writer struct {
	io.Writer
	conn.ICtxConn
}

func (w *writer) Write(b []byte) (n int, err error) {
	n, err = w.Writer.Write(b)
	if n > 0 {
		if id, ok := recordID(w); ok {
			err := SaveRequest(id, b[:n])
			if err != nil {
				logrus.WithField("record_id", id).WithError(err).Error("[data_dump] save request failed")
			}
		}
	}
	return n, err
}

type reader struct {
	io.Reader
	conn.ICtxConn
}

func (r *reader) Read(b []byte) (n int, err error) {
	n, err = r.Reader.Read(b)
	if n > 0 {
		if id, ok := recordID(r); ok {
			err := SaveResponse(id, b[:n])
			if err != nil {
				logrus.WithField("record_id", id).WithError(err).Error("[data_dump] save response failed")
			}
		}
	}
	return n, err
}

func recordID(ctx context.Context) (int64, bool) {
	req, ok := ctx.Value(constant.KeyRequestInfo).(typ.RequestInfo)
	if !ok {
		return 0, false
	}
	return req.ID(), true
}
