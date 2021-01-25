package xserverchan

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/Aoi-hosizora/go-serverchan"
	"github.com/sirupsen/logrus"
)

// ServerchanLogrus is a logrus logger used by serverchan.
type ServerchanLogrus struct {
	logger  *logrus.Logger
	logMode bool
}

// NewServerchanLogrus creates a new ServerchanLogrus with logrus.Logger.
func NewServerchanLogrus(logger *logrus.Logger, logMode bool) *ServerchanLogrus {
	return &ServerchanLogrus{logger: logger, logMode: logMode}
}

// Log does log using sckey, title, code and error.
func (s *ServerchanLogrus) Log(sckey string, title string, code int32, err error) {
	if !s.logMode {
		return
	}

	sckey = xstring.DefaultMaskToken(sckey)
	title = xstring.DefaultMaskToken(title)

	if err != nil {
		if err == serverchan.ErrNotSuccess {
			s.logger.Error(fmt.Sprintf("[Serverchan] Send to %s error: %v", sckey, err))
		} else {
			s.logger.WithFields(map[string]interface{}{
				"module":    "serverchan",
				"sckeyMask": sckey,
				"code":      code,
			}).Error(fmt.Sprintf("[Serverchan] Send to %s error: %v", sckey, err))
		}
	} else {
		s.logger.WithFields(map[string]interface{}{
			"module":    "serverchan",
			"sckeyMask": sckey,
			"titleMask": title,
			"code":      0,
		}).Info(fmt.Sprintf("[Serverchan] < | %s | %s", sckey, title))
	}
}
