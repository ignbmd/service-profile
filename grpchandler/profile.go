package grpchandler

import (
	"context"

	grpcprofile "github.com/btwedutech/grpc/service/profile"
	"google.golang.org/protobuf/types/known/timestamppb"
	"smartbtw.com/services/profile/lib"
)

type ProfileDelivery struct {
	grpcprofile.ProfileServer
}

func (s *ProfileDelivery) GetProfileElastic(context context.Context, req *grpcprofile.GetElasticRequest) (*grpcprofile.GetElasticResponse, error) {
	studentData, err := lib.GetStudentProfileElastic(int(req.SmartbtwId))
	if err != nil {
		return nil, err
	}
	//TODO: Update struct soon
	dom := &grpcprofile.GetElasticResponse{
		SmartbtwId:             int32(studentData.SmartbtwID),
		Name:                   studentData.Name,
		Email:                  studentData.Email,
		Phone:                  studentData.Phone,
		Gender:                 studentData.Gender,
		BirthDate:              timestamppb.New(studentData.BirthDate),
		Province:               studentData.Province,
		ProvinceId:             uint32(studentData.ProvinceID),
		Region:                 studentData.Region,
		RegionId:               uint32(studentData.RegionID),
		DomicileProvince:       studentData.DomicileProvince,
		DomicileProvinceId:     uint32(studentData.DomicileProvinceID),
		DomicileRegion:         studentData.DomicileRegion,
		DomicileRegionId:       uint32(studentData.RegionID),
		LastEdId:               studentData.LastEdID,
		LastEdName:             studentData.LastEdName,
		LastEdType:             studentData.LastEdMajor,
		LastEdMajor:            studentData.LastEdMajor,
		LastEdMajorId:          uint32(studentData.LastEdMajorID),
		LastEdRegion:           studentData.LastEdRegion,
		LastEdRegionId:         uint32(studentData.LastEdRegionID),
		EyeColorBlind:          studentData.EyeColorBlind,
		Height:                 studentData.Height,
		Weight:                 studentData.Weight,
		SchoolPtkId:            uint32(studentData.SchoolPTKID),
		SchoolNamePtk:          studentData.SchoolNamePTK,
		MajorNamePtk:           studentData.MajorNamePTK,
		MajorPtkId:             uint32(studentData.MajorPTKID),
		PolbitTypePtk:          studentData.PolbitTypePTK,
		PolbitCompetitionPtkId: uint32(studentData.PolbitCompetitionPTKID),
		PolbitLocationPtkId:    uint32(studentData.PolbitLocationPTKID),
		CreatedAtPtk:           timestamppb.New(studentData.CreatedAtPTK),
		SchoolPtnId:            uint32(studentData.SchoolPTNID),
		SchoolNamePtn:          studentData.SchoolNamePTN,
		MajorPtnId:             uint32(studentData.MajorPTNID),
		MajorNamePtn:           studentData.MajorNamePTN,
		CreatedAtPtn:           timestamppb.New(studentData.CreatedAtPTN),
		InstanceCpnsId:         uint32(studentData.InstanceCPNSID),
		InstaceCpnsName:        studentData.InstanceCPNSName,
		PositionCpnsId:         uint32(studentData.PositionCPNSID),
		PositionCpnsName:       studentData.PositionCPNSName,
		FormationCpnsType:      studentData.FormationCPNSType,
		CompetitionCpnsId:      uint32(studentData.CompetitionCPNSID),
		CreatedAtCpns:          timestamppb.New(studentData.CreatedAtCPNS),
		AccountType:            studentData.AccountType,
		CreatedAt:              timestamppb.New(studentData.CreatedAt),
	}

	if studentData.Photo != nil {
		dom.Photo = *studentData.Photo
	}
	if studentData.BranchCode != nil {
		dom.BranchCode = *studentData.BranchCode
	}
	if studentData.BranchName != nil {
		dom.BranchName = *studentData.BranchName
	}
	if studentData.Interest != nil {
		dom.Interest = *studentData.Interest
	}
	return dom, nil
}
