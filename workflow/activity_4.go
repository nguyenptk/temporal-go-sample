package workflow

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
)

const ACTIVITY4_NAME = "ACTIVITY_4"

type Activity4Response struct {
	Message string
	TS      time.Time
}

type Activity4 struct {
	TemporalClient    client.Client
	Activity4Response *Activity4Response
}

func NewActivity4(temporalClient client.Client, updateProcessingPayment interface{}) *Activity4 {
	return &Activity4{
		TemporalClient: temporalClient,
		Activity4Response: &Activity4Response{
			Message: "Activity4's message",
			TS:      time.Now(),
		},
	}
}

func (a *Activity4) Activity4(ctx context.Context, request interface{}) error {
	// Do logic
	fmt.Println("Do Activity4 logic...")
	time.Sleep(time.Second * 1)

	// Send message by signal to do_something workflow
	if err := a.TemporalClient.SignalWorkflow(
		ctx,
		GetDoSomethingWorkflowName(X_REQUEST_ID),
		"",
		GetActivity4ChanName(X_REQUEST_ID),
		a.Activity4Response,
	); err != nil {
		fmt.Println("Error sending update processing payment to workflow:", err)
	}

	fmt.Println("Completed Activity4 task. Signaling workflow...")
	return nil
}

func GetActivity4ChanName(paymentID string) string {
	return fmt.Sprintf("%s|%s",
		GetDoSomethingWorkflowName(paymentID),
		ACTIVITY4_NAME,
	)
}
