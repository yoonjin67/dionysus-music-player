package main

import (
	"errors"
	"log"
	"os"
	"strings"
//TODO: use beep library for audio play
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var supported_codecs [11]string = [11]string{
	"mp3", //working on it
	"flac",
	"ogg", //working on it
	"m4a",
	"wav",
	"wma",
	"aiff",
	"dsd",
	"alac",
	"pcm",
	"aac",
} //supported audio codecs

const CURRENT_LOGLEVEL = DEBUG

type songfileInfo struct {
	string name,
	string location 
}

func main() {
	mainApp, err := gtk.ApplicationNew(APPLICATION_ID, glib.APPLICATION_FLAGS_NONE)
	errorHandler(err, "Creating main app", CURRENT_LOGLEVEL, "error")
	mainApp.Connect("startup", func() {
		log.Println("Application Startup")
	})

	mainApp.Connect("activate", func() {
		uiStartUp(mainApp)
	})
	mainApp.Connect("shutdown", func() {
		log.Println("Shutdown")
	})

	os.Exit(mainApp.Run(os.Args))
}

func uiStartUp(mainApp *gtk.Application) {
	//get main ui
	log.Println("Application Activate")
	builder, err := gtk.BuilderNewFromFile(MAIN_WINDOW_XML)
	errorHandler(err, "Loading UI", CURRENT_LOGLEVEL, "error")

	songnameList := make([]songfileInfo, 500000)
	signals := map[string]interface{}{
		"on_main_window_destroy": func() {
			if CURRENT_LOGLEVEL == DEBUG {
				log.Println("Application Destroyed")
			}
		},
	}

	builder.ConnectSignals(signals) //log when main window destroyed. note: debug

	//get main window
	obj, err := builder.GetObject("Dionysus-ToplevelWindow")
	errorHandler(err, "loading toplevel Window", CURRENT_LOGLEVEL, "error")
	win, err := IS_APP_WINDOW(obj)
	errorHandler(err, "checking object type: *Gtk.ApplicationWindow", CURRENT_LOGLEVEL, "warn")
	if win == nil {
		if CURRENT_LOGLEVEL >= TEST {
			log.Println("Top Level Window is nil")
		}
		return
	}

	obj, err = builder.GetObject("Songlist")
	errorHandler(err, "Loading SongList Widget", CURRENT_LOGLEVEL, "warn")
	songlist, err := IS_TREE_VIEW(obj)
	songinfoArr := make([]string, 5)

	songinfoArr[0] = "Track No."
	songinfoArr[1] = "Artist"
	songinfoArr[2] = "Song Name"
	songinfoArr[3] = "Released"
	songinfoArr[4] = "Time"

	for _, info := range songinfoArr {
		col, err := gtk.TreeViewColumnNew()
		errorHandler(err, "Creating Song List Column Label", CURRENT_LOGLEVEL, "warn")
		col.SetTitle(info)
		col.SetSortOrder(gtk.SORT_ASCENDING)
		col.Connect("clicked", func() {
			if col.GetSortOrder() == gtk.SORT_ASCENDING {
				col.SetSortOrder(gtk.SORT_DESCENDING)
			} else {
				col.SetSortOrder(gtk.SORT_ASCENDING)
			}
		})
		songlist.AppendColumn(col)
	}

	fmenuList := make([]string, 2)

	fmenuList[0] = "FileMenuQuit"
	fmenuList[1] = "FileMenuOpen"
	for i, curr := range fmenuList {
		if curr == "" {
			log.Println("menuitem name is nil?")
			continue
		}
		obj, err = builder.GetObject(curr)
		errorHandler(err, "loading file menu(quit)", CURRENT_LOGLEVEL, "error")
		itm, err := IS_MENU_ITEM(obj)
		errorHandler(err, "checking object type: *Gtk.MenuItem", CURRENT_LOGLEVEL, "error")
		if itm == nil {
			log.Println("Menu Item (Quit) is nil")
			return
		}
		switch i {
		case 0:
			itm.Connect("activate", func() { os.Exit(0) })
		case 1:
			itm.Connect("activate", func() {
				dialog, err := gtk.FileChooserDialogNewWith1Button("Open Files...", win, gtk.FILE_CHOOSER_ACTION_OPEN, "Open", gtk.RESPONSE_ACCEPT)
				dialog.SetSelectMultiple(true)

				//open file chooser dialog, and retrieve file names as []string type
				errorHandler(err, "Creating File Chooser Dialog", CURRENT_LOGLEVEL, "warn")
				res := dialog.Run()
				if res == gtk.RESPONSE_ACCEPT {
					filenames, err := dialog.GetFilenames()
					errorHandler(err, "Loading File Names", CURRENT_LOGLEVEL, "info")
					for _, musicFile := range filenames {
						filenameSplit := strings.Split(musicFile, ".")
						extension := filenameSplit[len(filenameSplit)-1]
						if inCodecList(extension) {
							if CURRENT_LOGLEVEL == DEBUG {
								fnameSlice := filenameSplit[0 : len(filenameSplit)-1]
								fname := strings.Join(fnameSlice, ".")
								fnameSlice = strings.Split(fname, "/")
								fname = fnameSlice[len(fnameSlice)-1]
								log.Println("Music Name: ", fname)
								var currSongInfo songfileInfo
								currSongInfo.name = fname
								currsongInfo.location = musicFile
								songnameList = append(songnameList,currSongInfo)
							}
						}
						//todo: music file loading
					}
					dialog.Destroy()
					appendSongInfos(songlist,songnameList)
				}
			})
		}

	}

	win.Show()
	mainApp.AddWindow(win)

}

func appendSongInfos(view *Gtk.TreeViewm, songs []songFileInfo) {
//TODO: Get proper tags and display
}

func inCodecList(extension string) bool {
	for _, codec := range supported_codecs {
		if strings.Compare(extension, codec) == 0 {
			return true
		}
	}
	return false
}

func IS_TREE_VIEW(obj glib.IObject) (*gtk.TreeView, error) {
	if itm, ok := obj.(*gtk.TreeView); ok {
		return itm, nil
	}
	return nil, errors.New("Not a Gtk.TreeViwe")
}

func IS_MENU_ITEM(obj glib.IObject) (*gtk.MenuItem, error) {
	if itm, ok := obj.(*gtk.MenuItem); ok {
		return itm, nil
	}

	return nil, errors.New("Not a *Gtk.MenuItem")
}

func IS_APP_WINDOW(obj glib.IObject) (*gtk.ApplicationWindow, error) {
	if win, ok := obj.(*gtk.ApplicationWindow); ok {
		return win, nil
	}
	return nil, errors.New("Not a *Gtk.ApplicationWindow")
}
