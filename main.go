package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ledongthuc/pdf"
)

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
	lines []string
}

func NewPdf() *Pdf {
	return &Pdf{}
}

func (p *Pdf) Read() {
	// use pdf libs
	content, err := readPdf4(os.Args[1]) // Read local pdf file
	if err != nil {
		panic(err)
	}

	fmt.Println(content)
	// if strings.Contains(content, "")
	// s := strings.Split(content, " ")
	// fmt.Println(s, len(s))
	// for idx, text := range s {
	// 	fmt.Println(text, idx)
	// }
}

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

	excel := NewExcel()
	excel.Read()
}

func readPdf4(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		// 1 file -> N page -> ;
		// check inside page -> ОАК || single ? -> send page - another func ->

		rows, _ := p.GetTextByRow()
		println(pageIndex)

		for _, row := range rows {
			// println(">>>> row: ", row.Position, idx)
			fmt.Println(row.Content, row.Position)
			// if strings.Contains(row.Content, "Общий анализ крови") {
			// for _, word := range row.Content {
			// 	fmt.Println(word.S)
			// }
		}
	}
	return "", nil
}
