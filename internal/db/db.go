package db

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Структура, хранящая переменную для доступа к БД. Также реализует Querier.
type DB struct {
	db *gorm.DB
}

// Интерфейс, содержащий методы для работы с БД.
type Querier interface {
	New(username string, url string, shorty string) error
	ListShorties(username string) ([]Shorties, error)
	ShortyInfo(shortyURL string) (Shorties, error)
	DeleteShorty(shortyURL string) error
}

// Модель для работы с БД.
type Shorties struct {
	ID          uint      `gorm:"primaryKey;not null;autoIncrement"`
	Username    string    `gorm:"not null"`
	OriginalURL string    `gorm:"not null"`
	ShortyURL   string    `gorm:"not null;unique"`
	DateCreated time.Time `gorm:"not null"`
}

// Инициализирует БД, возвращает структуру *DB или одну из возможных ошибок.
func InitDB() (*DB, error) {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	err = db.Transaction(func(t *gorm.DB) error {
		ok := t.Migrator().HasTable(&Shorties{})
		if !ok {
			err := t.Migrator().CreateTable(&Shorties{})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

// Создает новую запись в БД. Используется в функции handlers.New.
func (d *DB) New(username string, url string, shorty string) error {
	var shortiesCount int64

	t := d.db.Table("shorties").Where("username = ?", username).Where("original_url = ?", url).Count(&shortiesCount)
	if shortiesCount != 0 {
		return gorm.ErrDuplicatedKey
	}
	if t.Error != nil && t.Error != gorm.ErrRecordNotFound {
		return t.Error
	}

	t = d.db.Table("shorties").Where("username = ?", username).Count(&shortiesCount)
	if shortiesCount >= 5 {
		return gorm.ErrCheckConstraintViolated
	}
	if t.Error != nil && t.Error != gorm.ErrRecordNotFound {
		return t.Error
	}

	t = d.db.Table("shorties").Create(&Shorties{
		Username:    username,
		OriginalURL: url,
		ShortyURL:   shorty,
		DateCreated: time.Now().UTC(),
	})
	if t.Error != nil {
		return t.Error
	}

	return nil
}

// Возвращает все доступные цитаты по имени пользователя. Используется в функциях handlers.ListShorties и handlers.DeleteShorty.
func (d *DB) ListShorties(username string) ([]Shorties, error) {
	var shorties []Shorties

	t := d.db.Table("shorties").Where("username = ?", username).Find(&shorties)
	if len(shorties) == 0 {
		return []Shorties{}, gorm.ErrRecordNotFound
	}
	if t.Error != nil && t.Error != gorm.ErrRecordNotFound {
		return []Shorties{}, t.Error
	}
	return shorties, nil
}

// Возвращает конкретную цитату по идентификатору сокращенной ссылки. Используется в функции handlers.ShortyInfo.
func (d *DB) ShortyInfo(shortyURL string) (Shorties, error) {
	var shorty Shorties

	t := d.db.Table("shorties").Where("shorty_url = ?", shortyURL).Take(&shorty)
	if shorty.OriginalURL == "" {
		return Shorties{}, gorm.ErrRecordNotFound
	}
	if t.Error != nil && t.Error != gorm.ErrRecordNotFound {
		return Shorties{}, t.Error
	}
	return shorty, nil
}

// Удаляет цитату из БД. Используется в функции handlers.DeleteSelectedShorty.
func (d *DB) DeleteShorty(shortyURL string) error {
	var shorty Shorties

	t := d.db.Table("shorties").Where("shorty_url = ?", shortyURL).Take(&shorty).Delete(&Shorties{})
	if shorty.OriginalURL == "" {
		return gorm.ErrRecordNotFound
	}
	if t.Error != nil && t.Error != gorm.ErrRecordNotFound {
		return t.Error
	}
	return nil
}
