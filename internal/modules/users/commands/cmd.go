package commands

type Cmd[I any, O any] interface {
	Exec(*I) (*O, error)
}

type CmdNoOutput[I any] interface {
	Exec(*I) error
}
