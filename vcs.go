package main

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"

	"github.com/motemen/ghq/cmdutil"
)

// A VCSBackend represents a VCS backend.
type VCSBackend struct {
	// Clones a remote repository to local path.
	Clone func(*url.URL, string, bool) error
	// Updates a cloned local repository.
	Update func(string) error
}

var GitBackend = &VCSBackend{
	Clone: func(remote *url.URL, local string, shallow bool) error {
		dir, _ := filepath.Split(local)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		args := []string{"clone"}
		if shallow {
			args = append(args, "--depth", "1")
		}
		args = append(args, remote.String(), local)

		return cmdutil.Run("git", args...)
	},
	Update: func(local string) error {
		return cmdutil.RunInDir(local, "git", "pull", "--ff-only")
	},
}

var SubversionBackend = &VCSBackend{
	Clone: func(remote *url.URL, local string, shallow bool) error {
		dir, _ := filepath.Split(local)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		args := []string{"checkout"}
		if shallow {
			args = append(args, "--depth", "1")
		}
		args = append(args, remote.String(), local)

		return cmdutil.Run("svn", args...)
	},
	Update: func(local string) error {
		return cmdutil.RunInDir(local, "svn", "update")
	},
}

var GitsvnBackend = &VCSBackend{
	// git-svn seems not supporting shallow clone currently.
	Clone: func(remote *url.URL, local string, ignoredShallow bool) error {
		dir, _ := filepath.Split(local)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		return cmdutil.Run("git", "svn", "clone", remote.String(), local)
	},
	Update: func(local string) error {
		return cmdutil.RunInDir(local, "git", "svn", "rebase")
	},
}

var MercurialBackend = &VCSBackend{
	// Mercurial seems not supporting shallow clone currently.
	Clone: func(remote *url.URL, local string, ignoredShallow bool) error {
		dir, _ := filepath.Split(local)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		return cmdutil.Run("hg", "clone", remote.String(), local)
	},
	Update: func(local string) error {
		return cmdutil.RunInDir(local, "hg", "pull", "--update")
	},
}

var DarcsBackend = &VCSBackend{
	Clone: func(remote *url.URL, local string, shallow bool) error {
		dir, _ := filepath.Split(local)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		args := []string{"get"}
		if shallow {
			args = append(args, "--lazy")
		}
		args = append(args, remote.String(), local)

		return cmdutil.Run("darcs", args...)
	},
	Update: func(local string) error {
		return cmdutil.RunInDir(local, "darcs", "pull")
	},
}

var cvsDummyBackend = &VCSBackend{
	Clone: func(remote *url.URL, local string, ignoredShallow bool) error {
		return errors.New("CVS clone is not supported")
	},
	Update: func(local string) error {
		return errors.New("CVS update is not supported")
	},
}

const fossilRepoName = ".fossil" // same as Go

var FossilBackend = &VCSBackend{
	Clone: func(remote *url.URL, local string, shallow bool) error {
		dir, _ := filepath.Split(local)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}

		err = cmdutil.Run("fossil", "clone", remote.String(), filepath.Join(dir, fossilRepoName))
		if err != nil {
			return err
		}

		err = os.Chdir(dir)
		if err != nil {
			return err
		}

		return cmdutil.Run("fossile", "open", fossilRepoName)
	},
	Update: func(local string) error {
		return cmdutil.RunInDir(local, "fossil", "update")
	},
}

var vcsRegistry = map[string]*VCSBackend{
	"git":        GitBackend,
	"github":     GitBackend,
	"svn":        SubversionBackend,
	"subversion": SubversionBackend,
	"git-svn":    GitsvnBackend,
	"hg":         MercurialBackend,
	"mercurial":  MercurialBackend,
	"darcs":      DarcsBackend,
	"fossil":     FossilBackend,
}
