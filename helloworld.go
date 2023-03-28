package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

func GetResp(url string) (error, string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	// log.Printf(sb)
	// fmt.Println(err)
	return err, sb
}

func ExecDownloadCmd(S3object string, directory string) {
	fmt.Println(S3object)
	out, err := exec.Command("s5cmd", "cp", S3object, directory).Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func ExecUploadCmd(object string, S3directory string) {
	fmt.Println(object)
	exec.Command("s5cmd", "rm", object).Output()
	out, err := exec.Command("s5cmd", "cp", object, S3directory).Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}
func main() {

	// var hostname string = "http://book-dev.ap-southeast-1.elasticbeanstalk.com/api/books/"
	var hostname string = "s3://elasticbeanstalk-ap-southeast-1-017761180421/"
	start := time.Now()
	var index int = 1
	for index <= 6 {
		// url := hostname + string(index)
		// err, sb := GetResp(url)

		videoPath := hostname + "Video/video_path_" + strconv.Itoa(index)
		txtPath := hostname + "Data/txt_path_" + strconv.Itoa(index)

		ExecDownloadCmd(videoPath, "F:\\project\\DATN\\DownloadS3\\S3Download\\Video")
		ExecDownloadCmd(txtPath, "F:\\project\\DATN\\DownloadS3\\S3Download\\Data")

		index++
		// if err != nil {
		// 	index++
		// }
	}
	end := time.Now()
	fmt.Println(start)
	fmt.Println(end)
}
