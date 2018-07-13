package relay

func GenerateConnectionTypes(baseType string) string {
	return `
	type ` + baseType + `Connection {
		edges: [` + baseType + `Edge]
		totalCount: Int
		pageInfo: PageInfo!
	}

	type ` + baseType + `Edge {
		node: ` + baseType + `
		cursor: ID!
	}
	`
}
