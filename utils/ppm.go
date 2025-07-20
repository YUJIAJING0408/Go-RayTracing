// 使用模版生成PPM文件

package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// 像素结构体
type Pixel struct {
	R, G, B int
}

func (pix Pixel) GetString() string {
	return fmt.Sprintf("%d %d %d ", pix.R, pix.G, pix.B)
}

// 定义 PPM 模板
const ppmTemplate = `P3
{{.Width}} {{.Height}}
{{.Max}}
{{range $i, $p := .Pixels}}
{{- if $i}}{{if mod $i $.Width}} {{else}}
{{end}}{{end -}}
{{$p.R}} {{$p.G}} {{$p.B -}}
{{end}}
`

// PPM 图像结构体
type PPMImage struct {
	Width  int     // 图像宽度
	Height int     // 图像高度
	Max    int     // 最大颜色值 (0-255)
	Pixels []Pixel // 像素数据 (按行存储)
}

func NewPPMImage(width, height, max int) *PPMImage {
	if max <= 0 {
		max = 255
	}
	if width <= 0 || height <= 0 {
		panic("宽高不能为负数")
	}
	return &PPMImage{
		Width:  width,
		Height: height,
		Max:    max,
		Pixels: make([]Pixel, 0, width*height),
	}
}

func (p *PPMImage) Full(pix []Pixel) {
	if len(pix) > p.Width*p.Height {
		p.Pixels = pix[0 : p.Width*p.Height]
	} else {
		p.Pixels = pix
	}
}

func (p *PPMImage) WriteAndSave(path string) {
	if len(p.Pixels) == 0 {
		panic("尚未填充像素无法保存")
	}
	// 注册自定义函数 (用于像素换行)
	funcMap := template.FuncMap{
		"mod": func(a, b int) int { return a % b },
	}

	// 解析模板
	tmpl, err := template.New("ppm").Funcs(funcMap).Parse(ppmTemplate)
	if err != nil {
		panic(err)
	}

	// 生成文件
	f, err := os.Create(path + "output.ppm")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	// 执行模板写入
	if err := tmpl.Execute(f, p); err != nil {
		panic(err)
	}
}

func (p *PPMImage) FastWriteAndSave(name string) {
	if len(p.Pixels) == 0 {
		panic("尚未填充像素无法保存")
	}

	// 生成文件
	f, err := os.Create(name + ".ppm")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	// 写入文件头
	fmt.Fprintf(f, "P3\n%d %d\n%d\n", p.Width, p.Height, p.Max)
	for j := 0; j < p.Height; j++ {
		var b strings.Builder
		b.Grow(p.Width * len(strconv.Itoa(p.Max)) * 4)
		for i := 0; i < p.Width; i++ {
			index := i + j*p.Width
			//fmt.Printf("w:%d,h%d index is %d\n", i, j, index)
			b.WriteString(p.Pixels[index].GetString())
		}
		fmt.Fprintf(f, "%s\n", b.String())
	}
	// 确保文件以换行符结束
	fmt.Fprint(f, "\n")
}
