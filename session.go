package session

import "context"

type SessionResolverBuilder interface {
	SetRepo(repo SessionRepo) SessionResolverBuilder
	Build() (SessionResolver,error)
}

type SessionResolver interface {
	GetOrderIdFromToken(context context.Context,token string) (string,error)
}

type SessionRepo interface {
	GetUserSessionByTokenToStruct(context context.Context,token string,dest interface{}) error
}