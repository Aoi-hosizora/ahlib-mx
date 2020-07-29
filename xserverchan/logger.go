package xserverchan

import (
	"github.com/Aoi-hosizora/go-serverchan"
	"github.com/sirupsen/logrus"
)

type ServerchanLogrus struct {
	logger  *logrus.Logger
	logMode bool
}

// noinspection GoUnusedExportedFunction
func NewServerChanLogrus(logger *logrus.Logger, logMode bool) *ServerchanLogrus {
	return &ServerchanLogrus{logger: logger, logMode: logMode}
}

func (s *ServerchanLogrus) Log(sckey string, title string, code int32, err error) {
	if !s.logMode {
		return
	}

	sckey = serverchan.Mask(sckey)
	title = serverchan.Mask(title)

	if err != nil {
		if !serverchan.IsResponseError(err) {
			s.logger.Errorf("[Serverchan] failed to send message to %s: %v", sckey, err)
		} else {
			s.logger.WithFields(map[string]interface{}{
				"module":    "serverchan",
				"sckeyMask": sckey,
				"code":      code,
			}).Errorf("[Serverchan] failed to send message to %s: %v", sckey, err)
		}
	} else {
		s.logger.WithFields(map[string]interface{}{
			"module":    "serverchan",
			"sckeyMask": sckey,
			"titleMask": title,
			"code":      0,
		}).Infof("[Serverchan] <- %s | %s", sckey, title)
	}
}

// Please use serverchan.DefaultLogger
// noinspection GoUnusedType
type serverchanLogger struct{}
