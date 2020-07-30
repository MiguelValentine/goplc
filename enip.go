package goplc

type TcpState struct {
	establishing bool
	established  bool
}

type SessionState struct {
	establishing bool
	established  bool
}

type ConnectionState struct {
	establishing bool
	established  bool
	seqNum       uint64
}

type ErrorState struct {
}

type state struct {
	_tcp        TcpState
	_session    SessionState
	_connection ConnectionState
	_error      ErrorState
}

type ENIP struct {
	state state
}

func (e *ENIP) error() ErrorState {
	return e.state._error
}

func (e *ENIP) establishing() bool {
	return e.state._session.establishing
}

func (e *ENIP) established() bool {
	return e.state._session.established
}
