package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/schollz/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/schollz/croc/v10/src/croc"
	"github.com/schollz/croc/v10/src/utils"
)

func sendTabItem(a fyne.App, w fyne.Window) *container.TabItem {
	status := widget.NewLabel("")
	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Sprint(r))
		}
	}()
	prog := widget.NewProgressBar()
	prog.Hide()
	topline := widget.NewLabel(lp("Pick a file to send"))
	randomCode := utils.GetRandomName()
	sendEntry := widget.NewEntry()
	sendEntry.SetText(randomCode)
	copyCodeButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		a.Clipboard().SetContent(sendEntry.Text)
	})
	copyCodeButton.Hide()

	sendDir, _ := os.MkdirTemp("", "crocgui-send")

	boxholder := container.NewVBox()
	senderScroller := container.NewVScroll(boxholder)
	fileentries := make(map[string]*fyne.Container)

	if len(os.Args) > 0 {
		for _, arg := range os.Args[1:] {
			if err := addPath(arg, sendDir, fileentries, boxholder, sendEntry); err != nil {
				log.Errorf(err.Error())
			}
		}
	}

	w.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		if len(uris) == 0 {
			return
		}
		for _, uri := range uris {
			if err := addPath(uri.Path(), sendDir, fileentries, boxholder, sendEntry); err != nil {
				log.Errorf(err.Error())
			}
		}
		SelectIndex(w, 0)
	})

	addFileButton := widget.NewButtonWithIcon("", theme.FileIcon(), func() {
		ShowFileOpen(func(f fyne.URIReadCloser, e error) {
			if e != nil {
				log.Errorf("Open dialog error: %s", e.Error())
				return
			}
			if f != nil {
				if err := addPath(f.URI().Path(), sendDir, fileentries, boxholder, sendEntry); err != nil {
					log.Errorf(err.Error())
				}
			}
		}, w)
	})

	debugBox := container.NewHBox(widget.NewLabel(lp("Debug log:")),
		layout.NewSpacer(),
		widget.NewButtonWithIcon(lp("Export full log"), theme.DocumentSaveIcon(), func() {
			savedialog := dialog.NewFileSave(func(f fyne.URIWriteCloser, e error) {
				if f != nil {
					logoutput.buf.WriteTo(f)
					f.Close()
				}
			}, w)
			savedialog.SetFileName("crocdebuglog.txt")
			savedialog.Resize(w.Canvas().Size())
			savedialog.Show()
		}))
	debugObjects = append(debugObjects, debugBox)

	cancelchan := make(chan bool)
	activeButtonHolder := container.NewVBox()
	var cancelButton, sendButton *widget.Button

	resetSender := func() {
		prog.Hide()
		prog.SetValue(0)
		for _, obj := range activeButtonHolder.Objects {
			activeButtonHolder.Remove(obj)
		}
		activeButtonHolder.Add(sendButton)

		for fpath, fe := range fileentries {
			boxholder.Remove(fe)
			os.Remove(fpath)
			log.Tracef("Removed file from internal cache: %s", fpath)
			delete(fileentries, fpath)
		}

		topline.SetText(lp("Pick a file to send"))
		addFileButton.Show()
		if sendEntry.Text == randomCode {
			randomCode = utils.GetRandomName()
			sendEntry.SetText(randomCode)
		}
		copyCodeButton.Hide()
		sendEntry.Enable()
	}

	sendButton = widget.NewButtonWithIcon(lp("Send"), theme.MailSendIcon(), func() {
		if len(sendEntry.Text) < 6 {
			log.Error("no receive code entered")
			dialog.ShowInformation(
				lp("Send"),
				lp("Enter code to download"),
				w,
			)
			return
		}

		// Only send if files selected
		if len(fileentries) < 1 {
			log.Error("no files selected")
			dialog.ShowInformation(
				lp("Send"),
				lp("Pick a file to send"),
				w,
			)
			return
		}

		addFileButton.Hide()
		sender, err := croc.New(croc.Options{
			IsSender:         true,
			SharedSecret:     sendEntry.Text,
			Debug:            crocDebugMode(),
			RelayAddress:     a.Preferences().String("relay-address"),
			RelayPorts:       strings.Split(a.Preferences().String("relay-ports"), ","),
			RelayPassword:    a.Preferences().String("relay-password"),
			Stdout:           false,
			NoPrompt:         true,
			DisableLocal:     a.Preferences().Bool("disable-local"),
			NoMultiplexing:   a.Preferences().Bool("disable-multiplexing"),
			OnlyLocal:        a.Preferences().Bool("force-local"),
			NoCompress:       a.Preferences().Bool("disable-compression"),
			Curve:            a.Preferences().String("pake-curve"),
			HashAlgorithm:    a.Preferences().String("croc-hash"),
			ThrottleUpload:   a.Preferences().String("upload-throttle"),
			ZipFolder:        false,
			GitIgnore:        false,
			MulticastAddress: a.Preferences().String("multicast-address"),
			Exclude:          []string{},
		})
		if err != nil {
			log.Errorf("croc error: %s\n", err.Error())
			return
		}
		log.SetLevel(crocDebugLevel())
		log.Trace("croc sender created")

		var filename string
		status.SetText(fmt.Sprintf("%s: %s", lp("Receive Code"), sendEntry.Text))
		copyCodeButton.Show()
		prog.Show()

		for _, obj := range activeButtonHolder.Objects {
			activeButtonHolder.Remove(obj)
		}
		activeButtonHolder.Add(cancelButton)

		donechan := make(chan bool)
		sendnames := make(map[string]int)
		go func() {
			ticker := time.NewTicker(time.Millisecond * 100)
			for {
				select {
				case <-ticker.C:
					if sender.Step2FileInfoTransferred {
						cnum := sender.FilesToTransferCurrentNum
						fi := sender.FilesToTransfer[cnum]
						filename = filepath.Base(fi.Name)
						sendnames[filename] = cnum
						fyne.Do(func() {
							topline.SetText(fmt.Sprintf("%s: %s(%d/%d)", lp("Sending file"), filename, cnum+1, len(sender.FilesToTransfer)))
							prog.Max = float64(fi.Size)
							prog.SetValue(float64(sender.TotalSent))
						})
					}
				case <-donechan:
					ticker.Stop()
					return
				}
			}
		}()
		go func() {
			var filepaths []string
			for fpath := range fileentries {
				filepaths = append(filepaths, fpath)
			}
			fyne.Do(sendEntry.Disable)
			fi, emptyfolders, numFolders, ferr := croc.GetFilesInfo(filepaths, false, false, []string{})
			if ferr != nil {
				log.Errorf("file info failed: %s\n", ferr)
			}
			serr := sender.Send(fi, emptyfolders, numFolders)
			donechan <- true
			if serr != nil {
				log.Errorf("Send failed: %s\n", serr)
			} else {
				fyne.Do(func() {
					status.SetText(fmt.Sprintf("%s: %s", lp("Sent file"), filename))
				})
			}
			fyne.Do(resetSender)
		}()
		go func() {
			<-cancelchan
			log.Warnf("Send cancelled. %s %v", sendDir, os.RemoveAll(sendDir))

			fyne.Do(func() {
				// status.SetText(lp("Send cancelled."))
				// resetSender()
				restart(a)
			})
		}()
	})

	cancelButton = widget.NewButtonWithIcon(lp("Cancel"), theme.CancelIcon(), func() {
		cancelchan <- true
	})

	activeButtonHolder.Add(sendButton)

	sendTop := container.NewVBox(
		container.NewHBox(topline, layout.NewSpacer(), addFileButton),
		widget.NewForm(&widget.FormItem{Text: lp("Send Code"), Widget: sendEntry}),
	)
	sendBot := container.NewVBox(
		activeButtonHolder,
		prog,
		container.NewHBox(status, copyCodeButton),
		debugBox,
	)

	return container.NewTabItemWithIcon(lp("Send"), theme.MailSendIcon(),
		container.NewBorder(sendTop, sendBot, nil, nil, senderScroller))
}

