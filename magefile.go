//go:build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
)

func step(msg string) {
	fmt.Printf("%s▶ %s%s\n", colorCyan, msg, colorReset)
}

func success(msg string) {
	fmt.Printf("%s✔ %s%s\n", colorGreen, msg, colorReset)
}

func fail(msg string) {
	fmt.Printf("%s✘ %s%s\n", colorRed, msg, colorReset)
}

var buf = func() string {
	if b := os.Getenv("BUF"); b != "" {
		return b
	}
	return "buf"
}()

func All() error {
	step("Running all: lint → generate")
	if err := Lint(); err != nil {
		return err
	}
	return Generate()
}

func Generate() error {
	step("Running buf generate...")
	if err := sh.RunV(buf, "generate"); err != nil {
		fail("Generate failed")
		return err
	}
	success("Generate complete")
	return nil
}

func Lint() error {
	step("Running buf lint...")
	if err := sh.RunV(buf, "lint"); err != nil {
		fail("Lint failed")
		return err
	}
	success("Lint passed")
	return nil
}

func Format() error {
	step("Running buf format...")
	if err := sh.RunV(buf, "format", "-w"); err != nil {
		fail("Format failed")
		return err
	}
	success("Format complete")
	return nil
}

func Clean() error {
	step("Cleaning generated .pb.go files...")
	if err := sh.RunV("find", "blanketops", "-name", "*.pb.go", "-delete"); err != nil {
		fail("Clean failed")
		return err
	}
	success("Clean complete")
	return nil
}

func Regen() error {
	step("Regenerating: clean → generate")
	if err := Clean(); err != nil {
		return err
	}
	return Generate()
}