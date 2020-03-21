package user

import "time"

type User struct {
	tableName struct{} `pg:"users"json:"-"`
	Id        string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time `pg:"default:now()"`
}

type Society struct {
	tableName struct{} `pg:"societies"`
	Id        string
	Name      string
	CreatedAt time.Time `pg:"default:now()"`
}

type Member struct {
	tableName  struct{} `pg:"societies_members"`
	UserId     string
	SocietyId  string
	Permission membership
}

type Applicant struct {
	tableName struct{} `pg:"societies_applicants"`
	UserId    string
	SocietyId string
	CreatedAt time.Time `pg:"default:now()"`
}

type membership string

type UserGroupRequest struct {
	UserId    string
	SocietyId string
}
