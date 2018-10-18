package graphql

import (
	"fmt"
	"reflect"

	"github.com/globalsign/mgo/bson"
)

func getType(f reflect.StructField) string {
	switch f.Type {
	case reflect.TypeOf((*string)(nil)):
		return "String"
	case reflect.TypeOf((*string)(nil)).Elem():
		return "String"
	case reflect.TypeOf((*int32)(nil)):
		return "Int"
	case reflect.TypeOf((*int32)(nil)).Elem():
		return "Int"
	case reflect.TypeOf((*bson.ObjectId)(nil)):
		return "ID"
	case reflect.TypeOf((*bson.ObjectId)(nil)).Elem():
		return "ID"
	case reflect.TypeOf((*bool)(nil)):
		return "Boolean"
	case reflect.TypeOf((*bool)(nil)).Elem():
		return "Boolean"
	default:
		return f.Type.Name()
	}
}

func getTypeFieldName(f reflect.StructField) string {
	if n, ok := f.Tag.Lookup("gql"); ok {
		return n
	}
	return f.Name
}

func buildTypeField(f reflect.StructField) string {
	return fmt.Sprintf("\t%s: %s\n", getTypeFieldName(f), getType(f))
}

func buildTypeBody(t reflect.Type) string {
	result := ""

	for i := 0; i < t.NumField(); i++ {
		result += buildTypeField(t.Field(i))
	}

	return result
}

func Build(t reflect.Type) string {
	result := fmt.Sprintf("type %s {\n%s\n}", t.Name(), buildTypeBody(t))

	return result
}
