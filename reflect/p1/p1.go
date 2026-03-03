package p1

type P1 struct {
	PublicField  string
	privateField float64
}

func New() *P1 {
	return &P1{
		PublicField:  "I am public",
		privateField: 3.14,
	}
}
