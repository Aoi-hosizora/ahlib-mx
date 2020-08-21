package xstatus

type JwtStatus int8

const (
	JwtExpired   JwtStatus = iota // exp, ExpiresAt
	JwtIssuer                     // iss, Issuer
	JwtNotIssued                  // iat, IssuedAt
	JwtNotValid                   // nbf, NotBefore
	JwtID                         // jti, Id
	JwtAudience                   // aud, Audience
	JwtSubject                    // sub, Subject
	JwtInvalid                    // generic
	JwtUserErr                    // user
	JwtFailed                     // server error
	JwtTagA                       // tag a
	JwtTagB                       // tag b
	JwtTagC                       // tag c
	JwtTagD                       // tag d
	JwtTagE                       // tag e
)

func (j JwtStatus) String() string {
	switch j {
	case JwtExpired:
		return "jwt-expired"
	case JwtIssuer:
		return "jwt-issuer"
	case JwtNotIssued:
		return "jwt-not-issued"
	case JwtNotValid:
		return "jwt-not-valid"
	case JwtInvalid:
		return "jwt-invalid"
	case JwtID:
		return "jwt-id"
	case JwtAudience:
		return "jwt-audience"
	case JwtSubject:
		return "jwt-subject"
	case JwtUserErr:
		return "jwt-user-err"
	case JwtFailed:
		return "jwt-failed"
	case JwtTagA:
		return "jwt-tag-a"
	case JwtTagB:
		return "jwt-tag-b"
	case JwtTagC:
		return "jwt-tag-c"
	case JwtTagD:
		return "jwt-tag-d"
	case JwtTagE:
		return "jwt-tag-e"
	default:
		return "jwt-?"
	}
}
