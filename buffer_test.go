package simple_buffer

// import (
// 	"fmt"
// 	"testing"
// )

// func TestUp(t *testing.T) {
// 	buf := New([]rune(code)...)

// 	buf.cursor = 0
// 	buf.line = 0

// 	fmt.Println(buf.GetRange(0, buf.cursor))

// 	for i := 0; i < 3; i++ {
// 		buf.CursorRight()
// 		fmt.Println(buf.GetRange(0, buf.cursor))
// 	}

// 	for i := 0; i < 3; i++ {
// 		buf.CursorDown()
// 		fmt.Println(buf.GetRange(0, buf.cursor))
// 	}

// 	// buf.CursorDown()
// 	fmt.Println(buf.GetRange(0, buf.cursor))
// 	// assert.NoError(t, buf.Append('a', '\n', 'b'))

// }

// const code = `package main

// import (
// 	_ "embed"

// 	"github.com/KlyuchnikovV/edicode/api"
// 	"github.com/KlyuchnikovV/edicode/core"
// 	"github.com/wailsapp/wails"
// 	"golang.org/x/net/context"
// )

// //go:embed ui/public/build/bundle.js
// var js string

// //go:embed ui/public/index.html
// var html string

// //go:embed ui/public/build/bundle.css
// var css string

// func main() {
// 	app := wails.CreateApp(&wails.AppConfig{
// 		Width:  1000,
// 		Height: 700,
// 		Title:  "edicode",
// 		JS:     js,
// 		CSS:    css,
// 		HTML:   html,
// 		Colour: "#FFFFFF",
// 	})

// 	c, err := core.New(context.Background(), "main.go", "main_test.go")
// 	if err != nil {
// 		panic(err)
// 	}

// 	a := api.New(c)
// 	a.Bind(app)

// 	if err := c.Start(); err != nil {
// 		panic(err)
// 	}

// 	if err := app.Run(); err != nil {
// 		panic(err)
// 	}
// }`
