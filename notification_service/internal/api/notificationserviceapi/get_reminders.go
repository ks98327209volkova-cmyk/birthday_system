package notificationserviceapi

import (
	"context"
	"log"

	"notification_service/internal/models"
	pbmodels "notification_service/internal/pb/models"
)

func (s *NotificationServiceAPI) GetReminders(ctx context.Context, req *pbmodels.ReminderRequest) (*pbmodels.ReminderResponse, error) {
	log.Printf("GetReminders request for month: %s", req.Month)

	reminders, err := s.notificationService.GetAllReminders(ctx)
	if err != nil {
		return nil, err
	}

	return &pbmodels.ReminderResponse{
		Reminders: mapReminders(reminders),
	}, nil
}

func mapReminders(reminders []*models.Reminder) []*pbmodels.Reminder {
	protoReminders := make([]*pbmodels.Reminder, 0, len(reminders))

	for _, r := range reminders {
		protoReminders = append(protoReminders, &pbmodels.Reminder{
			Id:           r.ID,
			PersonId:     uint64(r.PersonID),
			PersonName:   r.PersonName,
			Birthday:     r.Birthday,
			ReminderDate: r.ReminderDate,
			DaysUntil:    uint32(r.DaysUntil),
		})
	}

	return protoReminders
}