// Big File Dialog
func ShowFileOpen(callback func(reader fyne.URIReadCloser, err error), parent fyne.Window) {
	if mobile {
		dialog.ShowFileOpen(callback, parent)
		return
	}
	fd := dialog.NewFileOpen(callback, parent)
	fd.Resize(parent.Canvas().Size())
	fd.Show()
}

func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// Dirty refresh
func SelectIndex(window fyne.Window, index int) {
	var findTabs func(fyne.CanvasObject) *container.AppTabs
	findTabs = func(obj fyne.CanvasObject) *container.AppTabs {
		switch v := obj.(type) {
		case *container.AppTabs:
			return v
		case *fyne.Container:
			for _, child := range v.Objects {
				if tabs := findTabs(child); tabs != nil {
					return tabs
				}
			}
		}
		return nil
	}
	if tabs := findTabs(window.Content()); tabs != nil {
		tabs.SelectIndex(index)
		tabs.Refresh()
	}
}

func addPath(nfile, sendDir string,
	fileentries map[string]*fyne.Container,
	boxholder *fyne.Container,
	sendEntry *widget.Entry) error {
	fpath := filepath.Join(sendDir, filepath.Base(nfile))

	fi, err := os.Stat(nfile)
	if err != nil {
		return fmt.Errorf("URI (%s) %s", nfile, err.Error())
	} else if fi.IsDir() {
		log.Tracef("URI (%s), is dir", nfile)
		return nil
	}

	fi, err = os.Stat(fpath)
	if err == nil {
		log.Tracef("URI (%s), already in internal cache %s", nfile, fpath)
		return nil
	}

	if err := CopyFile(nfile, fpath); err != nil {
		return fmt.Errorf("Unable to copy file, error: %s - %s\n", sendDir, err.Error())
	}
	log.Tracef("URI (%s), copied to internal cache %s", nfile, fpath)

	if _, sterr := os.Stat(fpath); sterr != nil {
		return fmt.Errorf("Stat error: %s - %s\n", fpath, sterr.Error())
	}

	labelFile := widget.NewLabel(filepath.Base(nfile))
	newentry := container.NewHBox(
		labelFile,
		layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
			if !sendEntry.Disabled() {
				if fe, ok := fileentries[fpath]; ok {
					boxholder.Remove(fe)
					os.Remove(fpath)
					log.Tracef("Removed file from internal cache: %s", fpath)
					delete(fileentries, fpath)
				}
			}
		}),
	)

	fileentries[fpath] = newentry
	boxholder.Add(newentry)
	return nil
}
func restart(a fyne.App) {
	if !mobile {
		exec.Command(os.Args[0]).Start()
	}
	a.Quit()
}
