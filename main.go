package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

// 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
func main() {
	pdf := NewPdf()
	// read pdf; compasref dbName == line; collectTo; all pdf files; 2 variant; sortRemoveDuplicate = final json
	pdf.Read()
	// excel := NewExcel()
	// excel.Read()
}

type Reader interface {
	Read()
}

type Excel struct {
	data [][]string
}

func NewExcel() *Excel {
	return &Excel{}
}

func (e *Excel) Read() {
	log.Println("excel read")
}

type Pdf struct {
	pages []*Page
}

func NewPdf() *Pdf {
	return &Pdf{}
}

type Marker struct {
	id         int64
	date       time.Time
	value      float64
	abbr       string
	name       string
	refference string
	ids        []int64
	isGroup    bool
	markers    []Marker
	// groupMarkers []Marker
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
	listMarker []Marker
	Group      []ResultGroup `json:"group_markers"`
}

func NewGroup(markers []Marker) *Group {
	return &Group{listMarker: markers}
}

//isGroup ? -> fill; by ids -> for loop range - ids ->fill id; value;
//save - all; result - like cache ->
/* then final result -> {

dbSaveView : single : [], group[],
graphicView : single = group; remove - duplicate(by date; and markerid)
*/

// type Response struct {
// 	markersGraphic struct{}
// 	markersDb      struct{}
// }

// type Group struct {
// 	id      int64 //  43
// 	date    time.Time
// 	markers []Marker // 17, 44
// }

type markersDb struct {
	single []Marker
	group  []Group
}

type Response struct {
	User       *User          `json:"user"`
	ListGroup  []ResultGroup  `json:"group_markers"`
	ListSingle []ResultSingle `json:"markers"`
}

type Page struct {
	line       []string
	groupIndex int
	isGroup    bool
	name       string
	// totalPage  int
}

func NewPage() *Page {
	return &Page{}
}

func mockDbData() []Marker {
	var markersDb []Marker

	markersDb = append(markersDb, Marker{}, Marker{}, Marker{}, Marker{})

	markersDb[0].abbr = "oak"
	markersDb[0].id = 43
	markersDb[0].ids = []int64{17, 44}
	markersDb[0].name = "Общий анализ крови"
	m := Marker{id: 17, abbr: "HGB", name: "Гемоглобин"}
	m2 := Marker{id: 44, abbr: "", name: "Цветной показатель"}
	markersDb[0].markers = append(markersDb[0].markers, m, m2)
	// markersDb[0].isGroup = true\

	markersDb[1].abbr = "ttg"
	markersDb[1].id = 69
	markersDb[1].name = "ТТГ"

	markersDb[2].abbr = "HGB"
	markersDb[2].id = 23
	markersDb[2].name = "Гемоглобин"

	markersDb[3].abbr = ""
	markersDb[3].id = 18
	markersDb[3].name = "Железо"

	return markersDb
}

// coolect each; pdf file; result -> to response {}, return final - result
func (r Response) Collect() {}

func (p *Page) division(listRow []pdf.Rows) (pages []*Page) {
	//todo ; receive []string{oak, oam, etc}
	for _, rows := range listRow {
		for _, row := range rows {
			page := &Page{}
			for _, word := range row.Content {
				page.line = append(page.line, word.S)
				// общий анализ мочи/крови, etc
				if strings.Contains(word.S, "Общий анализ крови") {
					page.isGroup = true
					// page.name = word.S
					page.name = "Общий анализ крови" // || мочи, etc
				}
			}
			pages = append(pages, page)
		}
	}
	return pages
}

/*data json -> client -> ; 2 variant;
-register -> save - share|| save ->  register ->; 1.phone; -> 2.anketa; 3.send 1 variant;; save Db; -> call /share
*/
/*
1 read pdf files
2 register - by phone
3 request  -  save each single/group in Db
4 share - getRelId; markerId
5 new signin -> return like trends all markers


		if further
		check author file - by IIN || #order - in Db
		return response - /upload pdf || akgun - 1 variant
		//check username each page ? - 1 user - all pdf
		//nomer zayvki - passed_marker; if nil -> insert; else error

*/

func (m Marker) Division() (single, group []Marker) {

	/*
	  1 remove duplicate;
	  1.2 division by group/single
	  1.3 sort by date;
	  1.4 group - like share service;

	  prettier to single marker;
	  prettier to group marker
	  merge - single/group
	*/

	//division =  single  / group
	for _, marker := range m.markers {
		if marker.isGroup {
			//func1()
			group = append(group, marker)
		} else {
			single = append(single, marker)
			//func2()
			/*
				same marker; appedn 1 struct; by tag m.id
			*/
		}
	}
	return
}

