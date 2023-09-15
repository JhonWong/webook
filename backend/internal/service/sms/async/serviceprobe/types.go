package serviceprobe

import "context"

type ServiceProbe interface {
	Add(ctx context.Context, err error) bool
	IsCrashed(ctx context.Context) bool
}
