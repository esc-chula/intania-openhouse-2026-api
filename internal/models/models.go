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

type TransportMode string

const (
	TransportModePersonalCar         TransportMode = "personal_car"          // รถยนต์ส่วนบุคคล
	TransportModeDomesticFlight      TransportMode = "domestic_flight"       // เที่ยวบินในประเทศ
	TransportModePersonalPickupTruck TransportMode = "personal_pickup_truck" // รถกระบะส่วนบุคคล
	TransportModePublicVan           TransportMode = "public_van"            // รถตู้ประจำทาง
	TransportModeTaxi                TransportMode = "taxi"                  // แท็กซี่
	TransportModePublicBus           TransportMode = "public_bus"            // รถโดยสารประจำทาง
	TransportModePersonalElectricCar TransportMode = "personal_electric_car" // รถยนต์ไฟฟ้าส่วนบุคคล
	TransportModeDieselRailcar       TransportMode = "diesel_railcar"        // รถไฟดีเซลราง
	TransportModePersonalVan         TransportMode = "personal_van"          // รถตู้ส่วนบุคคล
	TransportModePublicBoat          TransportMode = "public_boat"           // เรือสาธารณะ
	TransportModeMotorcycle          TransportMode = "motorcycle"            // จักรยานยนต์
	TransportModeElectricTrain       TransportMode = "electric_train"        // รถไฟฟ้า
)

type OriginLocation string

