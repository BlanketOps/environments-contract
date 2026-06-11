//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/sh"
)

const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
	colorRed   = "\033[31m"
)

func step(msg string)    { fmt.Printf("%s▶ %s%s\n", colorCyan, msg, colorReset) }
func success(msg string) { fmt.Printf("%s✔ %s%s\n", colorGreen, msg, colorReset) }
func fail(msg string)    { fmt.Printf("%s✘ %s%s\n", colorRed, msg, colorReset) }

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
	if err := GenerateGo(); err != nil {
		return err
	}
	if err := GenerateCSharp(); err != nil {
		return err
	}
	if err := GenerateJava(); err != nil {
		return err
	}
	success("Generate complete")
	return nil
}

func GenerateGo() error {
	step("Generating Go...")
	if err := sh.RunV(buf, "generate", "--template", "buf.gen.yaml"); err != nil {
		fail("Go generate failed")
		return err
	}
	success("Go generation complete")
	return nil
}

func GenerateCSharp() error {
	step("Generating C#...")
	targets := []struct{ template, input, out string }{
		{"buf.gen.csharp.v1alpha1.yaml", "blanketops/environments/v1alpha1", "gen/csharp/environments/v1alpha1"},
		{"buf.gen.csharp.v1beta1.yaml", "blanketops/environments/v1beta1", "gen/csharp/environments/v1beta1"},
		{"buf.gen.csharp.v1.yaml", "blanketops/environments/v1", "gen/csharp/environments/v1"},
		{"buf.gen.csharp.events.v1alpha1.yaml", "blanketops/events/v1alpha1", "gen/csharp/events/v1alpha1"},
		{"buf.gen.csharp.events.v1beta1.yaml", "blanketops/events/v1beta1", "gen/csharp/events/v1beta1"},
		{"buf.gen.csharp.events.v1.yaml", "blanketops/events/v1", "gen/csharp/events/v1"},
		{"buf.gen.csharp.sources.v1alpha1.yaml", "blanketops/sources/v1alpha1", "gen/csharp/sources/v1alpha1"},
		{"buf.gen.csharp.sources.v1beta1.yaml", "blanketops/sources/v1beta1", "gen/csharp/sources/v1beta1"},
		{"buf.gen.csharp.sources.v1.yaml", "blanketops/sources/v1", "gen/csharp/sources/v1"},
	}
	for _, t := range targets {
		fmt.Printf("  %s▶%s %s\n", colorCyan, colorReset, t.input)
		if err := os.MkdirAll(t.out, 0755); err != nil {
			return err
		}
		cmd := exec.Command(buf, "generate", "--template", t.template, "--path", t.input)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fail(fmt.Sprintf("C# generate failed for %s", t.input))
			return err
		}
	}
	success("C# generation complete")
	return nil
}

