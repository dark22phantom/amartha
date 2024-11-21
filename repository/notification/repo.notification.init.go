package notification

type Repository struct{}

func New() (*Repository, error) {
	return &Repository{}, nil
}
