package main

import ("fmt"
	"errors"
	"net/http"
	"io"
	"strings"
	"os"
	"os/exec"
	"image"
	_ "image/png"
)

var download_dir string = "/tmp/shawel-go/parse"

func download_file(url string, to_disk bool) (string, error) {
	// Делаем http запрос
	resp, err := http.Get(url)
	if(err != nil) {
		panic(err) // Паникуем в случае ошибки
	}
	defer resp.Body.Close() // Освобождаем ресурсы после всех действий
	if(resp.StatusCode != 200) {
		return "", errors.New(url+": "+resp.Status)
	}
	if(to_disk == false) {
		bodyBytes, err := io.ReadAll(resp.Body) // Получаем body в виде байтов
		if err != nil {
			panic(err)
		}
		return string(bodyBytes), nil // Возвращаем body в виде строки
	} else {
		var file_split []string
		var file string
		file_split = strings.Split(url, "/")
		file = file_split[len(file_split)-1]
		out, err := os.Create(download_dir+"/"+file)
		defer out.Close()
		if(err != nil) {
			panic(err)
		}
		_, err = io.Copy(out, resp.Body)
		if(err != nil) {
			panic(err)
		}
		return file, nil
	}
}

func get_attr(body string, tag string, attr string) [3]string {
	iframes := [3]string{}
	tag = "<"+tag+" "
	attr = attr+"=\""
	for i := 0; i < len(iframes); i++ {
		var params string
		_, body, _ = strings.Cut(body, tag) // ищем нужный нам тег
		params, _, _ = strings.Cut(body, ">") // получаем все между < и >
		iframes[i] = params // добавляем тег в массив
		var link string
		_, link, _ = strings.Cut(params, attr) // ищем нужный нам атрибут
		link, _, _ = strings.Cut(link, "\"") // получаем его значение
		iframes[i] = link
	}
	return iframes
}

func check_img_white_percent(file string) {
	file_data, err := os.Open(file)
	defer file_data.Close()
	if(err != nil) {
		panic(err)
	}
	img, _, err := image.Decode(file_data)
	if(err != nil) {
		panic(err)
	}
    bounds := img.Bounds()
    var white_pixels int = 0
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if(r == 0 && g == 0 && b == 0 && a == 0) {
				white_pixels += 1
			}
		}
	}
	if(float64(white_pixels / (bounds.Max.Y*bounds.Max.X)) >= 0.995) {
		err = os.Remove(file)
		if(err != nil && !os.IsNotExist(err)) {
			panic(err)
		}
	}
}

func main() {
	prefixes := map[int]string{
		0: "rasp",
		1: "test",
		2: "zvonki",
	}
	err := os.MkdirAll(download_dir, 0750)
	if(err != nil && !os.IsExist(err)) {
		panic(err)
	}
	defer os.RemoveAll(download_dir)
	body, err := download_file("https://engschool9.ru/content/raspisanie.html", false)
	if(err != nil) {
		panic(err)
	}
	links := get_attr(body, "iframe", "src")
	var filename string
	for i, link := range links {
		prefix, ok := prefixes[i]
		if(ok == false || link == "") {
			continue
		}
		filename, err = download_file(link, true)
		if(err != nil) {
			fmt.Println(err)
			continue
		}
		err := os.Mkdir(download_dir+"/"+prefix, 0750)
		if(err != nil && !os.IsExist(err)) {
			panic(err)
		}
		cmd := exec.Command("convert", "-colorspace", "RGB", "-density", "200", download_dir+"/"+filename, download_dir+"/"+prefix+"/out.png")
		_, err = cmd.Output()
		if(err != nil) {
			panic(err)
		}
		imgs, err := os.ReadDir(download_dir+"/"+prefix)
		if(err != nil) {
			panic(err)
		}
		for _, img := range imgs {
			check_img_white_percent(download_dir+"/"+prefix+"/"+img.Name())
		}
	}
}