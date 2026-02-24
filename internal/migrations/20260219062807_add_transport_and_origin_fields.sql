-- +goose Up
-- +goose StatementBegin
CREATE TYPE transport_mode_enum AS ENUM (
    'personal_car', -- รถยนต์ส่วนบุคคล
    'domestic_flight', -- เที่ยวบินในประเทศ
    'personal_pickup_truck', -- รถกระบะส่วนบุคคล
    'public_van', -- รถตู้ประจำทาง
    'taxi', -- แท็กซี่
    'public_bus', -- รถโดยสารประจำทาง
    'personal_electric_car', -- รถยนต์ไฟฟ้าส่วนบุคคล
    'diesel_railcar', -- รถไฟดีเซลราง
    'personal_van', -- รถตู้ส่วนบุคคล
    'public_boat', -- เรือสาธารณะ
    'motorcycle', -- จักรยานยนต์
    'electric_train' -- รถไฟฟ้า
);
CREATE TYPE origin_location_enum AS ENUM (
    -- เขตกรุงเทพ
    'phra_nakhon', -- เขตพระนคร
    'dusit', -- เขตดุสิต
    'nong_chok', -- เขตหนองจอก
    'bang_rak', -- เขตบางรัก
    'bang_khen', -- เขตบางเขน
    'bang_kapi', -- เขตบางกะปิ
    'pathum_wan', -- เขตปทุมวัน
    'pom_prap_sattru_phai', -- เขตป้อมปราบศัตรูพ่าย
    'phra_khanong', -- เขตพระโขนง
    'min_buri', -- เขตมีนบุรี
    'lat_krabang', -- เขตลาดกระบัง
    'yan_nawa', -- เขตยานนาวา
    'khlong_san', -- เขตคลองสาน
    'bang_khae', -- เขตบางแค
    'bang_kho_laem', -- เขตบางคอแหลม
    'bang_sue', -- เขตบางซื่อ
    'bang_na', -- เขตบางนา
    'thawi_watthana', -- เขตทวีวัฒนา
    'thung_khru', -- เขตทุ่งครุ
    'bang_plad', -- เขตบางพลัด
    'bang_bon', -- เขตบางบอน
    'bang_khun_thian', -- เขตบางขุนเทียน
    'phasi_charoen', -- เขตภาษีเจริญ
    'taling_chan', -- เขตตลิ่งชัน
    'chatuchak', -- เขตจตุจักร
    'lak_si', -- เขตหลักสี่
    'sai_mai', -- เขตสายไหม
    'khlong_toei', -- เขตคลองเตย
    'suan_luang', -- เขตสวนหลวง
    'rat_burana', -- เขตราษฎร์บูรณะ
    'huai_khwang', -- เขตห้วยขวาง
    'khlong_sam_wa', -- เขตคลองสามวา
    'wang_thonglang', -- เขตวังทองหลาง
    'saphan_sung', -- เขตสะพานสูง
    'bangkok_yai', -- เขตบางกอกใหญ่
    'bangkok_noi', -- เขตบางกอกน้อย
    'samphanthawong', -- เขตสัมพันธวงศ์
    'phaya_thai', -- เขตพญาไท
    'ratchathewi', -- เขตราชเทวี
    'don_mueang', -- เขตดอนเมือง
    'prawet', -- เขตประเวศ
    'din_daeng', -- เขตดินแดง
    'bueng_kum', -- เขตบึงกุ่ม
    'sathon', -- เขตสาทร
    'chom_thong', -- เขตจอมทอง
    'watthana', -- เขตวัฒนา
    'kannayao', -- เขตคันนายาว

    -- จังหวัด
    'amnat_charoen', -- อำนาจเจริญ
    'ang_thong', -- อ่างทอง
    'bueng_kan', -- บึงกาฬ
    'buriram', -- บุรีรัมย์
    'chachoengsao', -- ฉะเชิงเทรา
    'chai_nat', -- ชัยนาท
    'chaiyaphum', -- ชัยภูมิ
    'chanthaburi', -- จันทบุรี
    'chiang_mai', -- เชียงใหม่
    'chiang_rai', -- เชียงราย
    'chonburi', -- ชลบุรี
    'chumphon', -- ชุมพร
    'kalasin', -- กาฬสินธุ์
    'kamphaeng_phet', -- กำแพงเพชร
    'kanchanaburi', -- กาญจนบุรี
    'khon_kaen', -- ขอนแก่น
    'krabi', -- กระบี่
    'lampang', -- ลำปาง
    'lamphun', -- ลำพูน
    'loei', -- เลย
    'lopburi', -- ลพบุรี
    'mae_hong_son', -- แม่ฮ่องสอน
    'maha_sarakham', -- มหาสารคาม
    'mukdahan', -- มุกดาหาร
    'nakhon_nayok', -- นครนายก
    'nakhon_pathom', -- นครปฐม
    'nakhon_phanom', -- นครพนม
    'nakhon_ratchasima', -- นครราชสีมา
    'nakhon_sawan', -- นครสวรรค์
    'nakhon_si_thammarat', -- นครศรีธรรมราช
    'nan', -- น่าน
    'narathiwat', -- นราธิวาส
    'nong_bua_lamphu', -- หนองบัวลำภู
    'nong_khai', -- หนองคาย
    'nonthaburi', -- นนทบุรี
    'pathum_thani', -- ปทุมธานี
    'pattani', -- ปัตตานี
    'phang_nga', -- พังงา
    'phatthalung', -- พัทลุง
    'phayao', -- พะเยา
    'phetchabun', -- เพชรบูรณ์
    'phetchaburi', -- เพชรบุรี
    'phichit', -- พิจิตร
    'phitsanulok', -- พิษณุโลก
    'phra_nakhon_si_ayutthaya', -- พระนครศรีอยุธยา
    'phrae', -- แพร่
    'phuket', -- ภูเก็ต
    'prachinburi', -- ปราจีนบุรี
    'prachuap_khiri_khan', -- ประจวบคีรีขันธ์
    'ranong', -- ระนอง
    'ratchaburi', -- ราชบุรี
    'rayong', -- ระยอง
    'roi_et', -- ร้อยเอ็ด
    'sa_kaeo', -- สระแก้ว
    'sakon_nakhon', -- สกลนคร
    'samut_prakan', -- สมุทรปราการ
    'samut_sakhon', -- สมุทรสาคร
    'samut_songkhram', -- สมุทรสงคราม
    'sara_buri', -- สระบุรี
    'satun', -- สตูล
    'sing_buri', -- สิงห์บุรี
    'sisaket', -- ศรีสะเกษ
    'songkhla', -- สงขลา
    'sukhothai', -- สุโขทัย
    'suphan_buri', -- สุพรรณบุรี
    'surat_thani', -- สุราษฎร์ธานี
    'surin', -- สุรินทร์
    'tak', -- ตาก
    'trang', -- ตรัง
    'trat', -- ตราด
    'ubon_ratchathani', -- อุบลราชธานี
    'udon_thani', -- อุดรธานี
    'uthai_thani', -- อุทัยธานี
    'uttaradit', -- อุตรดิตถ์
    'yala', -- ยะลา
    'yasothon' -- ยโสธร
);

ALTER TABLE users
ADD COLUMN transport_mode transport_mode_enum,
ADD COLUMN is_from_bangkok BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN origin_location origin_location_enum;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN IF EXISTS transport_mode,
DROP COLUMN IF EXISTS is_from_bangkok,
DROP COLUMN IF EXISTS origin_location;

DROP TYPE IF EXISTS transport_mode_enum;
DROP TYPE IF EXISTS origin_location_enum;

-- +goose StatementEnd
