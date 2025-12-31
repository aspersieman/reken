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

	btn := func(label string, fn func()) *widget.Button {
		return widget.NewButton(label, fn)
	}

	grid := container.NewGridWithColumns(4,
		btn("7", func() { appendText("7") }),
		btn("8", func() { appendText("8") }),
		btn("9", func() { appendText("9") }),
		btn("÷", func() { appendText("÷") }),

		btn("4", func() { appendText("4") }),
		btn("5", func() { appendText("5") }),
		btn("6", func() { appendText("6") }),
		btn("×", func() { appendText("×") }),

		btn("1", func() { appendText("1") }),
		btn("2", func() { appendText("2") }),
		btn("3", func() { appendText("3") }),
		btn("-", func() { appendText("-") }),

		btn("0", func() { appendText("0") }),
		btn(".", func() { appendText(".") }),
		btn("=", equal),
		btn("+", func() { appendText("+") }),
	)

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeyReturn, fyne.KeyEnter:
			equal()
		case fyne.KeyEscape:
			clear()
		}
	})

	w.SetContent(container.NewVBox(
		display,
		grid,
	))

	w.ShowAndRun()
}
