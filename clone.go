package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
)

func CloneProject(tenant, repoName, repoUrl, id string) error {
	// Check if the tenant folder exists
	if _, err := os.Stat("./" + tenant); os.IsNotExist(err) {
		err = os.Mkdir("./"+tenant, os.ModePerm)
		if err != nil {
			log.Error().Err(err).Msg("Error on creating folder")
			return err
		}
	}

	// Check if the project is already cloned, if not clone with depth -1
	if _, err := os.Stat("./" + repoName); os.IsNotExist(err) {
		run := exec.Command("git", "clone", "--depth", "1", repoUrl)
		run.Dir = "./" + tenant
		out, err := run.CombinedOutput()
		if err != nil {
			log.Error().Err(err).Msg("Error cloning project")
			return err
		}
		log.Info().Msg(string(out))
	}
	branch := randomString(10)
	// Fetch the pull request
	run := exec.Command("git", "fetch", "origin", "pull/"+id+"/head:"+branch)
	run.Dir = "./" + repoName
	out, err := run.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Msg("Error fetch project")
		return err
	}
	log.Info().Msg(string(out))

	// Checkout to the branch
	run = exec.Command("git", "checkout", branch)
	run.Dir = "./" + repoName
	out, err = run.CombinedOutput()
	if err != nil {
		log.Error().Err(err).Msg("Error checkout project")
		return err
	}
	log.Info().Msg(string(out))
	return nil
}
