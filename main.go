package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	colorPrint "github.com/fatih/color"
	colors "github.com/teacat/noire"
)

func main() {
	fmt.Println()

	flags := flag.NewFlagSet("Highlight colors (9-15)", flag.ContinueOnError)
	hisat :=
		flags.Float64("hisat", 0, "saturation for highlights")
	hidesat :=
		flags.Float64("hidesat", 0.1, "desaturation for highlights")
	hihue :=
		flags.Float64("hihue", -5, "hue for highlights")
	hilight :=
		flags.Float64("hilight", 0.25, "lighten for highlights")

	nsat :=
		flags.Float64("nsat", 0, "saturation for normal colors")
	ndesat :=
		flags.Float64("ndesat", 0, "desaturation for normal colors")
	nhue :=
		flags.Float64("nhue", 0, "hue for normal colors")
	nlight :=
		flags.Float64("nlight", 0, "lighten for normal colors")

	output_file :=
		flags.String("out", "", "different output than input file")

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: " + os.Args[0] + " [FILE]")
		fmt.Println("Target color.json generated by pywal to generate ")
		fmt.Println("highlight colors from the existing 8 colors.")
		fmt.Println()
		fmt.Println("default path: ~/.cache/wal/colors.json")
		fmt.Println()
		flags.Usage()
		return
	}
	flags.Parse(os.Args[2:])

	filePath := string(args[0])
	savePath := filePath
	jsonFile, err := os.Open(filePath)

	if err != nil {
		printError(err.Error())
		return
	}
	if *output_file != "" {
		savePath = *output_file
	}
	fmt.Println(savePath)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result colorScheme
	json.Unmarshal([]byte(byteValue), &result)

	wallpaper := strings.SplitAfter(result.Wallpaper, string('/'))
	printLine("Wallpaper: "+wallpaper[len(wallpaper)-1], colorPrint.FgGreen)
	printLine("Generating brighter colors", colorPrint.FgGreen)
	printError(fmt.Sprint(*nhue != 0))

	for key := range result.Colors {
		test := strings.SplitAfter(key, "color")
		n, _ := strconv.ParseInt(test[len(test)-1], 10, 64)
		if (*nhue != 0) || (*nsat != 0) || (*nlight != 0) || (*ndesat != 0) {
			result.Colors[key] = "#" + colors.
				NewHex(result.Colors["color"+fmt.Sprint(n)]).
				AdjustHue(*nhue).
				Lighten(*nlight).
				Desaturate(*ndesat).
				Saturate(*nsat).
				Hex()
			fmt.Println("     " + key)
		}

		if n > 8 {
			base_color := fmt.Sprint(n - 8)
			result.Colors[key] = "#" + colors.
				NewHex(result.Colors["color"+base_color]).
				AdjustHue(*hihue).
				Lighten(*hilight).
				Desaturate(*hidesat).
				Saturate(*hisat).
				Hex()
			fmt.Println("     " + key)
		}
	}

	result.Colors["color8"] = "#" + colors.NewHex(result.Colors["color0"]).Lighten(0.05).Hex()
	result.Special["cursor"] = "#" + colors.NewHex(result.Colors["color2"]).Lighten(0.10).Saturate(0.10).Hex()

	printLine("Done!", colorPrint.FgGreen)

	bytes, _ := json.MarshalIndent(result, "", "")
	_ = ioutil.WriteFile(savePath, bytes, 0644)
}

func printLine(text string, c colorPrint.Attribute) {
	col := colorPrint.New(c).Add(colorPrint.Bold)
	col.Print(":: ")
	col = colorPrint.New(colorPrint.FgWhite)
	col.Print(text)
	col.Println()
}

func printError(text string) {
	col := colorPrint.New(colorPrint.FgHiRed).Add(colorPrint.Bold)
	col.Print(":: ")
	col = colorPrint.New(colorPrint.FgRed).Add(colorPrint.Bold)
	col.Print(text)
	col.Println()
}

type colorScheme struct {
	Wallpaper string            `json:"wallpaper"`
	Special   map[string]string `json:"special"`
	Colors    map[string]string `json:"colors"`
}
