package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
)

func downloadFile(filepath string, url string) (err error) {

	fmt.Printf("Download file '%s' from URL %s\n", filepath, url)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Get the data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}

	return nil
}

func gunzipFile(gzFilePath, dstFilePath string) (int64, error) {
	fmt.Printf("Gunzip file '%s' to '%s'\n", gzFilePath, dstFilePath)

	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		return 0, fmt.Errorf("Failed to open file %s for unpack: %s", gzFilePath, err)
	}
	dstFile, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return 0, fmt.Errorf("Failed to create destination file %s for unpack: %s", dstFilePath, err)
	}

	ioReader, ioWriter := io.Pipe()

	go func() { // goroutine leak is possible here
		gzReader, _ := gzip.NewReader(gzFile)
		// it is important to close the writer or reading from the other end of the
		// pipe or io.copy() will never finish
		defer func(){
			gzFile.Close()
			gzReader.Close()
			ioWriter.Close()
		}()

		io.Copy(ioWriter, gzReader)
	}()

	written, err := io.Copy(dstFile, ioReader)
	if err != nil {
		return 0, err // goroutine leak is possible here
	}
	ioReader.Close()
	dstFile.Close()

	fi1, _ := os.Stat(gzFilePath)
	fi2, _ := os.Stat(dstFilePath)
	fmt.Printf("  %s (%dMB) -> %s (%dMB)\n", gzFilePath, int32(math.Round(float64(fi1.Size())/1024/1024)), dstFilePath, int32(math.Round(float64(fi2.Size())/1024/1024)))

	return written, nil
}