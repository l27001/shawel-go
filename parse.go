package main

import (//"fmt"
	"net/http"
	"io"
	"strings"
	"os"
)

func download_file(url string, to_disk bool) string {
	// Делаем http запрос
	resp, err := http.Get(url)
	if(err != nil) {
		panic(err) // Паникуем в случае ошибки
	}
	defer resp.Body.Close() // Освобождаем ресурсы после всех действий
	if(resp.StatusCode != 200) {
		panic("HTTP код не равен 200: "+resp.Status)
	}
	if(to_disk == false) {
		bodyBytes, err := io.ReadAll(resp.Body) // Получаем body в виде байтов
	    if err != nil {
	        panic(err)
	    }
	    return string(bodyBytes) // Возвращаем body в виде строки
	} else {
		var file_split []string
		var file string
    	file_split = strings.Split(url, "/")
    	file = file_split[len(file_split)-1]
		out, err := os.Create(file)
    	defer out.Close()
    	if(err != nil) {
    		panic(err)
    	}
    	_, err = io.Copy(out, resp.Body)
    	if(err != nil) {
    		panic(err)
    	}
    	return file
	}
}

func get_attr(body string, tag string, attr string) []string {
	var iframes []string
	tag = "<"+tag+" "
	attr = attr+"=\""
	tags_count := strings.Count(body, tag)
    for i := 0; i < tags_count; i++ {
    	var params string
    	_, body, _ = strings.Cut(body, tag) // ищем нужный нам тег
    	params, _, _ = strings.Cut(body, ">") // получаем все между < и >
    	iframes = append(iframes, params) // добавляем тег в массив
    	var link string
    	_, link, _ = strings.Cut(params, attr) // ищем нужный нам атрибут
    	link, _, _ = strings.Cut(link, "\"") // получаем его значение
    	iframes[i] = link
    }
    return iframes
}

func main() {
	err := os.Mkdir("/tmp/shawel-go", 750)
	if(err != nil && !os.IsExist(err)) {
		panic(err)
	}
	var body string
	var links []string
	body = download_file("https://engschool9.ru/content/raspisanie.html", false)
    links = get_attr(body, "iframe", "src")
    for _, link := range links {
    	download_file(link, true)
    }
}