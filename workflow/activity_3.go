package workflow

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
)

const ACTIVITY3_NAME = "ACTIVITY_3"

type Activity3Response struct {
	Message string
	TS      time.Time
}

type Activity3 struct {
	TemporalClient    client.Client
	Activity3Response *Activity3Response
}

func NewActivity3(temporalClient client.Client, updateProcessingPayment interface{}) *Activity3 {
	return &Activity3{
		TemporalClient: temporalClient,
		Activity3Response: &Activity3Response{
			Message: "Activity3's message",
			TS:      time.Now(),
		},
	}
}

func (a *Activity3) Activity3(ctx context.Context, request interface{}) error {
	// Do logic
	fmt.Println("Do Activity3 logic...")
	time.Sleep(time.Second * 1)

	// Send message by signal to do_something workflow
	if err := a.TemporalClient.SignalWorkflow(
		ctx,
		GetDoSomethingWorkflowName(X_REQUEST_ID),
		"",
		GetActivity3ChanName(X_REQUEST_ID),
		a.Activity3Response,
	); err != nil {
		fmt.Println("Error sending update processing payment to workflow:", err)
	}

	fmt.Println("Completed Activity3 task. Signaling workflow...")
	return nil
}

func GetActivity3ChanName(paymentID string) string {
	return fmt.Sprintf("%s|%s",
		GetDoSomethingWorkflowName(paymentID),
		ACTIVITY3_NAME,
	)
}
