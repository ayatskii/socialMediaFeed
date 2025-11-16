package hashtag

type Hashtag struct {
	ID          int64
	Tag         string
	usage_count int
}

func NewHashtag(ID int64, tag string) *Hashtag {
	return &Hashtag{
		ID:  ID,
		Tag: tag,
	}
}
