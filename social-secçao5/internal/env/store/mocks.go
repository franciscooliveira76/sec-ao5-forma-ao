package store


func NewMockStore() Storage {
	return Storage {
		Users: &MockUserStore{},
	}
}

type MockUserStore struck {
	users []User
}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {

	return nil
}

func (m *MockUserStore) GetByID(context.Context, int64) (*User, error) {
	return &User{
		Id: 42,
	}, nil
}

func (m *MockUserStore) GetByEmail(context.Context, string ) (&User, error) {

	return &User{}, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, t string) error {

	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, id int64) error {
	return nil
}