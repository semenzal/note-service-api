package note

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/semenzal/note-service-api/pkg/note_v1"
)

func (s *Service) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	_, err := s.noteRepository.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}