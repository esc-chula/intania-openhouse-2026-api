package myValidator

import (
	_ "embed"
	"encoding/json"
	"errors"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
)

//go:embed enums.json
var enumsJSON []byte

var (
	ErrInvalidGender           = errors.New("invalid gender")
	ErrInvalidParticipantType  = errors.New("invalid participant type")
	ErrInvalidTransportMode    = errors.New("invalid transport mode")
	ErrInvalidOriginLocation   = errors.New("invalid origin location")
	ErrInvalidWorkshopCategory = errors.New("invalid workshop category")
)

type EnumOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Enums struct {
	Genders            []EnumOption `json:"genders"`
	ParticipantTypes   []EnumOption `json:"participant_types"`
	TransportModes     []EnumOption `json:"transport_modes"`
	OriginLocations    []EnumOption `json:"origin_locations"`
	WorkshopCategories []EnumOption `json:"workshop_categories"`
}

var (
	validGenders            map[string]bool
	validParticipantTypes   map[string]bool
	validTransportModes     map[string]bool
	validOriginLocations    map[string]bool
	validWorkshopCategories map[string]bool
	loadedEnums             Enums
)

func init() {
	if err := json.Unmarshal(enumsJSON, &loadedEnums); err != nil {
		panic("failed to load enums.json: " + err.Error())
	}
	validGenders = buildValidMap(loadedEnums.Genders)
	validParticipantTypes = buildValidMap(loadedEnums.ParticipantTypes)
	validTransportModes = buildValidMap(loadedEnums.TransportModes)
	validOriginLocations = buildValidMap(loadedEnums.OriginLocations)
	validWorkshopCategories = buildValidMap(loadedEnums.WorkshopCategories)
}

func buildValidMap(options []EnumOption) map[string]bool {
	m := make(map[string]bool, len(options))
	for _, opt := range options {
		m[opt.Value] = true
	}
	return m
}
func ValidateUserEnums(user *models.User) error {
	if !validGenders[string(user.Gender)] {
		return ErrInvalidGender
	}
	if !validParticipantTypes[string(user.ParticipantType)] {
		return ErrInvalidParticipantType
	}
	if !validTransportModes[string(user.TransportMode)] {
		return ErrInvalidTransportMode
	}
	if !validOriginLocations[string(user.OriginLocation)] {
		return ErrInvalidOriginLocation
	}
	return nil
}

func ValidateWorkshopCategory(category string) error {
	if !validWorkshopCategories[category] {
		return ErrInvalidWorkshopCategory
	}
	return nil
}
