package gofumpt

import (
	"github.com/rs/zerolog/log"
	"os/exec"
)

func Format(repoName string) error {
	// Format already downloaded code
	cmd := exec.Command("gofumpt", "-l", "-w", repoName)
	cmd.Dir = "./" + repoName
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Msg("Failed to format code")
		return err
	}
	log.Info().Msg(string(out))
	return nil
}
