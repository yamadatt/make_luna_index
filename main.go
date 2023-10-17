package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// テンプレートへ引き渡すオブジェクト
// メンバはエクスポート(大文字)しておかないと参照できない
type ExportFormat struct {
	Title             string
	Filename          string
	Primarycategory   string
	Secondarycategory string
	Date              string
	Image             string
	Files             []string
	Inyou             string
	Body              string
	Slug              string
}

func main() {

	searchPath := "photo"

	// 構造体にアクセスするため、初期化。hはhtmlの略で使用した。
	var h ExportFormat

	filepath.WalkDir(searchPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("failed filepath.WalkDir")
		}

		if !info.IsDir() {
			if filepath.Base(path) == "index.html" {
				//debug ディレクトリの最後の文字列を取得したい
				fmt.Println("処理対象のディレクトリ", filepath.Dir(path))
				directory_name := strings.Split(filepath.Dir(path), "/")
				h.Slug = directory_name[len(directory_name)-1]

				fileInfos, _ := ioutil.ReadFile(path)
				stringReader := strings.NewReader(string(fileInfos))

				doc, err := goquery.NewDocumentFromReader(stringReader)

				if err != nil {
					fmt.Print("goquery failed")
				}

				// title
				doc.Find("title").Each(func(_ int, s *goquery.Selection) {
					h.Title = strings.TrimSuffix(s.Text(), " - Photo Gallery")
				})

				h.Date = time2str(time.Now())

				files, _ := ioutil.ReadDir(filepath.Dir(path))

				image_files := []string{}

				for _, f := range files {
					// fmt.Println(f.Name())
					if f.Name() != "index.html" {
						image_files = append(image_files, f.Name())
					}
				}
				fmt.Println(image_files)

				h.Image = image_files[0]
				h.Files = image_files

				//出力用の関数を呼ぶ
				h.make_export(filepath.Dir(path))
			}

		}

		return nil
	})

}

// time.Time型 -> 文字列
func time2str(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z07:00")
}

func (h *ExportFormat) make_export(path string) {

	// テンプレートオブジェクトを生成
	tmpl := template.Must(template.ParseFiles("./test.tmpl"))

	// テンプレートへ渡すデータを作る
	g := &ExportFormat{
		h.Title,
		h.Filename,
		h.Primarycategory,
		h.Secondarycategory,
		h.Date,
		h.Image,
		h.Files,
		h.Inyou,
		h.Body,
		h.Slug,
	}

	// create file
	fp, err := os.Create(filepath.Join(path, "index.md"))
	if err != nil {
		fmt.Println("error creating file", err)
	}
	defer fp.Close()

	// write file
	if err = tmpl.Execute(fp, g); err != nil {
		fmt.Println(err)
	}

}
