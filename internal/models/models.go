package models

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

type Gender string

const (
	GenderMale           Gender = "male"              // ชาย
	GenderFemale         Gender = "female"            // หญิง
	GenderPreferNotToSay Gender = "prefer_not_to_say" // ไม่ต้องการระบุ
	GenderOther          Gender = "other"             // อื่นๆ
)

type ParticipantType string

const (
	ParticipantTypeStudent                ParticipantType = "student"                  // นักเรียน/ผู้ที่สนใจศึกษาต่อ
	ParticipantTypeIntania                ParticipantType = "intania"                  // นิสิตปัจจุบัน/นิสิตเก่าวิศวะจุฬาฯ
	ParticipantTypeOtherUniversityStudent ParticipantType = "other_university_student" // นิสิตจากมหาลัยอื่น
	ParticipantTypeTeacher                ParticipantType = "teacher"                  // ครู
	ParticipantTypeOther                  ParticipantType = "other"                    // ผู้ปกครอง/บุคคลภายนอก
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID              int64           `bun:"id,pk,autoincrement" json:"id,omitempty"`
	FirstName       string          `bun:"first_name" json:"first_name,omitempty"`
	LastName        string          `bun:"last_name" json:"last_name,omitempty"`
	Gender          Gender          `bun:"gender" json:"gender,omitempty"`
	PhoneNumber     string          `bun:"phone_number" json:"phone_number,omitempty"`
	Email           string          `bun:"email" json:"email,omitempty"`
	ParticipantType ParticipantType `bun:"participant_type" json:"participant_type,omitempty"`

	AttendanceDates      []string        `bun:"attendance_dates,type:date,array" json:"attendance_dates,omitempty"` // Date in format `2024-12-31`
	InterestedActivities []string        `bun:"interested_activities,array" json:"interested_activities,omitempty"`
	DiscoveryChannel     []string        `bun:"discovery_channel,array" json:"discovery_channel,omitempty"`
	ExtraAttributes      json.RawMessage `bun:"extra_attributes,type:jsonb" json:"extra_attributes,omitempty"`

	CreatedAt time.Time `bun:"created_at,nullzero" json:"created_at,omitempty"`
	UpdatedAt time.Time `bun:"updated_at,nullzero" json:"updated_at,omitempty"`
}

type StudentExtraAttributes struct {
	EducationLevel   string   `json:"education_level"`
	SchoolName       string   `json:"school_name"`
	StudyPlan        string   `json:"study_plan"`
	Province         string   `json:"province"`
	TcasRank         string   `json:"tcas_rank"`
	InterestedMajors []string `json:"interested_majors"`
	EmergencyContact string   `json:"emergency_contact"`
}

type IntaniaExtraAttributes struct {
	IntaniaGeneration string `json:"intania_generation"`
}

type OtherUniversityStudentExtraAttributes struct {
	YearLevel  string `json:"year_level"`
	Faculty    string `json:"faculty"`
	University string `json:"university"`
}

type TeacherExtraAttributes struct {
	SchoolName    string `json:"school_name"`
	Province      string `json:"province"`
	SubjectTaught string `json:"subject_taught"`
}
