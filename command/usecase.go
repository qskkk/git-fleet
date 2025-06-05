package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
	"github.com/qskkk/git-fleet/config"
)

func ExecuteAll(args []string) (string, error) {
	out, err := ExecuteHandled(args)
	if err != nil {
		err = fmt.Errorf("error executing handled command: %w", err)
		return "", err
	}
	if out != "" {
		return out, nil
	}

	if len(args) < 2 {
		return "", nil
	}

	repos, ok := config.Cfg.Groups[args[1]]
	if !ok {
		log.Errorf("Error: group '%s' not found in configuration", args[1])
		os.Exit(1)
	}

	for _, repo := range repos {
		out, err := Execute(repo, args[2:])
		if err != nil {
			log.Errorf("Error executing command in '%s': %v", repo, err)
		}

		log.Info(out)
	}
	return out, nil
}

func Execute(repoName string, command []string) (string, error) {
	rc, ok := config.Cfg.Repositories[repoName]
	if !ok {
		err := fmt.Errorf("error: repository '%s' not found in configuration", repoName)
		return "", err
	}

	if info, err := os.Stat(rc.Path); err != nil || !info.IsDir() {
		err := fmt.Errorf("error: '%s' is not a valid directory: %w", rc.Path, err)
		return "", err
	}

	var out bytes.Buffer
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = rc.Path
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("error executing command in '%s': %w", rc.Path, err)
		return "", err
	}

	return fmt.Sprintf("Command output for '%s':\n%s", rc.Path, out.String()), nil
}

func ExecuteHandled(args []string) (string, error) {
	if len(args) < 2 {
		return "", nil
	}

	fmt.Printf("args: %v\n", args)

	if _, ok := GlobalHandled[args[1]]; ok {
		out, err := GlobalHandled[args[1]]()
		if err != nil {
			err = fmt.Errorf("error executing global command '%s': %w", args[1], err)
			return "", err
		}

		return out, nil
	}

	if len(args) < 3 {
		return "", nil
	}

	if _, ok := Handled[args[2]]; ok {
		out, err := Handled[args[2]]()
		if err != nil {
			err = fmt.Errorf("error executing command '%s': %w", args[2], err)
			return "", err
		}
		return out, nil
	}

	return "", nil
}

func ExecuteStatus() (string, error) {
	var result bytes.Buffer

	for repoName, rc := range config.Cfg.Repositories {
		if info, err := os.Stat(rc.Path); err != nil || !info.IsDir() {
			result.WriteString(fmt.Sprintf("Repository '%s': invalid directory '%s'\n", repoName, rc.Path))
			continue
		}

		cmd := exec.Command("git", "status", "--porcelain")
		cmd.Dir = rc.Path
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			result.WriteString(fmt.Sprintf("Repository '%s': error running git status: %v\n", repoName, err))
			continue
		}

		created, edited, deleted := 0, 0, 0
		for _, line := range bytes.Split(out.Bytes(), []byte("\n")) {
			if len(line) < 2 {
				continue
			}
			switch line[0] {
			case 'A', '?': // Added or untracked files
				created++
			case 'M': // Modified files
				edited++
			case 'D': // Deleted files
				deleted++
			}
		}

		result.WriteString(fmt.Sprintf("Repository '%s': %d created, %d edited, %d deleted files\n", repoName, created, edited, deleted))
	}

	return result.String(), nil
}
