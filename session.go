package session

import "context"

type SessionResolverBuilder interface {
	SetRepo(repo SessionRepo) SessionResolverBuilder
	Build() (SessionResolver,error)
}

type SessionResolver interface {
	GetActiveOrderByContext(context context.Context) (string,error)
}

type SessionRepo interface {
	GetUserSessionByTokenToStruct(context context.Context,token string,dest interface{}) error
}