func GenerateJava() error {
	step("Generating Java...")
	targets := []struct{ template, input, out string }{
		{"buf.gen.java.v1alpha1.yaml", "blanketops/environments/v1alpha1", "gen/java/environments/v1alpha1"},
		{"buf.gen.java.v1beta1.yaml", "blanketops/environments/v1beta1", "gen/java/environments/v1beta1"},
		{"buf.gen.java.v1.yaml", "blanketops/environments/v1", "gen/java/environments/v1"},
		{"buf.gen.java.events.v1alpha1.yaml", "blanketops/events/v1alpha1", "gen/java/events/v1alpha1"},
		{"buf.gen.java.events.v1beta1.yaml", "blanketops/events/v1beta1", "gen/java/events/v1beta1"},
		{"buf.gen.java.events.v1.yaml", "blanketops/events/v1", "gen/java/events/v1"},
		{"buf.gen.java.sources.v1alpha1.yaml", "blanketops/sources/v1alpha1", "gen/java/sources/v1alpha1"},
		{"buf.gen.java.sources.v1beta1.yaml", "blanketops/sources/v1beta1", "gen/java/sources/v1beta1"},
		{"buf.gen.java.sources.v1.yaml", "blanketops/sources/v1", "gen/java/sources/v1"},
	}
	for _, t := range targets {
		fmt.Printf("  %s▶%s %s\n", colorCyan, colorReset, t.input)
		if err := os.MkdirAll(t.out, 0755); err != nil {
			return err
		}
		cmd := exec.Command(buf, "generate", "--template", t.template, "--path", t.input)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fail(fmt.Sprintf("Java generate failed for %s", t.input))
			return err
		}
	}
	success("Java generation complete")
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
	step("Cleaning generated artifacts...")
	for _, d := range []string{"gen", "bin", "dist"} {
		if _, err := os.Stat(d); err == nil {
			fmt.Printf("  - removing %s\n", d)
			if err := os.RemoveAll(d); err != nil {
				return err
			}
		}
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

func ScaffoldTemplates() error {
	step("Scaffolding buf templates...")
	type tmpl struct {
		file string
		out  string
		lang string
	}
	templates := []tmpl{
		{"buf.gen.csharp.v1alpha1.yaml", "gen/csharp/environments/v1alpha1", "csharp"},
		{"buf.gen.csharp.v1beta1.yaml", "gen/csharp/environments/v1beta1", "csharp"},
		{"buf.gen.csharp.v1.yaml", "gen/csharp/environments/v1", "csharp"},
		{"buf.gen.csharp.events.v1alpha1.yaml", "gen/csharp/events/v1alpha1", "csharp"},
		{"buf.gen.csharp.events.v1beta1.yaml", "gen/csharp/events/v1beta1", "csharp"},
		{"buf.gen.csharp.events.v1.yaml", "gen/csharp/events/v1", "csharp"},
		{"buf.gen.csharp.sources.v1alpha1.yaml", "gen/csharp/sources/v1alpha1", "csharp"},
		{"buf.gen.csharp.sources.v1beta1.yaml", "gen/csharp/sources/v1beta1", "csharp"},
		{"buf.gen.csharp.sources.v1.yaml", "gen/csharp/sources/v1", "csharp"},
		{"buf.gen.java.v1alpha1.yaml", "gen/java/environments/v1alpha1", "java"},
		{"buf.gen.java.v1beta1.yaml", "gen/java/environments/v1beta1", "java"},
		{"buf.gen.java.v1.yaml", "gen/java/environments/v1", "java"},
		{"buf.gen.java.events.v1alpha1.yaml", "gen/java/events/v1alpha1", "java"},
		{"buf.gen.java.events.v1beta1.yaml", "gen/java/events/v1beta1", "java"},
		{"buf.gen.java.events.v1.yaml", "gen/java/events/v1", "java"},
		{"buf.gen.java.sources.v1alpha1.yaml", "gen/java/sources/v1alpha1", "java"},
		{"buf.gen.java.sources.v1beta1.yaml", "gen/java/sources/v1beta1", "java"},
		{"buf.gen.java.sources.v1.yaml", "gen/java/sources/v1", "java"},
	}
	for _, t := range templates {
		if _, err := os.Stat(t.file); err == nil {
			fmt.Printf("  [exists]  %s\n", t.file)
			continue
		}
		var pb, grpc string
		if t.lang == "csharp" {
			pb = "buf.build/protocolbuffers/csharp"
			grpc = "buf.build/grpc/csharp"
		} else {
			pb = "buf.build/protocolbuffers/java"
			grpc = "buf.build/grpc/java"
		}
		content := fmt.Sprintf("version: v1\nplugins:\n- plugin: %s\n  out: %s\n- plugin: %s\n  out: %s\n",
			pb, t.out, grpc, t.out)
		if err := os.WriteFile(t.file, []byte(content), 0644); err != nil {
			return err
		}
		fmt.Printf("  [created] %s\n", t.file)
	}
	success("Templates scaffolded — commit these files")
	return nil
}
