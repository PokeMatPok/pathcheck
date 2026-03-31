package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

var systemPath, userPath string
var systemPathEntries, userPathEntries []string

func expandPath(path string) string {
	return os.ExpandEnv(path)
}

func searchPath(needle string) ([]string, []string, error) {
	sysMatches := []string{}
	userMatches := []string{}
	for _, entry := range systemPathEntries {
		// Check if the entry is in the user PATH
		fullPath := path.Join(entry, needle+".exe")
		_, err := os.Stat(fullPath)
		if err == nil {
			sysMatches = append(sysMatches, entry)
		}
	}

	for _, entry := range userPathEntries {
		// Check if the entry is in the user PATH
		fullPath := path.Join(entry, needle+".exe")
		_, err := os.Stat(fullPath)
		if err == nil {
			userMatches = append(userMatches, entry)
		}
	}
	return sysMatches, userMatches, nil
}

func difference(a, b []string) []string {
	lookup := make(map[string]bool)
	for _, entry := range b {
		lookup[strings.ToLower(entry)] = true
	}
	result := []string{}
	for _, entry := range a {
		if !lookup[strings.ToLower(entry)] {
			result = append(result, entry)
		}
	}
	return result
}

func help() {
	println(Yellow + "Usage: pathcheck" + Red + " <command>" + Green + " [args]" + Reset)

	println(Red + "Commands:" + Reset)
	println(Red + "  audit" + Reset + "   - Audit the system and user PATH variables")
	println(Red + "  which" + Reset + "   - Show the location of a command")
	println(Red + "  cast" + Reset + "    - Run a command from the PATH")
	println(Red + "  diff" + Reset + "    - Compare system and user PATH variables")
	println(Red + "  unique" + Reset + "  - Show unique entries in PATH in machine readable format separated by newlines")
	println(Red + "  install" + Reset + " - Install path")
	println(Red + "  uninstall" + Reset + " - Uninstall path")

	println(Green + "Args:" + Reset)
	println(Green + "  <program> <passed_args>" + Reset)
	os.Exit(0)
}

