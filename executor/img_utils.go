package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/disintegration/gift"
)

func downloadImage(url string) (string, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return "", err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}

	fmt.Println(n, "bytes downloaded.")

	return fileName, nil
}

func uploadImage(url string, fileName string) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	fw, err := w.CreateFormFile("image", fileName)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return err
	}

	if fw, err = w.CreateFormField("name"); err != nil {
		return err
	}
	if _, err = fw.Write([]byte(fileName)); err != nil {
		return err
	}
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return err
	}

	return nil
}

func procImage(fileName string) (string, error) {
	reader, err := os.Open(fileName)
	if err != nil {
		fmt.Errorf("Failed to open image with error: %s\n", err)
		return "", err
	}
	defer reader.Close()

	src, _, err := image.Decode(reader)
	if err != nil {
		fmt.Errorf("Failed to decode image with error: %s\n", err)
		return "", err
	}

	g := gift.New(
		gift.Invert(),
	)

	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	var opt jpeg.Options
	opt.Quality = 100

	fileName = fileName + "_processed.jpg"

	out, _ := os.Create(fileName)
	err = jpeg.Encode(out, dst, &opt)
	if err != nil {
		fmt.Errorf("Failed to encode image with error: %s\n", err)
		return "", err
	}

	return fileName, nil
}
