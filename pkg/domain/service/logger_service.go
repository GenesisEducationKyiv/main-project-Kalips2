package service

type (
	Logger interface {
		loggerDebug
		loggerError
		loggerInfo
	}
	loggerDebug interface {
		LogDebug(v ...any)
	}

	loggerError interface {
		LogError(v ...any)
	}

	loggerInfo interface {
		LogInfo(v ...any)
	}
)
