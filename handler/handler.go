package handler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/vithsutra/biometric-message-consumer/models"
)

type consumerGroupHandler struct {
	dbRepo models.AttendanceInterface
}

func NewConsumerGroupHandler(dbRepo models.AttendanceInterface) *consumerGroupHandler {
	return &consumerGroupHandler{
		dbRepo,
	}
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		req := new(models.AttendanceLogRequest)

		if err := json.Unmarshal(msg.Value, req); err != nil {
			sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, "")
			sess.Commit()
			log.Println("error occurred while decoding the attendance log message request, Error: ", err.Error())
			continue
		}

		studentId, err := h.dbRepo.GetStudentId(req.DeviceId, req.StudentDeviceId)

		if err != nil {
			sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, "")
			sess.Commit()
			log.Println("error occurred with database while getting the student, Device Id: ", req.DeviceId, " Error: ", err.Error())
			continue
		}

		t, err := time.Parse("2006-01-02T15:04:05", req.TimeStamp)

		if err != nil {
			sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, "")
			sess.Commit()
			log.Println("error occurred with database while parsing the time stamp, Device Id: ", req.DeviceId, " Error: ", err.Error())
			continue
		}

		date := t.Format("2006-01-02")

		tm := t.Format("15:04")

		logStatus, err := h.dbRepo.CheckLoginOrLogout(studentId, date)

		if err != nil {
			sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, "")
			sess.Commit()
			log.Println("error occurred with database while checking login or logout, Device Id: ", req.DeviceId, " Error: ", err.Error())
			continue
		}

		if !logStatus {
			attendanceLog := models.AttendanceLog{
				StudentId:       studentId,
				StudentDeviceId: req.StudentDeviceId,
				DeviceId:        req.DeviceId,
				Date:            date,
				Login:           tm,
				Logout:          "25:00",
			}
			if err := h.dbRepo.InsertAttendanceLog(&attendanceLog); err != nil {
				sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, "")
				sess.Commit()
				log.Println("error occurred with database while inserting the attendance log, Device Id: ", req.DeviceId, " Error: ", err.Error())
				continue
			}
		} else {
			if err := h.dbRepo.UpdateAttendanceLog(studentId, date, tm); err != nil {
				sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, "")
				sess.Commit()
				log.Println("error occurred with database while updating the attendance log, Device Id: ", req.DeviceId, " Error: ", err.Error())
				continue
			}
		}

		sess.MarkOffset(msg.Topic, msg.Partition, msg.Offset+1, "")
		sess.Commit()
	}
	return nil
}
