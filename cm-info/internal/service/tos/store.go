package service

import (
	tosRepository "github.com/maciejas22/conference-manager-api/cm-info/internal/repository/tos"
)

type ToSService struct {
	tosRepo tosRepository.ToSRepoInterface
}

func NewToSService(tosRepo tosRepository.ToSRepoInterface) ToSServiceInterface {
	return &ToSService{
		tosRepo: tosRepo,
	}
}

func (s *ToSService) GetTermsOfService() (*TermsOfService, error) {
	tos, err := s.tosRepo.GetTermsOfService()
	if err != nil {
		return nil, err
	}

	tgtSections := make([]*Section, len(tos.Sections))
	for i, sec := range tos.Sections {
		if sec != nil {
			tgtSubsections := make([]*Subsection, len(sec.Subsections))
			for j, sub := range sec.Subsections {
				if sub != nil {
					tgtSubsections[j] = &Subsection{
						Id:        sub.Id,
						SectionId: sub.SectionId,
						Title:     sub.Title,
						Content:   sub.Content,
						CreatedAt: sub.CreatedAt,
					}
				}
			}

			tgtSections[i] = &Section{
				Id:               sec.Id,
				TermsOfServiceId: sec.TermsOfServiceId,
				Title:            sec.Title,
				Content:          sec.Content,
				CreatedAt:        sec.CreatedAt,
				Subsections:      tgtSubsections,
			}
		}
	}

	return &TermsOfService{
		Id:              tos.Id,
		CreatedAt:       tos.CreatedAt,
		UpdatedAt:       tos.UpdatedAt,
		Introduction:    tos.Introduction,
		Acknowledgement: tos.Acknowledgement,
		Sections:        tgtSections,
	}, nil
}
