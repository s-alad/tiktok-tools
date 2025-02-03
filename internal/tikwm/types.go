package tikwm

type MediaType string

const (
	MediaTypeNone   MediaType = ""
	MediaTypeVideo  MediaType = "video"
	MediaTypeSlides MediaType = "slides"
)

type TikwmRequestBody struct {
	URL    string `json:"url"`    // URL of the TikTok video
	Count  int    `json:"count"`  // Number of items to fetch
	Cursor int    `json:"cursor"` // Pagination cursor
	Web    int    `json:"web"`    // Web flag (1 or 0)
	HD     int    `json:"hd"`     // HD flag (1 or 0)
}

type TikwmRequestResponse struct {
	Code          int       `json:"code"`
	Msg           string    `json:"msg"`
	ProcessedTime float64   `json:"processed_time"`
	Data          tikwmData `json:"data"`
}

type tikwmData struct {
	ID                  string            `json:"id"`
	Region              string            `json:"region"`
	Title               string            `json:"title"`
	Cover               string            `json:"cover"`
	Duration            int               `json:"duration"`
	Play                string            `json:"play"`
	Wmplay              string            `json:"wmplay"`
	Hdplay              string            `json:"hdplay"`
	Size                int               `json:"size"`
	WmSize              int               `json:"wm_size"`
	HdSize              int               `json:"hd_size"`
	Music               string            `json:"music"`
	MusicInfo           tikwmMusicInfo    `json:"music_info"`
	PlayCount           int               `json:"play_count"`
	DiggCount           int               `json:"digg_count"`
	CommentCount        int               `json:"comment_count"`
	ShareCount          int               `json:"share_count"`
	DownloadCount       int               `json:"download_count"`
	CollectCount        int               `json:"collect_count"`
	CreateTime          int64             `json:"create_time"`
	Anchors             interface{}       `json:"anchors"`
	AnchorsExtras       string            `json:"anchors_extras"`
	IsAd                bool              `json:"is_ad"`
	CommerceInfo        tikwmCommerceInfo `json:"commerce_info"`
	CommercialVideoInfo string            `json:"commercial_video_info"`
	ItemCommentSettings int               `json:"item_comment_settings"`
	MentionedUsers      string            `json:"mentioned_users"`
	Author              tikwmAuthor       `json:"author"`
	Images              []string          `json:"images"`
}

type tikwmMusicInfo struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Play     string `json:"play"`
	Author   string `json:"author"`
	Original bool   `json:"original"`
	Duration int    `json:"duration"`
	Album    string `json:"album"`
}

type tikwmCommerceInfo struct {
	AdvPromotable          bool `json:"adv_promotable"`
	AuctionAdInvited       bool `json:"auction_ad_invited"`
	BrandedContentType     int  `json:"branded_content_type"`
	WithCommentFilterWords bool `json:"with_comment_filter_words"`
}

type tikwmAuthor struct {
	ID       string `json:"id"`
	UniqueID string `json:"unique_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
