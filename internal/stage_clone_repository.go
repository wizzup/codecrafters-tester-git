package internal

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"time"
)

type TestFile struct {
	path     string
	contents string
}

type TestRepo struct {
	url            string
	exampleCommits []string
	exampleFiles   []TestFile
}

func (r TestRepo) randomCommit() string {
	return r.exampleCommits[rand.Intn(len(r.exampleCommits))]
}

func (r TestRepo) randomFile() TestFile {
	return r.exampleFiles[rand.Intn(len(r.exampleFiles))]
}

var testRepos []TestRepo = []TestRepo{
	TestRepo{
		url: "https://github.com/codecrafters-io/git-sample-1",
		exampleCommits: []string{
			"3b0466d22854e57bf9ad3ccf82008a2d3f199550",
		},
		exampleFiles: []TestFile{
			TestFile{
				path:     "scooby/dooby/doo",
				contents: "dooby yikes dumpty scooby monkey donkey horsey humpty vanilla doo",
			},
		},
	},
	TestRepo{
		url: "https://github.com/codecrafters-io/git-sample-2",
		exampleCommits: []string{
			"b521b9179412d90a893bc36f33f5dcfd987105ef",
		},
		exampleFiles: []TestFile{
			TestFile{
				path:     "humpty/vanilla/yikes",
				contents: "scooby yikes dooby",
			},
		},
	},
	TestRepo{
		url: "https://github.com/codecrafters-io/git-sample-3",
		exampleCommits: []string{
			"23f0bc3b5c7c3108e41c448f01a3db31e7064bbb",
			"b521b9179412d90a893bc36f33f5dcfd987105ef",
		},
		exampleFiles: []TestFile{
			TestFile{
				path:     "donkey/donkey/monkey",
				contents: "monkey humpty doo scooby dumpty donkey vanilla horsey dooby",
			},
		},
	},
}

func randomRepo() TestRepo {
	rand.Seed(time.Now().UnixNano())
	return testRepos[rand.Intn(3)]
}

func testCloneRepository(executable *Executable, logger *customLogger) error {
	tempDir, err := ioutil.TempDir("", "worktree")
	if err != nil {
		return err
	}

	executable.WorkingDir = tempDir

	testRepo := randomRepo()

	logger.Debugf("Running ./your_git.sh clone %s <testDir>", testRepo.url)
	result, err := executable.Run("clone", testRepo.url, "test_dir")
	if err != nil {
		return err
	}

	if err = assertExitCode(result, 0); err != nil {
		return err
	}

	repoDir := path.Join(tempDir, "test_dir")

	// Test a commit
	commit_sha := testRepo.randomCommit()

	logger.Debugf("Running git cat-file commit %s", commit_sha)
	result, err = runGitCmdUnsafe(repoDir, "cat-file", "commit", commit_sha)
	if err != nil {
		return err
	}

	if err = assertExitCode(result, 0); err != nil {
		return err
	}

	if err = assertStdoutContains(result, "author Paul Kuruvilla"); err != nil {
		return err
	}
	logger.Successf("Commit contents verified")

	// Test a commit
	testFile := testRepo.randomFile()

	logger.Debugf("Reading contents of a sample file")
	bytes, err := ioutil.ReadFile(path.Join(repoDir, testFile.path))
	if err != nil {
		return err
	}

	expected := testFile.contents
	actual := string(bytes)

	if expected != actual {
		return fmt.Errorf("Expected '%s' as file contents, got: '%s'", expected, actual)
	}
	logger.Successf("File contents verified")

	return nil
}