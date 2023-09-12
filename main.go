package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// GetImageHistory reads image history from a file and returns it as a map.
func GetImageHistory(filename string)(map[int]bool, error){
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var imageData bytes.Buffer

	_,err = io.Copy(&imageData, file)
	if err != nil {
		return nil, err
	}
	
	ImageHistory := map[int] bool{}
	mapIndex := imageData.String()
	indexes := bytes.Fields([]byte(mapIndex))

	for _, number := range indexes {
		num, err := strconv.Atoi(string(number))
		if err == nil {
			ImageHistory[num] = true
		}
	}		

	return ImageHistory, nil
}

// GetPageContent sends an HTTP GET request and returns the response.
func GetPageContent(url string) (*http.Response, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("unexpected status code %v", response.Status)
	}

	return response, nil
}

// GetImageUrl extracts an image URL from a string.
func GetImageUrl(HTMLWebContent, imageUrlTarget string) (string, error) {
	imageIndex := strings.Index(HTMLWebContent, imageUrlTarget)
	if imageIndex == -1 {
		return "", fmt.Errorf("image URL not found")
	}

	urlImage := HTMLWebContent[imageIndex+len(imageUrlTarget):]
	urlEndIndex := strings.Index(urlImage, `"`)
	if urlEndIndex == -1 {
		return "", fmt.Errorf("url final part not found")
	}

	// Extract the URL
	imageUrl := urlImage[:urlEndIndex]
	return imageUrl, nil
}

// CalculateMD5 calculates the MD5 hash of data.
func CalculateMD5(data bytes.Buffer) (string, error) {
	hash := md5.New()
	_, err := io.Copy(hash, &data)
	if err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)	
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

// WriteBufferToFile writes the content of a buffer to a file.
func WriteBufferToFile(fileName string, buffer *bytes.Buffer) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	
	_, err = io.Copy(file, buffer)
	
	return err
}

func main() {
	filename := "Image_History.txt"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return 
	}
	defer file.Close()

	downloadedImages, err := GetImageHistory(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for id := 1; id <= 2826; id++{

		exist := downloadedImages[id]

		if exist {
			fmt.Printf("Image %v already downloaded.\n", id)
		} else {
			webPageUrl := "https://xkcd.com/" + strconv.Itoa(id) + "/info.0.json"
			imageUrlTarget := `"img": "`
			
			response, err := GetPageContent(webPageUrl)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			defer response.Body.Close()

			htmlWebPage, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			imageUrl, err := GetImageUrl(string(htmlWebPage), imageUrlTarget)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			response, err = GetPageContent(imageUrl)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			defer response.Body.Close()

			var imageData bytes.Buffer
			_, err = io.Copy(&imageData, response.Body)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			md5Hash, err := CalculateMD5(imageData)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fileName := md5Hash + ".png"

			err = WriteBufferToFile(fileName, &imageData)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			downloadedImages[id] = true
			file.WriteString(strconv.Itoa(id) + " ")
		}
	}
}
