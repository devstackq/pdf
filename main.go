package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

//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

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

//temp
// type DBMarkerIds struct {
// 	id      int64
// 	ids     []int64
// 	name    string
// 	abbr    string
// 	markers []Marker
// }

type Marker struct {
	id      int64
	date    time.Time
	value   float64
	abbr    string
	name    string
	ids     []int64
	markers []Marker
}

func mockDbData() []*Marker {

	var markersDb []*Marker

	markersDb = append(markersDb, &Marker{}, &Marker{}, &Marker{}, &Marker{})

	markersDb[0].abbr = "oak"
	markersDb[0].id = 43
	markersDb[0].ids = []int64{17, 44}
	markersDb[0].name = "Общий анализ крови"
	m := Marker{id: 17, abbr: "HGB", name: "Гемоглобин"}
	m2 := Marker{id: 44, abbr: "", name: "Цветной показатель"}

	markersDb[0].markers = append(markersDb[0].markers, m, m2)

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

//isGroup ? -> fill; by ids -> for loop range - ids ->fill id; value;
//save - all; result - like cache ->
/* then final result -> {

dbSaveView : single : [], group[],
graphicView : single = group; remove - duplicate(by date; and markerid)
*/

type Response struct {
	markersGraphic struct{}
	markersDb      struct{}
}

type Group struct {
	id      int64    //  43
	markers []Marker // 17, 44
}

type markersDb struct {
	single []Marker
	group  []Group
}

//coolect each; pdf file; result -> to response {}, return final - result
func (r Response) Collect() {}

func (p *Pdf) Read() {

	content, err := readPdf4(os.Args[1]) // Read local pdf file
	if err != nil {
		panic(err)
	}
	p.pages = content

	markersDb := mockDbData()

	//loop each page;
	for _, pg := range p.pages {
		pg.Compare(markersDb)
	}

	//prepareFinalResult; 1 variant; 2 variant;

}

func main() {
	pdf := NewPdf()
	//read pdf; compare dbName == line; collectTo; all pdf files; 2 variant; sortRemoveDuplicate = final json
	pdf.Read()
	// excel := NewExcel()
	// excel.Read()
}

type Page struct {
	line       []string
	groupIndex int
	isGroup    bool
	name       string
	// totalPage  int
}

//sort by markerId; []marker; []31: {val, date}

func (p *Page) getMrkId(dbMrks []*Marker) (int64, int, []Marker) {

	for _, marker := range dbMrks {
		for idx, text := range p.line {
			if marker.name == text {
				return marker.id, idx, marker.markers // return id, lineIdx
			}
			if strings.Contains(marker.name, text) {
				return marker.id, idx, marker.markers
			}
		}
	}
	return 0, 0, nil
}

func (p *Page) GetByIdMarker(id int64, idxLine int) Marker {

	m := Marker{}
	m.id = id

	s := strings.Split(p.line[idxLine+2], " ")
	ch := strings.Replace(s[0], ",", ".", 1)
	f, _ := strconv.ParseFloat(ch, 64)
	m.value = f

	chT := strings.Replace(p.line[idxLine+6], ".", "-", 2)
	t, _ := time.Parse("01-02-2006 15:04", chT)
	// log.Println(t, err)
	m.date = t
	return m
}

func (p *Page) helper(dbMrks []*Marker) (seqMarkers []Marker) {

	for _, marker := range dbMrks {
		for idx, text := range p.line {
			if marker.name == text {

				m := Marker{}
				m.id = marker.id

				s := strings.Split(p.line[idx+2], " ")
				ch := strings.Replace(s[0], ",", ".", 1)
				f, _ := strconv.ParseFloat(ch, 64)
				m.value = f

				chT := strings.Replace(p.line[idx+6], ".", "-", 2)

				t, err := time.Parse("01-02-2006 15:04", chT)

				log.Println(t, err)
				m.date = t
				seqMarkers = append(seqMarkers, m)

				log.Println("find marker ", p.line[idx+6], err)
				// collect to single(); 2 variant
			}
		}
	}
	return seqMarkers
}

//set each marker - markerId;
func (p *Page) Compare(dbMrks []*Marker) {
	/*compare; map[name]=[]int{}*/
	singleMrks := Marker{}
	// groupMrks := []Group{}
	g := Group{}

	if p.isGroup {
		// collectTo; plain -> how to compare data;
		log.Println("group", p.name, p.groupIndex)
		//get Id; & ids
		id, _, mrks := p.getMrkId(dbMrks)
		g.id = id
		// g.markers = mrks
		//fill each marker
		log.Println(id, mrks, "Group markers", len(g.markers), dbMrks)
		for _, l := range dbMrks {
			if l.id == g.id {
				//compare name; set id- to marker object
				for _, sm := range l.markers {
					for _, text := range p.line {
						if sm.name == text {
							fill to marker object; append - to groupMarker; time; value;
							log.Println(text, 99)
						}
					}
				}
			}
		}

		// for _, mk := range dbMrks
		//call func, concurrent ?
		log.Println(g.markers, "res")
	} else {
		//func()? gorutine ?
		singleMrks = p.helper(dbMrks)[0]
	}

	log.Println(singleMrks)
}

func NewPage() *Page {
	return &Page{}
}

func readPdf4(path string) ([]*Page, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}
	totalPage := r.NumPage()
	pages := []*Page{}

	//On(4)
	//each page - read; write to struct
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		// 1 file -> N page -> ;
		// check inside page -> ОАК || single ? -> send page - another func ->
		rows, _ := p.GetTextByRow()
		/* write to struct(then manipalate Data) */
		page := &Page{}
		for _, row := range rows {
			for _, word := range row.Content {
				page.line = append(page.line, word.S)
				// общий анализ мочи/крови, etc
				if strings.Contains(word.S, "Общий анализ крови") {
					page.groupIndex = pageIndex
					page.isGroup = true
					log.Println("send group function, this page", pageIndex)
					page.name = word.S
				}
				//else sne signleFunc()
			}
		}
		pages = append(pages, page)
	}
	// log.Println(pages[0], pages[1], pages[2], len(pages[2].line), len(pages), "len pages")
	// p.pages = pages
	return pages, nil
}
