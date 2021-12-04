package headers

// Headers are referred from https://github.com/go-http-utils/headers and https://en.wikipedia.org/wiki/List_of_HTTP_header_fields.

// Standard header fields.
const (
	AIM                           = "A-IM"                             // Used in requests
	Accept                        = "Accept"                           // Used in requests
	AcceptCH                      = "Accept-CH"                        // Used in responses
	AcceptCharset                 = "Accept-Charset"                   // Used in requests
	AcceptDatetime                = "Accept-Datetime"                  // Used in requests
	AcceptEncoding                = "Accept-Encoding"                  // Used in requests
	AcceptLanguage                = "Accept-Language"                  // Used in requests
	AcceptPatch                   = "Accept-Patch"                     // Used in responses
	AcceptRanges                  = "Accept-Ranges"                    // Used in responses
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials" // Used in responses
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"     // Used in responses
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"     // Used in responses
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"      // Used in responses
	AccessControlExposeHeaders    = "Access-Control-Expose-Headers"    // Used in responses
	AccessControlMaxAge           = "Access-Control-Max-Age"           // Used in responses
	AccessControlRequestHeaders   = "Access-Control-Request-Headers"   // Used in requests
	AccessControlRequestMethod    = "Access-Control-Request-Method"    // Used in requests
	Age                           = "Age"                              // Used in responses
	Allow                         = "Allow"                            // Used in responses
	AltSvc                        = "Alt-Svc"                          // Used in responses
	Authorization                 = "Authorization"                    // Used in requests
	CacheControl                  = "Cache-Control"                    // Used in requests, responses
	Connection                    = "Connection"                       // Used in requests, responses
	ContentDisposition            = "Content-Disposition"              // Used in responses
	ContentEncoding               = "Content-Encoding"                 // Used in requests, responses
	ContentLanguage               = "Content-Language"                 // Used in responses
	ContentLength                 = "Content-Length"                   // Used in requests, responses
	ContentLocation               = "Content-Location"                 // Used in responses
	ContentMD5                    = "Content-MD5"                      // Used in requests, responses
	ContentRange                  = "Content-Range"                    // Used in responses
	ContentType                   = "Content-Type"                     // Used in requests, responses
	Cookie                        = "Cookie"                           // Used in requests
	Date                          = "Date"                             // Used in requests, responses
	DeltaBase                     = "Delta-Base"                       // Used in responses
	ETag                          = "ETag"                             // Used in responses
	Expect                        = "Expect"                           // Used in requests
	Expires                       = "Expires"                          // Used in responses
	Forwarded                     = "Forwarded"                        // Used in requests
	From                          = "From"                             // Used in requests
	HTTP2Settings                 = "HTTP2-Settings"                   // Used in requests
	Host                          = "Host"                             // Used in requests
	IM                            = "IM"                               // Used in responses
	IfMatch                       = "If-Match"                         // Used in requests
	IfModifiedSince               = "If-Modified-Since"                // Used in requests
	IfNoneMatch                   = "If-None-Match"                    // Used in requests
	IfRange                       = "If-Range"                         // Used in requests
	IfUnmodifiedSince             = "If-Unmodified-Since"              // Used in requests
	LastModified                  = "Last-Modified"                    // Used in responses
	Link                          = "Link"                             // Used in responses
	Location                      = "Location"                         // Used in responses
	MaxForwards                   = "Max-Forwards"                     // Used in requests
	Origin                        = "Origin"                           // Used in requests
	P3P                           = "P3P"                              // Used in responses
	Pragma                        = "Pragma"                           // Used in requests, responses
	Prefer                        = "Prefer"                           // Used in requests
	PreferenceApplied             = "Preference-Applied"               // Used in responses
	ProxyAuthenticate             = "Proxy-Authenticate"               // Used in responses
	ProxyAuthorization            = "Proxy-Authorization"              // Used in requests
	PublicKeyPins                 = "Public-Key-Pins"                  // Used in responses
	Range                         = "Range"                            // Used in requests
	Referer                       = "Referer"                          // Used in requests
	RetryAfter                    = "Retry-After"                      // Used in responses
	Server                        = "Server"                           // Used in responses
	SetCookie                     = "Set-Cookie"                       // Used in responses
	StrictTransportSecurity       = "Strict-Transport-Security"        // Used in responses
	TE                            = "TE"                               // Used in requests
	Tk                            = "Tk"                               // Used in responses
	Trailer                       = "Trailer"                          // Used in requests, responses
	TransferEncoding              = "Transfer-Encoding"                // Used in requests, responses
	Upgrade                       = "Upgrade"                          // Used in requests, responses
	UserAgent                     = "User-Agent"                       // Used in requests
	Vary                          = "Vary"                             // Used in responses
	Via                           = "Via"                              // Used in requests, responses
	WWWAuthenticate               = "WWW-Authenticate"                 // Used in responses
	Warning                       = "Warning"                          // Used in requests, responses
	XFrameOptions                 = "X-Frame-Options"                  // Used in responses
)

// Common non-standard header fields.
const (
	ContentSecurityPolicy   = "Content-Security-Policy"   // Used in responses
	DNT                     = "DNT"                       // Used in requests
	ExpectCT                = "Expect-CT"                 // Used in responses
	FrontEndHttps           = "Front-End-Https"           // Used in requests
	NEL                     = "NEL"                       // Used in responses
	PermissionsPolicy       = "Permissions-Policy"        // Used in responses
	ProxyConnection         = "Proxy-Connection"          // Used in requests
	Refresh                 = "Refresh"                   // Used in responses
	ReportTo                = "Report-To"                 // Used in responses
	SaveData                = "Save-Data"                 // Used in requests
	Status                  = "Status"                    // Used in responses
	TimingAllowOrigin       = "Timing-Allow-Origin"       // Used in responses
	UpgradeInsecureRequests = "Upgrade-Insecure-Requests" // Used in requests
	XATTDeviceId            = "X-ATT-DeviceId"            // Used in requests
	XContentDuration        = "X-Content-Duration"        // Used in responses
	XContentSecurityPolicy  = "X-Content-Security-Policy" // Used in responses
	XContentTypeOptions     = "X-Content-Type-Options"    // Used in responses
	XCorrelationID          = "X-Correlation-ID"          // Used in requests, responses
	XCsrfToken              = "X-Csrf-Token"              // Used in requests
	XForwardedFor           = "X-Forwarded-For"           // Used in requests
	XForwardedHost          = "X-Forwarded-Host"          // Used in requests
	XForwardedProto         = "X-Forwarded-Proto"         // Used in requests
	XHttpMethodOverride     = "X-Http-Method-Override"    // Used in requests
	XPoweredBy              = "X-Powered-By"              // Used in responses
	XRateLimitLimit         = "X-RateLimit-Limit"         // Used in responses
	XRateLimitRemaining     = "X-RateLimit-Remaining"     // Used in responses
	XRateLimitReset         = "X-RateLimit-Reset"         // Used in responses
	XRealIP                 = "X-Real-IP"                 // Used in requests
	XRedirectBy             = "X-Redirect-By"             // Used in responses
	XRequestID              = "X-Request-ID"              // Used in requests, responses
	XRequestedWith          = "X-Requested-With"          // Used in requests
	XUACompatible           = "X-UA-Compatible"           // Used in responses
	XUIDH                   = "X-UIDH"                    // Used in requests
	XWapProfile             = "X-Wap-Profile"             // Used in requests
	XWebKitCSP              = "X-WebKit-CSP"              // Used in responses
	XXSSProtection          = "X-XSS-Protection"          // Used in responses
)
