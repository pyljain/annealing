package main

import (
	"annealing/pkg/config"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sync/errgroup"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "annealing.yaml", "Annealing YAML config file to identify services to build")
	flag.Parse()

	if configPath == "" {
		fmt.Fprintf(os.Stderr, "An Annealing config must be provided")
		os.Exit(-1)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "The Annealing config could not be loaded: %s", err)
		os.Exit(-1)
	}

	// Run Git diff
	output, err := exec.Command("git", "diff", "--name-only").CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running Git diff: %s", output)
		os.Exit(-1)
	}

	changedServices := map[string]*config.Service{}

	filesChanged := strings.Split(string(output), "\n")
	for _, file := range filesChanged {
		for _, svc := range cfg.Spec.Services {
			if strings.HasPrefix(file, svc.Path) {
				// cmd := exec.Command("/bin/bash", "-c", svc.Command)
				// cmd.Dir = svc.Path
				// cmd.Run()
				changedServices[svc.Path] = &svc
				break
			}
		}
	}

	ctx := context.Background()
	eg, egctx := errgroup.WithContext(ctx)
	for name, cs := range changedServices {
		fmt.Printf("Running build for service %s\n", name)
		eg.Go(func() error {
			for _, c := range cs.Commands {
				cmd := exec.CommandContext(egctx, "/bin/bash", "-c", c)
				cmd.Dir = cs.Path
				err = cmd.Run()
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running build command: %s", err)
		os.Exit(-1)
	}

}

/*
1. Run git diff with only file
2. Compare the containing directory to the directories in the config manifest
3. Run the respective build command in the correct directory
*/
