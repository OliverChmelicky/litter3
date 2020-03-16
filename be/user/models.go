package user

type UserModel struct {
	tableName struct{} `pg:"users"`
	Id        string
	FirstName string
	LastName  string
	Email     string
	Created   int64
}

type SocietyModel struct {
	tableName struct{} `pg:"societies"`
	Id        string
	Name      string
	Created   int64
}

type MemberModel struct {
	tableName  struct{} `pg:"societies_members"`
	UserId     string
	SocietyId  string
	Permission membership
}

type ApplicantModel struct {
	tableName struct{} `pg:"societies_applicants"`
	UserId    string
	SocietyId string
	Created   int64
}

type membership string

type UserGroupRequest struct {
	UserId    string
	SocietyId string
}
