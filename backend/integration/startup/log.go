package startup

import "github.com/johnwongx/webook/backend/pkg/logger"

func InitLog() logger.Logger {
	return logger.NewNopLogger()
}
