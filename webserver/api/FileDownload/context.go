package filedownload

import "context"

type idFIletype struct {
	id       int
	filetype string
}

type fileNameFIletype struct {
	filename string
	filetype string
}

func contextWriteIdFIleType(ctx context.Context, id int, filetype string) context.Context {
	tmp := idFIletype{id, filetype}
	return context.WithValue(ctx, idFIletype{}, tmp)
}

func contextReadIdFileType(ctx context.Context) (idFIletype, bool) {
	tmp, ok := ctx.Value(idFIletype{}).(idFIletype)
	return tmp, ok
}

func contextWriteFileNameFIleType(ctx context.Context, filename, filetype string) context.Context {
	tmp := fileNameFIletype{filename, filetype}
	return context.WithValue(ctx, fileNameFIletype{}, tmp)
}

func contextReadFileNameFIleType(ctx context.Context) (fileNameFIletype, bool) {
	tmp, ok := ctx.Value(fileNameFIletype{}).(fileNameFIletype)
	return tmp, ok
}
