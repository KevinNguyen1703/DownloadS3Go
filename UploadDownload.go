package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

type DBdata struct {
	Video       string `json:"Video"`
	Data        string `json:"Data"`
	Description string `json:"Description"`
}

type PostData struct {
	Video       string
	Data        string
	Description string
}

func GetResp(url string) (DBdata, error) {
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
	// sb := string(body)
	var dbData = &DBdata{}
	json.Unmarshal(body, dbData)
	return *dbData, err
}

func Post(url string, data DBdata) error {
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println("ERROR: wrong psot data")
	}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	return err
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

func ExecUploadCmd(object string, S3directory string) error {
	fmt.Println(object)
	fmt.Println(S3directory)
	exec.Command("s5cmd", "rm", S3directory).Output()
	out, err := exec.Command("s5cmd", "cp", object, S3directory).Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	return err
}

func ServerDownload(url string, args ...string) error {
	videoLocalLocation := "F:\\project\\DATN\\DownloadS3\\S3Download\\Video\\"
	dataLocalLocation := "F:\\project\\DATN\\DownloadS3\\S3Download\\Data\\"

	if len(args) == 2 {
		videoLocalLocation = args[0]
		dataLocalLocation = args[1]
	}
	sb, err := GetResp(url)
	if err != nil {
		var videoPath string = sb.Video
		var dataPath string = sb.Video

		ExecDownloadCmd(videoPath, videoLocalLocation)
		ExecDownloadCmd(dataPath, dataLocalLocation)
	}
	return err
}

func ServerUpload(url string, videoObject string, dataObject string, args ...string) error {
	var video string = ""
	var data string = ""
	var description string = ""

	if len(args) == 3 {
		video = args[0]
		data = args[1]
		description = args[2]
	}

	postData := DBdata{
		Video:       video,
		Data:        data,
		Description: description,
	}

	var videoS3Location = video
	var dataS3Location = data
	if len(args) == 5 {
		videoS3Location = args[3]
		dataS3Location = args[4]
	}
	errVideo := ExecUploadCmd(videoObject, videoS3Location)
	errData := ExecUploadCmd(dataObject, dataS3Location)
	if errVideo == nil && errData == nil {
		err := Post(url, postData)
		return err
	}

	return nil
}

func main() {

	var hostname string = "http://ec2-3-88-173-160.compute-1.amazonaws.com:9010/Data/"

	// var index int = 1
	// for {
	// 	url := hostname + strconv.Itoa(index)
	// 	err := ServerDownload(url,"F:\\project\\DATN\\DownloadS3\\S3Download\\Video\\","F:\\project\\DATN\\DownloadS3\\S3Download\\Data\\")
	// 	if err != nil {
	// 		index++
	// 	}
	// }

	var s3Storage string = "s3://lvtn-data/"
	var localStorage string = "F:\\project\\GolangPrj\\UploadDownloadS3\\S3Download\\"
	var index int = 1
	for {
		localVideoLocation := localStorage + "Video\\video_path_" + strconv.Itoa(index)
		localDataLocation := localStorage + "Data\\txt_path_" + strconv.Itoa(index)

		s3VideoLocation := s3Storage + "video/video_server_" + strconv.Itoa(index) + ".avi"
		s3DataLocation := s3Storage + "data/dataserver_" + strconv.Itoa(index) + ".txt"

		err := ServerUpload(
			hostname,
			localVideoLocation,
			localDataLocation,
			s3VideoLocation,
			s3DataLocation,
			"test description",
		)
		if err == nil {
			fmt.Printf("status done: %s\n", strconv.Itoa(index))
			index++
		}
	}

}
