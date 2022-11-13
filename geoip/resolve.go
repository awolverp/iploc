package geoip

import (
	"net"
	"sort"
)

// compile-time type checking
var (
	_ Process = (*countryProcess)(nil)
	_ Process = (*ipIsBetweenProcess)(nil)
	_ Process = (*allProcess)(nil)
)

type countryProcess struct {
	Name string

	Offset int
	Limit  int

	Objects []*GeoIPObject

	length int
	found  int
	err    error
}

func (p *countryProcess) Init() error { p.length = len(p.Name); return nil }

func (p countryProcess) Where(count int, obj *GeoIPObject) bool {
	switch p.length {
	case 2:
		return obj.Country[0] == p.Name

	case 3:
		return obj.Country[1] == p.Name

	default:
		return obj.Country[2] == p.Name
	}
}

func (p *countryProcess) Handle(count int, obj *GeoIPObject) bool {
	if p.Offset > 0 {
		p.Offset--
		return true
	}

	p.Objects = append(p.Objects, obj)

	p.found++
	return (p.Limit < 1 || p.Limit > p.found)
}

func (p *countryProcess) OnError(count int, err error) bool {
	p.err = err
	return false
}

func (p countryProcess) Done() error {
	return p.err
}

func ResolveCountry(stream *Stream, name string, offset, limit int) ([]*GeoIPObject, error) {
	p := countryProcess{Name: name, Limit: limit, Offset: offset}
	err := stream.HandleProcess(&p, 0)
	return p.Objects, err
}

type ipIsBetweenProcess struct {
	IP net.IP
	q  uint64

	Objects []*GeoIPObject

	err error
}

func (p *ipIsBetweenProcess) Init() error {
	p.q = IPToUint64(p.IP)
	return nil
}

func (p ipIsBetweenProcess) Where(count int, obj *GeoIPObject) bool {
	return obj.IsBetween(p.q)
}

func (p *ipIsBetweenProcess) Handle(count int, obj *GeoIPObject) bool {
	p.Objects = append(p.Objects, obj)
	return false
}

func (p *ipIsBetweenProcess) OnError(count int, err error) bool {
	p.err = err
	return false
}

func (p ipIsBetweenProcess) Done() error {
	return p.err
}

func ResolveIP(stream *Stream, ip net.IP) ([]*GeoIPObject, error) {
	p := ipIsBetweenProcess{IP: ip}
	err := stream.HandleProcess(&p, 0)
	return p.Objects, err
}

type allProcess struct {
	Offset int
	Limit  int

	Objects []*GeoIPObject

	found int
	err   error
}

func (p *allProcess) Init() error { return nil }

func (p allProcess) Where(_ int, _ *GeoIPObject) bool {
	return true
}

func (p *allProcess) Handle(_ int, obj *GeoIPObject) bool {
	if p.Offset > 0 {
		p.Offset--
		return true
	}

	p.Objects = append(p.Objects, obj)

	p.found++
	return (p.Limit < 1 || p.Limit > p.found)
}

func (p *allProcess) OnError(_ int, err error) bool {
	p.err = err
	return false
}

func (p allProcess) Done() error {
	return p.err
}

func ResolveAll(stream *Stream, offset int, limit int) ([]*GeoIPObject, error) {
	p := allProcess{Offset: offset, Limit: limit}
	err := stream.HandleProcess(&p, 0)
	return p.Objects, err
}

type listCountriesProcess struct {
	Objects [][3]string
	err     error
}

func (p *listCountriesProcess) Init() error { return nil }

func (p listCountriesProcess) Where(_ int, obj *GeoIPObject) bool {
	for _, v := range p.Objects {
		if v[0] == obj.Country[0] {
			return false
		}
	}
	return true
}

func (p *listCountriesProcess) Handle(_ int, obj *GeoIPObject) bool {
	p.Objects = append(p.Objects, obj.Country)
	return true
}

func (p *listCountriesProcess) OnError(_ int, err error) bool {
	p.err = err
	return false
}

func (p listCountriesProcess) Done() error {
	if p.err != nil {
		return p.err
	}
	sort.Slice(p.Objects, func(i, j int) bool {
		return p.Objects[i][0] > p.Objects[j][0]
	})
	return nil
}

func ListCountries(stream *Stream) ([][3]string, error) {
	p := listCountriesProcess{}
	err := stream.HandleProcess(&p, 0)
	return p.Objects, err
}

func Resolve(stream *Stream, query string, offset, limit int) ([]*GeoIPObject, error) {
	if IP := net.ParseIP(query); IP != nil {
		return ResolveIP(stream, IP)
	} else {
		return ResolveCountry(stream, query, offset, limit)
	}
}
