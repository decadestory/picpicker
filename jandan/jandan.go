package jandan

import (
	"fmt"
	"io"
	"jerry/guid"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	url    string = "http://jandan.net"
	dirNv  string = "./download/ooxx/"
	dirbor string = "./download/pic/"
)

var (
	pageUrls    chan string
	imgUrls     chan string
	downloadCnt chan int
)

func init() {
	pageUrls = make(chan string, 100)
	imgUrls = make(chan string, 1000)
	downloadCnt = make(chan int, 1000)

}

//GetNv get meizi
func GetNv(cnt int) {
	fmt.Println("开始获取", cnt, "张妹子图")
	go queryPages("ooxx")
	go queryImgUrls()
	go download(cnt, dirNv)

	for item := range downloadCnt {
		fmt.Println("已经下载完成第" + strconv.Itoa(item) + "张图片")
		if item == cnt {
			fmt.Println("下载图片完成")
			close(downloadCnt)
			close(imgUrls)
			close(pageUrls)
		}
	}
}

//GetPic get pic
func GetPic(cnt int) {
	fmt.Println("开始获取", cnt, "张无聊图")
	go queryPages("pic")
	go queryImgUrls()
	go download(cnt, dirbor)

	for item := range downloadCnt {
		fmt.Println("已经下载完成第" + strconv.Itoa(item) + "张图片")
		if item == cnt {
			fmt.Println("下载图片完成")
			close(imgUrls)
			close(pageUrls)
			close(downloadCnt)
		}
	}
}

func queryPages(cat string) {
	index := 1
	for {
		curl := url + "/" + cat + "/page-" + strconv.Itoa(index) + "#comments"
		pageUrls <- curl
		time.Sleep(1 * 1e9)
		index++
	}
}

func queryImgUrls() {
	for item := range pageUrls {
		doc, err := goquery.NewDocument(item)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find(".commentlist img").Each(func(i int, s *goquery.Selection) {
			imgURLjpg, exist := s.Attr("src")
			imgURLGif, existgif := s.Attr("org_src")

			if existgif {
				imgUrls <- imgURLGif
			} else if exist {
				imgUrls <- imgURLjpg
			}
		})
	}
}

func download(cnt int, dir string) {
	basePath, err := os.Stat("./download/")
	if err != nil || basePath.IsDir() == false {
		err := os.Mkdir("./download/", os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败:", err)
		}
	}

	info, err := os.Stat(dir)
	if err != nil || info.IsDir() == false {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败:", err)
		}
	}

	dlcnt := 1
	for item := range imgUrls {
		fileName := guid.New().String()
		fileExt := path.Ext(item)
		filePath := dir + fileName + fileExt

		data, err := http.Get("http:" + item)
		if err != nil {
			fmt.Println("下载图片失败:", err)
		}

		f, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		io.Copy(f, data.Body)
		downloadCnt <- dlcnt
		dlcnt++
	}
}
