package agent

type Session struct {
	previousResponseID string
}

func NewSession() *Session {
	return &Session{}
}

func (s *Session) PreviousResponseID() string {
	return s.previousResponseID
}

func (s *Session) SetPreviousResponseID(responseID string) {
	s.previousResponseID = responseID
}

func (s *Session) Reset() {
	s.previousResponseID = ""
}
