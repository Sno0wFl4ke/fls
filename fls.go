package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	_ "time"
)

// Converts bytes to a human-readable format (KB, MB, GB, etc.)
func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case size >= TB:
		return color.HiRedString(fmt.Sprintf("%.2fTB", float64(size)/float64(TB)))
	case size >= GB:
		return color.HiMagentaString(fmt.Sprintf("%.2fGB", float64(size)/float64(GB)))
	case size >= MB:
		return color.HiBlueString(fmt.Sprintf("%.2fMB", float64(size)/float64(MB)))
	case size >= KB:
		return color.HiGreenString(fmt.Sprintf("%.2fKB", float64(size)/float64(KB)))
	default:
		return fmt.Sprintf("%dB", size)
	}
}

func main() {
	// Define flags
	deep := flag.Bool("d", false, "List files recursively")
	allInfo := flag.Bool("a", false, "Show detailed file info (size, permissions, modified date)")

	flag.Parse()

	// Run the appropriate function
	if *deep {
		fancyLsDeep(*allInfo)
	} else if *allInfo {
		fancyLs(*allInfo)
	} else {
		simpleLs()
	}
}

func simpleLs() {
	dir, err := os.ReadDir("./")
	if err != nil {
		return
	}

	fmt.Println(color.HiWhiteString("┤ " + filepath.Dir(os.Args[0])))
	for _, entry := range dir {
		if entry.IsDir() {
			fmt.Println(color.HiBlueString("├ 📁" + entry.Name()))
		} else if entry.Type().IsRegular() {
			fmt.Println(color.HiMagentaString("│ 📝" + entry.Name()))
		} else {
			fmt.Println(color.HiRedString("│ " + entry.Type().String()))
		}
	}
}

func fancyLs(showInfo bool) {
	dir, err := os.ReadDir("./")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	fmt.Println(color.HiWhiteString("┤ " + filepath.Dir(os.Args[0])))

	for _, entry := range dir {
		printFileInfo(entry, showInfo, "│ ")
	}
}

func fancyLsDeep(showInfo bool) {
	dir, err := os.ReadDir("./")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	fmt.Println(color.HiWhiteString("┤ " + filepath.Dir(os.Args[0])))

	for _, entry := range dir {
		printFileInfo(entry, showInfo, "│ ")

		// If it's a directory, list its contents
		if entry.IsDir() {
			subDir, _ := os.ReadDir(entry.Name())
			for _, subEntry := range subDir {
				printFileInfo(subEntry, showInfo, "│ ├ ")
			}
		}
	}
}

// Function to print file info based on the -a flag
func printFileInfo(entry os.DirEntry, showInfo bool, prefix string) {
	if showInfo {
		info, err := entry.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			return
		}

		modTime := info.ModTime().Format("2006-01-02 15:04:05") // Format date
		size := formatSize(info.Size())                         // Convert size
		perms := info.Mode()

		if entry.IsDir() {
			fmt.Printf("%s📁 %s\t[%s]  %s  %s\n", color.HiBlueString(prefix), entry.Name(), perms, size, modTime)
		} else {
			fmt.Printf("%s📝 %s\t[%s]  %s  %s\n", color.HiMagentaString(prefix), entry.Name(), perms, size, modTime)
		}
	} else {
		// Simple output (no extra info)
		if entry.IsDir() {
			fmt.Println(color.HiBlueString(prefix + "📁 " + entry.Name()))
		} else {
			fmt.Println(color.HiMagentaString(prefix + "📝 " + entry.Name()))
		}
	}
}
