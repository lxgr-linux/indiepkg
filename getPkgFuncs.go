package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func writeLoadPkg(pkgName string, pkgFile string, load bool) Package {
	newFile(infoPath+pkgName+".json", pkgFile, "An error occurred while writing package information for %s", pkgName)

	if load {
		var pkg Package

		if load {
			pkg = readLoad(pkgName)
		}

		return pkg
	}

	return Package{}
}

func findPkg(pkgName string) string {
	log(1, "Finding package %s...", bolden(pkgName))
	urls := parseSources()
	var validInfos, validUrls []string

	log(1, "Checking urls length...")
	urlsLen := len(urls)

	if urlsLen <= 0 {
		log(4, "You don't have any sources defined in %s.", bolden(configPath+"sources.txt"))
		os.Exit(1)
	} else if urlsLen == 1 {
		debugLog("Only one source defined in %s. Using that source.", bolden(configPath+"sources.txt"))
		log(1, "Getting package info for %s...", bolden(pkgName))
		pkgFile, err := viewFile(urls[0]+pkgName+".json", "An error occurred while getting package information for %s", pkgName)

		if errIs404(err) {
			log(4, "Package %s not found.", bolden(pkgName))
			os.Exit(1)
		}

		errorLog(err, 4, "An error occurred while getting package information for %s", bolden(pkgName))

		return pkgFile
	}

	for _, url := range urls {
		log(1, "Parsing URL...")
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}
		pkgUrl := url + pkgName + ".json"

		log(1, "Checking %s for package info...", bolden(url))
		debugLog("URL: %s", pkgUrl)
		infoFile, err := viewFile(pkgUrl, "An error occurred while getting package information for %s", pkgName)

		if errIs404(err) {
			continue
		}

		errorLog(err, 4, "An error occurred while getting package information for %s", bolden(pkgName))

		log(0, "Found %s in %s!", bolden(pkgName), bolden(url))
		log(1, "Saving valid info & url...")
		validInfos = append(validInfos, infoFile)
		validUrls = append(validUrls, url)
	}

	log(1, "Checking valid info length...")
	lenValidInfos := len(validInfos)

	if lenValidInfos < 1 {
		log(4, "Package %s not found in any repo.", bolden(pkgName))
		os.Exit(1)
	} else if lenValidInfos == 1 {
		return validInfos[0]
	} else {
		fmt.Printf("\n")
		log(1, "Multiple packages found. Please choose one:")
		for i, url := range validUrls {
			log(1, "%d) %s", i, bolden(url+pkgName))
		}

		choice := input("0", "Number between 0 and %d or q to quit", lenValidInfos-1)

		if strings.Contains(choice, "q") {
			os.Exit(1)
		} else {
			convChoice, err := strconv.Atoi(choice)
			errorLog(err, 4, "An error occurred while converting choice to int")
			return validInfos[convChoice]
		}
	}
	return ""
}

func getPkgFromNet(pkgName string) (Package, string) {
	packageFile := findPkg(pkgName)

	pkg := loadPkg(packageFile, pkgName)

	return pkg, packageFile
}

func downloadPkg(pkgName string, load bool) Package {
	log(1, "Downloading package info for %s...", bolden(pkgName))

	return writeLoadPkg(pkgName, findPkg(pkgName), load)
}

func doDirectDownload(pkg Package, pkgName string, srcPath string) {
	pkgDispName := bolden(pkgName)

	log(1, "Making sure %s is not already downloaded...", pkgDispName)
	delPath(3, srcPath+pkg.Name, "An error occurred while deleting temporary downloaded files for %s", pkgName)

	log(1, "Getting download URL for %s", pkgDispName)
	url := getDownloadUrl(pkg)

	log(1, "Making directory for %s...", pkgDispName)
	newDir(srcPath+pkg.Name, "An error occurred while creating temporary directory for %s", pkgName)

	log(1, "Downloading file for %s from %s...", pkgDispName, bolden(url))
	nameOfFile := srcPath + pkg.Name + "/" + pkg.Name

	debugLog("Downloading and saving to %s", bolden(nameOfFile))
	downloadFileWithProg(nameOfFile, url, "An error occurred while downloading file for %s", pkgName)
}
