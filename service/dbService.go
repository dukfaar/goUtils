package service

import "github.com/globalsign/mgo/bson"

type DBService interface {
	HasElementBeforeID(id string) (bool, error)
	HasElementAfterID(id string) (bool, error)

	Count() (int, error)

	HasElementBeforeIDWithQuery(query bson.M, id string) (bool, error)
	HasElementAfterIDWithQuery(query bson.M, id string) (bool, error)

	CountWithQuery(query bson.M) (int, error)
}
