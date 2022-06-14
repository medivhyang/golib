package http

import "net/http"

const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderXRequestedWith      = "X-Requested-With"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"
)

const (
	StatusContinue                      = http.StatusContinue
	StatusSwitchingProtocols            = http.StatusSwitchingProtocols
	StatusProcessing                    = http.StatusProcessing
	StatusEarlyHints                    = http.StatusEarlyHints
	StatusOK                            = http.StatusOK
	StatusCreated                       = http.StatusCreated
	StatusAccepted                      = http.StatusAccepted
	StatusNonAuthoritativeInfo          = http.StatusNonAuthoritativeInfo
	StatusNoContent                     = http.StatusNoContent
	StatusResetContent                  = http.StatusResetContent
	StatusPartialContent                = http.StatusPartialContent
	StatusMultiStatus                   = http.StatusMultiStatus
	StatusAlreadyReported               = http.StatusAlreadyReported
	StatusIMUsed                        = http.StatusIMUsed
	StatusMultipleChoices               = http.StatusMultipleChoices
	StatusMovedPermanently              = http.StatusMovedPermanently
	StatusFound                         = http.StatusFound
	StatusSeeOther                      = http.StatusSeeOther
	StatusNotModified                   = http.StatusNotModified
	StatusUseProxy                      = http.StatusUseProxy
	StatusTemporaryRedirect             = http.StatusTemporaryRedirect
	StatusPermanentRedirect             = http.StatusPermanentRedirect
	StatusBadRequest                    = http.StatusBadRequest
	StatusUnauthorized                  = http.StatusUnauthorized
	StatusPaymentRequired               = http.StatusPaymentRequired
	StatusForbidden                     = http.StatusForbidden
	StatusNotFound                      = http.StatusNotFound
	StatusMethodNotAllowed              = http.StatusMethodNotAllowed
	StatusNotAcceptable                 = http.StatusNotAcceptable
	StatusProxyAuthRequired             = http.StatusProxyAuthRequired
	StatusRequestTimeout                = http.StatusRequestTimeout
	StatusConflict                      = http.StatusConflict
	StatusGone                          = http.StatusGone
	StatusLengthRequired                = http.StatusLengthRequired
	StatusPreconditionFailed            = http.StatusPreconditionFailed
	StatusRequestEntityTooLarge         = http.StatusRequestEntityTooLarge
	StatusRequestURITooLong             = http.StatusRequestURITooLong
	StatusUnsupportedMediaType          = http.StatusUnsupportedMediaType
	StatusRequestedRangeNotSatisfiable  = http.StatusRequestedRangeNotSatisfiable
	StatusExpectationFailed             = http.StatusExpectationFailed
	StatusTeapot                        = http.StatusTeapot
	StatusMisdirectedRequest            = http.StatusMisdirectedRequest
	StatusUnprocessableEntity           = http.StatusUnprocessableEntity
	StatusLocked                        = http.StatusLocked
	StatusFailedDependency              = http.StatusFailedDependency
	StatusTooEarly                      = http.StatusTooEarly
	StatusUpgradeRequired               = http.StatusUpgradeRequired
	StatusPreconditionRequired          = http.StatusPreconditionRequired
	StatusTooManyRequests               = http.StatusTooManyRequests
	StatusRequestHeaderFieldsTooLarge   = http.StatusRequestHeaderFieldsTooLarge
	StatusUnavailableForLegalReasons    = http.StatusUnavailableForLegalReasons
	StatusInternalServerError           = http.StatusInternalServerError
	StatusNotImplemented                = http.StatusNotImplemented
	StatusBadGateway                    = http.StatusBadGateway
	StatusServiceUnavailable            = http.StatusServiceUnavailable
	StatusGatewayTimeout                = http.StatusGatewayTimeout
	StatusHTTPVersionNotSupported       = http.StatusHTTPVersionNotSupported
	StatusVariantAlsoNegotiates         = http.StatusVariantAlsoNegotiates
	StatusInsufficientStorage           = http.StatusInsufficientStorage
	StatusLoopDetected                  = http.StatusLoopDetected
	StatusNotExtended                   = http.StatusNotExtended
	StatusNetworkAuthenticationRequired = http.StatusNetworkAuthenticationRequired
)

func StatusText(code int) string {
	return http.StatusText(code)
}

const (
	MIMEText = "text/plain"
	MIMEHTML = "text/html"
	MIMEXML  = "application/xml"
	MIMEJSON = "application/json"
)

var (
	ErrResponseBodyHasRead = errorf("response body has read")
	ErrUnknownError        = errorf("unknown error")
	ErrRequirePointerType  = errorf("require pointer type")
	ErrRequireSliceType    = errorf("require slice type")
	ErrRequireStructType   = errorf("require struct type")
	ErrCannotSetValue      = errorf("can not set value")
)

var BindingTagKey = "binding"
