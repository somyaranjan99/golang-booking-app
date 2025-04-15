package dbrepo

import (
	"context"
	"fmt"
	"github/somyaranjan99/basic-go-project/internal/reservationtypes"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (m *RepositoryDBHandler) AllUsers() bool {
	return true
}
func (m *RepositoryDBHandler) BookReservation(res *reservationtypes.Reservation) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	_, err := m.Db.ExecContext(ctx, `INSERT INTO reservations(first_name, 
last_name, 
email,
phone ,
start_date, 
end_date , 
room_id,
created_at,
updated_at
) VALUES(?,?,?,?,?,?,?,?,?)`, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, res.CreatedAt, res.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return "success", nil
}
func (m *RepositoryDBHandler) SearchAvailabilityForAllRooms(start_date, end_date time.Time) ([]reservationtypes.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	rows, err := m.Db.QueryContext(ctx, `SELECT r.id,r.room_name FROM rooms r WHERE r.id IN (
	      SELECT rr.room_id FROM room_restrictions rr WHERE  ? NOT BETWEEN rr.start_date AND rr.end_date
		   AND ? NOT BETWEEN rr.start_date AND rr.end_date
	) `, start_date, end_date)
	if err != nil {
		return nil, err
	}
	var rooms []reservationtypes.Room
	for rows.Next() {
		var room reservationtypes.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil

}
func (m *RepositoryDBHandler) GetRoomByID(roomid int) (reservationtypes.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var room reservationtypes.Room
	getRoom := m.Db.QueryRowContext(ctx, "SELECT id,room_name,created_at,updated_at FROM rooms WHERE id=?", roomid)
	err := getRoom.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return room, err
	}
	return room, nil
}

func (m *RepositoryDBHandler) IsAuthenticatedUser(email, password string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rows := m.Db.QueryRowContext(ctx, `SELECT Id
															first_name,
															last_name,
															email,password,
															access_level FROM users WHERE email=? and password=?`, email, password)
	var user reservationtypes.User
	err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.AccessLevel)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil

}

func (m *RepositoryDBHandler) SignupUser(user *reservationtypes.User) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 	id int AI PK
	// first_name varchar(255)
	// last_name varchar(255)
	// email varchar(255)
	// password varchar(60)
	// access_level int
	// created_at datetime
	// updated_at datetime
	res, err := m.Db.ExecContext(ctx, `INSERT INTO users
	(first_name,last_name,email,password,access_level,created_at,updated_at) VALUES
	(?,?,?,?,?,?,?)`, user.FirstName, user.LastName, user.Email, user.Password, user.AccessLevel, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return false, err
	}
	fmt.Println(res)
	return true, nil

}
func (m *RepositoryDBHandler) GetUserByEmail(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	row := m.Db.QueryRowContext(ctx, "SELECT email FROM users WHERE email=?", email)

	user := reservationtypes.User{}

	err := row.Scan(&user.Email)
	if err != nil {
		return false, err
	}
	if user.Email == "" {
		return false, nil
	}
	return true, nil
}
