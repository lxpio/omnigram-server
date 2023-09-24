package schema

type FavoriteBook struct {
	ID        int64 `json:"id" gorm:"primary_key,comment:进度ID"`
	BookID    int64 `json:"book_id" gorm:"uniqueIndex:uni_idx_fav_bookid_userid,comment:书籍ID"`
	UserID    int64 `json:"user_id" gorm:"uniqueIndex:uni_idx_fav_bookid_userid,comment:用户ID"`
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime"`

	// BookTitle string `json:"book_title"`
	// Book      *Book  `json:"book"`
}
