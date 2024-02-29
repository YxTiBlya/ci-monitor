package service

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pingcap/errors"
	"go.uber.org/zap"
)

func (svc *Service) Monitor(ch chan string) {
	for {
		for _, repo := range svc.cfg.Repositories {
			updated, err := svc.check(repo)
			if err != nil {
				svc.log.Errorw("failed to check new commits",
					zap.String("repo", repo),
					zap.Error(err),
				)
			}

			if updated {
				ch <- repo
			}
		}

		time.Sleep(5 * time.Second) // TODO: mov to cfg?
	}
}

func (svc *Service) check(repo string) (bool, error) {
	var firstHash, secondHash string

	if _, err := os.Stat(repo); err != nil {
		return false, errors.Wrap(err, "directory is not exists")
	}

	if _, err := os.Stat(fmt.Sprintf("%s/pipeline.yaml", repo)); err != nil {
		return false, errors.Wrap(err, "pipeline.yaml is not exists")
	}

	cmd := exec.Command("git", "reset", "--hard")
	cmd.Dir = repo
	_, err := cmd.Output()
	if err != nil {
		return false, errors.Wrap(err, "failed to exec git reset")
	}

	cmd = exec.Command("git", "log", "-n1")
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		return false, errors.Wrap(err, "failed to exec git log")
	}

	firstHash = getHash(&out)

	cmd = exec.Command("git", "pull")
	cmd.Dir = repo
	_, err = cmd.Output()
	if err != nil {
		return false, errors.Wrap(err, "failed to exec git pull")
	}

	cmd = exec.Command("git", "log", "-n1")
	cmd.Dir = repo
	out, err = cmd.Output()
	if err != nil {
		return false, errors.Wrap(err, "failed to exec git log")
	}

	secondHash = getHash(&out)

	if firstHash != secondHash {
		return true, nil
	}

	return false, nil
}

func getHash(b *[]byte) string {
	return strings.Split(strings.Split(string(*b), " ")[1], "\n")[0]
}
