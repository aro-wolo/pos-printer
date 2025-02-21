//go:build darwin
// +build darwin

package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func printToDefaultPrinter(receipt string) {
	cmd := exec.Command("lp")
	cmd.Stdin = bytes.NewBufferString(receipt)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error printing receipt:", err)
	}
}
