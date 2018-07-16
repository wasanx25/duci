package github_test

import (
	"encoding/json"
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/duck8823/duci/service/github"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type MockHandler struct {
	Body   interface{}
	Status int
}

func (h *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(h.Body)
	w.Write(payload)
	w.WriteHeader(h.Status)
}

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
	mux := http.NewServeMux()
	mux.Handle("/repos/duck8823/duci/pulls/5", &MockHandler{
		Body: struct {
			Id int64 `json:"id"`
		}{Id: 19},
		Status: 200,
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	baseUrl, err := url.Parse(ts.URL + "/")
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	s, err := github.NewWithEnv()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}
	s.Client.BaseURL = baseUrl

	repo := &MockRepo{
		FullName: "duck8823/duci",
		SSHURL:   "git@github.com:duck8823/duci.git",
	}
	pr, err := s.GetPullRequest(context.New("test/task"), repo, 5)
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	actual := pr.GetID()
	expected := 19
	t.Logf("%+v", pr)
	if pr.GetID() != 19 {
		t.Errorf("id must be equal %+v, but got %+v. \npr=%+v", expected, actual, pr)
	}
}

func TestService_CreateCommitStatus(t *testing.T) {
	t.Run("when github server returns status ok", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.Handle("/repos/duck8823/duci/statuses/0000000000000000000000000000000000000000", &MockHandler{
			Status: 200,
		})

		ts := httptest.NewServer(mux)
		defer ts.Close()

		baseUrl, err := url.Parse(ts.URL + "/")
		if err != nil {
			t.Fatalf("error occured. %+v", err)
		}

		s, err := github.NewWithEnv()
		if err != nil {
			t.Fatalf("error occured. %+v", err)
		}
		s.Client.BaseURL = baseUrl

		repo := &MockRepo{
			FullName: "duck8823/duci",
			SSHURL:   "git@github.com:duck8823/duci.git",
		}
		if err := s.CreateCommitStatus(context.New("test/task"), repo, plumbing.Hash{}, github.SUCCESS, ""); err != nil {
			t.Errorf("error must not occured: but got %+v", err)
		}
	})

	t.Run("when github server returns status not found", func(t *testing.T) {
		mux := http.NewServeMux()
		ts := httptest.NewServer(mux)
		defer ts.Close()

		baseUrl, err := url.Parse(ts.URL + "/")
		if err != nil {
			t.Fatalf("error occured. %+v", err)
		}

		s, err := github.NewWithEnv()
		if err != nil {
			t.Fatalf("error occured. %+v", err)
		}
		s.Client.BaseURL = baseUrl

		repo := &MockRepo{
			FullName: "duck8823/duci",
			SSHURL:   "git@github.com:duck8823/duci.git",
		}
		if err := s.CreateCommitStatus(context.New("test/task"), repo, plumbing.Hash{}, github.SUCCESS, ""); err == nil {
			t.Error("errot must occred. but got nil")
		}
	})
}