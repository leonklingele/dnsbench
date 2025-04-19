package dnsbench

type Domain struct {
	domain string
}

func NewDomain(d string) Domain {
	return Domain{
		domain: d,
	}
}

func (d Domain) String() string {
	return d.domain
}
