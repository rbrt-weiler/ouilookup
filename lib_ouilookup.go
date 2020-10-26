package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func compressData(content bytes.Buffer) (bytes.Buffer, error) {
	var buf bytes.Buffer

	devMessage("Entering compressData()")

	gzipWriter, gzipWriterErr := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if gzipWriterErr != nil {
		return buf, fmt.Errorf("Could not create gzip stream: %s", gzipWriterErr)
	}
	gzipWriter.ModTime = time.Now()
	gzipWriter.Comment = fmt.Sprintf("created with %s", toolID)
	_, writeErr := gzipWriter.Write(content.Bytes())
	if writeErr != nil {
		return buf, fmt.Errorf("Could not write to gzip buffer: %s", writeErr)
	}
	flushErr := gzipWriter.Flush()
	if flushErr != nil {
		return buf, fmt.Errorf("Could not flush gzip buffer: %s", flushErr)
	}
	closeErr := gzipWriter.Close()
	if closeErr != nil {
		return buf, fmt.Errorf("Could not close gzip stream: %s", closeErr)
	}

	return buf, nil
}

func storeData(fileName string, content bytes.Buffer) error {
	devMessage("Entering storeData()")

	fileHandle, fileErr := os.Create(fileName)
	if fileErr != nil {
		return fmt.Errorf("Could not create outfile: %s", fileErr)
	}
	defer fileHandle.Close()

	fileWriter := bufio.NewWriter(fileHandle)
	_, writeErr := fileWriter.WriteString(content.String())
	if writeErr != nil {
		return fmt.Errorf("Could not write to outfile: %s", writeErr)
	}

	flushErr := fileWriter.Flush()
	if flushErr != nil {
		return fmt.Errorf("Could not flush file buffer: %s", flushErr)
	}

	return nil
}

func loadData(fileName string) (bytes.Buffer, error) {
	var retVal bytes.Buffer

	devMessage("Entering loadData()")

	content, contentErr := ioutil.ReadFile(fileName)
	if contentErr != nil {
		return retVal, fmt.Errorf("Could not read file: %s", contentErr)
	}
	retVal = *bytes.NewBuffer(content)

	return retVal, nil
}

func decompressData(content bytes.Buffer) (bytes.Buffer, error) {
	var buf bytes.Buffer

	devMessage("Entering decompressData()")

	reader, readerErr := gzip.NewReader(&content)
	if readerErr != nil {
		return buf, fmt.Errorf("Could not read compressed data: %s", readerErr)
	}

	_, bufErr := buf.ReadFrom(reader)
	if bufErr != nil {
		return buf, fmt.Errorf("Could not read compressed data: %s", bufErr)
	}

	return buf, nil
}
