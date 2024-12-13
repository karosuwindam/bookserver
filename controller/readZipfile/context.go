package readzipfile

import "context"

type zipId struct{}

func ContextWriteZipId(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, zipId{}, id)
}

func contextReadZipId(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(zipId{}).(int)
	return id, ok
}
