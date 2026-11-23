package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/repository"
	staffMembersModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	uuid "github.com/satori/go.uuid"
)

func SaveStaffScope(req *requests.StaffScopeRequest, addedBy string) error {
	var deployment staffMembersModels.StaffDeployment
	userNumber := req.StaffNumber
	if err := config.DB.Where("number = ? OR user_number = ?", req.StaffNumber, req.StaffNumber).First(&deployment).Error; err == nil && deployment.UserNumber != "" {
		userNumber = deployment.UserNumber
	}

	var rows []models.StaffScope

	if len(req.Campuses) == 0 && req.SchoolNumber != "" {
		rows = append(rows, models.StaffScope{
			Number:       uuid.NewV4().String(),
			StaffNumber:  req.StaffNumber,
			UserNumber:   userNumber,
			SchoolNumber: req.SchoolNumber,
			AddedBy:      addedBy,
		})
	}

	for _, campus := range req.Campuses {
		if len(campus.Curriculums) == 0 {
			rows = append(rows, models.StaffScope{
				Number:       uuid.NewV4().String(),
				StaffNumber:  req.StaffNumber,
				UserNumber:   userNumber,
				SchoolNumber: req.SchoolNumber,
				CampusNumber: campus.CampusNumber,
				AddedBy:      addedBy,
			})
			continue
		}
		for _, curr := range campus.Curriculums {
			if len(curr.Levels) == 0 {
				rows = append(rows, models.StaffScope{
					Number:           uuid.NewV4().String(),
					StaffNumber:      req.StaffNumber,
					UserNumber:       userNumber,
					SchoolNumber:     req.SchoolNumber,
					CampusNumber:     campus.CampusNumber,
					CurriculumNumber: curr.CurriculumNumber,
					AddedBy:          addedBy,
				})
				continue
			}
			for _, level := range curr.Levels {
				if len(level.Classes) == 0 {
					rows = append(rows, models.StaffScope{
						Number:            uuid.NewV4().String(),
						StaffNumber:       req.StaffNumber,
						UserNumber:        userNumber,
						SchoolNumber:      req.SchoolNumber,
						CampusNumber:      campus.CampusNumber,
						CurriculumNumber:  curr.CurriculumNumber,
						CampusLevelNumber: level.CampusLevelNumber,
						AddedBy:           addedBy,
					})
					continue
				}
				for _, classItem := range level.Classes {
					if len(classItem.Streams) == 0 {
						rows = append(rows, models.StaffScope{
							Number:            uuid.NewV4().String(),
							StaffNumber:       req.StaffNumber,
							UserNumber:        userNumber,
							SchoolNumber:      req.SchoolNumber,
							CampusNumber:      campus.CampusNumber,
							CurriculumNumber:  curr.CurriculumNumber,
							CampusLevelNumber: level.CampusLevelNumber,
							ClassNumber:       classItem.ClassNumber,
							AddedBy:           addedBy,
						})
						continue
					}
					for _, stream := range classItem.Streams {
						if len(stream.AssessmentAreas) == 0 {
							rows = append(rows, models.StaffScope{
								Number:            uuid.NewV4().String(),
								StaffNumber:       req.StaffNumber,
								UserNumber:        userNumber,
								SchoolNumber:      req.SchoolNumber,
								CampusNumber:      campus.CampusNumber,
								CurriculumNumber:  curr.CurriculumNumber,
								CampusLevelNumber: level.CampusLevelNumber,
								ClassNumber:       classItem.ClassNumber,
								StreamNumber:      stream.StreamNumber,
								AddedBy:           addedBy,
							})
							continue
						}
						for _, area := range stream.AssessmentAreas {
							rows = append(rows, models.StaffScope{
								Number:               uuid.NewV4().String(),
								StaffNumber:          req.StaffNumber,
								UserNumber:           userNumber,
								SchoolNumber:         req.SchoolNumber,
								CampusNumber:         campus.CampusNumber,
								CurriculumNumber:     curr.CurriculumNumber,
								CampusLevelNumber:    level.CampusLevelNumber,
								ClassNumber:          classItem.ClassNumber,
								StreamNumber:         stream.StreamNumber,
								AssessmentAreaNumber: area,
								AddedBy:              addedBy,
							})
						}
					}
				}
			}
		}
	}

	repo := repository.NewStaffScopeRepository(config.DB)
	return repo.ReplaceStaffScopes(req.StaffNumber, req.SchoolNumber, rows)
}

