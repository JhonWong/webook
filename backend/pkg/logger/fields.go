package logger

func Error(err error) Field {
	return Field{
		Key:   "Error",
		Value: err,
	}
}

func Int64(key string, err int64) Field {
	return Field{
		Key:   key,
		Value: err,
	}
}

func String(key string, err string) Field {
	return Field{
		Key:   key,
		Value: err,
	}
}
