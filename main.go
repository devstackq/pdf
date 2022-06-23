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

type Group struct {
	id      int64 //  43
	date    time.Time
	markers []Marker // 17, 44
}

type markersDb struct {
	single []Marker
	group  []Group
}

type Response struct {
	listGroup  []Group
	listSingle []Marker
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

func (m Marker) Prettier(markers []Marker) {
	/*1 remove duplicate;
	  1.2 division by group/single
	  1.3 sort by date;
	*/
	res := Response{}

	for _, m := range markers {
		log.Println(m)

		if m.isGroup {
			/*
				write to []GroupMarkers
			*/
			res.listGroup = append(res.listGroup)
			// append to []groupStruct;
		} else {
			// append single markers
		}
	}
}

func (p *Pdf) Read() {
	// gorutine - read list pdfs

	markers := []Marker{}
	// read pdf files
	for i := 1; i <= len(os.Args[1:]); i++ {
		rows, err := p.readPdf(os.Args[i]) // Read local pdf files
		if err != nil {
			panic(err)
		}
		page := NewPage()
		// division by single/group marker
		p.pages = page.division(rows)
		// get mock db markers
		markersDb := mockDbData()

		// compare each page; pdf.page.row.name & db.marker name
		for _, pg := range p.pages {
			markers = append(markers, pg.Compare(markersDb))
		}
	}
	// log.Println(markers, len(markers))
	// write/sort/duplicate to final struct
	m := Marker{}
	m.Prettier(markers)
}

// set each marker - markerId;
func (p *Page) Compare(dbMrks []Marker) (result Marker) {
	// func()? gorutine ?
	if p.isGroup {
		// get Id; & ids
		id, mrks := p.getGroupId(dbMrks)
		result.id = id
		// call func, concurrent ?
		result.markers = p.helper(mrks, true)
		result.date = result.markers[0].date
		result.isGroup = true
		// regexp.Match("/./")
	} else {
		result = p.helper(dbMrks, false)[0]
	}
	// log.Println(result, "res")
	return
}

// sort by markerId; []marker; []31: {val, date}

func (p *Page) getGroupId(dbMrks []Marker) (int64, []Marker) {
	for _, marker := range dbMrks {
		for _, text := range p.line {
			if strings.Contains(marker.name, text) {
				return marker.id, marker.markers
			}
		}
	}
	return 0, nil
}

// func (p *Page) getDate(dbMrks []Marker) (date time.Time) {
// }
// todo refatror
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
					chT = strings.Replace(p.line[idx+8], ".", "-", 2)
				} else {
					chT = strings.Replace(p.line[idx+6], ".", "-", 2)
				}
				t, _ := time.Parse("01-02-2006 15:04", chT)
				m.date = t

				m.refference = p.line[idx+4]
				// log.Println(t, err)
				seqMarkers = append(seqMarkers, m)

				// log.Println("find marker ", p.line[idx+6], p.line[idx+8])
				// collect to single(); 2 variant
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