func GetStaffScopeTree(staffNumber, schoolNumber string) (*requests.StaffScopeRequest, error) {
	repo := repository.NewStaffScopeRepository(config.DB)
	rows, err := repo.GetStaffScopes(staffNumber, schoolNumber)
	if err != nil {
		return nil, err
	}

	result := &requests.StaffScopeRequest{
		StaffNumber:  staffNumber,
		SchoolNumber: schoolNumber,
		Campuses:     make([]requests.StaffScopeCampus, 0),
	}

	type streamNode struct {
		areas map[string]bool
	}
	type classNode struct {
		streams map[string]*streamNode
	}
	type levelNode struct {
		classes map[string]*classNode
	}
	type currNode struct {
		levels map[string]*levelNode
	}
	type campusNode struct {
		curriculums map[string]*currNode
	}

	tree := make(map[string]*campusNode)

	for _, row := range rows {
		if row.CampusNumber == "" {
			continue
		}
		cNode, ok := tree[row.CampusNumber]
		if !ok {
			cNode = &campusNode{curriculums: make(map[string]*currNode)}
			tree[row.CampusNumber] = cNode
		}
		if row.CurriculumNumber == "" {
			continue
		}
		cuNode, ok := cNode.curriculums[row.CurriculumNumber]
		if !ok {
			cuNode = &currNode{levels: make(map[string]*levelNode)}
			cNode.curriculums[row.CurriculumNumber] = cuNode
		}
		if row.CampusLevelNumber == "" {
			continue
		}
		lNode, ok := cuNode.levels[row.CampusLevelNumber]
		if !ok {
			lNode = &levelNode{classes: make(map[string]*classNode)}
			cuNode.levels[row.CampusLevelNumber] = lNode
		}
		if row.ClassNumber == "" {
			continue
		}
		clNode, ok := lNode.classes[row.ClassNumber]
		if !ok {
			clNode = &classNode{streams: make(map[string]*streamNode)}
			lNode.classes[row.ClassNumber] = clNode
		}
		if row.StreamNumber == "" {
			continue
		}
		sNode, ok := clNode.streams[row.StreamNumber]
		if !ok {
			sNode = &streamNode{areas: make(map[string]bool)}
			clNode.streams[row.StreamNumber] = sNode
		}
		if row.AssessmentAreaNumber != "" {
			sNode.areas[row.AssessmentAreaNumber] = true
		}
	}

	for cID, cNode := range tree {
		campusDTO := requests.StaffScopeCampus{
			CampusNumber: cID,
			Curriculums:  make([]requests.StaffScopeCurriculum, 0),
		}
		for cuID, cuNode := range cNode.curriculums {
			currDTO := requests.StaffScopeCurriculum{
				CurriculumNumber: cuID,
				Levels:           make([]requests.StaffScopeLevel, 0),
			}
			for lID, lNode := range cuNode.levels {
				lvlDTO := requests.StaffScopeLevel{
					CampusLevelNumber: lID,
					Classes:           make([]requests.StaffScopeClass, 0),
				}
				for clID, clNode := range lNode.classes {
					clsDTO := requests.StaffScopeClass{
						ClassNumber: clID,
						Streams:     make([]requests.StaffScopeStream, 0),
					}
					for sID, sNode := range clNode.streams {
						areas := make([]string, 0, len(sNode.areas))
						for areaID := range sNode.areas {
							areas = append(areas, areaID)
						}
						clsDTO.Streams = append(clsDTO.Streams, requests.StaffScopeStream{
							StreamNumber:    sID,
							AssessmentAreas: areas,
						})
					}
					lvlDTO.Classes = append(lvlDTO.Classes, clsDTO)
				}
				currDTO.Levels = append(currDTO.Levels, lvlDTO)
			}
			campusDTO.Curriculums = append(campusDTO.Curriculums, currDTO)
		}
		result.Campuses = append(result.Campuses, campusDTO)
	}

	return result, nil
}

func ClearStaffScope(staffNumber, schoolNumber string) error {
	repo := repository.NewStaffScopeRepository(config.DB)
	return repo.ClearStaffScope(staffNumber, schoolNumber)
}

