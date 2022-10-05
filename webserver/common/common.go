package common

type UserType int

const (
	ADMIN UserType = 1
	USER  UserType = 1 << 1
	GUEST UserType = 1 << 2
)
