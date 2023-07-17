package service

type (
	Logger interface {
		loggerDebug
		loggerError
		loggerInfo
	}
	loggerDebug interface {
		LogDebug(v ...any) error
	}

	loggerError interface {
		LogError(v ...any) error
	}

	loggerInfo interface {
		LogInfo(v ...any) error
	}
)
