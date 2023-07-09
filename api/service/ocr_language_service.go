package service

import (
	"sort"
	"strings"
	"voltaserve/model"
	"voltaserve/repo"

	"github.com/reactivex/rxgo/v2"
)

type OCRLanguage struct {
	ID        string `json:"id"`
	ISO639Pt3 string `json:"iso639pt3"`
	Name      string `json:"name"`
}

type OCRLanguageList struct {
	Data          []*OCRLanguage `json:"data"`
	TotalPages    uint           `json:"totalPages"`
	TotalElements uint           `json:"totalElements"`
	Page          uint           `json:"page"`
	Size          uint           `json:"size"`
}

type OCRLanguageListOptions struct {
	Query     string
	Page      uint
	Size      uint
	SortBy    string
	SortOrder string
}

type OCRLanguageService struct {
	ocrLanguageRepo   repo.OCRLanguageRepo
	ocrLanguageMapper *ocrLanguageMapper
}

func NewOCRLanguageService() *OCRLanguageService {
	return &OCRLanguageService{
		ocrLanguageRepo:   repo.NewOCRLanguageRepo(),
		ocrLanguageMapper: newOCRLanguageMapper(),
	}
}

func (svc *OCRLanguageService) List(opts OCRLanguageListOptions) (*OCRLanguageList, error) {
	all, err := svc.ocrLanguageRepo.FindAll()
	if err != nil {
		return nil, err
	}
	filtered, err := svc.doFiltering(opts.Query, all)
	if err != nil {
		return nil, err
	}
	sorted := svc.doSorting(filtered, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	mapped, err := svc.ocrLanguageMapper.mapMany(paged)
	if err != nil {
		return nil, err
	}
	return &OCRLanguageList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint(len(mapped)),
	}, nil
}

func (svc *OCRLanguageService) doSorting(data []model.OCRLanguage, sortBy string, sortOrder string) []model.OCRLanguage {
	if sortBy == SortByID {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetID() > data[j].GetID()
			} else {
				return data[i].GetID() < data[j].GetID()
			}
		})
		return data
	} else if sortBy == SortByISO639Pt3 {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetISO639Pt3() > data[j].GetISO639Pt3()
			} else {
				return data[i].GetISO639Pt3() < data[j].GetISO639Pt3()
			}
		})
		return data
	}
	return data
}

func (svc *OCRLanguageService) doPagination(data []model.OCRLanguage, page, size uint) ([]model.OCRLanguage, uint, uint) {
	totalElements := uint(len(data))
	totalPages := (totalElements + size - 1) / size
	if page > totalPages {
		return nil, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	pageData := data[startIndex:endIndex]
	return pageData, totalElements, totalPages
}

func (svc *OCRLanguageService) doFiltering(query string, data []model.OCRLanguage) ([]model.OCRLanguage, error) {
	filtered, _ := rxgo.Just(data)().
		Filter(func(v interface{}) bool {
			if query != "" {
				ocrLanguage := v.(model.OCRLanguage)
				return strings.Contains(strings.ToLower(ocrLanguage.GetID()), strings.ToLower(query)) || strings.Contains(strings.ToLower(ocrLanguage.GetISO639Pt3()), strings.ToLower(query)) || strings.Contains(strings.ToLower(ocrLanguage.GetName()), strings.ToLower(query))
			} else {
				return true
			}
		}).
		ToSlice(0)
	var res []model.OCRLanguage
	for _, v := range filtered {
		res = append(res, v.(model.OCRLanguage))
	}
	return res, nil
}

type ocrLanguageMapper struct {
}

func newOCRLanguageMapper() *ocrLanguageMapper {
	return &ocrLanguageMapper{}
}

func (mp *ocrLanguageMapper) mapOne(ocrLanguage model.OCRLanguage) *OCRLanguage {
	return &OCRLanguage{
		ID:        ocrLanguage.GetID(),
		ISO639Pt3: ocrLanguage.GetISO639Pt3(),
		Name:      ocrLanguage.GetName(),
	}
}

func (mp *ocrLanguageMapper) mapMany(ocrLanguages []model.OCRLanguage) ([]*OCRLanguage, error) {
	res := []*OCRLanguage{}
	for _, u := range ocrLanguages {
		res = append(res, mp.mapOne(u))
	}
	return res, nil
}
