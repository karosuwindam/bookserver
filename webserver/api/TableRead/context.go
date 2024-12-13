package tableread

import "context"

type tableName struct{}

type tableIdName struct {
	tableName string
	id        int
}

func contextWriteTableName(ctx context.Context, table string) context.Context {
	return context.WithValue(ctx, tableName{}, table)
}

func contextReadTableName(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(tableName{}).(string)
	return v, ok
}

func contextWriteTableIdName(ctx context.Context, table string, id int) context.Context {
	return context.WithValue(ctx, tableIdName{}, tableIdName{table, id})
}

func contextReadTableIdName(ctx context.Context) (tableIdName, bool) {
	v, ok := ctx.Value(tableIdName{}).(tableIdName)
	return v, ok
}
