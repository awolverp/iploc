package geoip

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

// compile-time type checking
var (
	_ fmt.Stringer = (*GeoIPObject)(nil)
	_ io.Closer    = (*Stream)(nil)
)

func IPToUint64(__ip net.IP) uint64 {
	v := __ip.To4()
	return uint64(v[0])<<24 + uint64(v[1])<<16 + uint64(v[2])<<8 + uint64(v[3])
}

func Uint64ToIP(__num uint64) net.IP {
	return net.IPv4(byte((__num>>24)&0xFF), byte((__num>>16)&0xFF), byte((__num>>8)&0xFF), byte(__num&0xFF))
}

type Stream struct {
	fd    *os.File
	csv_R *csv.Reader
}

func OpenStream(filename string) (*Stream, error) {
	s := Stream{}

	fd, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	s.fd = fd
	s.csv_R = csv.NewReader(fd)
	return &s, nil
}

func (s *Stream) Close() error { return s.fd.Close() }

func (s *Stream) Back() {
	s.fd.Seek(0, io.SeekStart)
	s.csv_R = csv.NewReader(s.fd)
}

type Process interface {
	// Initialize process to start
	Init() error
	// HandleProcess calls Process.Handle if Process.Where return true, otherwise continue
	Where(count int, obj *GeoIPObject) bool
	// It used after Where
	//
	// If returns false, HandleProcess break and calls Process.Done()
	Handle(count int, obj *GeoIPObject) bool
	// It used when an error retuned from other operations
	//
	// If returns false, HandleProcess break and calls Process.Done()
	OnError(count int, err error) bool
	// It used after all process
	//
	// It can returns error
	Done() error
}

func (s *Stream) HandleProcess(p Process, limit int) error {
	if err := p.Init(); err != nil {
		return err
	}

	for i := 0; (limit < 1) || (i < limit); i++ {
		record, err := s.csv_R.Read()

		if err != nil {
			if err == io.EOF || !p.OnError(i, err) {
				break
			}
			continue
		}

		obj, err := parseRecordToObject(record)

		if err != nil {
			if !p.OnError(i, err) {
				break
			}
			continue
		}

		if p.Where(i, obj) {
			if !p.Handle(i, obj) {
				break
			}
		}
	}

	s.Back()
	return p.Done()
}

type GeoIPObject struct {
	From     uint64    `json:"from"`
	To       uint64    `json:"to"`
	Registry string    `json:"registery"`
	Num      uint64    `json:"num"`
	Country  [3]string `json:"country"`

	__rec []string
}

func parseRecordToObject(record []string) (*GeoIPObject, error) {
	if len(record) != 7 {
		return nil, errors.New("invalid CSV record")
	}

	from, err := strconv.ParseUint(record[0], 10, 0)
	if err != nil {
		return nil, err
	}

	to, err := strconv.ParseUint(record[1], 10, 0)
	if err != nil {
		return nil, err
	}

	num, err := strconv.ParseUint(record[3], 10, 0)
	if err != nil {
		return nil, err
	}

	gp := GeoIPObject{
		From: from, To: to, Registry: record[2], Num: num,
		Country: [3]string{record[4], record[5], record[6]},
	}
	return &gp, nil
}

/*
is __n between obj.From and obj.To?
*/
func (obj GeoIPObject) IsBetween(__n uint64) bool {
	return (obj.From <= __n) && (__n <= obj.To)
}

/*
is __n between obj.From and obj.To?
*/
func (obj GeoIPObject) IsIPBetween(__n net.IP) bool {
	return (obj.From <= IPToUint64(__n)) && (IPToUint64(__n) <= obj.To)
}

func (obj GeoIPObject) JSON() string {
	b, _ := json.Marshal(map[string]interface{}{
		"from":     Uint64ToIP(obj.From).String(),
		"to":       Uint64ToIP(obj.To).String(),
		"registry": obj.Registry,
		"num":      obj.Num,
		"country":  obj.Country,
	})
	if b == nil {
		return ""
	}

	return string(b)
}

func (obj GeoIPObject) CSV() string {
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)
	if len(obj.__rec) < 1 {
		obj.__rec = make([]string, 7)
		obj.__rec[0] = Uint64ToIP(obj.From).String()
		obj.__rec[1] = Uint64ToIP(obj.To).String()
		obj.__rec[2] = obj.Registry
		obj.__rec[3] = strconv.FormatUint(obj.Num, 10)
		obj.__rec[4] = obj.Country[0]
		obj.__rec[5] = obj.Country[1]
		obj.__rec[6] = obj.Country[2]
	}
	w.Write(obj.__rec)
	w.Flush()
	return b.String()[:b.Len()-1]
}

func (obj GeoIPObject) String() string {
	return fmt.Sprintf(
		".GREEN%v.RESET-.RED%v\t\t.YELLOW%s.RESET\t%d\t.BOLD%s.RESET [%s / %s]",
		Uint64ToIP(obj.From), Uint64ToIP(obj.To), obj.Registry, obj.Num, obj.Country[2], obj.Country[0], obj.Country[1],
	)
}
