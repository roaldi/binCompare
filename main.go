package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image"
	"image/color"
	"io/ioutil"
	"github.com/roaldi/gobinviz"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	dia "github.com/sqweek/dialog"
	"strconv"
)

var (
	fileA string
	fileB string
	imgA image.Image
	imgB image.Image
	update = make(chan bool,1024)
	distances string
)

func main() {

	imgPlaceholder := image.NewRGBA(image.Rect(0, 0, 100, 100))
		// Colors are defined by Red, Green, Blue, Alpha uint8 values.
		cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			switch {
			case x < 100/2 && y <100/2: // upper left quadrant
				imgPlaceholder.Set(x, y, cyan)
			case x >= 100/2 && y >= 100/2: // lower right quadrant
				imgPlaceholder.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	fileAWidget := widget.NewFormItem("", canvas.NewImageFromImage(imgPlaceholder))
	fileBWidget := widget.NewFormItem("", canvas.NewImageFromImage(imgPlaceholder))
	fileAButton := widget.NewButton("File A", func() {
		fileA, _ = dia.File().Load()
	})
	fileBButton := widget.NewButton("File B", func() {
		fileB, _ = dia.File().Load()
	})
	processButton := widget.NewButton("Process", func() {
		process()
	})




	a := app.New()
	w := a.NewWindow("Binary Image Comparison")
	go func(){
		for {
			localUpdate := <-update
			if localUpdate {
				fileAImage := canvas.NewImageFromImage(imgA)
				fileBImage := canvas.NewImageFromImage(imgB)
				fileAImage.SetMinSize(fyne.NewSize(float32(imgA.Bounds().Size().X), float32(imgA.Bounds().Size().Y)))
				fileBImage.SetMinSize(fyne.NewSize(float32(imgB.Bounds().Size().X), float32(imgB.Bounds().Size().Y)))
				renderA := widget.NewFormItem("", fileAImage)
				renderB := widget.NewFormItem("", fileBImage)
				w.SetContent(container.NewVBox(widget.NewForm(renderA,renderB),fileAButton,fileBButton,processButton,widget.NewLabel(distances)))
			}
		}
	}()
	w.SetContent(container.NewVBox(widget.NewLabel("test"),
		widget.NewForm(fileAWidget,fileBWidget),
		fileAButton,
		fileBButton,
		processButton,))
	w.ShowAndRun()
}

func process() {
	byteData, _ := ioutil.ReadFile(fileA)
	r, _ := binviz.ProcessBinary(byteData)
	imgA = r.Image
	altData, _ := ioutil.ReadFile(fileB)
	a, _ := binviz.ProcessBinary(altData)
	imgB = a.Image

	distances = "Average Distance: " + strconv.Itoa(r.AverageDistance(a)) + "\n"
	distances += "Block Distance: " + strconv.Itoa(r.BlockHashDistance(a)) + "\n"
	distances += "Difference Distance: " + strconv.Itoa(r.DifferenceDistance(a)) + "\n"
	distances += "Marr-Hilde Distance: " + strconv.Itoa(r.MarrHildeDistance(a)) + "\n"
	distances += "Median Distance: " + strconv.Itoa(r.MedianDistance(a)) + "\n"
	update <- true
}