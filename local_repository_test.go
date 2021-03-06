package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

import (
	"net/url"
	"os"
)

func TestNewLocalRepository(t *testing.T) {
	RegisterTestingT(t)

	_localRepositoryRoots = []string{"/repos"}

	r, err := LocalRepositoryFromFullPath("/repos/github.com/motemen/ghq")
	Expect(err).To(BeNil())
	Expect(r.NonHostPath()).To(Equal("motemen/ghq"))
	Expect(r.Subpaths()).To(Equal([]string{"ghq", "motemen/ghq", "github.com/motemen/ghq"}))

	r, err = LocalRepositoryFromFullPath("/repos/stash.com/scm/motemen/ghq")
	Expect(err).To(BeNil())
	Expect(r.NonHostPath()).To(Equal("scm/motemen/ghq"))
	Expect(r.Subpaths()).To(Equal([]string{"ghq", "motemen/ghq", "scm/motemen/ghq", "stash.com/scm/motemen/ghq"}))

	githubURL, _ := url.Parse("ssh://git@github.com/motemen/ghq.git")
	r, err = LocalRepositoryFromURL(githubURL)
	Expect(err).To(BeNil())
	Expect(r.FullPath).To(Equal("/repos/github.com/motemen/ghq"))

	stashURL, _ := url.Parse("ssh://git@stash.com/scm/motemen/ghq")
	r, _ = LocalRepositoryFromURL(stashURL)
	Expect(r.FullPath).To(Equal("/repos/stash.com/scm/motemen/ghq"))

	svnSourceforgeURL, _ := url.Parse("http://svn.code.sf.net/p/ghq/code/trunk")
	r, _ = LocalRepositoryFromURL(svnSourceforgeURL)
	Expect(r.FullPath).To(Equal("/repos/svn.code.sf.net/p/ghq/code/trunk"))

	gitSourceforgeURL, _ := url.Parse("http://git.code.sf.net/p/ghq/code")
	r, _ = LocalRepositoryFromURL(gitSourceforgeURL)
	Expect(r.FullPath).To(Equal("/repos/git.code.sf.net/p/ghq/code"))

	svnSourceforgeJpURL, _ := url.Parse("http://scm.sourceforge.jp/svnroot/ghq/")
	r, _ = LocalRepositoryFromURL(svnSourceforgeJpURL)
	Expect(r.FullPath).To(Equal("/repos/scm.sourceforge.jp/svnroot/ghq"))

	gitSourceforgeJpURL, _ := url.Parse("http://scm.sourceforge.jp/gitroot/ghq/ghq.git")
	r, _ = LocalRepositoryFromURL(gitSourceforgeJpURL)
	Expect(r.FullPath).To(Equal("/repos/scm.sourceforge.jp/gitroot/ghq/ghq"))

	svnAssemblaURL, _ := url.Parse("https://subversion.assembla.com/svn/ghq/")
	r, _ = LocalRepositoryFromURL(svnAssemblaURL)
	Expect(r.FullPath).To(Equal("/repos/subversion.assembla.com/svn/ghq"))

	gitAssemblaURL, _ := url.Parse("https://git.assembla.com/ghq.git")
	r, _ = LocalRepositoryFromURL(gitAssemblaURL)
	Expect(r.FullPath).To(Equal("/repos/git.assembla.com/ghq"))
}

func TestLocalRepositoryRoots(t *testing.T) {
	RegisterTestingT(t)

	defer func(orig string) { os.Setenv("GHQ_ROOT", orig) }(os.Getenv("GHQ_ROOT"))

	_localRepositoryRoots = nil
	os.Setenv("GHQ_ROOT", "/path/to/ghqroot")
	Expect(localRepositoryRoots()).To(Equal([]string{"/path/to/ghqroot"}))

	_localRepositoryRoots = nil
	os.Setenv("GHQ_ROOT", "/path/to/ghqroot1"+string(os.PathListSeparator)+"/path/to/ghqroot2")
	Expect(localRepositoryRoots()).To(Equal([]string{"/path/to/ghqroot1", "/path/to/ghqroot2"}))
}

// https://gist.github.com/kyanny/c231f48e5d08b98ff2c3
func TestList_Symlink(t *testing.T) {
	RegisterTestingT(t)

	root, err := ioutil.TempDir("", "")
	Expect(err).To(BeNil())

	symDir, err := ioutil.TempDir("", "")
	Expect(err).To(BeNil())

	_localRepositoryRoots = []string{root}

	err = os.MkdirAll(filepath.Join(root, "github.com", "atom", "atom", ".git"), 0777)
	Expect(err).To(BeNil())

	err = os.MkdirAll(filepath.Join(root, "github.com", "zabbix", "zabbix", ".git"), 0777)
	Expect(err).To(BeNil())

	err = os.Symlink(symDir, filepath.Join(root, "github.com", "ghq"))
	Expect(err).To(BeNil())

	paths := []string{}
	walkLocalRepositories(func(repo *LocalRepository) {
		paths = append(paths, repo.RelPath)
	})

	Expect(paths).To(HaveLen(2))
}
