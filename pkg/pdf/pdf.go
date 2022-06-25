package pdf

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

type Reader interface {
	Read([]Marker) Response
}
type Prettier interface {
	Pretty()
}

type Pdf struct {
	pages []*Page
}

func NewPdf() *Pdf {
	return &Pdf{}
}

func (p *Pdf) readPdf(path string) (listRow []pdf.Rows, err error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}
	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		// 1 file -> N page -> ;
		rows, err := p.GetTextByRow()
		if err != nil {
			return nil, err
		}
		listRow = append(listRow, rows)
	}
	return listRow, nil
}

type Response struct {
	User       User           `json:"user"`
	ListGroup  []Marker       `json:"group_markers"`
	ListSingle []ResultSingle `json:"markers"`
}

type ResultSingle struct {
	MarkerId  int64    `json:"marker_id"`
	Reference string   `json:"reference"`
	Name      string   `json:"name"`
	Result    []Values `json:"values"`
}
type Values struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

type Single struct {
	listMarker []Marker
	Markers    []ResultSingle `josn:"markers"`
	mapSingle  map[int64]interface{}
}

func NewSingle(markers []Marker) *Single {
	return &Single{listMarker: markers}
}

type Group struct {
	Id         int64  `json:"group_id"`
	Name       string `josn:"name"`
	IsGroup    bool   `json:"is_group"`
	listMarker []Marker
}

func NewGroup(markers []Marker) *Group {
	return &Group{
		listMarker: markers,
	}
}

// data - getFrom -mockDbData()
func mockDbGroupData() []Marker {
	var markersDb []Marker
	markersDb = append(markersDb, Marker{}, Marker{})

	markersDb[0].abbr = "oak"
	markersDb[0].Id = 43
	markersDb[0].Name = "Общий анализ крови"

	markersDb[1].abbr = "oam"
	markersDb[1].Id = 50
	markersDb[1].Name = "Общий анализ мочи"
	return markersDb
}

func mockDbData() []Marker {
	var markersDb []Marker
	markersDb = append(markersDb, Marker{}, Marker{}, Marker{}, Marker{}, Marker{})

	// temp oak
	markersDb[0].abbr = "oak"
	markersDb[0].Id = 43
	markersDb[0].ids = []int64{17, 44}
	markersDb[0].Name = "Общий анализ крови"
	m := Marker{Id: 17, abbr: "HGB", Name: "Гемоглобин"}
	m2 := Marker{Id: 44, abbr: "", Name: "Цветной показатель"}
	markersDb[0].Markers = append(markersDb[0].Markers, m, m2)
	// temp oam
	markersDb[1].abbr = "oam"
	markersDb[1].Id = 50
	markersDb[1].ids = []int64{12, 19, 33}
	markersDb[1].Name = "Общий анализ мочи"
	m3 := Marker{Id: 12, abbr: "SLR", Name: "Сангрия"}
	m4 := Marker{Id: 19, abbr: "GOL", Name: "Лайм брю"}
	m5 := Marker{Id: 33, abbr: "SMT", Name: "Сайм нит"}
	markersDb[1].Markers = append(markersDb[0].Markers, m3, m4, m5)

	markersDb[2].abbr = "ttg"
	markersDb[2].Id = 69
	markersDb[2].Name = "ТТГ"

	markersDb[3].abbr = "HGB"
	markersDb[3].Id = 23
	markersDb[3].Name = "Гемоглобин"

	markersDb[4].abbr = "ferr"
	markersDb[4].Id = 18
	markersDb[4].Name = "Железо"

	return markersDb
}

type Marker struct {
	Id         int64     `json:"id"`
	Date       time.Time `json:"date"`
	Value      float64   `json:"value,omitempty"`
	Name       string    `json:"name"`
	Refference string    `json:"reference,omitempty"`
	abbr       string
	ids        []int64
	IsGroup    bool     `json:"is_group,omitempty"`
	Markers    []Marker `json:"markers,omitempty"`
}

func NewMarker(markers []Marker) *Marker {
	return &Marker{
		Markers: markers,
	}
}

// division markers  to  single  / group structs
func (m Marker) Division() (single, group []Marker) {
	for _, marker := range m.Markers {
		if marker.IsGroup {
			group = append(group, marker)
		} else {
			single = append(single, marker)
		}
	}
	return
}

func (s *Single) Unboxing() {
	res := []ResultSingle{}
	for _, m := range s.mapSingle {
		res = append(res, m.(ResultSingle))
	}
	s.Markers = res
}

