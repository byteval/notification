package common

import (
	"notification/internal/domain/notification"
)

func ToResponse(n *notification.Notification) *Response {
	if n == nil {
		return nil
	}

	response := &Response{
		ID:        n.ID,
		CreatedAt: n.CreatedAt,
	}

	for _, receiver := range n.NotificationReceivers {
		response.Receivers = append(response.Receivers, ReceiverResponse{
			ID:     receiver.ID,
			Email:  receiver.Email,
			Status: string(receiver.Status),
		})
	}

	for _, attachment := range n.Attachments {
		response.Attachments = append(response.Attachments, AttachmentResponse{
			ID:          attachment.ID,
			FileName:    attachment.FileName,
			Size:        attachment.Size,
			ContentType: attachment.ContentType,
		})
	}

	return response
}
