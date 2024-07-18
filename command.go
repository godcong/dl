package dl

import (
	"log"
	"os"
	"os/exec"
)

// ExecGoImports Executes the goimports command on the given file
func ExecGoImports(dir, moduleName, name string) error {
	localPath, err := exec.LookPath("goimports")
	if err != nil {
		if err := ExecGoInstall(dir, "golang.org/x/tools/cmd/goimports@latest"); err != nil {
			log.Println("not found goimports command, try to start with go run mode...")
		}
	}
	args := []string{"-local", dir, "-w", name}
	if localPath == "" {
		localPath = "go"
		args = append([]string{"run", "-mod=mod", "golang.org/x/tools/cmd/goimports"}, args...)
	}
	cmd := exec.Command(localPath, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	log.Println("Command description", "#run", cmd.String())
	return cmd.Run()
}

func ExecGoInstall(dir, path string) error {
	localPath, err := exec.LookPath("go")
	if err != nil {
		log.Println("not found go command, please install go first")
		return nil
	}

	cmd := exec.Command(localPath, "install", path)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	log.Println("Command description", "#run", cmd.String())
	return cmd.Run()
}
