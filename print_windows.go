//go:build windows
// +build windows

package main

import (
	"fmt"

	"github.com/alexbrainman/printer"
)

func printToDefaultPrinter(receipt string) {
	printerName, err := printer.Default()
	if err != nil {
		fmt.Println("Error getting default printer:", err)
		return
	}

	h, err := printer.Open(printerName)
	if err != nil {
		fmt.Println("Error opening printer:", err)
		return
	}
	defer h.Close()

	err = h.StartDocument("Receipt", "RAW")
	if err != nil {
		fmt.Println("Error starting document:", err)
		return
	}
	defer h.EndDocument()

	err = h.StartPage()
	if err != nil {
		fmt.Println("Error starting page:", err)
		return
	}
	defer h.EndPage()

	_, err = h.Write([]byte(receipt))
	if err != nil {
		fmt.Println("Error writing to printer:", err)
		return
	}
}
