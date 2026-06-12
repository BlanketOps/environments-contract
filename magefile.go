//go:build mage

package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// Lint runs buf lint across the proto tree.
func Lint() error {
	step("Linting protos...")
	if err := sh.RunV("buf", "lint"); err != nil {
		fail("Lint failed")
		return err
	}
	success("Lint complete")
	return nil
}

// Gen runs buf generate for all targets (Go, C#, Java, TypeScript).
func Gen() error {
	step("Generating contracts...")
	if err := sh.RunV("buf", "generate"); err != nil {
		fail("Gen failed")
		return err
	}
	success("Gen complete")
	return nil
}

// Verify runs lint and generate.
func Verify() error {
	step("Verifying contracts...")
	if err := Lint(); err != nil {
		return err
	}
	return Gen()
}

// Bundle assembles the OCI contract bundle: buf image, proto sources,
// generated code, and metadata.
func Bundle() error {
	if err := Verify(); err != nil {
		return err
	}
	step("Building contract bundle...")
	if err := os.MkdirAll("bundle", 0755); err != nil {
		return err
	}
	// Canonical compiled contract — language-neutral binary descriptor.
	if err := sh.RunV("buf", "build", "-o", "bundle/image.binpb"); err != nil {
		fail("Bundle failed")
		return err
	}
	if err := tarGz("blanketops", "bundle/protos.tar.gz"); err != nil {
		fail("Bundle failed")
		return err
	}
	if err := tarGz("gen", "bundle/gen.tar.gz"); err != nil {
		fail("Bundle failed")
		return err
	}
	if err := writeMetadata("bundle/metadata.yaml"); err != nil {
		fail("Bundle failed")
		return err
	}
	success("Bundle complete")
	return nil
}

// Publish builds the bundle and pushes it to the OCI registry via ORAS.
func Publish() error {
	if err := Bundle(); err != nil {
		return err
	}
	ref, err := contractRef()
	if err != nil {
		fail("Publish failed")
		return err
	}
	step(fmt.Sprintf("Publishing %s...", ref))
	err = sh.RunV("oras", "push", ref,
		"bundle/image.binpb:application/vnd.bufbuild.buf.image.v1+binary",
		"bundle/protos.tar.gz:application/gzip",
		"bundle/gen.tar.gz:application/gzip",
		"bundle/metadata.yaml:application/yaml",
	)
	if err != nil {
		fail("Publish failed")
		return err
	}
	success(fmt.Sprintf("Published %s", ref))
	return nil
}

// Clean removes generated and bundled artifacts.
func Clean() error {
	step("Cleaning...")
	for _, dir := range []string{"gen", "bundle"} {
		if err := os.RemoveAll(dir); err != nil {
			fail("Clean failed")
			return err
		}
	}
	success("Clean complete")
	return nil
}

func tarGz(srcDir, dst string) error {
	if _, err := os.Stat(srcDir); err != nil {
		return fmt.Errorf("stat %s: %w", srcDir, err)
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	gz := gzip.NewWriter(out)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(path)
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(tw, f)
		return err
	})
}

func writeMetadata(dst string) error {
	meta := fmt.Sprintf(
		"version: %s\nrepository: %s\npublishedAt: %s\n",
		version(),
		envOr("GITHUB_REPOSITORY", "ntlaletsi70/blanketops-environments-contract"),
		time.Now().UTC().Format(time.RFC3339),
	)
	return os.WriteFile(dst, []byte(meta), 0644)
}

func contractRef() (string, error) {
	v := version()
	if v == "" {
		return "", fmt.Errorf("no version: set VERSION or GITHUB_REF_NAME, or tag the commit")
	}
	registry := envOr("CONTRACT_REGISTRY", "ghcr.io")
	owner := envOr("GITHUB_REPOSITORY_OWNER", "ntlaletsi70")
	name := envOr("CONTRACT_PACKAGE", "environments-contract")
	return fmt.Sprintf("%s/%s/%s:%s", registry, owner, name, v), nil
}

func version() string {
	if v := os.Getenv("VERSION"); v != "" {
		return v
	}
	if v := os.Getenv("GITHUB_REF_NAME"); v != "" {
		return v
	}
	if v, err := sh.Output("git", "describe", "--tags", "--abbrev=0"); err == nil {
		return strings.TrimSpace(v)
	}
	return ""
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
