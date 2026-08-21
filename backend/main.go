package main

import (
	"fmt"
	"os"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/cmd"
)

func main() {
	if err := cmd.Root().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
