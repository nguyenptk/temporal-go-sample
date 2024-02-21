package workflow

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	DO_SOMETHING_WORKFLOW_NAME = "DO_SOMETHING"
	X_REQUEST_ID               = "uuid"
)

func GetDoSomethingWorkflowName(paymentID string) string {
	return fmt.Sprintf("%s[%s]",
		DO_SOMETHING_WORKFLOW_NAME,
		paymentID,
	)
}

type DoSomethingInput struct {
	PaymentID string
}

func DoSomething(ctx workflow.Context, input DoSomethingInput) error {
	// Step 1: Activity1
	fmt.Println("Doing step 1...")
	var step1 *Activity1
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 15 * time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    time.Second,
				BackoffCoefficient: 2.0,
				MaximumAttempts:    3,
			},
		}),
		step1.Activity1,
	).Get(ctx, nil); err != nil {
		return err
	}
	fmt.Println("Done step 1")

	fmt.Println("Waiting for step 1 signal")
	signalStep1 := waitSignal[Activity1](ctx, GetDoSomethingWorkflowName(input.PaymentID), GetActivity1ChanName(input.PaymentID))

	// Step 2: Activity2
	fmt.Println("Doing step 2...")
	var step2 *Activity2
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 15 * time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    time.Second,
				BackoffCoefficient: 2.0,
				MaximumAttempts:    3,
			},
		}),
		step2.Activity2,
		signalStep1.Activity1Response,
	).Get(ctx, nil); err != nil {
		return err
	}
	fmt.Println("Done step 2")

	fmt.Println("Waiting for step 2 signal")
	signalStep2 := waitSignal[Activity2](ctx, GetDoSomethingWorkflowName(input.PaymentID), GetActivity2ChanName(input.PaymentID))

	// Step 3: Activity3
	fmt.Println("Doing step 3...")
	var step3 *Activity3
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 15 * time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    time.Second,
				BackoffCoefficient: 2.0,
				MaximumAttempts:    3,
			},
		}),
		step3.Activity3,
		signalStep2.Activity2Response,
	).Get(ctx, nil); err != nil {
		return err
	}
	fmt.Println("Done step 3")

	// Wait for step 3
	fmt.Println("Waiting for step 3 signal")
	signalStep3 := waitSignal[Activity3](ctx, GetDoSomethingWorkflowName(input.PaymentID), GetActivity3ChanName(input.PaymentID))

	// Step 4: Update the status=complete
	fmt.Println("Doing step 4...")
	var step4 *Activity4
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 15 * time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    time.Second,
				BackoffCoefficient: 2.0,
				MaximumAttempts:    3,
			},
		}),
		step4.Activity4,
		signalStep3.Activity3Response,
	).Get(ctx, nil); err != nil {
		return err
	}
	fmt.Println("Done step 4")

	// Wait for step 4
	fmt.Println("Waiting for step 4 signal")
	signalStep4 := waitSignal[Activity4](ctx, GetDoSomethingWorkflowName(input.PaymentID), GetActivity4ChanName(input.PaymentID))

	// Step 5: Call Merchant System
	fmt.Println("Doing step 5...")
	var step5 *Activity5
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 15 * time.Minute,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval:    time.Second,
				BackoffCoefficient: 2.0,
				MaximumAttempts:    3,
			},
		}),
		step5.Activity5,
		signalStep4.Activity4Response,
	).Get(ctx, nil); err != nil {
		return err
	}
	fmt.Println("Done step 5")

	// Wait for step 5
	fmt.Println("Finish workflow...")
	workflow.CompleteSession(ctx)

	return nil
}

func waitSignal[T any](ctx workflow.Context, workflowID string, channel string) T {
	println("WaitSignal:", channel)

	var signal T
	workflow.WithWorkflowID(ctx, workflowID)
	signalChan := workflow.GetSignalChannel(ctx, channel)

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &signal)
	})

	selector.Select(ctx)

	return signal
}
