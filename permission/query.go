package permission

import "github.com/dukfaar/goUtils/graphql"

const Query string = `query {
	users {
		edges {
			node {
				_id
				roles {
					edges {
						node {
							_id
						}
					}
				}
			}
		}
	}
	roles {
		edges {
			node {
				_id
				permissions {
					edges {
						node {
							_id
						}
					}
				}
			}
		}
	}
	permissions {
		edges {
			node {
				_id
				name
			}
		}
	}
	tokens {
		edges {
			node {
				accessToken
				accessTokenExpiresAt
				userId
			}
		}
	}
}`

func ParseQueryResponse(queryResult graphql.Response, permissionService *Service) {
	tokenEdges := queryResult.GetObject("tokens").GetArray("edges")
	for j := 0; j < tokenEdges.Len(); j++ {
		tokenEdge := tokenEdges.Get(j)
		token := tokenEdge.GetObject("node")
		accessTokenExpiresAt, _ := token.GetInt64("accessTokenExpiresAt")
		expiresAt := graphql.JSTimestampToTime(accessTokenExpiresAt)

		permissionService.SetToken(token.GetString("accessToken"), token.GetString("userId"), expiresAt)
	}

	userEdges := queryResult.GetObject("users").GetArray("edges")
	for i := 0; i < userEdges.Len(); i++ {
		userEdge := userEdges.Get(i)
		user := userEdge.GetObject("node")
		id := user.GetString("_id")

		roleEdges := user.GetObject("roles").GetArray("edges")
		userRoles := make([]string, roleEdges.Len())
		for j := 0; j < roleEdges.Len(); j++ {
			roleEdge := roleEdges.Get(j)
			role := roleEdge.GetObject("node")
			userRoles[j] = role.GetString("_id")
		}

		permissionService.SetUser(id, userRoles)
	}

	roleEdges := queryResult.GetObject("roles").GetArray("edges")
	for i := 0; i < roleEdges.Len(); i++ {
		roleEdge := roleEdges.Get(i)
		role := roleEdge.GetObject("node")
		id := role.GetString("_id")

		permissionEdges := role.GetObject("permissions").GetArray("edges")
		rolePermissions := make([]string, permissionEdges.Len())
		for j := 0; j < permissionEdges.Len(); j++ {
			permissionEdge := permissionEdges.Get(j)
			permission := permissionEdge.GetObject("node")
			rolePermissions[j] = permission.GetString("_id")
		}

		permissionService.SetRole(id, rolePermissions)
	}

	permissionEdges := queryResult.GetObject("permissions").GetArray("edges")
	for j := 0; j < permissionEdges.Len(); j++ {
		permissionEdge := permissionEdges.Get(j)
		permission := permissionEdge.GetObject("node")
		permissionService.SetPermission(permission.GetString("_id"), permission.GetString("name"))
	}
}
