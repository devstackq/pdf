package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
type DBMarkerIds struct {
	id      int
	ids     []int
	isGroup bool
	name    string
	abbr    string
	marker  []Marker
}
type Marker struct {
	id    int
	value float64
}

func (p *Pdf) Read() {
	content, err := readPdf4(os.Args[1]) // Read local pdf file
	if err != nil {
		panic(err)
	}
	p.pages = content

	var markersDb []*DBMarkerIds

	markersDb = append(markersDb, &DBMarkerIds{}, &DBMarkerIds{})

	markersDb[0].abbr = "oak"
	markersDb[0].id = 43
	markersDb[0].ids = []int{17, 44}
	markersDb[0].name = "Общий анализ крови"

	markersDb[1].abbr = "ttg"
	markersDb[1].id = 69
	markersDb[1].name = "ТТГ"

	markersDb[0].abbr = "HGB"
	markersDb[0].id = 23
	markersDb[0].name = "Гемоглобин"

	for _, page := range p.pages {
		if page.groupIndex > 0 {
			//group logic;
		} else {
			for idx, text := range page.line {
				for _, dbMrk := range markersDb {
					//here write value
					// 	// dbMrk.name == text || if abbr empty ? -> comapre by name
					if dbMrk.abbr == text {
						log.Println(page.line[idx+2], "val", idx, text)
					}
				}
			}
		}

	}

	is group ? -> fill; by ids -> for loop range - ids ->fill id; value;
	//save - all; result - like cache ->
	// then final result ->

	for _, marker := range markersDb {
		//1: for - analyses
		//2 for - save Db
		if len(marker.ids) > 0 {
			//fill group marker
		} else {
			//fill single marker
		}
	}

	fmt.Println(content)

}

// func(p *Pdf compchan)

/*
if user -> want auth ? -> call
select * names; abbr -> set in Map[name]=id...
hash sum - of file; compare -> next file is  = uniq ?

readPdf -> save || share -> isAuth -> authByPhone -> fill relative;
case1: share -> client json -> []groups; []single -> request - fetchs(relId, m_id, value) (save in db markers)
client -> sent url\share -> markerId; relId; - receive link

readPdf solve:
return  {
1 variant: []marker -> for graphic
2 variant: []group markers; by date;  []single markers; by date; -> for save data Db
}

save || share :
case1: send request -> for loop []group => {
	fetch(groupId : 43, relId : 22,  []{id: 3, val: 2.1; id : 15 : val: 45})..
}


loop - pdfFiles {
	loop - pdfPages;
	prepare1Variant() {
		getDbMarkerId() // name || abbr
		merge single/group ->
	}
	prepare2Variant() {
		filter by date;
		devide single & group markers
		getDbMarkerId() // name || abbr

	}
}
return client
*/

func main() {
	pdf := NewPdf()
	pdf.Read()
	// excel := NewExcel()
	// excel.Read()
}

type Page struct {
	line       []string
	groupIndex int
	name       string
	totalPage  int
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
		// println(pageIndex)

		/* write to struct(then manipalate Data) */
		page := &Page{}

		for _, row := range rows {
			// println(">>>> row: ", row.Position, idx)
			for _, word := range row.Content {
				page.line = append(page.line, word.S)

				//Қан  талдауы / Общий анализ крови || общий анализ мочи, etc
				if strings.Contains(word.S, "Общий анализ крови") {
					page.groupIndex = pageIndex
					log.Println("send group function, this page", pageIndex)
					page.name = word.S
				}
				//else sne signleFunc()
			}
		}
		pages = append(pages, page)
	}
	log.Println(pages[0], pages[1], pages[2], len(pages))
	// p.pages = pages
	return pages, nil
}
