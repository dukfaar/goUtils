package eventbus

type ResolveType struct {
	//Which query to use to resolve this
	By string

	//What to use from the types fields for the query arguments
	FieldArguments map[string]string
}

type FieldType struct {
	//What is the fields name that gets added
	Name string

	//Which Type this will resolve to add this
	Type string

	//How the field will be resolved
	Resolve ResolveType
}

/*
Example (in json format) in itemService
{
	"type": "Item",
	"fields": [
		{
			"type": "Namespace",
			"name": "namespace",
			"resolve": {
				"by": "namespace",
				"fieldArguments": {
					"id": "namespaceId"
				}
			}
		}
	]
}

Example (in json format) in namespaceService
{
	"type": "Namespace",
	"fields": [
		{
			"type": "ItemConnection",
			"name": "items",
			"resolve": {
				"by": "itemsByNamespaceId",
				"fieldArguments": {
					"namespaceId": "_id"
				}
			}
		}
	]
}
*/
type SchemaExtension struct {
	//Which type to extend
	Type string

	//Fields to add to the type
	Fields []FieldType
}

type ServiceInfo struct {
	Name                  string
	Hostname              string
	Port                  string
	GraphQLHttpEndpoint   string
	GraphQLSchemaEndpoint string
	GraphQLSocketEndpoint string

	SchemaExtensions []SchemaExtension
}