const (
	// Bangkok Districts
	OriginLocationPhraNakhon        OriginLocation = "phra_nakhon"
	OriginLocationDusit             OriginLocation = "dusit"
	OriginLocationNongChok          OriginLocation = "nong_chok"
	OriginLocationBangRak           OriginLocation = "bang_rak"
	OriginLocationBangKhen          OriginLocation = "bang_khen"
	OriginLocationBangKapi          OriginLocation = "bang_kapi"
	OriginLocationPathumWan         OriginLocation = "pathum_wan"
	OriginLocationPomPrapSattruPhai OriginLocation = "pom_prap_sattru_phai"
	OriginLocationPhraKhanong       OriginLocation = "phra_khanong"
	OriginLocationMinBuri           OriginLocation = "min_buri"
	OriginLocationLatKrabang        OriginLocation = "lat_krabang"
	OriginLocationYanNawa           OriginLocation = "yan_nawa"
	OriginLocationKhlongSan         OriginLocation = "khlong_san"
	OriginLocationBangKhae          OriginLocation = "bang_khae"
	OriginLocationBangKhoLaem       OriginLocation = "bang_kho_laem"
	OriginLocationBangSue           OriginLocation = "bang_sue"
	OriginLocationBangNa            OriginLocation = "bang_na"
	OriginLocationThawiWatthana     OriginLocation = "thawi_watthana"
	OriginLocationThungKhru         OriginLocation = "thung_khru"
	OriginLocationBangPlad          OriginLocation = "bang_plad"
	OriginLocationBangBon           OriginLocation = "bang_bon"
	OriginLocationBangKhunThian     OriginLocation = "bang_khun_thian"
	OriginLocationPhasiCharoen      OriginLocation = "phasi_charoen"
	OriginLocationTalingChan        OriginLocation = "taling_chan"
	OriginLocationChatuchak         OriginLocation = "chatuchak"
	OriginLocationLakSi             OriginLocation = "lak_si"
	OriginLocationSaiMai            OriginLocation = "sai_mai"
	OriginLocationKhlongToei        OriginLocation = "khlong_toei"
	OriginLocationSuanLuang         OriginLocation = "suan_luang"
	OriginLocationRatBurana         OriginLocation = "rat_burana"
	OriginLocationHuaiKhwang        OriginLocation = "huai_khwang"
	OriginLocationKhlongSamWa       OriginLocation = "khlong_sam_wa"
	OriginLocationWangThonglang     OriginLocation = "wang_thonglang"
	OriginLocationSaphanSung        OriginLocation = "saphan_sung"
	OriginLocationBangkokYai        OriginLocation = "bangkok_yai"
	OriginLocationBangkokNoi        OriginLocation = "bangkok_noi"
	OriginLocationSamphanthawong    OriginLocation = "samphanthawong"
	OriginLocationPhayaThai         OriginLocation = "phaya_thai"
	OriginLocationRatchathewi       OriginLocation = "ratchathewi"
	OriginLocationDonMueang         OriginLocation = "don_mueang"
	OriginLocationPrawet            OriginLocation = "prawet"
	OriginLocationDinDaeng          OriginLocation = "din_daeng"
	OriginLocationBuengKum          OriginLocation = "bueng_kum"
	OriginLocationSathon            OriginLocation = "sathon"
	OriginLocationChomThong         OriginLocation = "chom_thong"
	OriginLocationWatthana          OriginLocation = "watthana"
	OriginLocationKannayao          OriginLocation = "kannayao"

	// Provinces
	OriginLocationAmnatCharoen          OriginLocation = "amnat_charoen"
	OriginLocationAngThong              OriginLocation = "ang_thong"
	OriginLocationBuengKan              OriginLocation = "bueng_kan"
	OriginLocationBuriram               OriginLocation = "buriram"
	OriginLocationChachoengsao          OriginLocation = "chachoengsao"
	OriginLocationChaiNat               OriginLocation = "chai_nat"
	OriginLocationChaiyaphum            OriginLocation = "chaiyaphum"
	OriginLocationChanthaburi           OriginLocation = "chanthaburi"
	OriginLocationChiangMai             OriginLocation = "chiang_mai"
	OriginLocationChiangRai             OriginLocation = "chiang_rai"
	OriginLocationChonburi              OriginLocation = "chonburi"
	OriginLocationChumphon              OriginLocation = "chumphon"
	OriginLocationKalasin               OriginLocation = "kalasin"
	OriginLocationKamphaengPhet         OriginLocation = "kamphaeng_phet"
	OriginLocationKanchanaburi          OriginLocation = "kanchanaburi"
	OriginLocationKhonKaen              OriginLocation = "khon_kaen"
	OriginLocationKrabi                 OriginLocation = "krabi"
	OriginLocationLampang               OriginLocation = "lampang"
	OriginLocationLamphun               OriginLocation = "lamphun"
	OriginLocationLoei                  OriginLocation = "loei"
	OriginLocationLopburi               OriginLocation = "lopburi"
	OriginLocationMaeHongSon            OriginLocation = "mae_hong_son"
	OriginLocationMahaSarakham          OriginLocation = "maha_sarakham"
	OriginLocationMukdahan              OriginLocation = "mukdahan"
	OriginLocationNakhonNayok           OriginLocation = "nakhon_nayok"
	OriginLocationNakhonPathom          OriginLocation = "nakhon_pathom"
	OriginLocationNakhonPhanom          OriginLocation = "nakhon_phanom"
	OriginLocationNakhonRatchasima      OriginLocation = "nakhon_ratchasima"
	OriginLocationNakhonSawan           OriginLocation = "nakhon_sawan"
	OriginLocationNakhonSiThammarat     OriginLocation = "nakhon_si_thammarat"
	OriginLocationNan                   OriginLocation = "nan"
	OriginLocationNarathiwat            OriginLocation = "narathiwat"
	OriginLocationNongBuaLamphu         OriginLocation = "nong_bua_lamphu"
	OriginLocationNongKhai              OriginLocation = "nong_khai"
	OriginLocationNonthaburi            OriginLocation = "nonthaburi"
	OriginLocationPathumThani           OriginLocation = "pathum_thani"
	OriginLocationPattani               OriginLocation = "pattani"
	OriginLocationPhangNga              OriginLocation = "phang_nga"
	OriginLocationPhatthalung           OriginLocation = "phatthalung"
	OriginLocationPhayao                OriginLocation = "phayao"
	OriginLocationPhetchabun            OriginLocation = "phetchabun"
	OriginLocationPhetchaburi           OriginLocation = "phetchaburi"
	OriginLocationPhichit               OriginLocation = "phichit"
	OriginLocationPhitsanulok           OriginLocation = "phitsanulok"
	OriginLocationPhraNakhonSiAyutthaya OriginLocation = "phra_nakhon_si_ayutthaya"
	OriginLocationPhrae                 OriginLocation = "phrae"
	OriginLocationPhuket                OriginLocation = "phuket"
	OriginLocationPrachinburi           OriginLocation = "prachinburi"
	OriginLocationPrachuapKhiriKhan     OriginLocation = "prachuap_khiri_khan"
	OriginLocationRanong                OriginLocation = "ranong"
	OriginLocationRatchaburi            OriginLocation = "ratchaburi"
	OriginLocationRayong                OriginLocation = "rayong"
	OriginLocationRoiEt                 OriginLocation = "roi_et"
	OriginLocationSaKaeo                OriginLocation = "sa_kaeo"
	OriginLocationSakonNakhon           OriginLocation = "sakon_nakhon"
	OriginLocationSamutPrakan           OriginLocation = "samut_prakan"
	OriginLocationSamutSakhon           OriginLocation = "samut_sakhon"
	OriginLocationSamutSongkhram        OriginLocation = "samut_songkhram"
	OriginLocationSaraBuri              OriginLocation = "sara_buri"
	OriginLocationSatun                 OriginLocation = "satun"
	OriginLocationSingBuri              OriginLocation = "sing_buri"
	OriginLocationSisaket               OriginLocation = "sisaket"
	OriginLocationSongkhla              OriginLocation = "songkhla"
	OriginLocationSukhothai             OriginLocation = "sukhothai"
	OriginLocationSuphanBuri            OriginLocation = "suphan_buri"
	OriginLocationSuratThani            OriginLocation = "surat_thani"
	OriginLocationSurin                 OriginLocation = "surin"
	OriginLocationTak                   OriginLocation = "tak"
	OriginLocationTrang                 OriginLocation = "trang"
	OriginLocationTrat                  OriginLocation = "trat"
	OriginLocationUbonRatchathani       OriginLocation = "ubon_ratchathani"
	OriginLocationUdonThani             OriginLocation = "udon_thani"
	OriginLocationUthaiThani            OriginLocation = "uthai_thani"
	OriginLocationUttaradit             OriginLocation = "uttaradit"
	OriginLocationYala                  OriginLocation = "yala"
	OriginLocationYasothon              OriginLocation = "yasothon"
)

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
