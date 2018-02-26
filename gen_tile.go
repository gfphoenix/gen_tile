package main

import (
	"flag"
	"image"
	"image/draw"
	"image/png"
_	"image/jpeg"
	"log"
	"os"
	"strconv"
	"strings"
)
var input_file string
var output_file string
var inch_size string // 4x6 5x7
var sheet_size string // in pixels
var layout string // 4x3 rows x cols
var padding string
//var paddingX int
//var paddingY int
var dpi int
func init() {
	flag.StringVar(&input_file, "i", "", "input single photo")
	flag.StringVar(&output_file, "o", "a.out.png", "output png name")
	flag.StringVar(&inch_size, "inch_size", "4x6", "available size is `1-in`")
	flag.StringVar(&sheet_size, "sheet_size", "", "sheet size in pixels") 
	flag.StringVar(&layout, "layout", "4x3", "照片布局，行数x列数")
	flag.StringVar(&padding, "padding", "60x20", "<paddingX>x<paddingY>")
	flag.IntVar(&dpi, "dpi", 300, "DPI, default is 300dpi")
	flag.Parse()
}

func parseInt2(axb string) (int, int) {
	sp := strings.Split(axb, "x")
	if len(sp)!=2 {
		log.Fatalln("need <a>x<b>, got ", axb)
	}
	a, err := strconv.Atoi(strings.TrimSpace(sp[0]))
	if err != nil {
		log.Fatalln(err)
	}
	b, err := strconv.Atoi(strings.TrimSpace(sp[1]))
	if err != nil {
		log.Fatalln(err)
	}
	return a, b
}

type TileSheet struct {
	sheet_width int
	sheet_height int
	padding_x int
	padding_y int
	rows int
	cols int
}

func parse() TileSheet {
	if ! flag.Parsed() {
		flag.Parse()
	}
	var sheet_w, sheet_h int
	if sheet_size != "" {
		sheet_w, sheet_h = parseInt2(sheet_size)
	}else {
		sheet_w, sheet_h = parseInt2(inch_size)
		sheet_w *= dpi
		sheet_h *= dpi
	}
	rows, cols := parseInt2(layout)
	paddingX, paddingY := parseInt2(padding)
	return TileSheet {
		sheet_width : sheet_w,
		sheet_height: sheet_h,
		padding_x : paddingX,
		padding_y : paddingY,
		rows : rows,
		cols : cols,
	}
}

func (ts *TileSheet) Gen(img image.Image) image.Image {
	zp := image.ZP
	r := image.Rectangle{Min:zp, Max:image.Point{ts.sheet_width, ts.sheet_height}}
	dst := image.NewNRGBA(r)
	for i:=0; i<len(dst.Pix); i++ {
		dst.Pix[i]=uint8(255)
	}
	
	r0 := img.Bounds()
	var dx, dy int
	off_x, off_y := (r.Dx()-r0.Dx())/2, (r.Dy()-r0.Dy())/2
	if ts.cols >=2 {
		dx = (r.Max.X - ts.padding_x*2 - r0.Dx())/(ts.cols-1)
		off_x = ts.padding_x
	}
	if ts.rows >=2 {
		dy = (r.Max.Y - ts.padding_y*2 - r0.Dy())/(ts.rows-1)
		off_y = ts.padding_y
	}
	r0 = r0.Add(image.Point{off_x, off_y})
	for x:=0; x<ts.cols; x++ {
		for y:=0; y<ts.rows; y++ {
			draw.Draw(dst, r0.Add(image.Point{dx*x, dy*y}), img, zp, draw.Over)
		}
	}
	return dst
}
func handle(in, out string) {
	fin, err := os.Open(in)
	if err!=nil {
		log.Fatalln(err)
	}
	defer fin.Close()
	fout, err := os.OpenFile(out, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err!=nil {
		log.Fatalln(err)
	}
	defer fout.Close()
	
	img, tag, err :=  image.Decode(fin)
	if err!=nil {
		log.Fatalln(err)
	}
	log.Println("tag = ", tag)
	log.Println("input bounds = ", img.Bounds())
	
	ts := parse()
	dst := ts.Gen(img)
	err = png.Encode(fout, dst)
	
	if err == nil {
		log.Printf("write to '%s' OK!\n", output_file)
	}else {
		log.Println("encode to file failed:", err)
	}
}
func main() {
	if input_file=="" {
		log.Println("input file must not be empty")
		flag.PrintDefaults()
		os.Exit(1)
	}
	handle(input_file, output_file)
}
