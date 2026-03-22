package origin

type UnknownOrigin struct {
	origin
}

func (uo *UnknownOrigin) Definition() []string {
	return []string{}
}

func NewUnkownOrigin() *UnknownOrigin {
	return &UnknownOrigin{
		origin: origin{},
	}
}
