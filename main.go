package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	assets "github.com/aspersieman/reken/assets"
)

var version = "dev"

func eval(expr string) (string, error) {
	e, err := parser.ParseExpr(expr)
	if err != nil {
		return "", err
	}

	var evalNode func(ast.Expr) (float64, error)

	evalNode = func(n ast.Expr) (float64, error) {
		switch v := n.(type) {

		case *ast.BinaryExpr:
			l, err := evalNode(v.X)
			if err != nil {
				return 0, err
			}
			r, err := evalNode(v.Y)
			if err != nil {
				return 0, err
			}

			switch v.Op {
			case token.ADD:
				return l + r, nil
			case token.SUB:
				return l - r, nil
			case token.MUL:
				return l * r, nil
			case token.QUO:
				return l / r, nil
			default:
				return 0, fmt.Errorf("unsupported operator")
			}

		case *ast.BasicLit:
			return strconv.ParseFloat(v.Value, 64)

		default:
			return 0, fmt.Errorf("invalid expression")
		}
	}

	val, err := evalNode(e)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(val, 'f', -1, 64), nil
}

func main() {
	fmt.Println("Reken Calculator v" + version)
	a := app.New()
	a.SetIcon(assets.ResourceIconPng)

	w := a.NewWindow("Reken Calculator")
	w.Resize(fyne.NewSize(260, 360))

	display := widget.NewEntry()
	display.Disable()
	display.SetText("0")

	appendText := func(s string) {
		if display.Text == "0" {
			display.SetText(s)
		} else {
			display.SetText(display.Text + s)
		}
	}

	backspace := func() {
		if len(display.Text) <= 1 {
			display.SetText("0")
			return
		}
		display.SetText(display.Text[:len(display.Text)-1])
	}

	clear := func() {
		display.SetText("0")
	}

	equal := func() {
		expr := strings.ReplaceAll(display.Text, "×", "*")
		expr = strings.ReplaceAll(expr, "÷", "/")

		result, err := eval(expr)
		if err != nil {
			display.SetText("Error")
			return
		}
		display.SetText(result)
	}

	handleInput := func(s string) {
		switch s {
		case "C":
			clear()
		case "=":
			equal()
		case "⌫":
			backspace()
		default:
			appendText(s)
		}
	}

	btn := func(label string) *widget.Button {
		return widget.NewButton(label, func() {
			handleInput(label)
		})
	}

	grid := container.NewGridWithColumns(4,
		btn("C"), btn("⌫"), btn("÷"), btn("×"),
		btn("7"), btn("8"), btn("9"), btn("-"),
		btn("4"), btn("5"), btn("6"), btn("+"),
		btn("1"), btn("2"), btn("3"), btn("="),
		btn("0"), btn("."), widget.NewLabel(""), widget.NewLabel(""),
	)

	w.Canvas().SetOnTypedRune(func(r rune) {
		fmt.Printf("Typed rune: %v\n", r)
		switch r {

			// Digits (top row + numpad)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			handleInput(string(r))

			// Operators (top row + numpad)
		case '+', '-':
			handleInput(string(r))
		case '*':
			handleInput("×")
		case '/':
			handleInput("÷")

			// Decimal
		case '.':
			handleInput(".")
		}
	})

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		fmt.Printf("Typed key: %v\n", k)
		switch k.Name {

			// Equals (Enter key, including numpad Enter)
		case fyne.KeyReturn, fyne.KeyEnter:
			handleInput("=")

			// Backspace
		case fyne.KeyBackspace:
			handleInput("⌫")

			// Clear
		case fyne.KeyEscape, fyne.KeyDelete:
			handleInput("C")
		}
	})

	w.SetContent(container.NewVBox(
		display,
		grid,
	))

	w.ShowAndRun()
}
