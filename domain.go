package dnsbench

type Domain struct {
	domain string
}

func (d Domain) String() string {
	return d.domain
}

func NewDomain(d string) Domain {
	return Domain{
		domain: d,
	}
}
