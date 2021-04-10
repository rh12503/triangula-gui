package main

import (
  _ "embed"
  "github.com/wailsapp/wails"
  _ "image/jpeg"
  _ "image/png"
)

//go:embed frontend/build/main.js
var js string

//go:embed frontend/build/main.css
var css string

func main() {

  app := wails.CreateApp(&wails.AppConfig{
		Width:  980,
		Height: 720,
		Title:  "Triangula",
		JS:     js,
		CSS:    css,
		Colour: "rgb(225,225,225)",
		Resizable: true,
		MinWidth: 580,
		MinHeight: 426,
  })
  
  app.Bind(&Runner{})
  
  app.Run()
}
