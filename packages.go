package main

import (
	"os"
	"strings"
)

var home string = os.Getenv("HOME")

var srcPath string = home + "/.local/share/indiepkg/package_src/"
var installedPath string = home + "/.local/share/indiepkg/installed_packages/"

type Package struct {
	Name         string
	Author       string
	Description  string
	Url          string
	Install      []string
	Uninstall    []string
	upgrade      []string
	Config_paths []string
}

var environmentVariables = map[string]string{
	"PATH": home + "/.local",
}

func installPackage(pkgName string) {
	pkgInfoPath := installedPath + pkgName + ".json"
	url := "https://raw.githubusercontent.com/talwat/indiepkg/main/packages/" + pkgName + ".json"
	var err error

	log(1, "Making required directories...")
	newDir(srcPath)       //nolint:errcheck
	newDir(installedPath) //nolint:errcheck

	log(1, "Downloading package info...")
	log(1, "URL: %s", url)
	err = downloadFile(pkgInfoPath, url)
	errorLog(err, 4, "An error occurred while getting package information for %s.", pkgName)

	log(1, "Reading package info...")
	pkgFile, err := readFile(pkgInfoPath)
	errorLog(err, 4, "An error occurred while reading package information for %s.", pkgName)

	log(1, "Loading package info...")
	pkg, err := loadPackage(pkgFile)
	errorLog(err, 4, "An error occurred while loading package information for %s.", pkgName)

	log(1, "Cloning source code...")
	runCommand(srcPath, "git", "clone", pkg.Url)

	log(1, "Running install commands...")
	for _, command := range pkg.Install {
		runCommand(srcPath+pkg.Name, strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)
	}

	log(0, "Installed %s successfully!", pkgName)
}

func uninstallPackage(pkgName string) {
	pkgInfoPath := installedPath + pkgName + ".json"
	var err error

	installed, err := pathExists(pkgInfoPath)
	errorLog(err, 4, "An error occurred while checking if package %s exists.", pkgName)
	if !installed {
		log(4, "%s is not installed, so it can't be uninstalled.", pkgName)
		os.Exit(1)
	}

	log(1, "Reading package...")
	pkgFile, err := readFile(pkgInfoPath)
	errorLog(err, 4, "An error occurred while reading package %s.", pkgName)

	log(1, "Loading package info...")
	pkg, err := loadPackage(pkgFile)
	errorLog(err, 4, "An error occurred while loading package information for %s.", pkgName)

	log(1, "Running uninstall commands...")
	for _, command := range pkg.Uninstall {
		runCommand(srcPath+pkg.Name, strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)
	}

	log(1, "Deleting source files for %s...", pkgName)
	err = delDir(srcPath + pkgName)
	errorLog(err, 4, "An error occurred while deleting source files for %s.", pkgName)

	log(1, "Deleting info file for %s...", pkgName)
	err = delFile(pkgInfoPath)
	errorLog(err, 4, "An error occurred while deleting info file for package %s.", pkgName)

	log(0, "Successfully uninstalled %s.", pkgName)
}

func infoPackage(pkgName string) {
	packageFile, err := viewFile("https://raw.githubusercontent.com/talwat/indiepkg/main/packages/" + pkgName + ".json")
	errorLog(err, 4, "An error occurred while getting package info for %s.", pkgName)

	pkgInfo, err := loadPackage(packageFile)
	errorLog(err, 4, "An error occurred while loading package information for %s.", pkgName)

	log(1, "Name: %s", pkgInfo.Name)
	log(1, "Author: %s", pkgInfo.Author)
	log(1, "Description: %s", pkgInfo.Description)
	log(1, "Git URL: %s", pkgInfo.Url)
}