func audit() {
	// find duplicates
	pathSources := make(map[string][]string)
	for _, entry := range systemPathEntries {
		pathSources[entry] = append(pathSources[entry], "System")
	}
	for _, entry := range userPathEntries {
		pathSources[entry] = append(pathSources[entry], "User")
	}

	duplicates := make([]string, 0)

	// Print duplicates
	for entry, sources := range pathSources {
		if len(sources) < 2 {
			continue
		}
		duplicates = append(duplicates, entry)
	}

	if len(duplicates) == 0 {
		println(Green + "No duplicates found in PATH" + Reset)
	} else {
		println(Yellow + strconv.Itoa(len(duplicates)) + " duplicates found in PATH:" + Reset)

		// Print duplicates
		for _, entry := range duplicates {
			println(Red + "	" + entry + Reset)
		}
	}

	//check for invalid paths
	invalidPaths := make([]string, 0)

	for _, entry := range systemPathEntries {
		if _, err := os.Stat(entry); err != nil {
			invalidPaths = append(invalidPaths, entry+" (System)")
		}
	}

	for _, entry := range userPathEntries {
		if _, err := os.Stat(entry); err != nil {
			invalidPaths = append(invalidPaths, entry+" (User)")
		}
	}

	if len(invalidPaths) == 0 {
		println(Green + "No invalid paths found in PATH" + Reset)
		return
	} else {
		println(Yellow + strconv.Itoa(len(invalidPaths)) + " invalid paths found in PATH:" + Reset)
	}

	// Print invalid paths
	for _, entry := range invalidPaths {
		println(Red + "	" + entry + Reset)
	}

	// Check for potential path hijacks (e.g. C:\Windows\System32 in user PATH)
	potentialHijacks := make([]string, 0)

	for _, entry := range userPathEntries {
		if strings.EqualFold(entry, "C:\\Windows\\System32") {
			potentialHijacks = append(potentialHijacks, entry+" (C:\\Windows\\System32 in User PATH)")
		}

		if strings.EqualFold(entry, "C:\\Windows") {
			potentialHijacks = append(potentialHijacks, entry+" (C:\\Windows in User PATH)")
		}
	}

	for _, entry := range systemPathEntries {
		if strings.EqualFold(entry, "C:\\Windows\\System32") {
			potentialHijacks = append(potentialHijacks, entry+" (C:\\Windows\\System32 in System PATH)")
		}

		if strings.EqualFold(entry, "C:\\Windows") {
			potentialHijacks = append(potentialHijacks, entry+" (C:\\Windows in System PATH)")
		}
	}

	if len(potentialHijacks) == 0 {
		println(Green + "No potential path hijacks found in user PATH" + Reset)
	} else {
		println(Yellow + strconv.Itoa(len(potentialHijacks)) + " potential path hijacks found in user PATH:" + Reset)
	}

	// Print potential hijacks
	for _, entry := range potentialHijacks {
		println(Red + "	" + entry + Reset)
	}

	// print summary
	println(Blue + "Summary:" + Reset)
	println(Blue + "  Total entries in System PATH: " + strconv.Itoa(len(systemPathEntries)) + Reset)
	println(Blue + "  Total entries in User PATH: " + strconv.Itoa(len(userPathEntries)) + Reset)
	println(Blue + "  Total duplicates: " + strconv.Itoa(len(duplicates)) + Reset)
	println(Blue + "  Total invalid paths: " + strconv.Itoa(len(invalidPaths)) + Reset)
	println(Blue + "  Total potential hijacks: " + strconv.Itoa(len(potentialHijacks)) + Reset)
	println("")

	// print recommendations
	println(Blue + "Recommendations:" + Reset)
	println(Blue + "  - Remove duplicate entries from PATH" + Reset)
	println(Blue + "  - Remove invalid paths from PATH" + Reset)
	println(Blue + "  - Avoid adding critical system directories (e.g. C:\\Windows\\System32) to user PATH" + Reset)
}

func which() {

	sysMatches, userMatches, _ := searchPath(os.Args[2])

	matchesCount := len(sysMatches) + len(userMatches)

	if matchesCount > 0 {
		println(Green + strconv.Itoa(matchesCount) + " matches found" + Reset)
		if len(sysMatches) > 0 {
			println(Yellow + "System PATH matches:" + Reset)
			for _, match := range sysMatches {
				println(White + "  " + match + Reset + Magenta + " → " + strings.Join(strings.Split(path.Join(match, os.Args[2]+".exe"), "/"), "\\") + Reset)
			}
		}
		if len(userMatches) > 0 {
			println(Yellow + "User PATH matches:" + Reset)
			for _, match := range userMatches {
				println(White + "  " + match + Reset + Magenta + " → " + strings.Join(strings.Split(path.Join(match, os.Args[2]+".exe"), "/"), "\\") + Reset)
			}
		}
	} else {
		println(Red + "No matches found in PATH" + Reset)
	}
}

