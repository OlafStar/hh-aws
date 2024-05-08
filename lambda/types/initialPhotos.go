package types

type PhotoType string

const (
	FullFace  PhotoType = "fullFace"
	Forehead  PhotoType = "forehead"
	LeftSide  PhotoType = "left side"
	RightSide PhotoType = "right side"
	Nose      PhotoType = "nose"
	Chin      PhotoType = "chin"
)

type InitialPhotosStruct struct {
	Image string `json:"image"`
	Type  PhotoType `json:"type"`
}

type InitialPhotosOfClientPostBody struct {
	Images []InitialPhotosStruct `json:"images"`
}

type InitialPhotoRecord struct {
	ID string `json:"id"`
	ClientID string        `json:"clientId"`
	Images   []InitialPhotosStruct `json:"images"`
}