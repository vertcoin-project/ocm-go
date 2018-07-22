package main

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/andlabs/ui"
)

func startupWindow() {
	var GPUType string
	err := ui.Main(func() {
		box := ui.NewVerticalBox()
		status := ui.NewLabel("Starting OCM")
		box.Append(status, false)
		window := ui.NewWindow("OCM-go", 200, 50, false)
		window.SetChild(box)
		window.SetMargined(true)

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
			if _, err := os.Stat("./miners/AMD"); os.IsNotExist(err) {
				status.SetText("Downloading AMD miner")
				err := DownloadFile("https://github.com/CryptoGraphics/lyclMiner/releases/download/untagged-95777e4326ae4e5ccdb5/lyclMiner015.zip", "./miners/AMD.zip")
				if err != nil {
					panic(err)
				}

				err = UnzipFile("./miners/AMD.zip", "./miners/AMD")
				if err != nil {
					panic(err)
				}

				err = os.Remove("./miners/AMD.zip")
				if err != nil {
					panic(err)
				}
			}

			err := os.Chdir("./miners/AMD/lyclMiner015")
			if err != nil {
				panic(err)
			}

			if _, err := os.Stat("lycl.conf"); err == nil {
				err := os.Remove("lycl.conf")
				if err != nil {
					panic(err)
				}
			}

			cmd := exec.Command("lyclMiner.exe", "-g", "lycl.conf")
			err = cmd.Run()
			if err != nil {
				panic(err)
			}

			err = ReplaceInFile("lycl.conf", "stratum+tcp://example.com:port", "stratum+tcp://p2proxy.vertcoin.org:9171")
			if err != nil {
				panic(err)
			}

			GPUType = "Radeon"
		case strings.Contains(gpuVendor, "NVIDIA"):
			if _, err := os.Stat("./miners/NVIDIA"); os.IsNotExist(err) {
				status.SetText("Downloading NVIDIA miner")
				err := DownloadFile("https://vtconline.org/downloads/ccminer.zip", "./miners/NVIDIA.zip")
				if err != nil {
					panic(err)
				}

				err = UnzipFile("./miners/NVIDIA.zip", "./miners/NVIDIA")
				if err != nil {
					panic(err)
				}

				err = os.Remove("./miners/NVIDIA.zip")
				if err != nil {
					panic(err)
				}
			}

			GPUType = "NVIDIA"
		default:
			panic("Neither AMD or nVidia GPU found")
		}

		window.Hide()

		go mainWindow(GPUType)
	})

	if err != nil {
		panic(err)
	}
}

func startAMD(address string) (*exec.Cmd, io.ReadCloser, error) {
	err := ReplaceInFile("lycl.conf", "user", address)
	if err != nil {
		return nil, nil, err
	}

	cmd := exec.Command("cmd", "/K", "lyclMiner.exe", "lycl.conf")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	return cmd, stdout, nil
}

func startNVIDIA(address string) (*exec.Cmd, io.ReadCloser, error) {
	cmd := exec.Command("cmd", "/K", "./miners/NVIDIA/ccminer-x64.exe", "-a", "lyra2v2", "-o", "stratum+tcp://p2proxy.vertcoin.org:9171", "-u", address, "-p", "x")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	return cmd, stdout, nil
}

func mainWindow(GPUType string) {
	ui.QueueMain(func() {
		var cmd *exec.Cmd
		input := ui.NewEntry()
		button := ui.NewButton("Mine")
		status := ui.NewLabel("")
		box := ui.NewVerticalBox()
		box.Append(ui.NewLabel("VTC address:"), false)
		box.Append(input, false)
		box.Append(button, false)
		box.Append(status, false)
		mWindow := ui.NewWindow("OCM-go", 100, 50, false)
		mWindow.SetMargined(true)
		mWindow.SetChild(box)

		button.OnClicked(func(*ui.Button) {
			var err error
			if GPUType == "NVIDIA" {
				cmd, _, err = startNVIDIA(input.Text())
			} else {
				cmd, _, err = startAMD(input.Text())
			}

			if err != nil {
				panic(err)
			}

			status.SetText("Started mining")
		})

		mWindow.OnClosing(func(*ui.Window) bool {
			if cmd != nil {
				cmd.Process.Kill()
			}

			ui.Quit()
			return true
		})
		mWindow.Show()
	})
}

func main() {
	if runtime.GOOS != "windows" {
		panic("Only Windows is supported at present")
	}

	startupWindow()
}
