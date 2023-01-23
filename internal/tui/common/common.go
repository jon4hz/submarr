package common

const Ellipsis = "…"

type ErrMsg struct {
	Err error
}

func NewErrMsg(err error) ErrMsg {
	return ErrMsg{Err: err}
}

func (e ErrMsg) Error() string {
	return e.Err.Error()
}
