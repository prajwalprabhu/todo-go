package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type App struct {
	config_file string
	datalist    []data
	labellist   []*gtk.Label
	labelsGrid  *gtk.Grid
}

type data struct {
	Date  string
	Label string
}

func main() {
	app := App{config_file: "data.json"}
	app.run()
}
func (main *App) run() {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("ToDo")
	win.Connect("destroy", func() {
		main.save_data()
		gtk.MainQuit()
	})
	main.get_data()
	win.Add(main.windowWidget())
	win.ShowAll()
	gtk.Main()
}
func (main *App) windowWidget() *gtk.Widget {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	sw, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}
	grid.Attach(sw, 0, 0, 2, 1)
	sw.SetHExpand(true)
	sw.SetVExpand(true)
	main.labelsGrid, err = gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	main.labelsGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	sw.Add(main.labelsGrid)
	main.labelsGrid.SetHExpand(true)
	insertBtn, err := gtk.ButtonNewWithLabel("Add a label")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	removeBtn, err := gtk.ButtonNewWithLabel("Remove a label")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	insertBtn.Connect("clicked", func() {
		main.new_todo()
	})
	removeBtn.Connect("clicked", func() {
		main.rm_todo()
	})
	for _, data := range main.datalist {
		newlabel, err := gtk.LabelNew(fmt.Sprint("Date : ", data.Date, "\t Label : ", data.Label))
		if err != nil {
			log.Print("Unable to create label:", err)
		}
		main.labelsGrid.Add(newlabel)
		main.labellist = append(main.labellist, newlabel)
		newlabel.SetHExpand(true)
		main.labelsGrid.ShowAll()
	}
	grid.Attach(insertBtn, 0, 1, 1, 1)
	grid.Attach(removeBtn, 1, 1, 1, 1)
	return &grid.Container.Widget
}
func (main *App) new_todo() {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("ToDo")
	win.Connect("destroy", func() {
		win.Destroy()
	})
	labelsGrid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	labelsGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	entry_buffer, err := gtk.EntryBufferNew("Enter Date", 10)
	check_err(err)
	entry_date, err := gtk.EntryNewWithBuffer(entry_buffer)
	check_err(err)
	labelsGrid.Add(entry_date)
	entry_buffer2, err := gtk.EntryBufferNew("Enter Label", 15)
	check_err(err)
	entry_label, err := gtk.EntryNewWithBuffer(entry_buffer2)
	check_err(err)
	insertBtn, err := gtk.ButtonNewWithLabel("Create")
	check_err(err)
	insertBtn.Connect("clicked", func() {
		fmt.Println("Created")
		date, err := entry_buffer.GetText()
		check_err(err)
		label, err := entry_buffer2.GetText()
		check_err(err)
		new_data := data{date, label}
		main.datalist = append(main.datalist, new_data)
		main.restart()
	})
	labelsGrid.Add(entry_label)
	labelsGrid.Add(insertBtn)

	win.Add(labelsGrid)
	win.ShowAll()
}
func (main *App) rm_todo() {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	check_err(err)
	grid, err := gtk.GridNew()
	check_err(err)
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	list_box, err := gtk.ListBoxNew()
	check_err(err)
	for i, data := range main.datalist {
		newlabel, err := gtk.LabelNew(fmt.Sprint("Date : ", data.Date, "\t Label : ", data.Label))
		check_err(err)
		list_box.Insert(newlabel, i)
	}
	rm_button, err := gtk.ButtonNewWithLabel("Remove")
	check_err(err)
	rm_button.Connect("clicked", func() {
		fmt.Println("Clicked")
		list := list_box.GetSelectedRow()
		index := list.GetIndex()
		fmt.Print(index)
		copy(main.datalist[index:], main.datalist[index+1:]) // Shift a[i+1:] left one index.
		main.datalist[len(main.datalist)-1] = data{}         // Erase last element (write zero value).
		main.datalist = main.datalist[:len(main.datalist)-1]

	})
	grid.Add(list_box)
	grid.Add(rm_button)
	win.Connect("destroy", func() {
		win.Destroy()
	})
	win.Add(grid)
	win.ShowAll()
}
func (main *App) get_data() {
	file_data, err := ioutil.ReadFile(main.config_file)
	if err != nil {
		return
	}
	err = json.Unmarshal(file_data, &main.datalist)
	if err != nil {
		return
	}
}
func (main *App) save_data() {
	json_data, err := json.Marshal(main.datalist)
	if err != nil {
		log.Fatal("Error :", err)
	}
	err = ioutil.WriteFile(main.config_file, json_data, 0777)
	check_err(err)
}
func (main *App) restart() {
	for _, val := range main.labellist {
		val.Destroy()
	}
	for _, data := range main.datalist {
		newlabel, err := gtk.LabelNew(fmt.Sprint("Date : ", data.Date, "\t Label : ", data.Label))
		if err != nil {
			log.Print("Unable to create label:", err)
		}
		main.labelsGrid.Add(newlabel)
		main.labellist = append(main.labellist, newlabel)
		newlabel.SetHExpand(true)
		main.labelsGrid.ShowAll()
	}
}
func check_err(err error) {
	if err != nil {
		log.Fatal("Error :", err.Error())
	}
}
