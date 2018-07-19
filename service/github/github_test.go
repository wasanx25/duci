package github_test

import (
	"fmt"
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/duck8823/duci/service/github"
	"gopkg.in/h2non/gock.v1"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io/ioutil"
	"testing"
)

type MockRepo struct {
	FullName string
	SSHURL   string
}

func (r *MockRepo) GetFullName() string {
	return r.FullName
}

func (r *MockRepo) GetSSHURL() string {
	return r.SSHURL
}

func TestService_GetPullRequest(t *testing.T) {
	// setup
	s, err := github.NewWithEnv()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	t.Run("when github server returns status ok", func(t *testing.T) {
		// given
		repo := &MockRepo{
			FullName: "duck8823/duci",
			SSHURL:   "git@github.com:duck8823/duci.git",
		}
		num := 5
		var id int64 = 19

		// and
		gock.New("https://api.github.com").
			Get(fmt.Sprintf("/repos/%s/pulls/%d", repo.FullName, num)).
			Reply(200).
			JSON(&github.PullRequest{
				ID: &id,
			})

		// when
		pr, err := s.GetPullRequest(context.New("test/task"), repo, num)

		// then
		if err != nil {
			t.Fatalf("error occured. %+v", err)
		}

		if pr.GetID() != id {
			t.Errorf("id must be equal %+v, but got %+v. \npr=%+v", id, pr.GetID(), pr)
		}
	})

	t.Run("when github server returns status not found", func(t *testing.T) {
		// given
		repo := &MockRepo{
			FullName: "duck8823/duci",
			SSHURL:   "git@github.com:duck8823/duci.git",
		}
		num := 5

		// and
		gock.New("https://api.github.com").
			Get(fmt.Sprintf("/repos/%s/pulls/%d", repo.FullName, num)).
			Reply(404)

		// when
		pr, err := s.GetPullRequest(context.New("test/task"), repo, num)

		// then
		if err == nil {
			t.Error("error must occur")
		}

		if pr != nil {
			t.Errorf("pr must nil, but got %+v", pr)
		}
	})

	t.Run("with invalid repository", func(t *testing.T) {
		// given
		repo := &MockRepo{
			FullName: "",
		}
		num := 5

		// expect
		if _, err := s.GetPullRequest(context.New("test/task"), repo, num); err == nil {
			t.Error("errot must occred. but got nil")
		}
	})
}

func TestService_CreateCommitStatus(t *testing.T) {
	// setup
	s, err := github.NewWithEnv()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	t.Run("when github server returns status ok", func(t *testing.T) {
		// given
		repo := &MockRepo{
			FullName: "duck8823/duci",
		}

		// and
		gock.New("https://api.github.com").
			Post(fmt.Sprintf("/repos/%s/statuses/%s", repo.FullName, "0000000000000000000000000000000000000000")).
			Reply(200)

		// expect
		if err := s.CreateCommitStatus(context.New("test/task"), repo, plumbing.Hash{}, github.SUCCESS, ""); err != nil {
			t.Errorf("error must not occured: but got %+v", err)
		}

		// cleanup
		gock.Clean()
	})

	t.Run("when github server returns status not found", func(t *testing.T) {
		// given
		repo := &MockRepo{
			FullName: "duck8823/duci",
		}

		// and
		gock.New("https://api.github.com").
			Post(fmt.Sprintf("/repos/%s/statuses/%s", repo.FullName, "0000000000000000000000000000000000000000")).
			Reply(404)

		// expect
		if err := s.CreateCommitStatus(context.New("test/task"), repo, plumbing.Hash{}, github.SUCCESS, ""); err == nil {
			t.Error("errot must occred. but got nil")
		}

		// cleanup
		gock.Clean()
	})

	t.Run("with invalid repository", func(t *testing.T) {
		// given
		repo := &MockRepo{
			FullName: "",
		}

		// expect
		if err := s.CreateCommitStatus(context.New("test/task"), repo, plumbing.Hash{}, github.SUCCESS, ""); err == nil {
			t.Error("errot must occred. but got nil")
		}
	})

	t.Run("with long description", func(t *testing.T) {
		// given
		repo := &MockRepo{
			FullName: "duck8823/duci",
		}

		// and
		taskName := "test/task"
		description := "123456789012345678901234567890123456789012345678901234567890"
		malformedDescription := "1234567890123456789012345678901234567890123456..."
		state := github.SUCCESS

		gock.New("https://api.github.com").
			Post(fmt.Sprintf("/repos/%s/statuses/%s", repo.FullName, "0000000000000000000000000000000000000000")).
			MatchType("json").
			JSON(&github.Status{
				Context:     &taskName,
				Description: &malformedDescription,
				State:       &state,
			}).
			Reply(404)

		// expect
		if err := s.CreateCommitStatus(
			context.New(taskName),
			repo,
			plumbing.Hash{},
			state,
			description,
		); err == nil {
			t.Error("errot must occred. but got nil")
		}

		if !gock.IsDone() {
			t.Error("request missing...")
			for _, req := range gock.GetUnmatchedRequests() {
				bytes, _ := ioutil.ReadAll(req.Body)
				t.Logf("%+v", string(bytes))
			}
		}

		// cleanup
		gock.Clean()
	})
}
