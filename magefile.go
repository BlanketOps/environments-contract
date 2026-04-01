//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/sh"
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
	fmt.Println("▶ Cleaning generated artifacts...")

	dirs := []string{
		"gen",
		"bin",
		"dist",
	}

	for _, d := range dirs {
		if _, err := os.Stat(d); err == nil {
			fmt.Printf("  - removing %s\n", d)
			if err := os.RemoveAll(d); err != nil {
				return err
			}
		}
	}

	fmt.Println("✔ Clean complete")
	return nil
}

func GenerateCSharp() error {
	fmt.Println("▶ Generating C# SDKs per version...")

	targets := []struct {
		path string
		out  string
	}{
		{"blanketops/environments/v1alpha1", "gen/csharp/v1alpha1"},
		{"blanketops/environments/v1beta1", "gen/csharp/v1beta1"},
		{"blanketops/environments/v1", "gen/csharp/v1"},
	}

	for _, t := range targets {
		fmt.Println("▶", t.path)

		cmd := exec.Command("buf", "generate", t.path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}

		// move output
		os.RemoveAll(t.out)
		if err := os.Rename("gen/csharp", t.out); err != nil {
			return err
		}
	}

	fmt.Println("✔ C# generation complete")
	return nil
}

func Regen() error {
	step("Regenerating: clean → generate")
	if err := Clean(); err != nil {
		return err
	}
	return Generate()
}
