package dbrepo

import (
	"context"
	"fmt"
	"github/somyaranjan99/basic-go-project/internal/reservationtypes"
	"time"
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
