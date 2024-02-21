package workflow

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
)

const ACTIVITY5_NAME = "ACTIVITY_5"

type Activity5Response struct {
	Message string
	TS      time.Time
}

type Activity5 struct {
	TemporalClient    client.Client
	Activity5Response *Activity5Response
}

func NewActivity5(temporalClient client.Client, updateProcessingPayment interface{}) *Activity5 {
	return &Activity5{
		TemporalClient: temporalClient,
		Activity5Response: &Activity5Response{
			Message: "Activity5's message",
			TS:      time.Now(),
		},
	}
}

func (a *Activity5) Activity5(ctx context.Context, request interface{}) error {
	// Do logic
	fmt.Println("Do Activity5 logic...")
	time.Sleep(time.Second * 1)

	// Send message by signal to do_something workflow
	if err := a.TemporalClient.SignalWorkflow(
		ctx,
		GetDoSomethingWorkflowName(X_REQUEST_ID),
		"",
		GetActivity5ChanName(X_REQUEST_ID),
		a.Activity5Response,
	); err != nil {
		fmt.Println("Error sending update processing payment to workflow:", err)
	}

	fmt.Println("Completed Activity5 task. Signaling workflow...")
	return nil
}

func GetActivity5ChanName(paymentID string) string {
	return fmt.Sprintf("%s|%s",
		GetDoSomethingWorkflowName(paymentID),
		ACTIVITY5_NAME,
	)
}
