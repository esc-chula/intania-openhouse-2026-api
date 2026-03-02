package myValidator

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/go-playground/validator/v10"
)

var (
	ErrExtraAttributesRequired = errors.New("extra attributes required")
	ErrExtraAttributesInvalid  = errors.New("extra attributes invalid")
	ErrInvalidEventDate        = errors.New("invalid event date format, expected YYYY-MM-DD")
)

var validate = validator.New()

func ValidateExtraAttributes(participantType models.ParticipantType, fields *models.ExtraAttributesFields) (extraAttributes json.RawMessage, err error) {
	var chosenField any

	switch participantType {
	case models.ParticipantTypeStudent:
		if fields.StudentExtraAttributes == nil {
			return nil, ErrExtraAttributesRequired
		}
		chosenField = fields.StudentExtraAttributes

	case models.ParticipantTypeIntania:
		if fields.IntaniaExtraAttributes == nil {
			return nil, ErrExtraAttributesRequired
		}
		chosenField = fields.IntaniaExtraAttributes

	case models.ParticipantTypeOutsideStudent:
		if fields.OutsideStudentExtraAttributes == nil {
			return nil, ErrExtraAttributesRequired
		}
		chosenField = fields.OutsideStudentExtraAttributes

	case models.ParticipantTypeAlumni:
		if fields.AlumniExtraAttributes == nil {
			return nil, ErrExtraAttributesRequired
		}
		chosenField = fields.AlumniExtraAttributes

	case models.ParticipantTypeTeacher:
		if fields.TeacherExtraAttributes == nil {
			return nil, ErrExtraAttributesRequired
		}
		chosenField = fields.TeacherExtraAttributes

	case models.ParticipantTypeOther:
		return nil, nil
	default:
		return nil, ErrInvalidParticipantType
	}

	extraAttributes, err = json.Marshal(chosenField)
	if err != nil {
		return nil, err
	}

	return extraAttributes, err
}

func ValidateAttendanceDate(user *models.User) error {
	for _, date := range user.AttendanceDates {
		if _, err := time.Parse("2006-01-02", date); err != nil {
			return err
		}
	}
	return nil
}

func ValidateEventDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return ErrInvalidEventDate
	}
	return nil
}
