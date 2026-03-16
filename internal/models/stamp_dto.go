package models

type UserStamps struct {
	TotalCount           int64
	DepartmentStampCount int64
	ClubStampCount       int64
	ExhibitionStampCount int64
	DepartmentStamps     []StampItem
	ClubStamps           []StampItem
	ExhibitionStamps     []StampItem
}

type StampRedemptionStatus struct {
	DepartmentIsRedeemed bool
	ClubIsRedeemed       bool
	ExhibitionIsRedeemed bool
	DepartmentRedeemable bool
	ClubRedeemable       bool
	ExhibitionRedeemable bool
}
