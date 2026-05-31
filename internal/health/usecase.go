package health

type Usecase interface {
	Check() (Result, error)
}

type Result struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type usecase struct {
	service string
	// dbPing optional in Phase 1; for now, keep it simple and fast.
}

func New(service string) Usecase {
	return &usecase{service: service}
}

func (u *usecase) Check() (Result, error) {
	return Result{Status: "ok", Service: u.service}, nil
}

