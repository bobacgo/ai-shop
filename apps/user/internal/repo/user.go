package repo

type UserRepo struct {
	data *Data
}

func NewUserRepo(data *Data) *UserRepo {
	return &UserRepo{
		data: data,
	}
}