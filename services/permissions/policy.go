package permissions

var Policy = &PolicyConfiguration{
	Permissions: map[string]*Permission{
		"updateOwn": {
			Action: "update",
			Conditions: Conditions{
				func(request *Request) error {
					if request.Ressource.GetOwner() == string(request.User.ID) {
						return nil
					}
					return newAccessDeniedError()
				},
			},
		},
	},
	Roles: Roles{
		"admin": {
			Name:  "Admin",
			Color: "red",
			Permissions: map[string][]*Permission{
				"user": {
					{Action: "create"},
					{Action: "update"},
					{Action: "delete"},
				},
			},
			Parents: []string{"default", "developer"},
		},
		"developer": {
			Name:    "developer",
			Color:   "blue",
			Parents: []string{"default"},
		},
		"default": {
			Name:  "",
			Color: "gray",
		},
	},
}
