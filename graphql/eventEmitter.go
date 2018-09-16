package graphql

import (
	"github.com/dukfaar/goUtils/eventbus"
	graphqlIntrospection "github.com/graph-gophers/graphql-go/introspection"
)

func EmitRegisterEvents(eventName string, gqlType *graphqlIntrospection.Type, ebus eventbus.EventBus) {
	if gqlType == nil {
		return
	}

	typeFields := gqlType.Fields(nil)

	if typeFields == nil {
		return
	}

	for i := range *typeFields {
		field := (*typeFields)[i]

		ebus.Emit(eventName, field.Name())
	}
}

func isBuiltInType(name string) bool {
	switch name {
	case "String":
		return true
	case "ID":
		return true
	case "Int":
		return true
	case "Float":
		return true
	case "Date":
		return true
	case "Boolean":
		return true
	case "__Directive":
		return true
	case "__DirectiveLocation":
		return true
	case "__EnumValue":
		return true
	case "__TypeKind":
		return true
	case "__Type":
		return true
	case "__Schema":
		return true
	case "__InputValue":
		return true
	case "__Field":
		return true
	default:
		return false
	}
}

type TypeEvent struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields"`
}

func EmitRegisterTypeEvents(eventName string, gqlTypes []*graphqlIntrospection.Type, ebus eventbus.EventBus) {
	if gqlTypes != nil {
		for i := range gqlTypes {
			gqlType := gqlTypes[i]
			name := *gqlType.Name()

			if isBuiltInType(name) {
				continue
			}

			if name == "Query" || name == "Mutation" || name == "Subscription" {
				continue
			}

			fields := gqlType.Fields(nil)
			if fields != nil {
				fieldsArray := make([]string, 0)

				for fieldsIndex := range *fields {
					field := (*fields)[fieldsIndex]
					fieldsArray = append(fieldsArray, field.Name())
				}

				ebus.Emit(eventName, TypeEvent{
					Name:   name,
					Fields: fieldsArray,
				})
			}
		}
	}
}
