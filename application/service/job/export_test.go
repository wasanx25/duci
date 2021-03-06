package job

import "github.com/duck8823/duci/domain/model/job"

type StubService struct {
	ID string
}

func (s *StubService) FindBy(_ job.ID) (*job.Job, error) {
	return nil, nil
}

func (s *StubService) Start(_ job.ID) error {
	return nil
}

func (s *StubService) Append(_ job.ID, _ job.LogLine) error {
	return nil
}

func (s *StubService) Finish(_ job.ID) error {
	return nil
}

type ServiceImpl = serviceImpl

func (s *ServiceImpl) SetRepo(repo job.Repository) (reset func()) {
	tmp := s.repo
	s.repo = repo
	return func() {
		s.repo = tmp
	}
}
