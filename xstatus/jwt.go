package xstatus

type JwtStatus int8

const (
	JwtSuccess   JwtStatus = iota // success
	JwtExpired                    // exp, ExpiresAt
	JwtNotIssued                  // iat, IssuedAt
	JwtNotValid                   // nbf, NotBefore
	JwtIssuer                     // iss, Issuer
	JwtSubject                    // sub, Subject
	JwtAudience                   // aud, Audience
	JwtInvalid                    // generic
	JwtBlank                      // blank
	JwtNotFound                   // not found
	JwtUserErr                    // user
	JwtFailed                     // server error
	JwtTagA                       // tag a
	JwtTagB                       // tag b
	JwtTagC                       // tag c
)

func (j JwtStatus) String() string {
	switch j {
	case JwtSuccess:
		return "jwt-success"
	case JwtExpired:
		return "jwt-expired"
	case JwtNotIssued:
		return "jwt-not-issued"
	case JwtNotValid:
		return "jwt-not-valid"
	case JwtIssuer:
		return "jwt-issuer"
	case JwtSubject:
		return "jwt-subject"
	case JwtAudience:
		return "jwt-audience"
	case JwtInvalid:
		return "jwt-invalid"
	case JwtBlank:
		return "jwt-blank"
	case JwtNotFound:
		return "jwt-not-found"
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
	default:
		return "jwt-?"
	}
}
