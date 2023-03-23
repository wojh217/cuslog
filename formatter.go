package cuslog

type Formatter interface {
	Format(* Entry) error
}