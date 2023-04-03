package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Slyk Fractions Generator")

	maxSizeLabel := widget.NewLabel("Max Size                                                      ")
	maxSizeEntry := widget.NewEntry()
	maxSizeEntry.SetPlaceHolder("Enter Max Size")
	unitSelect := widget.NewSelect([]string{"In", "Ft"}, func(selected string) {
		fmt.Println("Selected:", selected)
	})
	unitSelect.SetSelectedIndex(0)

	smallestFractionLabel := widget.NewLabel("Smallest Fraction")
	smallestFractionEntry := widget.NewEntry()
	smallestFractionEntry.SetPlaceHolder("Enter Smallest Fraction")

	outputLocationLabel := widget.NewLabel("Output Location")
	outputLocationEntry := widget.NewEntry()

	openFileExplorerButton := widget.NewButton("...", func() {

		openDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()

			outputLocationLabel.SetText(fmt.Sprintf("File Path: %s", writer.URI().Path()))
			// Add your logic for saving or processing the selected file
		}, w)
		openDialog.SetFileName("example.txt")
		openDialog.Show()
	})

	outputLocationEntry.SetPlaceHolder("Enter Output Location")

	generateBtn := widget.NewButton("Generate", func() {
		fmt.Println("Generate button clicked")
		MA, err := strconv.Atoi(maxSizeEntry.Text)
		if err != nil {
			fmt.Println("Error converting Max Size to int")
		}
		MB := unitSelect.Selected
		SF, err := strconv.Atoi(smallestFractionEntry.Text)
		if err != nil {
			fmt.Println("Error converting Smallest Fraction to int")
		}
		FO := outputLocationEntry.Text
		Generate(MA, MB, SF, FO)
	})

	sectionA := container.NewVBox(maxSizeLabel, container.NewBorder(nil, nil, nil, unitSelect, maxSizeEntry))
	sectionB := container.NewVBox(smallestFractionLabel, smallestFractionEntry)
	sectionC := container.NewVBox(outputLocationLabel, container.NewBorder(nil, nil, nil, openFileExplorerButton, outputLocationEntry))
	sectionD := container.NewVBox(generateBtn)

	content := container.NewVBox(
		sectionA,
		sectionB,
		sectionC,
		sectionD,
	)

	outerVBox := container.NewVBox(content)

	w.SetContent(outerVBox)
	w.Resize(fyne.NewSize(360, 200))
	w.ShowAndRun()
}

func Generate(MaxSize int, SmallestFraction string, MinFraction int, OutputLocation string) {
	fmt.Println("Generate function called")
	var Lines int
	switch SmallestFraction {
	case "In":
		Lines = MinFraction * MaxSize
	case "Ft":
		Lines = MinFraction * MaxSize * 12
	}
	file, err := os.Create(OutputLocation)
	if err != nil {
		fmt.Println("Error creating file")
	}
	defer file.Close()
	var Numerator int = 1
	var WholeInch int = 0
	var WholeFoot int = 0
	file.WriteString("pK,Decimal,Fraction\n")
	for i := 1; i <= Lines; i++ {
		ThisDecimal := float64(i) / float64(MinFraction)

		var ThisFraction, PA, PB, PC string = "", "", "", ""
		if WholeFoot > 0 {
			spacing := ""
			if WholeInch > 0 {
				spacing = " "
			} else if Numerator > 0 {
				spacing = " "
			}
			PC = fmt.Sprintf("%d'%s", WholeFoot, spacing)
		}
		if WholeInch > 0 {
			dash := ""
			if Numerator > 0 {
				dash = "-"
			} else {
				dash = "\""
			}
			PB = fmt.Sprintf("%d%s", WholeInch, dash)
		}
		if Numerator > 0 {
			PA = FractionLogic(Numerator, MinFraction)
		}
		ThisFraction = PC + PB + PA
		Numerator++
		if Numerator == MinFraction {
			WholeInch++
			Numerator = 0
		}
		if WholeInch == 12 {
			WholeFoot++
			WholeInch = 0
		}

		file.WriteString(fmt.Sprintf("%d,%f,%s\n", i, ThisDecimal, ThisFraction))
	}
	fmt.Println("Lines Total: ", Lines)
}

func FractionLogic(Nu int, De int) string {
	for math.Remainder(float64(Nu), 2) == 0 {
		Nu = Nu / 2
		De = De / 2
	}

	FractionIs := fmt.Sprintf("%d/%d\"", Nu, De)
	return FractionIs
}
