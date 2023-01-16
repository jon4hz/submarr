package common

type ErrMsg struct {
	Err error
}

func NewErrMsg(err error) ErrMsg {
	return ErrMsg{Err: err}
}

func (e ErrMsg) Error() string {
	return e.Err.Error()
}
