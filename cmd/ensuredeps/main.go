package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/golangci/golangci-api/internal/shared/logutil"
	"github.com/golangci/golangci-api/pkg/goenvbuild/command"
	"github.com/golangci/golangci-api/pkg/goenvbuild/result"

	"github.com/golangci/golangci-api/pkg/goenvbuild/ensuredeps"
)

func main() {
	repoName := flag.String("repo", "", "repo name or path")
	flag.Parse()
	if *repoName == "" {
		log.Fatalf("Repo name must be set: use --repo")
	}

	loggger := logutil.NewStderrLog("")
	resLog := result.NewLog(log.New(os.Stdout, "", 0))
	runner := command.NewStreamingRunner(resLog)

	r := ensuredeps.NewRunner(loggger, runner)
	ret := r.Run(context.Background(), *repoName)
	if err := json.NewEncoder(os.Stdout).Encode(ret); err != nil {
		log.Fatalf("Failed to JSON output result: %s", err)
	}
}
