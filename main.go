package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/godarkthrt/filenameconverter/translator"
)

const DEFAULT_MAPPING_FILENAME = "filename_mapping.csv"

type mappingData struct {
	currentFileName string
	newFileName     string
}

func NewMappingData(mdata []string) mappingData {
	return mappingData{
		currentFileName: mdata[0],
		newFileName:     mdata[1],
	}
}

// csvRow returns
func (m mappingData) csvRow() []string {
	return []string{m.currentFileName, m.newFileName}
}

func (m mappingData) csvHeader() []string {
	return []string{"CURRENT_FILE_NAME", "NEW_FILE_NAME"}
}

type MappingFileData []mappingData

func NewMappingFileData(csvData [][]string) MappingFileData {
	var mfd MappingFileData = make([]mappingData, 0)
	for i, v := range csvData {
		if i == 0 {
			// skip the header row
			continue
		}
		mdata := NewMappingData(v)
		mfd.add(mdata)
	}

	return mfd
}

func NewMappingFileDataWithJapaneseFilenames(filenames []string) MappingFileData {
	var mfd MappingFileData = make([]mappingData, 0)
	for _, fname := range filenames {
		englishFilename := convertJapaneseToEnglish(fname)
		mdata := mappingData{currentFileName: fname, newFileName: sanatizeFileName(englishFilename)}
		mfd.add(mdata)
	}

	return mfd
}

func (m *MappingFileData) add(data mappingData) {
	*m = append(*m, data)
}

func (m *MappingFileData) save(filename string) error {
	mappingFileName := filename
	if len(filename) == 0 {
		// default mapping filename
		mappingFileName = DEFAULT_MAPPING_FILENAME
	}
	csvFile, err := os.Create(mappingFileName)
	if err != nil {
		return fmt.Errorf("error while creating mapping file : %s , error  : %s", mappingFileName, err)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
	for i, v := range *m {
		if i == 0 {
			csvwriter.Write(v.csvHeader())
		}
		csvwriter.Write(v.csvRow())
	}
	csvwriter.Flush()

	return nil
}

// load mapping data from given filename (if not provided use default mapping filename)
func loadMappingDataAndRenameFiles(filename string) error {
	mappingFileName := filename
	if len(filename) == 0 {
		// default mapping filename
		mappingFileName = DEFAULT_MAPPING_FILENAME
	}

	csvfile, err := os.Open(mappingFileName)
	if err != nil {
		return fmt.Errorf("error while trying to read mapping file : %s , error : %s", mappingFileName, err)
	}

	csvReader := csv.NewReader(csvfile)
	csvdata, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("error while reading mapping data from file : %s , error : %s", mappingFileName, err)
	}

	mfd := NewMappingFileData(csvdata)

	for _, md := range mfd {
		err := os.Rename(md.currentFileName, md.newFileName)
		if err != nil {
			fmt.Printf("error while renaming file : %s , error : %s , will continue renaming other files \n", md.currentFileName, err)
		}
	}

	return nil
}

func getAllFileNameFromDirectory(dir string) ([]string, error) {
	var fileNames []string = make([]string, 0, 10)
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Error while getting file names for dir : %s , error : %s \n", dir, err)
	}
	for _, dirEntry := range dirEntries {
		// commenting since we need to now convert directories too
		// if dirEntry.IsDir() {
		// 	continue
		// }
		fileNames = append(fileNames, dirEntry.Name())
	}
	return fileNames, nil
}

func convertJapaneseToEnglish(filenameInJapanese string) string {
	englishFilename, err := translator.Translate(filenameInJapanese, "ja", "en")
	fmt.Println("english names of file ", filenameInJapanese, " is : ", englishFilename)
	if err != nil {
		return ""
	}
	return englishFilename
}

func sanatizeFileName(rawFilename string) string {
	newFileName := rawFilename

	if strings.Contains(newFileName, "/") {
		newFileName = strings.ReplaceAll(newFileName, "/", "_")
	}
	if strings.Contains(newFileName, "\\") {
		newFileName = strings.ReplaceAll(newFileName, "\\", "_")
	}
	if strings.Contains(newFileName, ":") {
		newFileName = strings.ReplaceAll(newFileName, ":", "_")
	}
	if strings.Contains(newFileName, "*") {
		newFileName = strings.ReplaceAll(newFileName, "*", "_")
	}
	if strings.Contains(newFileName, "?") {
		newFileName = strings.ReplaceAll(newFileName, "?", "_")
	}
	if strings.Contains(newFileName, "\"") {
		newFileName = strings.ReplaceAll(newFileName, "\"", "_")
	}
	if strings.Contains(newFileName, "<") {
		newFileName = strings.ReplaceAll(newFileName, "<", "_")
	}
	if strings.Contains(newFileName, ">") {
		newFileName = strings.ReplaceAll(newFileName, ">", "_")
	}
	if strings.Contains(newFileName, "|") {
		newFileName = strings.ReplaceAll(newFileName, "|", "_")
	}

	return newFileName
}

func main() {

	createMappingFile := flag.Bool("m", false, "create a mapping file containing current and suggested new file names")
	renameStep := flag.Bool("r", false, "rename file name based on mapping file created")

	flag.Parse()

	switch {
	case *createMappingFile == true:
		fmt.Printf("*** creating a mapping file : %s \n", DEFAULT_MAPPING_FILENAME)
		files, _ := getAllFileNameFromDirectory(".")
		mfd := NewMappingFileDataWithJapaneseFilenames(files)
		mfd.save("")
	case *renameStep == true:
		fmt.Printf("*** loading mapping file : %s  to rename files \n", DEFAULT_MAPPING_FILENAME)
		// TODO handle a scenario if mappingfile was not created or somehow deleted
		loadMappingDataAndRenameFiles("")
	default:
		fmt.Println("unexpected error, use -h to see provided functionality")
	}

	files, _ := getAllFileNameFromDirectory(".")
	fmt.Println(len(files), files)

}
