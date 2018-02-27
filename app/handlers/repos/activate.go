package repos

import (
	"strings"

	"github.com/golangci/golangci-api/app/internal/auth/user"
	"github.com/golangci/golangci-api/app/internal/repos"
	"github.com/golangci/golangci-api/app/models"
	"github.com/golangci/golangci-api/app/returntypes"
	"github.com/golangci/golib/server/context"
	"github.com/golangci/golib/server/handlers/herrors"
	"github.com/golangci/golib/server/handlers/manager"
)

func changeRepo(ctx context.C) error {
	ga, err := user.GetGithubAuth(&ctx)
	if err != nil {
		return herrors.New(err, "can't get github auth")
	}

	repoOwner := ctx.URLVar("repoOwner")
	if !strings.EqualFold(ga.Login, repoOwner) {
		return herrors.New403Errorf("invalid repo owner: %q != %q", ga.Login, repoOwner)
	}

	repoName := ctx.URLVar("repoName")

	var gr *models.GithubRepo
	var activate = ctx.R.Method == "PUT"
	switch ctx.R.Method {
	case "PUT":
		gr, err = repos.ActivateRepo(&ctx, ga, repoOwner, repoName)
		if err != nil {
			return herrors.New(err, "can't activate repo")
		}
	case "DELETE":
		gr, err = repos.DeactivateRepo(&ctx, repoOwner, repoName)
		if err != nil {
			return herrors.New(err, "can't deactivate repo")
		}
	default:
		return herrors.New404Errorf("unallowed method")
	}

	ri := returntypes.RepoInfo{
		Name:        gr.Name,
		IsActivated: activate,
		HookID:      gr.HookID,
	}
	ctx.ReturnJSON(map[string]returntypes.RepoInfo{
		"repo": ri,
	})
	return nil
}

func init() {
	manager.Register("/v1/repos/{repoOwner}/{repoName}", changeRepo)
}
