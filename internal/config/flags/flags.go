package flags

import "flag"

type Flags struct {
	A string
	B string
	F string
	D string
	S string
}

func ProcessingFlags() Flags {
	f := Flags{}
	flag.StringVar(&f.A, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&f.B, "b", "http://localhost:8080", "base address")
	flag.StringVar(&f.F, "f", "storage/storage.json", "file storage path")
	flag.StringVar(&f.D, "d", "", "database address")
	flag.StringVar(&f.S, "s", "secretkey", "JWT secret key")
	flag.Parse()
	return f
}
