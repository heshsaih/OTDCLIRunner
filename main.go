package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const PATH = "D:\\otd\\"
const OSU_CONFIG = "config_osu.json"
const DRAWING_CONFIG = "config_drawing.json"

func runDaemon() *exec.Cmd {
	command := exec.Command(PATH + "OpenTabletDriver.Daemon.exe")
	err := command.Start()
	if err != nil {
		panic(err)
	}
	proccessPid := command.Process.Pid
	log.Println("Daemon is running with pid: ", proccessPid)
	return command
}

func determineConfigAndLoad() string {
	args := os.Args[1:]
	var loadedConfig string
	//App will run with osu config to load as default
	if len(args) == 0 || args[0] == "osu" {
		//I really don't know what's the issue, but for some reason daemon crahses when
		//I load osu config first and then I want to load drawing config,
		//but when I load drawing config first, everything's working
		loadConfig(DRAWING_CONFIG)
		loadConfig(OSU_CONFIG)
		loadedConfig = OSU_CONFIG
	} else if args[0] == "drawing" {
		loadConfig(DRAWING_CONFIG)
		loadedConfig = DRAWING_CONFIG
	} else {
		panic("Wrong config specified (osu | drawing only)")
	}
	return loadedConfig
}

func loadConfig(config string) {
	out, err := exec.Command(PATH+"OpenTabletDriver.Console.exe", "loadsettings", PATH+config).Output()
	if err != nil {
		panic(err)
	}
	log.Println(out)
	log.Println("Config \"" + config + "\" has been loaded successfully")

}

func printOptions() {
	fmt.Println("Available options:")
	fmt.Println("~ h/H -> Display this menu")
	fmt.Println("~ o/O -> Load osu config")
	fmt.Println("~ d/D -> Load drawing config")
	fmt.Println("~ q/Q -> Exit")
}

func main() {
	log.Println("Program is running!")

	command := runDaemon()
	defer command.Process.Kill()

	//Daemon loads for a bit, so that's why there is a sleep
	time.Sleep(1 * time.Second)

	loadedConfig := determineConfigAndLoad()
	printOptions()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.ToLower(strings.TrimSpace(text))
		switch text {
		case "h":
			printOptions()
		case "o":
			if loadedConfig == OSU_CONFIG {
				fmt.Println("Config already loaded")
			} else {
				loadConfig(OSU_CONFIG)
				loadedConfig = OSU_CONFIG
            }
		case "d":
			if loadedConfig == DRAWING_CONFIG {
				fmt.Println("Config already loaded")
			} else {
				loadConfig(DRAWING_CONFIG)
                loadedConfig = DRAWING_CONFIG
			}
		case "q":
			return
		default:
			fmt.Println("Wrong option (use \"h\" for help)")
		}
	}
}
