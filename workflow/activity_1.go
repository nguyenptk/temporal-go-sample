package workflow

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
)

const ACTIVITY1_NAME = "ACTIVITY_1"

type Activity1Response struct {
	Message string
	TS      time.Time
}

type Activity1 struct {
	TemporalClient    client.Client
	Activity1Response *Activity1Response
}

func NewActivity1(temporalClient client.Client, updateProcessingPayment interface{}) *Activity1 {
	return &Activity1{
		TemporalClient: temporalClient,
		Activity1Response: &Activity1Response{
			Message: "Activity1's message",
			TS:      time.Now(),
		},
	}
}

func (a *Activity1) Activity1(ctx context.Context, request interface{}) error {
	// Do logic
	fmt.Println("Do Activity1 logic...")
	time.Sleep(time.Second * 1)

	// Send message by signal to do_something workflow
	if err := a.TemporalClient.SignalWorkflow(
		ctx,
		GetDoSomethingWorkflowName(X_REQUEST_ID),
		"",
		GetActivity1ChanName(X_REQUEST_ID),
		a.Activity1Response,
	); err != nil {
		fmt.Println("Error sending update processing payment to workflow:", err)
	}

	fmt.Println("Completed Activity1 task. Signaling workflow...")
	return nil
}

func GetActivity1ChanName(paymentID string) string {
	return fmt.Sprintf("%s|%s",
		GetDoSomethingWorkflowName(paymentID),
		ACTIVITY1_NAME,
	)
}
