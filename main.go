package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
)

const (
    colorRed   = "\033[31m"
    colorGreen = "\033[32m"
    colorBlue  = "\033[34m"
    colorReset = "\033[0m"
)

var (
    wipeCache   bool
    wipeDalvik  bool
    showHelp    bool
    showVersion bool
)

func init() {
    // Options for wiping
    flag.BoolVar(&wipeCache, "c", false, "\tWipe cache only")
    flag.BoolVar(&wipeDalvik, "d", false, "\tWipe dalvik only")
    // Options for help and versio
    flag.BoolVar(&showHelp, "h", false, "\tShow help message")
    flag.BoolVar(&showVersion, "v", false, "\tShow version information")
}

func twrpwipe() {
    var scriptContent string

    if wipeCache && wipeDalvik {
        scriptContent = "wipe cache\nwipe dalvik\nreboot system"
    } else if wipeCache {
        scriptContent = "wipe cache\nreboot system"
    } else if wipeDalvik {
        scriptContent = "wipe dalvik\nreboot system"
    } else {
        // Default behavior if neither cache nor dalvik is selected
        scriptContent = "wipe cache\nwipe dalvik\nreboot system"
    }

    if _, err := os.Stat("/cache/recovery/openrecoveryscript"); err == nil {
        if err := os.Remove("/cache/recovery/openrecoveryscript"); err != nil {
            printError("Error removing openrecoveryscript file:\n", err)
            return
        }
    }
    err := ioutil.WriteFile("/cache/recovery/openrecoveryscript", []byte(scriptContent), 0644)
    if err != nil {
        printError("Error writing openrecoveryscript file:\n", err)
        return
    }
	cmd := exec.Command("su", "-c", "reboot recovery")
	err = cmd.Run()
    if err != nil {
        printError("Error running reboot recovery command:\n", err)
        return
    }
}

func printError(message string, err error) {
    fmt.Printf("%s%s%s%s\n", colorRed, message, err, colorReset)
}

func main() {
    flag.Parse()

    if showHelp {
        fmt.Printf("%stwrpwipe%s\nreboot to recovery, wipe partitions, and reboot back\n(default wipe cache and dalvik)\n\n", colorGreen, colorReset)
        flag.PrintDefaults()
        fmt.Println("\n")
        return
    }

    if showVersion {
    	fmt.Printf("%stwrpwipe%s v1.2.0\nby PhateValleyman\nJonas.Ned@outlook.com\n\n", colorGreen, colorReset)
        return
    }

    twrpwipe()
}
