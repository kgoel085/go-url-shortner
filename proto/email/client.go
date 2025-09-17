package email

import (
	"context"
	"errors"

	email "kgoel085.com/url-shortner/grpc/email" // Generated via 'protoc --go_out=grpc/email --go-grpc_out=grpc/email  proto/email/email.proto'
	"kgoel085.com/url-shortner/proto"
	"kgoel085.com/url-shortner/utils"
)

type GrpcSendEmailRequest struct {
	ToEmail   string
	Subject   string
	Content   string
	ProjectId string
}

func SendEmailViaGRPC(req GrpcSendEmailRequest) error {
	client, _ := proto.ClientManager.Get(string(proto.EmailServiceClientType))
	if client == nil {
		return errors.New("email service client not available")
	}

	utils.Log.Info("Using gRPC client to send email to ", req.ToEmail)
	emailClient := email.NewEmailServiceClient(client)

	utils.Log.Info("Sending email via gRPC to ", emailClient)
	resp, err := emailClient.SendEmail(context.Background(), &email.SendEmailRequest{
		ToEmail:   req.ToEmail,
		Subject:   req.Subject,
		Content:   req.Content,
		ProjectId: req.ProjectId,
	})
	utils.Log.Info("gRPC response: ", resp, err, req)
	if err != nil {
		return err
	}
	if resp.Id == "" {
		return errors.New("failed to send email")
	}

	return nil
}
