package models

import (
  "encoding/json"
  "time"

  "github.com/uptrace/bun"
)

type Gender string

type ParticipantType string

const (
  ParticipantTypeStudent                ParticipantType = "student"                  // นักเรียน/ผู้ที่สนใจศึกษาต่อ
  ParticipantTypeIntania                ParticipantType = "intania"                  // นิสิตปัจจุบัน/นิสิตเก่าวิศวะจุฬาฯ
  ParticipantTypeOtherUniversityStudent ParticipantType = "other_university_student" // นิสิตจากมหาลัยอื่น
  ParticipantTypeTeacher                ParticipantType = "teacher"                  // ครู
  ParticipantTypeOther                  ParticipantType = "other"                    // ผู้ปกครอง/บุคคลภายนอก
)

type TransportMode string

type OriginLocation string

type User struct {
  bun.BaseModel `bun:"table:users,alias:u"`

  ID              int64           `bun:"id,pk,autoincrement" json:"id"`
  FirstName       string          `bun:"first_name" json:"first_name"`
  LastName        string          `bun:"last_name" json:"last_name"`
  Gender          Gender          `bun:"gender" json:"gender"`
  PhoneNumber     string          `bun:"phone_number" json:"phone_number"`
  Email           string          `bun:"email" json:"email"`
  ParticipantType ParticipantType `bun:"participant_type" json:"participant_type"`
  TransportMode   TransportMode   `bun:"transport_mode" json:"transport_mode"`
  IsFromBangkok   bool            `bun:"is_from_bangkok" json:"is_from_bangkok"`
  OriginLocation  OriginLocation  `bun:"origin_location" json:"origin_location"`

  AttendanceDates      []string        `bun:"attendance_dates,type:date,array" json:"attendance_dates"` // Date in format `2024-12-31`
  InterestedActivities []string        `bun:"interested_activities,array" json:"interested_activities"`
  DiscoveryChannel     []string        `bun:"discovery_channel,array" json:"discovery_channel"`
  ExtraAttributes      json.RawMessage `bun:"extra_attributes,type:jsonb" json:"extra_attributes"`

  CreatedAt time.Time `bun:"created_at,nullzero" json:"created_at"`
  UpdatedAt time.Time `bun:"updated_at,nullzero" json:"updated_at"`
}

type StudentExtraAttributes struct {
  EducationLevel   string `json:"education_level"`
  SchoolName       string `json:"school_name"`
  StudyPlan        string `json:"study_plan"`
  Province         string `json:"province"`
  TcasRank         string `json:"tcas_rank"`
  InterestedMajor  string `json:"interested_major"`
  EmergencyContact string `json:"emergency_contact"`
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


