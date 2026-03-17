//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
	"os"
)

var buf = func() string {
	if b := os.Getenv("BUF"); b != "" {
		return b
	}
	return "buf"
}()

func All() error {
	if err := Lint(); err != nil {
		return err
	}
	return Generate()
}

func Generate() error {
	return sh.Run(buf, "generate")
}

func Lint() error {
	return sh.Run(buf, "lint")
}

func Format() error {
	return sh.Run(buf, "format", "-w")
}

func Clean() error {
	return sh.RunV("find", "blanketops", "-name", "*.pb.go", "-delete")
}

func Regen() error {
	if err := Clean(); err != nil {
		return err
	}
	return Generate()
}