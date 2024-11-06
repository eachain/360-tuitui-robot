package client_test

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/eachain/360-tuitui-robot/client"
)

func ExampleClient_UploadImage() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	image := flag.String("image", "", "image file path")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	fp, err := os.Open(*image)
	if err != nil {
		log.Println(err)
	}
	defer fp.Close()

	mediaId, err := cli.UploadImage(fp, filepath.Base(*image))
	if err != nil {
		log.Printf("upload image: %v", err)
		return
	}
	log.Printf("upload image, media id: %v", mediaId)
}

func ExampleClient_UploadFile() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	file := flag.String("file", "", "file path")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	fp, err := os.Open(*file)
	if err != nil {
		log.Println(err)
	}
	defer fp.Close()

	mediaId, err := cli.UploadImage(fp, filepath.Base(*file))
	if err != nil {
		log.Printf("upload file: %v", err)
		return
	}
	log.Printf("upload file, media id: %v", mediaId)
}

func ExampleClient_UploadFromDisk() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	file := flag.String("file", "", "file path")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	mediaId, isImage, err := cli.UploadFromDisk(*file)
	if err != nil {
		log.Printf("upload file: %v", err)
		return
	}
	log.Printf("upload file, media id: %v, is image: %v", mediaId, isImage)
}

func ExampleClient_UploadFromURL() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	url := flag.String("url", "https://so1.360tres.com/t012cdb572f41b93733.png", "file url path")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	mediaId, isImage, err := cli.UploadFromURL(*url)
	if err != nil {
		log.Printf("upload file from url: %v", err)
		return
	}
	log.Printf("upload file from url, media id: %v, is image: %v", mediaId, isImage)
}

func ExampleClient_FetchMediaTemporaryURL() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	mediaId := flag.String("media", "993a9e16b365f4529fc9ccfa", "file media id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	urlOf, warn, err := cli.FetchMediaTemporaryURL([]string{*mediaId})
	if err != nil {
		log.Printf("fetch media temporary url: %v", err)
		return
	}
	if warn != nil {
		log.Printf("fetch media temporary url failed media_ids: %v", warn.Explains)
	}

	log.Printf("fetch media temporary url: %v", urlOf)
}

func ExampleClient_GetMediaTemporaryURL() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	mediaId := flag.String("media", "993a9e16b365f4529fc9ccfa", "file media id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	url, err := cli.GetMediaTemporaryURL(*mediaId)
	if err != nil {
		log.Printf("get media temporary url: %v", err)
		return
	}

	log.Printf("get media temporary url: %v", url)
}
