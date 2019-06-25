package domain

type RegistrationData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreateGroupData struct {
	Name    string `json:"name"`
	SchoolID string `json:"schoolcode"`
}

type CreateChildData struct {
	Name    string `json:"name"`
	GroupID string `json:"groupid"`
}
