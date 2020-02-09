package constant

const (
	KeyProtocol    = "__protocol"
	KeyRequestInfo = "__request"
	KeyRule        = "__rule"
	KeyNamespace   = "__namespace"
	KeyProfile     = "__profile"
	KeyUseTLS      = "__tls" // 当开启mitm时，存在ctx中，type：bool

	ProtocolHTTP  = "http"
	ProtocolHTTPS = "https"

	ModeDirect = "ModeDirect"
	ModeGlobal = "ModeGlobal"
	ModeRule   = "ModeRule"
)
