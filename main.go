package main

import (
	"context"
	"fmt"
	"net/http"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/nguyenptk/temporal-go-sample/workflow"
)

func NewTemporalClient() (client.Client, error) {
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		fmt.Println("\nfailed to create temporal client")
		return nil, err
	}
	return temporalClient, nil
}

func main() {
	// Init Temporal
	temporalClient, err := NewTemporalClient()
	temporalWorker := worker.New(temporalClient, "all-tasks", worker.Options{})
	temporalWorker.RegisterWorkflow(workflow.DoSomething)
	temporalWorker.RegisterActivity(workflow.NewActivity1(temporalClient, nil).Activity1)
	temporalWorker.RegisterActivity(workflow.NewActivity2(temporalClient, nil).Activity2)
	temporalWorker.RegisterActivity(workflow.NewActivity3(temporalClient, nil).Activity3)
	temporalWorker.RegisterActivity(workflow.NewActivity4(temporalClient, nil).Activity4)
	temporalWorker.RegisterActivity(workflow.NewActivity5(temporalClient, nil).Activity5)

	// Handle HTTP endpoints
	http.HandleFunc("/do-something", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/do-something endpoint is called, triggering the workflow...")

		// Submit DoSomething Workflow to Temporal
		_, err = temporalClient.ExecuteWorkflow(
			context.TODO(),
			client.StartWorkflowOptions{
				ID:        workflow.GetDoSomethingWorkflowName(workflow.X_REQUEST_ID),
				TaskQueue: "all-tasks",
			},
			workflow.DoSomething,
			workflow.DoSomethingInput{
				PaymentID: workflow.X_REQUEST_ID,
			},
		)
		if err != nil {
			fmt.Println("Error submitting DoSomething Workflow to Temporal:", err)
		}
	})

	// Start Temporal
	go func() {
		err = temporalWorker.Run(worker.InterruptCh())
		if err != nil {
			fmt.Println("Error running temporal worker:", err)
		}
	}()

	// Start HTTP Go server
	port := 6000
	fmt.Printf("Server listening on :%d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
