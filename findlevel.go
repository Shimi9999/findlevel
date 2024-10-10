package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Shimi9999/gobms"
)

func main() {
	flag.Parse()

	if len(flag.Args()) >= 2 {
		fmt.Println("Usage: findlevel [dirpath]")
		os.Exit(1)
	}

	var path string
	if len(flag.Args()) == 0 {
		path = "./"
	} else {
		path = flag.Arg(0)
	}
	fInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("Path is wrong:", err.Error())
		os.Exit(1)
	}
	if !fInfo.IsDir() {
		fmt.Println("The entered path is not a directory")
		os.Exit(1)
	}

	bmsdirs := make([]gobms.BmsDirectory, 0)
	err = gobms.FindBmsInDirectory(path, &bmsdirs)
	if err != nil {
		fmt.Println("FindBmsInDirectory:", err.Error())
		os.Exit(1)
	}

	if len(bmsdirs) == 0 {
		fmt.Println("No BMS files")
		os.Exit(1)
	}

	noDifText := ""
	onlyOneText := ""
	okcount := 0
	ngcount := 0
	for _, dir := range bmsdirs {
		for _, bmsfile := range dir.BmsDataSet {
			difficulty, _ := strconv.Atoi(bmsfile.Difficulty)
			if difficulty >= 1 && difficulty <= 5 {
				okcount++
			} else {
				log := fmt.Sprintf("%s, #DIFFICULTY = %s\n", bmsfile.Path, bmsfile.Difficulty)
				if len(dir.BmsDataSet) == 1 {
					onlyOneText += "Difficulty is missing, but this has only one chart: " + log
				} else {
					noDifText += "Difficulty is missing: " + log
				}
				ngcount++
			}
		}
	}
	fmt.Printf("%s", noDifText)
	fmt.Printf("%s", onlyOneText)
	fmt.Printf("Difficulties OK: %d, NG: %d\n", okcount, ngcount)

	err = makeCsv(bmsdirs, path)
	if err != nil {
		fmt.Println("makeCsv:", err.Error())
		os.Exit(1)
	}
}

func makeCsv(bmsdirs []gobms.BmsDirectory, rootDirPath string) error {
	records := [][]string{}

	records = append(records, []string{"directory_path", "song_name", "for_aviutl", "level_text"})

	for _, bmsdir := range bmsdirs {
		baseDirPath := strings.Replace(filepath.Clean(bmsdir.Path), filepath.Clean(rootDirPath)+"\\", "", 1)
		records = append(records, []string{baseDirPath, bmsdir.Name, makeLevelTextForAviutl(bmsdir), makeSimpleLevelText(bmsdir)})
	}

	csvbuf := new(bytes.Buffer)
	w := csv.NewWriter(csvbuf)
	if err := w.WriteAll(records); err != nil {
		return fmt.Errorf("CSV text WriteAll: %w", err)
	}

	outputFilename := "findlevel_output.csv"
	file, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("CSV file Create: %w", err)
	}
	defer file.Close()
	_, err = file.Write(csvbuf.Bytes())
	if err != nil {
		return fmt.Errorf("CSV file Write: %w", err)
	}
	fmt.Println("Output file created:", outputFilename)

	return nil
}

func makeLevelTextForAviutl(bmsdir gobms.BmsDirectory) string {
	keys, alllevel := getAllLevel(bmsdir.BmsDataSet)

	difcolor := []string{"b0b0b0", "24fb5d", "00ccff", "f6962b", "ff3014", "a049ff"}
	indentset := []int{64, 76}
	text := ""
	for keyindex, keylevel := range alllevel {
		keytext := ""
		starstr := "☆"
		if keys[keyindex] == 9 {
			starstr = "Lv"
		}
		indent := indentset[0]
		if keys[keyindex] >= 10 {
			indent = indentset[1]
		}
		indentstr := "<p" + strconv.Itoa(indent) + ",+0>"
		levelcount := 0
		newlineflag := false
		for i := 0; i < 6; i++ {
			difficulty := (i + 1) % 6 // 0(unknown)を最後に
			leveltext := ""
			for index, level := range keylevel[difficulty] {
				if newlineflag {
					if index == 0 {
						keytext += "\n" + indentstr
					} else {
						leveltext += "\n" + indentstr
					}
					newlineflag = false
				}
				leveltext += starstr + strconv.Itoa(level) + " "
				levelcount++
				if levelcount%3 == 0 {
					newlineflag = true
				}
			}
			if leveltext != "" {
				keytext += "<#" + difcolor[difficulty] + ">" + leveltext
			}
		}
		if keytext != "" {
			keytext = strconv.Itoa(keys[keyindex]) + "KEY " + indentstr + keytext + "<#>\n"
		}
		text += keytext
	}
	text = strings.TrimRight(text, "\n")

	if strings.Count(text, "\n") >= 10 {
		fmt.Println("This level text is too long:", bmsdir.Name)
	}

	return text
}

func makeSimpleLevelText(bmsdir gobms.BmsDirectory) string {
	keys, alllevel := getAllLevel(bmsdir.BmsDataSet)

	difficultyChar := []string{"U", "B", "N", "H", "A", "I"}
	text := ""
	for keyindex, keylevel := range alllevel {
		leveltexts := ""
		for i := 0; i < 6; i++ {
			difficulty := (i + 1) % 6 // 0(unknown)を最後に
			for _, level := range keylevel[difficulty] {
				leveltexts += difficultyChar[difficulty] + strconv.Itoa(level) + " "
			}
		}
		if leveltexts != "" {
			text += strconv.Itoa(keys[keyindex]) + "K: " + leveltexts + "\n"
		}
	}
	text = strings.TrimRight(text, "\n")

	return text
}

func getAllLevel(bmsDataSet []gobms.BmsData) (keys [7]int, alllevel [7][6][]int) {
	keys = [7]int{5, 10, 7, 14, 9, 24, 48}
	// unknown, beginner, normal, hyper, another, insane
	alllevel = [7][6][]int{}
	for _, bmsfile := range bmsDataSet {
		var index int
		for i, k := range keys {
			if k == bmsfile.Keymode {
				index = i
			}
		}
		var d int
		if bmsfile.Difficulty == "" {
			d = 0
		} else {
			d, _ = strconv.Atoi(bmsfile.Difficulty) // not integer → 0
			if d < 0 || d >= 6 {
				d = 0
			}
		}
		level, _ := strconv.Atoi(bmsfile.Playlevel) // not integer → 0
		alllevel[index][d] = append(alllevel[index][d], level)
	}

	for index, keybms := range alllevel {
		for difficulty := range keybms {
			sort.Ints(alllevel[index][difficulty])
		}
	}

	return keys, alllevel
}
