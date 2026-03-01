package myValidator

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/go-playground/validator/v10"
)

var (
	ErrExtraAttributesInvalid = errors.New("extra attributes invalid")
	ErrInvalidEventDate       = errors.New("invalid event date format, expected YYYY-MM-DD")
)

var validate = validator.New()

func ValidateExtraAttributes(user *models.User) error {
	switch user.ParticipantType {
	case models.ParticipantTypeStudent:
		var studentExtraAttributes models.StudentExtraAttributes
		if err := validateRawMessage(user.ExtraAttributes, &studentExtraAttributes); err != nil {
			return ErrExtraAttributesInvalid
		}
		return nil
	case models.ParticipantTypeIntania:
		var intaniaExtraAttributes models.IntaniaExtraAttributes
		if err := validateRawMessage(user.ExtraAttributes, &intaniaExtraAttributes); err != nil {
			return ErrExtraAttributesInvalid
		}
		return nil
	case models.ParticipantTypeOutsideStudent:
		var outsideStudentExtraAttributes models.OutsideStudentExtraAttributes
		if err := validateRawMessage(user.ExtraAttributes, &outsideStudentExtraAttributes); err != nil {
			return ErrExtraAttributesInvalid
		}
		return nil
	case models.ParticipantTypeAlumni:
		var alumniExtraAttributes models.IntaniaExtraAttributes
		if err := validateRawMessage(user.ExtraAttributes, &alumniExtraAttributes); err != nil {
			return ErrExtraAttributesInvalid
		}
		return nil
	case models.ParticipantTypeTeacher:
		var teacherExtraAttributes models.TeacherExtraAttributes
		if err := validateRawMessage(user.ExtraAttributes, &teacherExtraAttributes); err != nil {
			return ErrExtraAttributesInvalid
		}
		return nil
	case models.ParticipantTypeOther:
		return nil
	default:
		return ErrInvalidParticipantType
	}
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

func validateRawMessage(raw json.RawMessage, myStruct interface{}) error {

	if len(raw) == 0 {
		return errors.New("extra_attributes is required")
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(myStruct); err != nil {
		return err
	}

	if err := validate.Struct(myStruct); err != nil {
		return err
	}

	return nil
}
