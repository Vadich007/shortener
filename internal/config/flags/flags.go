package flags

import "flag"

type Flags struct {
	A string
	B string
}

func ProcessingFlags() Flags {
	f := Flags{}
	flag.StringVar(&f.A, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&f.B, "b", "", "base address")
	flag.Parse()
	return f
}
