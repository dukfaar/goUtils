package service

import (
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type DBService interface {
	HasElementBeforeID(id string) (bool, error)
	HasElementAfterID(id string) (bool, error)

	Count() (int, error)
}

type DBServiceWithQuery interface {
	MakeBaseQuery() bson.M
	MakeListQuery(query bson.M, before *string, after *string)

	HasElementBeforeIDWithQuery(query bson.M, id string) (bool, error)
	HasElementAfterIDWithQuery(query bson.M, id string) (bool, error)

	CountWithQuery(query bson.M) (int, error)
}

type BaseMgoServiceWithQuery struct {
	Collection *mgo.Collection
}

func (s *BaseMgoServiceWithQuery) MakeBaseQuery() bson.M {
	return bson.M{}
}

func (s *BaseMgoServiceWithQuery) MakeListQuery(query bson.M, before *string, after *string) {
	if after != nil {
		query["_id"] = bson.M{
			"$gt": bson.ObjectIdHex(*after),
		}
	}

	if before != nil {
		query["_id"] = bson.M{
			"$lt": bson.ObjectIdHex(*before),
		}
	}
}

func (s *BaseMgoServiceWithQuery) HasElementBeforeIDWithQuery(inquery bson.M, id string) (bool, error) {
	query := bson.M{}

	for k, v := range inquery {
		query[k] = v
	}

	query["_id"] = bson.M{
		"$lt": bson.ObjectIdHex(id),
	}

	count, err := s.Collection.Find(query).Count()
	return count > 0, err
}

func (s *BaseMgoServiceWithQuery) HasElementAfterIDWithQuery(inquery bson.M, id string) (bool, error) {
	query := bson.M{}

	for k, v := range inquery {
		query[k] = v
	}

	query["_id"] = bson.M{
		"$gt": bson.ObjectIdHex(id),
	}

	count, err := s.Collection.Find(query).Count()
	return count > 0, err
}

func (s *BaseMgoServiceWithQuery) CountWithQuery(query bson.M) (int, error) {
	count, err := s.Collection.Find(query).Count()
	return count, err
}

func (s *BaseMgoServiceWithQuery) HasElementBeforeID(id string) (bool, error) {
	query := bson.M{}

	query["_id"] = bson.M{
		"$lt": bson.ObjectIdHex(id),
	}

	count, err := s.Collection.Find(query).Count()
	return count > 0, err
}

func (s *BaseMgoServiceWithQuery) HasElementAfterID(id string) (bool, error) {
	query := bson.M{}

	query["_id"] = bson.M{
		"$gt": bson.ObjectIdHex(id),
	}

	count, err := s.Collection.Find(query).Count()
	return count > 0, err
}

func (s *BaseMgoServiceWithQuery) Count() (int, error) {
	count, err := s.Collection.Find(bson.M{}).Count()

	return count, err
}

func (s *BaseMgoServiceWithQuery) GetSkipLimit(query bson.M, first *int32, last *int32) (int, int) {
	var (
		skip  int
		limit int
	)

	if first != nil {
		limit = int(*first)
	}

	if last != nil {
		count, _ := s.Collection.Find(query).Count()

		limit = int(*last)
		skip = count - limit
	}

	return skip, limit
}
