package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/andlabs/ui"
)

func startupWindow() {
	err := ui.Main(func() {
		box := ui.NewVerticalBox()
		status := ui.NewLabel("Starting OCM")
		box.Append(status, false)
		window := ui.NewWindow("OCM 2", 100, 50, false)
		window.SetChild(box)
		window.SetMargined(true)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})

		window.Show()

		status.SetText("Checking for miners dir")
		if _, err := os.Stat("./miners"); os.IsNotExist(err) {
			err := os.Mkdir("./miners", 0755)
			if err != nil {
				panic(err)
			}
		}

		gpuVendor := GetGPU()
		switch {
		case strings.Contains(gpuVendor, "Radeon"):
			status.SetText("Downloading AMD miner")
			err := DownloadFile("https://github.com/CryptoGraphics/lyclMiner/releases/download/untagged-95777e4326ae4e5ccdb5/lyclMiner015.zip", "./miners/AMD.zip")
			if err != nil {
				panic(err)
			}

			err = UnzipFile("./miners/AMD.zip", "./miners/AMD")
			if err != nil {
				panic(err)
			}

		case strings.Contains(gpuVendor, "NVIDIA"):
			status.SetText("Downloading NVIDIA miner")
			err := DownloadFile("https://vtconline.org/downloads/ccminer.zip", "./miners/NVIDIA.zip")
			if err != nil {
				panic(err)
			}

			err = UnzipFile("./miners/NVIDIA.zip", "./miners/NVIDIA")
			if err != nil {
				panic(err)
			}

		default:
			panic("Neither AMD or nVidia GPU found")
		}

		status.SetText(gpuVendor)
	})

	if err != nil {
		panic(err)
	}
}

func main() {
	if runtime.GOOS != "windows" {
		panic("Only Windows is supported at present")
	}

	startupWindow()
}