type Prettier interface {
	Pretty()
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

		if mapSng[marker.id] != nil {
			curr := mapSng[marker.id].(ResultSingle)
			v := Values{}
			v.Date = marker.date
			v.Value = marker.value
			curr.Result = append(curr.Result, v)
			mapSng[marker.id] = curr
		} else {
			//dry, func()
			rs := ResultSingle{}
			rs.MarkerId = marker.id
			rs.Name = marker.name
			rs.Reference = marker.refference
			v := Values{}
			v.Date = marker.date
			v.Value = marker.value
			rs.Result = append(rs.Result, v)
			mapSng[marker.id] = rs
		}
	}
	s.mapSingle = mapSng
	s.Unboxing()
}

type ResultGroup struct {
	Id      int64  `json:"group_id"`
	Name    string `josn:"name"`
	IsGroup bool   `json:"is_group"`
	markers []ResultSingle
}

func (g *Group) Pretty() {
	mapGroup := make(map[int64]interface{})
	single := Single{}
	// log.Println(len(g.listMarker))
	todo - append child markers - each group markers; define - why !read = all groupPage?

	for _, mg := range g.listMarker {
		if mapGroup[mg.id] != nil {
			curr := mapGroup[mg.id].(ResultGroup)
			//?
			single.listMarker = mg.markers
			log.Println(mg.markers, "group child markers")
			single.Pretty()
			single.Unboxing()
			curr.markers = single.Markers

		} else {
			rg := ResultGroup{}
			rg.Id = mg.id
			rg.IsGroup = true
			rg.Name = mg.name
			single.listMarker = mg.markers
			single.Pretty()
			single.Unboxing()

			rg.markers = single.Markers
			mapGroup[mg.id] = rg
		}

	}
	log.Println(mapGroup, "group final")
	// todo Unboxing()
	// g.Group = mapGroup
}

func (p *Pdf) Read() {
	// gorutine - read list pdfs
	markers := []Marker{}
	// read pdf files
	u := NewUser()
	// get mock db markers
	markersDb := mockDbData()

	for i := 1; i <= len(os.Args[1:]); i++ {
		rows, err := p.readPdf(os.Args[i]) // Read local pdf files
		if err != nil {
			panic(err)
		}
		page := NewPage()
		// division by single/group marker page
		p.pages = page.division(rows)
		// compare each page; pdf.page.row.name & db.marker name
		for _, pg := range p.pages {
			markers = append(markers, pg.Compare(markersDb))
		}
	}

	u.getUserCreds(p.pages[0])
	log.Println(*u, "user data")
	res := Response{
		User: u,
	}

	m := Marker{}
	m.markers = markers
	markers, groupMarkers := m.Division()

	s := NewSingle(markers)
	g := NewGroup(groupMarkers)

	s.Pretty()
	g.Pretty()
	res.ListSingle = s.Markers
	res.ListGroup = g.Group
	//todo  sort by date /duplicate remove

	log.Println(res.ListGroup, "final result")
}

// set each marker - markerId;
func (p *Page) Compare(dbMrks []Marker) (result Marker) {
	// func()? gorutine ?
	if p.isGroup {
		// get group Ids; from dbData
		id, mrks, name := p.getGroupId(dbMrks)
		result.id = id
		result.isGroup = true
		result.name = name
		// call func, concurrent ?
		result.markers = p.helper(mrks, true)
		result.date = result.markers[0].date
	} else {
		result = p.helper(dbMrks, false)[0]
	}
	return
}

// sort by markerId; []marker; []31: {val, date}
func (p *Page) getGroupId(dbMrks []Marker) (int64, []Marker, string) {
	for _, marker := range dbMrks {
		for _, text := range p.line {
			if strings.Contains(marker.name, text) {
				return marker.id, marker.markers, marker.name
			}
		}
	}
	return 0, nil, ""
}

type User struct {
	fullName  string
	birthDate string
	gender    string
}

func NewUser() *User {
	return &User{}
}

func (u *User) getUserCreds(page *Page) {

	for idx, text := range page.line {
		if text == "Дата рождения:" {
			u.fullName = page.line[idx+37]
			u.birthDate = page.line[idx+41] // convert method todo;
			u.gender = page.line[idx+47]
			break
		}
	}
}

//todo refactor
func (p *Page) helper(dbMrks []Marker, isGroup bool) (seqMarkers []Marker) {
	for _, marker := range dbMrks {
		for idx, text := range p.line {
			if marker.name == text {

				m := Marker{}
				m.id = marker.id

				s := strings.Split(p.line[idx+2], " ")
				ch := strings.Replace(s[0], ",", ".", 1)
				f, _ := strconv.ParseFloat(ch, 64)
				m.value = f
				chT := ""

				if isGroup {
					// m.isGroup = true
					chT = strings.Replace(p.line[idx+8], ".", "-", 2)
				} else {
					// m.isGroup=false
					chT = strings.Replace(p.line[idx+6], ".", "-", 2)
				}
				t, _ := time.Parse("01-02-2006 15:04", chT)
				m.date = t
				m.name = text
				m.refference = p.line[idx+4]
				seqMarkers = append(seqMarkers, m)
			}
		}
	}
	return seqMarkers
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
