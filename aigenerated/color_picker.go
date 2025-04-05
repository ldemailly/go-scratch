package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
)

// ColorData struct to pass color information to the template
type ColorData struct {
	Red   int
	Green int
	Blue  int
	Hex   string
}

func generateRandomColor() ColorData {
	r := rand.IntN(256)
	g := rand.IntN(256)
	b := rand.IntN(256)

	hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)

	return ColorData{
		Red:   r,
		Green: g,
		Blue:  b,
		Hex:   hex,
	}
}

func colorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		redStr := r.Form.Get("red")
		greenStr := r.Form.Get("green")
		blueStr := r.Form.Get("blue")

		red, err := strconv.Atoi(redStr)
		if err != nil || red < 0 || red > 255 {
			http.Error(w, "Invalid red value", http.StatusBadRequest)
			return
		}
		green, err := strconv.Atoi(greenStr)
		if err != nil || green < 0 || green > 255 {
			http.Error(w, "Invalid green value", http.StatusBadRequest)
			return
		}
		blue, err := strconv.Atoi(blueStr)
		if err != nil || blue < 0 || blue > 255 {
			http.Error(w, "Invalid blue value", http.StatusBadRequest)
			return
		}

		hex := fmt.Sprintf("#%02x%02x%02x", red, green, blue)
		colorData := ColorData{
			Red:   red,
			Green: green,
			Blue:  blue,
			Hex:   hex,
		}
		renderColorSliders(w, colorData)
		return
	}

	// Handle GET request (initial load or after submitting)
	colorData := generateRandomColor()
	renderColorSliders(w, colorData)
}

func renderColorSliders(w http.ResponseWriter, data ColorData) {
	tmpl, err := template.New("colorSliders").Parse(colorSlidersTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", colorHandler)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


// HTML Template with Sliders (Revised for Dynamic Updates)
const colorSlidersTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Color Picker with Sliders</title>
    <style>
        body {
            font-family: sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            margin: 0;
            background-color: #f0f0f0;
        }
        .color-container {
            text-align: center;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }
        #colorPreview { /* Added ID selector for specificity */
            width: 150px;
            height: 150px;
            border-radius: 50%;
            margin: 20px auto;
            border: 2px solid #ccc;
            background-color: rgb({{ .Red }}, {{ .Green }}, {{ .Blue }});
        }
        .color-info {
            margin-top: 15px;
        }
        .slider-group {
            margin-bottom: 15px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
        }
        input[type="range"] {
            width: 200px;
        }
        .slider-value {
            font-size: 0.9em;
            color: #555;
        }
        button {
            padding: 10px 20px;
            font-size: 16px;
            cursor: pointer;
            border: none;
            border-radius: 5px;
            background-color: #007bff;
            color: white;
            margin-top: 15px;
        }
        button:hover {
            background-color: #0056b3;
        }
        .controls {
            margin-top: 20px;
        }
    </style>
    <script>
        function updateColor() {
            const red = document.getElementById('red').value;
            const green = document.getElementById('green').value;
            const blue = document.getElementById('blue').value;

            const colorPreview = document.getElementById('colorPreview');
            if (colorPreview) {
                colorPreview.style.backgroundColor = 'rgb(' + red + ', ' + green + ', ' + blue + ')';
            }

            const redValue = document.getElementById('redValue');
            if (redValue) {
                redValue.textContent = red;
            }
            const greenValue = document.getElementById('greenValue');
            if (greenValue) {
                greenValue.textContent = green;
            }
            const blueValue = document.getElementById('blueValue');
            if (blueValue) {
                blueValue.textContent = blue;
            }
            const hexValue = document.getElementById('hexValue');
            if (hexValue) {
                hexValue.textContent = rgbToHex(parseInt(red), parseInt(green), parseInt(blue));
            }
        }

        function rgbToHex(r, g, b) {
            return "#" + componentToHex(r) + componentToHex(g) + componentToHex(b);
        }

        function componentToHex(c) {
            const hex = c.toString(16);
            return hex.length == 1 ? "0" + hex : hex;
        }

        // Call updateColor on initial load to set the preview based on Go template values
        window.onload = updateColor;
    </script>
</head>
<body>
    <div class="color-container">
        <h1>Color Picker</h1>
        <div id="colorPreview" class="color-preview"></div>
        <div class="color-info">
            <p>RGB: <span id="redValue">{{ .Red }}</span>, <span id="greenValue">{{ .Green }}</span>, <span id="blueValue">{{ .Blue }}</span></p>
            <p>Hex: <span id="hexValue">{{ .Hex }}</span></p>
        </div>
        <div class="controls">
            <h2>Adjust Color</h2>
            <div class="slider-group">
                <label for="red">Red:</label>
                <input type="range" id="red" name="red" min="0" max="255" value="{{ .Red }}" oninput="updateColor()">
                <span class="slider-value">{{ .Red }}</span>
            </div>
            <div class="slider-group">
                <label for="green">Green:</label>
                <input type="range" id="green" name="green" min="0" max="255" value="{{ .Green }}" oninput="updateColor()">
                <span class="slider-value">{{ .Green }}</span>
            </div>
            <div class="slider-group">
                <label for="blue">Blue:</label>
                <input type="range" id="blue" name="blue" min="0" max="255" value="{{ .Blue }}" oninput="updateColor()">
                <span class="slider-value">{{ .Blue }}</span>
            </div>
            <button onclick="location.reload()">Generate Random Color</button>
        </div>
    </div>
</body>
</html>
`
