package models

type AttendanceLogRequest struct {
	DeviceId        string `json:"did"`
	StudentDeviceId string `json:"sid"`
	TimeStamp       string `json:"tmstmp"`
}

type AttendanceLog struct {
	StudentId       string
	StudentDeviceId string
	DeviceId        string
	Date            string
	Login           string
	Logout          string
}

type AttendanceInterface interface {
	GetStudentId(deviceId string, studenDeviceId string) (string, error)
	CheckLoginOrLogout(studentId string, date string) (bool, error)
	InsertAttendanceLog(attendanceLog *AttendanceLog) error
	UpdateAttendanceLog(studentId string, date string, logout string) error
}
