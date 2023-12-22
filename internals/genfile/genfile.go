package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	directory := "/Users/pp/downloads/single pose/y2k pose01"

	// Read the specified directory
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Open output file
	outFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".ycd") {
			baseName := strings.TrimSuffix(f.Name(), ".ycd")
			strippedName := strings.TrimPrefix(baseName, "789@")                     // remove the prefix
			displayName := strings.Title(strings.ReplaceAll(strippedName, "@", " ")) // Capitalize and replace "@" with space for display

			// Construct the desired string
			content := fmt.Sprintf("[\"%s\"] = {\"789@%s\", \"%s\", \"789Store %s\", AnimationOptions =\n   {\n       EmoteLoop = true,\n       EmoteMoving = true,\n   }},\n", strippedName, strippedName, strippedName, displayName)
			// Write to file
			outFile.WriteString(content)
		}
	}

	// Move output file to the specified directory
	err = os.Rename("output.txt", filepath.Join(directory, "output.txt"))
	if err != nil {
		fmt.Println("Error moving the output file:", err)
		return
	}

	fmt.Println("File generated successfully in", directory)
}
