package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/mitchellh/go-homedir"
	"os/exec"
)

func main() {
	onExit := func() {
		fmt.Println("Starting onExit")
	}
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("BrokenWall")
	systray.SetTooltip("Lantern")
	mQuitOrig := systray.AddMenuItem("退出", "Quit the whole app")
	mQuitOrig.SetIcon(icon.Data)
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()

	// We can manipulate the systray in other goroutines
	go func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle("BrokenWall")
		systray.SetTooltip("Pretty awesome棒棒嗒")
		//mChecked := systray.AddMenuItem("runfaild", "run me")
		systray.AddSeparator()
		//port := "1080"
		for {
			select {
			//case <-mChecked.ClickedCh:
			//	if mChecked.Checked() {
			//		mChecked.Uncheck()
			//		mChecked.SetTitle("runfaild")
			//	} else {
			//		mChecked.Check()
			//		mChecked.SetTitle("running:1080")
			//	}
			default:
				home,_ := homedir.Dir()
				configPath := home+"/Desktop/ss.json"
				//_, err := net.Listen("tcp", ":"+port)
				//if err != nil {
				//	time.Sleep(5 * time.Second)
				//	continue
				//}
				cmd := exec.Command("/bin/bash","-c",fmt.Sprintf("%s %s",home+"/Desktop/"+"mac_client-x64",configPath))
				cmd.Run()
					//bytes,_ := cmd.Output()
					//fmt.Println(string(bytes))
			}
		}
	}()
}
