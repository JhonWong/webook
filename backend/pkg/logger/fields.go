package logger

func Error(err error) Field {
	return Field{
		Key:   "Error",
		Value: err,
	}
}
