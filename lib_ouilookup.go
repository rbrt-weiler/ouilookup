package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	resty "github.com/go-resty/resty/v2"
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

	devMessage("Entering ouiDatabase.ToText()")

	for oui, data := range db.OUIDatabase {
		lines = append(lines, fmt.Sprintf("%s\t%s", oui, data.VendorName))
	}

	devMessage("Leaving ouiDatabase.ToText()")
	return strings.Join(lines, "\n")
}

func (db *ouiDatabase) ToCSV() string {
	var lines []string

	devMessage("Entering ouiDatabase.ToCSV()")

	for oui, data := range db.OUIDatabase {
		lines = append(lines, fmt.Sprintf(`"%s","%s"`, oui, data.VendorName))
	}

	devMessage("Leaving ouiDatabase.ToCSV()")
	return strings.Join(lines, "\n")
}

func (db *ouiDatabase) ToJSON() string {
	devMessage("Entering ouiDatabase.ToJSON()")
	json, _ := json.MarshalIndent(db, "", "    ")
	devMessage("Leaving ouiDatabase.ToJSON()")
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

func fetchOnlineDatabase(url string) (bytes.Buffer, error) {
	var buf bytes.Buffer

	devMessage("Entering fetchOnlineDatabase()")

	client := resty.New()
	client.SetTimeout(time.Duration(config.Update.HTTPTimeoutSeconds) * time.Second)
	resp, respErr := client.R().Get(url)
	if respErr != nil {
		return buf, respErr
	}
	devMessage(fmt.Sprintf("Status Code   : %v", resp.StatusCode()))
	devMessage(fmt.Sprintf("Response Size : %v", resp.Size()))

	devMessage("Leaving fetchOnlineDatabase()")
	return *bytes.NewBuffer(resp.Body()), nil
}

func storeOnlineDatabase(url string, fileName string) error {
	var database *bytes.Buffer

	devMessage("Entering storeOnlineDatabase()")

	onlineDatabase, onlineDatabaseErr := fetchOnlineDatabase(url)
	if onlineDatabaseErr != nil {
		return fmt.Errorf("Error fetching online OUI database: %s", onlineDatabaseErr)
	}
	database = &onlineDatabase

	if strings.HasSuffix(fileName, ".gz") {
		compressedData, compressedDataErr := compressData(onlineDatabase)
		if compressedDataErr != nil {
			return fmt.Errorf("Error compressing OUI database: %s", compressedDataErr)
		}
		database = &compressedData
	}

	storeErr := storeData(fileName, *database)
	if storeErr != nil {
		return fmt.Errorf("Error storing local OUI database: %s", storeErr)
	}

	devMessage("Leaving storeOnlineDatabase()")
	return nil
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

	devMessage("Leaving storeData()")
	return nil
}

func loadData(fileName string) (bytes.Buffer, error) {
	var retVal bytes.Buffer

	devMessage("Entering loadData()")

	content, contentErr := os.ReadFile(fileName)
	if contentErr != nil {
		return retVal, fmt.Errorf("Could not read file: %s", contentErr)
	}
	retVal = *bytes.NewBuffer(content)

	devMessage("Leaving loadData()")
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

	devMessage("Leaving decompressData()")
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

	devMessage("Leaving compressData()")
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

	if strings.HasSuffix(fileName, ".gz") {
		db, err = decompressData(db)
		if err != nil {
			return db, fmt.Errorf("Could not decompress database: %s", err)
		}
	}

	devMessage("Leaving loadRawDatabase()")
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

	devMessage("Leaving parseRawDatabase()")
	return
}

func ouiToVendorDatabase(ouiDB ouiDatabase) (vendorDB map[string][]string, err error) {
	devMessage("Entering ouiToVendorDatabase()")
	if vendorDB == nil {
		vendorDB = make(map[string][]string)
	}
	for oui, data := range ouiDB.OUIDatabase {
		vendorDB[data.VendorName] = append(vendorDB[data.VendorName], oui)
	}
	for vendor := range vendorDB {
		sort.Strings(vendorDB[vendor])
	}
	devMessage("Leaving ouiToVendorDatabase()")
	return
}

func loadDatabase(fileName string) (db ouiDatabase, err error) {
	devMessage("Entering loadDatabase()")

	rawDB, rawDBErr := loadRawDatabase(config.DatabaseFile)
	if rawDBErr != nil {
		return db, fmt.Errorf("Error reading local OUI database: %s", rawDBErr)
	}
	db, err = parseRawDatabase(rawDB)
	if err != nil {
		return db, fmt.Errorf("Error parsing local OUI database: %s", err)
	}

	devMessage("Leaving loadDatabase()")
	return
}

/*
#### ##    ## ########  ##     ## ########    ##     ##    ###    ##    ## ########  ##       #### ##    ##  ######
 ##  ###   ## ##     ## ##     ##    ##       ##     ##   ## ##   ###   ## ##     ## ##        ##  ###   ## ##    ##
 ##  ####  ## ##     ## ##     ##    ##       ##     ##  ##   ##  ####  ## ##     ## ##        ##  ####  ## ##
 ##  ## ## ## ########  ##     ##    ##       ######### ##     ## ## ## ## ##     ## ##        ##  ## ## ## ##   ####
 ##  ##  #### ##        ##     ##    ##       ##     ## ######### ##  #### ##     ## ##        ##  ##  #### ##    ##
 ##  ##   ### ##        ##     ##    ##       ##     ## ##     ## ##   ### ##     ## ##        ##  ##   ### ##    ##
#### ##    ## ##         #######     ##       ##     ## ##     ## ##    ## ########  ######## #### ##    ##  ######
*/

func filterHexChars(r rune) rune {
	devMessage("Entering filterHexChars()")

	switch {
	case r >= '0' && r <= '9':
		return r
	case r >= 'A' && r <= 'F':
		return r
	case r >= 'a' && r <= 'f':
		return r
	}

	devMessage("Leaving filterHexChars()")
	return -1
}

func isValidMAC(mac string) bool {
	var reFullMAC = regexp.MustCompile(`^([0-9A-Fa-f]{2}[-:]?){5}[0-9A-Fa-f]{2}$`)
	var reOUIOnly = regexp.MustCompile(`^([0-9A-Fa-f]{2}[-:]?){2}[0-9A-Fa-f]{2}$`)

	devMessage("Entering isValidMAC()")

	if reFullMAC.Match([]byte(mac)) {
		return true
	}
	if reOUIOnly.Match([]byte(mac)) {
		return true
	}

	devMessage("Leaving isValidMAC()")
	return false
}

func extractOUI(mac string) (oui string, err error) {
	devMessage("Entering extractOUI()")

	hexOnly := strings.Map(filterHexChars, mac)

	if len(hexOnly) < 6 {
		err = fmt.Errorf("Not enough characters to extract an OUI")
		return
	}
	oui = hexOnly[:6]

	devMessage("Leaving extractOUI()")
	return
}

func normalizeMAC(mac string) (normalizedMAC string, err error) {
	var normalized []rune
	var charCount uint

	devMessage("Entering normalizeMAC()")

	mac = strings.Map(filterHexChars, mac)
	if len(mac) > 12 {
		err = fmt.Errorf("MAC too long")
		return
	}
	for len(mac) < 12 {
		mac = mac + "0"
	}

	charCount = 0
	for _, c := range mac {
		charCount++
		normalized = append(normalized, c)
		if charCount%2 == 0 {
			normalized = append(normalized, ':')
		}
	}
	normalizedMAC = strings.Trim(string(normalized), ":")

	devMessage("Leaving normalizeMAC()")
	return
}
