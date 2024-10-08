package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

//go:embed bin/rqlited
var binaryFile []byte

func main() {
	fmt.Println("Hello")
	binaryPath := "/tmp"
	binaryName := "rqlited"

	fullBinaryPath := fmt.Sprintf("%s"+"/"+"%s", binaryPath, binaryName)

	// Write the embedded binary to a file (optional)
	err := os.WriteFile(fullBinaryPath, binaryFile, 0777)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Printf("Binary file has been written: %s\n", fullBinaryPath)

	// Create a context for the children
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure the context is cleaned up

	// Create an errgroup to manage concurrent processes
	// var g errgroup.Group
	g, _ := errgroup.WithContext(ctx)

	// Example processes to run concurrently
	processes := []struct {
		binaryPath string
		args       []string
	}{
		{fullBinaryPath, []string{"~/node.1"}},
		// {"/path/to/binary2", []string{"arg1", "arg2"}},
	}

	// Loop over the processes and start them concurrently
	for _, process := range processes {
		// Capture the loop variables to avoid closure issues
		binaryPath := process.binaryPath
		args := process.args

		// Launch each process in a separate goroutine
		g.Go(func() error {
			return runProcessWithContext(ctx, binaryPath, args...)
		})
	}

	// Wait for all processes to finish
	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("All processes finished successfully.")
	}

}

func runProcessWithContext(ctx context.Context, binaryPath string, args ...string) error {
	// Create a command using the context for cancellation or timeout
	cmd := exec.CommandContext(ctx, binaryPath, args...)
	cmd.Stdout = os.Stdout // Redirect the process output to the terminal
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin // Pass input from the terminal to the process

	// Run the process and wait for it to complete
	err := cmd.Run()

	// 	// Start the process (non-blocking)
	// 	err := cmd.Start()
	// 	if err != nil {
	// 		fmt.Printf("Error starting process: %v\n", err)
	// 		return
	// 	}

	// 	// You can do other tasks while the process runs in the background
	// 	fmt.Println("Process started. PID:", cmd.Process.Pid)

	// 	// Wait for the process to finish (optional)
	// 	err = cmd.Wait()
	// 	if err != nil {
	// 		fmt.Printf("Process finished with error: %v\n", err)
	// 	}

	// Check for context timeout or cancellation
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("process %s timed out", binaryPath)
	}

	// Check for other errors
	if err != nil {
		return fmt.Errorf("error running process %s: %v", binaryPath, err)
	}

	fmt.Printf("Process %s finished successfully\n", binaryPath)
	return nil
}
