package domain

type SMSInfo struct {
	Tpl        string
	Args       []string
	Numbers    []string
	RetryTimes int
}

type SMSAsyncInfo struct {
	Id            int64
	Tpl           string
	Args          []string
	Numbers       []string
	MaxRetryCount int
}
