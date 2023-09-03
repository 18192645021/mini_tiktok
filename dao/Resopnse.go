package dao

type FavoriteListResponse struct {
	StatusCode int32    `json:"status_code"`
	StatusMsg  string   `json:"status_msg,omitempty"`
	VideoList  []*Video `json:"video_list,omitempty"`
}

type CommentActionResponse struct {
	StatusCode int32    `json:"status_code"`
	StatusMsg  string   `json:"status_msg,omitempty"`
	Comment    *Comment `json:"comment,omitempty"`
}
type CommentListResponse struct {
	StatusCode  int32      `json:"status_code"`
	StatusMsg   string     `json:"status_msg,omitempty"`
	CommentList []*Comment `json:"comment_list,omitempty"`
}
