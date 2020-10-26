package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

/*
######## ##    ## ########  ########  ######
   ##     ##  ##  ##     ## ##       ##    ##
   ##      ####   ##     ## ##       ##
   ##       ##    ########  ######    ######
   ##       ##    ##        ##             ##
   ##       ##    ##        ##       ##    ##
   ##       ##    ##        ########  ######
*/

type ouiEntry struct {
	VendorName    string   `json:"vendorName"`
	VendorAddress []string `json:"vendorAddress"`
}

type ouiDatabase struct {
	OUIDatabase map[string]ouiEntry `json:"ouiDatabase"`
}

func (db *ouiDatabase) ToText() string {
	var lines []string

	for oui, data := range db.OUIDatabase {
		lines = append(lines, fmt.Sprintf("%s\t%s", oui, data.VendorName))
	}

	return strings.Join(lines, "\n")
}

func (db *ouiDatabase) ToCSV() string {
	var lines []string

	for oui, data := range db.OUIDatabase {
		lines = append(lines, fmt.Sprintf(`"%s","%s"`, oui, data.VendorName))
	}

	return strings.Join(lines, "\n")
}

func (db *ouiDatabase) ToJSON() string {
	json, _ := json.MarshalIndent(db, "", "    ")
	return string(json)
}

/*
######## #### ##       ########    ##     ##    ###    ##    ## ########  ##       #### ##    ##  ######
##        ##  ##       ##          ##     ##   ## ##   ###   ## ##     ## ##        ##  ###   ## ##    ##
##        ##  ##       ##          ##     ##  ##   ##  ####  ## ##     ## ##        ##  ####  ## ##
######    ##  ##       ######      ######### ##     ## ## ## ## ##     ## ##        ##  ## ## ## ##   ####
##        ##  ##       ##          ##     ## ######### ##  #### ##     ## ##        ##  ##  #### ##    ##
##        ##  ##       ##          ##     ## ##     ## ##   ### ##     ## ##        ##  ##   ### ##    ##
##       #### ######## ########    ##     ## ##     ## ##    ## ########  ######## #### ##    ##  ######
*/

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

/*
########     ###    ########    ###    ########     ###     ######  ########    ##     ##    ###    ##    ## ########  ##       #### ##    ##  ######
##     ##   ## ##      ##      ## ##   ##     ##   ## ##   ##    ## ##          ##     ##   ## ##   ###   ## ##     ## ##        ##  ###   ## ##    ##
##     ##  ##   ##     ##     ##   ##  ##     ##  ##   ##  ##       ##          ##     ##  ##   ##  ####  ## ##     ## ##        ##  ####  ## ##
##     ## ##     ##    ##    ##     ## ########  ##     ##  ######  ######      ######### ##     ## ## ## ## ##     ## ##        ##  ## ## ## ##   ####
##     ## #########    ##    ######### ##     ## #########       ## ##          ##     ## ######### ##  #### ##     ## ##        ##  ##  #### ##    ##
##     ## ##     ##    ##    ##     ## ##     ## ##     ## ##    ## ##          ##     ## ##     ## ##   ### ##     ## ##        ##  ##   ### ##    ##
########  ##     ##    ##    ##     ## ########  ##     ##  ######  ########    ##     ## ##     ## ##    ## ########  ######## #### ##    ##  ######
*/

func loadRawDatabase(fileName string) (db bytes.Buffer, err error) {
	devMessage("Entering loadRawDatabase()")

	db, err = loadData(fileName)
	if err != nil {
		return db, fmt.Errorf("Could not load database: %s", err)
	}

	db, err = decompressData(db)
	if err != nil {
		return db, fmt.Errorf("Could not decompress database: %s", err)
	}

	return
}

func parseRawDatabase(content bytes.Buffer) (ouiDB ouiDatabase, err error) {
	var reBase16 = regexp.MustCompile(`^(.+?)\s+\(base 16\)\s+(.+?)$`)
	var inVendorBlock bool
	var vendorOUI string
	var vendorName string
	var vendorAddress []string

	devMessage("Entering parseRawDatabase()")

	if ouiDB.OUIDatabase == nil {
		ouiDB.OUIDatabase = make(map[string]ouiEntry)
	}

	inVendorBlock = false
	fs := bufio.NewScanner(bytes.NewReader(content.Bytes()))
	for fs.Scan() {
		if inVendorBlock {
			trimmed := strings.TrimSpace(fs.Text())
			if trimmed == "" {
				ouiDB.OUIDatabase[vendorOUI] = ouiEntry{VendorName: vendorName, VendorAddress: vendorAddress}
				inVendorBlock = false
				vendorAddress = []string{}
				continue
			}
			vendorAddress = append(vendorAddress, trimmed)
		} else if reBase16.Match(fs.Bytes()) {
			inVendorBlock = true
			ouiData := reBase16.FindStringSubmatch(fs.Text())
			vendorOUI = strings.ToLower(ouiData[1])
			vendorName = ouiData[2]
		}
	}

	return
}

func loadDatabase(fileName string) (db ouiDatabase, err error) {
	rawDB, rawDBErr := loadRawDatabase(config.DatabaseFile)
	if rawDBErr != nil {
		return db, fmt.Errorf("Error reading local OUI database: %s", rawDBErr)
	}
	db, err = parseRawDatabase(rawDB)
	if err != nil {
		return db, fmt.Errorf("Error parsing local OUI database: %s", err)
	}
	return
}
