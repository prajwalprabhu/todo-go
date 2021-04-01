package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

var labelList = list.New()
var config_file = "data.json"

// var data = make(map[string]string)
type data_ struct {
	Date  string
	Label string
}

// Field int `json:"myName"`
var data_list = []data_{}

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("ToDo")
	win.Connect("destroy", func() {
		save_data()
		gtk.MainQuit()
	})
	get_data()
	win.Add(windowWidget())
	win.ShowAll()

	gtk.Main()
}

func windowWidget() *gtk.Widget {
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	unused, err := gtk.LabelNew("This label is never used")
	if err != nil {

		unused.Destroy()
	}

	sw, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create scrolled window:", err)
	}

	grid.Attach(sw, 0, 0, 2, 1)
	sw.SetHExpand(true)
	sw.SetVExpand(true)

	labelsGrid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	labelsGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	sw.Add(labelsGrid)
	labelsGrid.SetHExpand(true)

	insertBtn, err := gtk.ButtonNewWithLabel("Add a label")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	removeBtn, err := gtk.ButtonNewWithLabel("Remove a label")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}

	insertBtn.Connect("clicked", func() {
		new_todo()

	})

	removeBtn.Connect("clicked", func() {
		rm_todo()
	})
	for _, data := range data_list {
		newlabel, err := gtk.LabelNew(fmt.Sprint("Date : ", data.Date, "\t Label : ", data.Label))
		if err != nil {
			log.Print("Unable to create label:", err)
		}

		labelList.PushBack(newlabel)
		labelsGrid.Add(newlabel)
		newlabel.SetHExpand(true)
		labelsGrid.ShowAll()

	}
	grid.Attach(insertBtn, 0, 1, 1, 1)
	grid.Attach(removeBtn, 1, 1, 1, 1)

	return &grid.Container.Widget
}
func new_todo() {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("ToDo")
	win.Connect("destroy", func() {
		// save_data()
		// gtk.MainQuit()
		win.Destroy()
	})
	labelsGrid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	labelsGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	a, e := gtk.LabelNew("Hello World")
	fmt.Println(e)
	labelsGrid.Add(a)
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
		new_data := data_{date, label}
		data_list = append(data_list, new_data)
	})
	labelsGrid.Add(entry_label)
	labelsGrid.Add(insertBtn)
	win.Add(labelsGrid)
	win.ShowAll()
	// win.Destroy()
}
func rm_todo() {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	check_err(err)
	grid, err := gtk.GridNew()
	check_err(err)
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	list_box, err := gtk.ListBoxNew()
	// list_box.SetSelectionMode(gtk.SELECTION_MULTIPLE)
	check_err(err)
	for i, data := range data_list {
		newlabel, err := gtk.LabelNew(fmt.Sprint("Date : ", data.Date, "\t Label : ", data.Label))
		check_err(err)
		labelList.PushBack(newlabel)
		list_box.Insert(newlabel, i)
	}
	rm_button, err := gtk.ButtonNewWithLabel("Remove")
	check_err(err)
	rm_button.Connect("clicked", func() {
		fmt.Println("Clicked")
		list := list_box.GetSelectedRow()
		index := list.GetIndex()
		fmt.Print(index)
		copy(data_list[index:], data_list[index+1:]) // Shift a[i+1:] left one index.
		data_list[len(data_list)-1] = data_{}     // Erase last element (write zero value).
		data_list = data_list[:len(data_list)-1]    
		
	})

	grid.Add(list_box)
	grid.Add(rm_button)
	win.Connect("destroy", func() {
		win.Destroy()
	})
	win.Add(grid)
	win.ShowAll()
}
func get_data() {
	file_data, err := ioutil.ReadFile(config_file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file_data, &data_list)
	if err != nil {
		log.Fatal(err)
	}

}
func save_data() {
	json_data, err := json.Marshal(data_list)
	if err != nil {
		log.Fatal("Error :", err)
	}
	fmt.Println(string(json_data))
	fmt.Println(config_file)
	err = ioutil.WriteFile(config_file, json_data, 0777)
	check_err(err)
}
func check_err(err error) {
	if err != nil {
		log.Fatal("Error :", err)
	}
}
