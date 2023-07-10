package builder

type Link struct {
	Title  string `json:"title"`
	Action string `json:"action"`
	View   string `json:"view"`
}
type BuilderResponse struct {
	Links []Link `json:"links"`
	Key   string `json:"key"`
}

func NewBuilderResponse() *BuilderResponse {
	return &BuilderResponse{}
}

func (b *BuilderResponse) AddLink(title, action, view string) *BuilderResponse {
	b.Links = append(b.Links, Link{
		Title:  title,
		Action: action,
		View:   view,
	})
	return b
}
func (b *BuilderResponse) SetKey(value string) *BuilderResponse {
	b.Key = value
	return b
}
