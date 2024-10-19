package database

import (
	"SkinRest/config"
	"SkinRest/pkg/models"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

//go:generate mockgen -source=db.go -destination=mocks/mock.go

type ApiHandler interface {
	CreateNewUser(user *models.User) error
	UpdateUserToken(user *models.User) (string, error)
	GetInfoUser(user *models.User) (*models.UserData, error)
	GetUserFromToken(token string) (*models.UserData, error)
	AddNewSkin(userData *models.UserData, skin *models.Skin) (*models.SkinData, error)
	GetUserSkins(userData *models.UserData) ([]models.SkinData, error)
	GetUserSkin(userData *models.UserData, id int) (*models.SkinData, error)
	DeleteUserSkin(userData *models.UserData, id int) error
}

type AppContext struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func New() *sql.DB {
	cfg := config.GetConfig()

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.SSL) // set database connection string
	fmt.Println(connStr)

	db, err := sql.Open(cfg.Database.Driver, connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS public.userstable (
        user_id SERIAL PRIMARY KEY,
        login VARCHAR(20) NOT NULL,
        password VARCHAR(255) NOT NULL,
        token VARCHAR(255) NOT NULL,
        CONSTRAINT userstable_login_key UNIQUE (login)
    )`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS public.skinstable (
        skin_id SERIAL PRIMARY KEY,
		owner_name VARCHAR(20) NOT NULL,
        skin_name VARCHAR(30) NOT NULL,
        skin_type VARCHAR(10) NOT NULL,
        skin_src VARCHAR(255) NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (m *AppContext) CreateNewUser(user *models.User) error {

	var exists int
	if err := m.DB.QueryRow("SELECT COUNT(1) FROM userstable WHERE login = $1", user.Login).Scan(&exists); err != nil {
		return err
	}

	if exists > 0 {
		return models.ErrAlrRegistered
	}

	token := GenerateNewToken(user)

	passwordHash, err := GetPasswordHash(user.Password)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec("INSERT INTO userstable (login, password, token) VALUES ($1, $2, $3)", user.Login, passwordHash, token)

	if err != nil {
		return err
	}

	return nil

}

func (m *AppContext) UpdateUserToken(user *models.User) (string, error) {
	newToken := GenerateNewToken(user)
	res, err := m.DB.Exec("UPDATE userstable SET token = $1 WHERE login = $2", newToken, user.Login)

	if err != nil {
		return "", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", models.ErrUserNotFound
	}

	return newToken, nil
}

func (m *AppContext) GetInfoUser(user *models.User) (*models.UserData, error) {
	var userData models.UserData

	err := m.DB.QueryRow("SELECT user_id, login, password, token FROM userstable WHERE login = $1", user.Login).Scan(&userData.Id, &userData.Login, &userData.Password, &userData.Token)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	if ValidatePasswordHash(user.Password, userData.Password) {
		return &userData, nil
	}

	return nil, models.ErrUserNotFound

}

func (m *AppContext) GetUserFromToken(token string) (*models.UserData, error) {
	var userData models.UserData

	err := m.DB.QueryRow("SELECT user_id, login, password, token FROM userstable WHERE token = $1", token).Scan(&userData.Id, &userData.Login, &userData.Password, &userData.Token)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &userData, nil
}

func (m *AppContext) AddNewSkin(userData *models.UserData, skin *models.Skin) (*models.SkinData, error) {
	var skin_id int

	err := m.DB.QueryRow("INSERT INTO skinstable (owner_name, skin_name, skin_type, skin_src) VALUES ($1, $2, $3, $4) RETURNING skin_id", userData.Login, skin.Name, skin.Type, skin.Src).Scan(&skin_id)

	if err != nil {
		return nil, err
	}

	skinData := &models.SkinData{
		Id:   skin_id,
		Name: skin.Name,
		Type: skin.Type,
		Src:  skin.Src,
	}

	return skinData, nil
}

func (m *AppContext) GetUserSkins(userData *models.UserData) ([]models.SkinData, error) {
	var skins []models.SkinData

	rows, err := m.DB.Query("SELECT skin_id, skin_name, skin_type, skin_src FROM skinstable WHERE owner_name = $1", userData.Login)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var skin models.SkinData
		if err := rows.Scan(&skin.Id, &skin.Name, &skin.Type, &skin.Src); err != nil {
			return nil, err
		}
		skins = append(skins, skin)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return skins, nil
}

func (m *AppContext) GetUserSkin(userData *models.UserData, id int) (*models.SkinData, error) {
	var skinData models.SkinData

	err := m.DB.QueryRow("SELECT skin_id, skin_name, skin_type, skin_src FROM skinstable WHERE skin_id = $1", id).Scan(&skinData.Id, &skinData.Name, &skinData.Type, &skinData.Src)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrSkinNotFound
		}
		return nil, err
	}

	return &skinData, nil

}

func (m *AppContext) DeleteUserSkin(userData *models.UserData, id int) error {
	res, err := m.DB.Exec("DELETE FROM skinstable WHERE skin_id = $1", id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrSkinNotFound
	}

	return nil
}
