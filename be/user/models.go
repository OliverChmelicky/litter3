package user

type UserModel struct {
	tableName struct{} `pg:"users"`
	Id        string
	FirstName string
	LastName  string
	Email     string
	Created   int64
}

type Society struct {
	tableName struct{} `pg:"societies"`
	Id        string
	Name      string
	Created   int64
}
