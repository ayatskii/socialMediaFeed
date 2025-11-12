package user

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func NewUserFactory(id, name, role string) *User {
	normalizedRole := "Standard"
	if role == "Admin" || role == "Verified" {
		normalizedRole = role
	}

	return &User{
		ID:   id,
		Name: name,
		Role: normalizedRole,
	}
}
