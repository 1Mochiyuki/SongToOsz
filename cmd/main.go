package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"slices"
	"strings"
)

const BASE = "C:\\Users\\Administrator\\AppData\\Local\\osu!\\Songs\\"

func main() {
	/* fileName := fmt.Sprintf("%s\\5551 Kawada Mami - Joint (TV Size)\\Kawada Mami - Joint (TV Size) (tom800510) [Deepsea's Streams].osu", BASE)
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	} */
	songDir, err := os.ReadDir(BASE)
	if err != nil {
		panic(err)
	}
	for _, song := range songDir {
		fmt.Printf("Creating zip for: %s\n", song.Name())
		createZip(song)
		break
	}
}

func readOsuMetadata(file *os.File) {
	fmt.Printf("File: %s\n", file.Name())
	text := bufio.NewReader(file)

	for {
		line, err := text.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("No more data")
				break
			}
			log.Fatalf("read file line error: %v", err)
			return
		}
		if strings.Contains(line, "Mode") {
			fmt.Printf("Meow: %s\n", line)
		}
	}

}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	fmt.Printf("\nBase in zip: %s\nBase path: %s", baseInZip, basePath)
	// Open the Directory
	files, err := os.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		fmt.Printf("Current file: %s\n", file.Name())
		absPath := fmt.Sprintf("%s%c%s", basePath, os.PathSeparator, file.Name())
		zipDir := fmt.Sprintf("%s%c%s", basePath, os.PathSeparator, baseInZip)
		fmt.Println("Opening zip directory")
		openedZipDir, err := os.Open(zipDir)
		if err != nil {
			panic(err)
		}
		zipDirEntries, err := openedZipDir.Readdirnames(len(files))
		if err != nil {
			panic(err)
		}
		if strings.Contains(file.Name(), ".osz") {
			fmt.Println("\nSkipping .osz file")
			continue
		}

		
		if strings.Contains(file.Name(), ".osu") {
			if slices.Contains(zipDirEntries, file.Name()) {
				fmt.Printf("Skipping: %s\n", file.Name())
				continue
			}
			fmt.Println("Found .osu file, reading metadata")

			fmt.Printf("Reading: %s\n", file.Name())

			file, err := os.Open(absPath)
			if err != nil {
				fmt.Println("open file err")
				panic(err)
			}
			readOsuMetadata(file)
		}

		path := fmt.Sprintf("%s%c%s", basePath, os.PathSeparator, file.Name())
		fmt.Printf("Found: %s\n", path)
		if !file.IsDir() {

			dat, err := os.ReadFile(path)

			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			fileInZip := fmt.Sprintf("%s%s", baseInZip, file.Name())
			f, err := w.Create(fileInZip)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("File in zip: %s\n", fileInZip)
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			newZipBase := fmt.Sprintf("%s%c", baseInZip, os.PathSeparator)
			fmt.Printf("Zip Base: %s\n", newZipBase)
			addFiles(w, newBase, newZipBase)
		}
	}
}

func songToOsz() {
	baseDir, err := os.Open(BASE)
	if err != nil {
		panic(err)
	}
	test, err := baseDir.ReadDir(1)
	if err != nil {
		panic(err)
	}

	for _, dir := range test {
		fmt.Printf("Found: %d\nName: %s", len(test), dir.Name())
		if err := createZip(dir); err != nil {
			panic(err)
		}
	}
}

func createZip(file fs.DirEntry) error {

	fmt.Println("Creating " + file.Name() + ".osz")
	archive, err := os.Create(BASE + file.Name() + ".osz")
	if err != nil {
		return err
	}
	defer archive.Close()

	writer := zip.NewWriter(archive)

	zipBasePath := fmt.Sprintf("%s%c", file.Name(), os.PathSeparator)

	addFiles(writer, BASE+file.Name(), zipBasePath)

	fmt.Println("closing zip archive")
	writer.Close()
	return nil
}
