package repository

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vithsutra/biometric-message-consumer/models"
)

type postgresRepository struct {
	dbConn *pgxpool.Pool
}

func NewPostgresRepository(dbConn *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{
		dbConn,
	}
}

func (repo *postgresRepository) GetStudentId(deviceId string, studenDeviceId string) (string, error) {
	query := `SELECT student_id FROM ` + strings.ToLower(deviceId) + ` WHERE student_unit_id=$1`
	var studentId string
	err := repo.dbConn.QueryRow(context.Background(), query, studenDeviceId).Scan(&studentId)
	return studentId, err
}

func (repo *postgresRepository) CheckLoginOrLogout(studentId string, date string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM attendance WHERE date=$1 AND student_id=$2 and logout=$3 )`
	var logStatus bool
	err := repo.dbConn.QueryRow(context.Background(), query, date, studentId, "25:00").Scan(&logStatus)
	return logStatus, err
}

func (repo *postgresRepository) InsertAttendanceLog(attendanceLog *models.AttendanceLog) error {
	query := `INSERT INTO attendance (student_id,student_unit_id,unit_id,date,login,logout) VALUES ($1,$2,$3,$4,$5,$6)`

	_, err := repo.dbConn.Exec(
		context.Background(),
		query,
		attendanceLog.StudentId,
		attendanceLog.StudentDeviceId,
		attendanceLog.DeviceId,
		attendanceLog.Date,
		attendanceLog.Login,
		attendanceLog.Logout,
	)

	return err
}

func (repo *postgresRepository) UpdateAttendanceLog(studentId string, date string, logout string) error {
	query := `UPDATE attendance SET logout=$4 WHERE student_id=$1 AND date=$2 AND logout=$3`
	_, err := repo.dbConn.Exec(
		context.Background(),
		query,
		studentId,
		date,
		"25:00",
		logout,
	)
	return err
}