func cast() {
	sysMatches, userMatches, _ := searchPath(os.Args[2])

	matchesCount := len(sysMatches) + len(userMatches)

	if matchesCount > 1 {
		println(Yellow + "multiple matches found")
		println("please select your command:" + Reset)
		if len(sysMatches) > 0 {
			println(Yellow + "System PATH matches:" + Reset)
			for i, match := range sysMatches {
				println("  " + strconv.Itoa(i+1) + ". " + match)
			}
		}
		if len(userMatches) > 0 {
			println(Yellow + "User PATH matches:" + Reset)
			for i, match := range userMatches {
				println("  " + strconv.Itoa(i+1+len(sysMatches)) + ". " + match)
			}
		}

		var choice int
		print("Enter the number of the command to run: ")
		_, err := fmt.Scan(&choice)
		if err != nil {
			println(Red + "Invalid input" + Reset)
			os.Exit(1)
		}
		if choice < 1 || choice > matchesCount {
			println(Red + "Invalid choice" + Reset)
			os.Exit(1)
		}

		var match string
		if choice <= len(sysMatches) {
			match = sysMatches[choice-1]
		} else {
			match = userMatches[choice-1-len(sysMatches)]
		}

		match = strings.TrimSpace(match) + "\\" + os.Args[2] + ".exe"
		println(Yellow + "running " + match + "..." + Reset)

		proc := exec.Command(match, os.Args[3:]...)
		proc.Stdout = os.Stdout
		proc.Stderr = os.Stderr
		proc.Stdin = os.Stdin
		proc.Run()
	} else {
		if matchesCount == 1 {
			match := strings.TrimSpace(strings.Join(sysMatches, "") + strings.Join(userMatches, ""))
			println(Yellow + "running " + match + "..." + Reset)

			proc := exec.Command(match, os.Args[3:]...)
			proc.Stdout = os.Stdout
			proc.Stderr = os.Stderr
			proc.Stdin = os.Stdin
			proc.Run()
		}
	}
}

func diff() {
	shellEntries := strings.Split(os.Getenv("PATH"), ";")
	registryEntries := []string{}
	for _, entry := range append(systemPathEntries, userPathEntries...) {
		registryEntries = append(registryEntries, os.ExpandEnv(entry))
	}

	missingFromShell := difference(registryEntries, shellEntries)
	missingFromRegistry := difference(shellEntries, registryEntries)

	if len(missingFromShell) == 0 {
		println(Green + "Shell PATH is up to date" + Reset)
	} else {
		println(Yellow + strconv.Itoa(len(missingFromShell)) + " entries in registry but missing from shell (stale):" + Reset)
		for _, entry := range missingFromShell {
			println(Red + "  " + entry + Reset)
		}
	}

	if len(missingFromRegistry) > 0 {
		println(Yellow + strconv.Itoa(len(missingFromRegistry)) + " entries in shell but removed from registry:" + Reset)
		for _, entry := range missingFromRegistry {
			println(Red + "  " + entry + Reset)
		}
	}
}

func unique() {
	uniqueMatches := make(map[string]bool)

	for _, entry := range systemPathEntries {
		uniqueMatches[entry] = true
	}
	for _, entry := range userPathEntries {
		uniqueMatches[entry] = true
	}

	for entry := range uniqueMatches {
		println(entry)
	}
}

func install() {
	println("Installing path...")
}

func uninstall() {
	println("Uninstalling path...")
}

func init() {
	key, _ := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE)
	systemPath, _, _ = key.GetStringValue("PATH")
	systemPath = expandPath(systemPath)
	systemPathEntries = strings.Split(systemPath, ";")

	key2, _ := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	userPath, _, _ = key2.GetStringValue("PATH")
	userPath = expandPath(userPath)
	userPathEntries = strings.Split(userPath, ";")

	/*

		verb := windows.StringToUTF16Ptr("runas")
		exe := windows.StringToUTF16Ptr(os.Args[0])
		params := windows.StringToUTF16Ptr(strings.Join(os.Args[1:], " "))

		err := windows.ShellExecute(0, verb, exe, params, nil, windows.SW_NORMAL)
		if err != nil {
			if err.(windows.Errno) == 1223 {
				println(Red + "Error: This program must be run as administrator" + Reset)
				os.Exit(1)
			} else {
				println(Red + "Error: " + err.Error() + Reset)
				os.Exit(1)
			}
		}

	*/
}

func main() {

	_, no_color := os.LookupEnv("NO_COLOR")

	if no_color {
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Magenta = ""
		Cyan = ""
		Gray = ""
		White = ""
		Reset = ""
	}

	if len(os.Args) < 2 {
		help()
	}

	switch os.Args[1] {
	case "audit":
		audit()
	case "which":
		which()
	case "cast":
		cast()
	case "diff":
		diff()
	case "unique":
		unique()
	case "install":
		install()
	case "uninstall":
		uninstall()
	default:
		help()
	}
}
