package pdf

type PdfReader interface {
	Read()
}

type Pdf struct{}

func NewPdf() *Pdf {
	return &Pdf{}
}

func (p *Pdf) Read() {
}
