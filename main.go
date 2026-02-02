package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "strings"
)

const (
    colorRed   = "\033[31m"
    colorGreen = "\033[32m"
    colorBlue  = "\033[34m"
    colorReset = "\033[0m"
)

var (
    // Flag: wipe cache partition
    wipeCache bool

    // Flag: wipe dalvik/art cache
    wipeDalvik bool

    // Flag: show help
    showHelp bool

    // Flag: show version
    showVersion bool
)

func init() {
    // Register flag for cache wipe
    flag.BoolVar(&wipeCache, "c", false, "Wipe cache only")

    // Register flag for dalvik wipe
    flag.BoolVar(&wipeDalvik, "d", false, "Wipe dalvik only")

    // Register flag for help
    flag.BoolVar(&showHelp, "h", false, "Show help message")

    // Register flag for version
    flag.BoolVar(&showVersion, "v", false, "Show version information")
}

func isRecoveryMode() bool {
    // Check for recovery binary path (TWRP specific)
    if _, err := os.Stat("/sbin/recovery"); err == nil {
        return true
    }

    // Execute getprop to read boot mode
    cmd := exec.Command("getprop", "ro.bootmode")
    output, err := cmd.Output()
    if err != nil {
        return false
    }

    // Check if bootmode equals recovery
    return strings.TrimSpace(string(output)) == "recovery"
}

func buildRecoveryScript() string {
    // Build OpenRecoveryScript based on flags
    if wipeCache && wipeDalvik {
        return "wipe cache\nwipe dalvik\nreboot system\n"
    }

    if wipeCache {
        return "wipe cache\nreboot system\n"
    }

    if wipeDalvik {
        return "wipe dalvik\nreboot system\n"
    }

    // Default behavior: wipe both
    return "wipe cache\nwipe dalvik\nreboot system\n"
}

func writeRecoveryScript(script string) error {
    // Ensure old script is removed
    _ = os.Remove("/cache/recovery/openrecoveryscript")

    // Write OpenRecoveryScript file
    return ioutil.WriteFile(
        "/cache/recovery/openrecoveryscript",
        []byte(script),
        0644,
    )
}

func rebootToRecovery() error {
    // Execute reboot to recovery using root
    cmd := exec.Command("su", "-c", "reboot recovery")
    return cmd.Run()
}

func runInRecovery() error {
    // Execute TWRP wipe commands directly
    var commands []string

    // Append cache wipe command
    if wipeCache || (!wipeCache && !wipeDalvik) {
        commands = append(commands, "wipe cache")
    }

    // Append dalvik wipe command
    if wipeDalvik || (!wipeCache && !wipeDalvik) {
        commands = append(commands, "wipe dalvik")
    }

    // Append reboot command
    commands = append(commands, "reboot system")

    // Join commands with newline
    script := strings.Join(commands, "\n")

    // Execute commands via recovery shell
    cmd := exec.Command("/sbin/sh", "-c", script)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    return cmd.Run()
}

func printError(message string, err error) {
    // Print error message in red color
    fmt.Printf("%s%s%s%s\n", colorRed, message, err, colorReset)
}

func main() {
    // Parse CLI flags
    flag.Parse()

    // Show help and exit
    if showHelp {
        fmt.Printf(
            "%stwrpwipe%s\nAutomatically detects system or recovery mode\n"+
                "Normal mode: writes script and reboots to recovery\n"+
                "Recovery mode: wipes cache/dalvik and reboots\n\n",
            colorGreen,
            colorReset,
        )
        flag.PrintDefaults()
        return
    }

    // Show version and exit
    if showVersion {
        fmt.Printf(
            "%stwrpwipe%s v1.3.0\nby PhateValleyman\nJonas.Ned@outlook.com\n\n",
            colorGreen,
            colorReset,
        )
        return
    }

    // Detect if device is running in recovery mode
    if isRecoveryMode() {
        // Run wipe directly inside recovery
        if err := runInRecovery(); err != nil {
            printError("Error running wipe in recovery:\n", err)
        }
        return
    }

    // Build OpenRecoveryScript content
    script := buildRecoveryScript()

    // Write OpenRecoveryScript to recovery cache
    if err := writeRecoveryScript(script); err != nil {
        printError("Error writing OpenRecoveryScript:\n", err)
        return
    }

    // Reboot device into recovery
    if err := rebootToRecovery(); err != nil {
        printError("Error rebooting to recovery:\n", err)
        return
    }
}
