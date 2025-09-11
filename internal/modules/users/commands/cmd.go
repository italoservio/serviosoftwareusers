package commands

type Cmd[I any, O any] interface {
	Exec(*I) (*O, error)
}
