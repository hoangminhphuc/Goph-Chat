package business

type RegisterRepo interface {
	FindUser
	CreateUser
}

type Hasher interface {
	Hash
	Compare
}

type registerBusiness struct {
	repo RegisterRepo
	hasher Hasher
}

func NewRegisterBusiness(repo RegisterRepo, hasher Hasher) registerBusiness {
	return registerBusiness{repo: repo, hasher: hasher}	
}

func (rb *registerBusiness) Register() error {
	
}

