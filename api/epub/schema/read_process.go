package schema

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserBookShip 用户阅读进度表
type ReadProgress struct {
	ID        int64   `json:"id" gorm:"primary_key,comment:进度ID"`
	BookID    int64   `json:"book_id" gorm:"uniqueIndex:uni_idx_read_bookid_userid,comment:书籍ID"`
	UserID    int64   `json:"user_id" gorm:"uniqueIndex:uni_idx_read_bookid_userid,comment:用户ID"`
	Progress  float32 `json:"progress" gorm:"comment:阅读进度百分比"`
	StartDate int64   `json:"start_date" gorm:"comment:阅读开始日期"`

	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime,comment:阅读更新时间"`
	//预计完成日期
	ExptEndDate int64 `json:"expt_end_date" gorm:"comment:预计完成日期"`
	EndDate     int64 `json:"end_date" gorm:"comment:阅读结束日期"`

	// ChapterID int64 `json:"chapter_id" gorm:"comment:当前章节"`
	//章节定位
	ChapterPos string `json:"chapter_pos" gorm:"comment:章节定位"`
}

func (p *ReadProgress) BeforeCreate(db *gorm.DB) error {
	p.StartDate = time.Now().Unix()
	p.ExptEndDate = time.Now().AddDate(0, 0, 15).Unix() //根据数据推测完成时间，
	// if p.

	if p.BookID == 0 || p.UserID == 0 {
		return errors.New("book_id or user_id not set")
	}

	return nil
}

func (p *ReadProgress) First(db *gorm.DB) error {
	return db.First(p).Error
}

func (b *Book) CreateReadProcess(db *gorm.DB, userID int64) (*ReadProgress, error) {

	proc := &ReadProgress{
		BookID:      b.ID,
		UserID:      userID,
		Progress:    0,
		StartDate:   time.Now().Unix(),
		ExptEndDate: time.Now().AddDate(0, 0, 15).Unix(), //根据数据推测完成时间，
		EndDate:     0,
		ChapterPos:  "",
	}

	err := db.Create(proc).Error

	return proc, err
}

func (p *ReadProgress) Update(db *gorm.DB) error {

	return db.Table(`read_progresses`).Where(`user_id = ? AND book_id = ?`, p.UserID, p.BookID).Select("progress", "update_at", "chapter_pos").Updates(p).Error

}

func (p *ReadProgress) Upsert(db *gorm.DB) error {

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "book_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"progress", "updated_at", "chapter_pos"}),
	}).Create(p).Error

}
