package posthook

type PostInterface interface {
	PostExec() (err error)
}

type PostHook struct {
	Post PostInterface
}

func (p *PostHook) PostExec() (err error) {
	return p.Post.PostExec()
}
