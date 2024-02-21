package workflow

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
)

const ACTIVITY2_NAME = "ACTIVITY_2"

type Activity2Response struct {
	Message string
	TS      time.Time
}

type Activity2 struct {
	TemporalClient    client.Client
	Activity2Response *Activity2Response
}

func NewActivity2(temporalClient client.Client, updateProcessingPayment interface{}) *Activity2 {
	return &Activity2{
		TemporalClient: temporalClient,
		Activity2Response: &Activity2Response{
			Message: "Activity2's message",
			TS:      time.Now(),
		},
	}
}

func (a *Activity2) Activity2(ctx context.Context, request interface{}) error {
	// Do logic
	fmt.Println("Do Activity2 logic...")
	time.Sleep(time.Second * 1)

	// Send message by signal to do_something workflow
	if err := a.TemporalClient.SignalWorkflow(
		ctx,
		GetDoSomethingWorkflowName(X_REQUEST_ID),
		"",
		GetActivity2ChanName(X_REQUEST_ID),
		a.Activity2Response,
	); err != nil {
		fmt.Println("Error sending update processing payment to workflow:", err)
	}

	fmt.Println("Completed Activity2 task. Signaling workflow...")
	return nil
}

func GetActivity2ChanName(paymentID string) string {
	return fmt.Sprintf("%s|%s",
		GetDoSomethingWorkflowName(paymentID),
		ACTIVITY2_NAME,
	)
}
