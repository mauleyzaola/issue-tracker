package domain

type DocumentMetadata struct {
	DocumentType string `json:"documentType" db:"-"`
	FriendName   string `json:"friendName" db:"-"`
	Id           string `json:"id" db:"-"`
}