func (s *Single) Pretty() {
	mapSng := make(map[int64]interface{})

	for _, marker := range s.listMarker {
		if mapSng[marker.Id] != nil {
			curr := mapSng[marker.Id].(ResultSingle)
			v := Values{}
			v.Date = marker.Date
			v.Value = marker.Value
			curr.Result = append(curr.Result, v)
			mapSng[marker.Id] = curr
		} else {
			// dry, func()
			rs := ResultSingle{}
			rs.MarkerId = marker.Id
			rs.Name = marker.Name
			rs.Reference = marker.Refference
			v := Values{}
			v.Date = marker.Date
			v.Value = marker.Value
			rs.Result = append(rs.Result, v)
			mapSng[marker.Id] = rs
		}
	}
	s.mapSingle = mapSng
	s.Unboxing()
}

// func (g *Group) Pretty() {
// 	mapGroup := make(map[int64]interface{})
// }

// then use dbData
func (p *Pdf) Read(dbData []Marker) Response {
	// gorutine - read list pdfs
	listSingleGroupMarkers := []Marker{}
	// read pdf files
	u := NewUser()
	// get mock db markers
	markersDb := mockDbData()
	groupMarkersDb := mockDbGroupData()

	for i := 1; i <= len(os.Args[1:]); i++ {
		rows, err := p.readPdf(os.Args[i]) // Read local pdf files
		if err != nil {
			panic(err)
		}
		page := NewPage()
		// division by single/group marker page
		p.pages = page.division(rows, groupMarkersDb)
		// compare each page; pdf.page.row.name & db.marker name
		for _, pg := range p.pages {
			listSingleGroupMarkers = append(listSingleGroupMarkers, pg.Compare(markersDb))
		}
	}

	u.getUserCreds(p.pages[0])

	res := Response{
		User: *u,
	}

	m := NewMarker(listSingleGroupMarkers)

	markers, groupMarkers := m.Division()

	s := NewSingle(markers)
	s.Pretty()
	// g := NewGroup(groupMarkers)
	// g.Pretty()

	res.ListSingle = s.Markers
	res.ListGroup = groupMarkers
	// todo  sort by date /duplicate remove
	return res
}

type Page struct {
	line    []string
	isGroup bool
	name    string
}

func NewPage() *Page {
	return &Page{}
}

// set each marker - markerId;
func (p *Page) Compare(dbMrks []Marker) (result Marker) {
	// func()? gorutine ?
	if p.isGroup {
		// get group Ids; from dbData
		id, mrks, name := p.getGroupId(dbMrks)
		result.Id = id
		result.IsGroup = true
		result.Name = name
		// call func, concurrent ?
		result.Markers = p.fillFields(mrks, true)
		result.Date = result.Markers[0].Date
	} else {
		result = p.fillFields(dbMrks, false)[0]
	}
	return
}

// sort by markerId; []marker; []31: {val, date}
func (p *Page) getGroupId(dbMrks []Marker) (int64, []Marker, string) {
	for _, marker := range dbMrks {
		for _, text := range p.line {
			if strings.Contains(marker.Name, text) {
				return marker.Id, marker.Markers, marker.Name
			}
		}
	}
	return 0, nil, ""
}

func (p *Page) fillFields(dbMrks []Marker, isGroup bool) (seqMarkers []Marker) {
	var (
		m            Marker
		splitted     []string
		changedValue string
		value        float64
		changedTime  string
		date         time.Time
	)

	for _, marker := range dbMrks {
		for idx, text := range p.line {
			if marker.Name == text {
				m.Id = marker.Id
				splitted = strings.Split(p.line[idx+2], " ")
				changedValue = strings.Replace(splitted[0], ",", ".", 1)
				value, _ = strconv.ParseFloat(changedValue, 64)
				m.Value = value

				if isGroup {
					changedTime = strings.Replace(p.line[idx+8], ".", "-", 2)
				} else {
					changedTime = strings.Replace(p.line[idx+6], ".", "-", 2)
				}
				date, _ = time.Parse("01-02-2006 15:04", changedTime)
				m.Date = date
				m.Name = text
				m.Refference = p.line[idx+4]
				seqMarkers = append(seqMarkers, m)
			}
		}
	}
	return seqMarkers
}

func (p *Page) division(listRow []pdf.Rows, dbGroupMarkers []Marker) (pages []*Page) {
	for _, rows := range listRow {
		for _, row := range rows {
			page := &Page{}
			for _, word := range row.Content {
				page.line = append(page.line, word.S)
				// общий анализ мочи/крови, etc
				for _, gm := range dbGroupMarkers {
					if strings.Contains(word.S, gm.Name) {
						page.isGroup = true
						// page.name = word.S
						page.name = gm.Name
						break
					}
				}
			}
			pages = append(pages, page)
		}
	}
	return pages
}

type User struct {
	FullName  string `json:"full_name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) getUserCreds(page *Page) {
	for idx, text := range page.line {
		if text == "Дата рождения:" {
			u.FullName = page.line[idx+37]
			u.BirthDate = page.line[idx+41] // convert method todo;
			u.Gender = page.line[idx+47]
			break
		}
	}
}
