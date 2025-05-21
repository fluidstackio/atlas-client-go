package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	atlas "github.com/fluidstackio/atlas-client-go/v1alpha1"
	"github.com/google/uuid"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
)

func main() {
	client := newClient()

	if os.Getenv("ATLAS_PROJECT_ID") == "" {
		log.Fatal("Missing required environment variable ATLAS_PROJECT_ID")
	}
	projectID := uuid.MustParse(os.Getenv("ATLAS_PROJECT_ID"))

	ctx := context.Background()

	instance, err := createInstance(ctx, client, projectID, "example-instance-01", "cpu.2x")
	if err != nil {
		log.Fatalf("Failed to create instance: %v", err)
	}

	if err := stopInstance(ctx, client, projectID, instance.Id); err != nil {
		log.Fatalf("Failed to stop instance: %v", err)
	}

	if err := startInstance(ctx, client, projectID, instance.Id); err != nil {
		log.Fatalf("Failed to start instance: %v", err)
	}

	if err := deleteInstance(ctx, client, projectID, instance.Id); err != nil {
		log.Fatalf("Failed to delete instance: %v", err)
	}
}

func newClient() *atlas.ClientWithResponses {
	token := os.Getenv("ATLAS_API_TOKEN")
	if token == "" {
		log.Fatal("Missing required environment variable ATLAS_API_TOKEN")
	}

	url := os.Getenv("ATLAS_REGION_URL")
	if url == "" {
		log.Fatal("Missing required environment variable ATLAS_REGION_URL")
	}

	bearerAuth, err := securityprovider.NewSecurityProviderBearerToken(token)
	if err != nil {
		log.Fatal(err)
	}

	client, err := atlas.NewClientWithResponses(url+"/api/v1alpha1/", atlas.WithRequestEditorFn(bearerAuth.Intercept))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func createInstance(ctx context.Context, client *atlas.ClientWithResponses, projectID uuid.UUID, name, instanceType string) (*atlas.Instance, error) {
	log.Printf("Creating instance")

	resp, err := client.PostInstancesWithResponse(ctx, &atlas.PostInstancesParams{
		XPROJECTID: projectID,
	}, atlas.InstancesPostRequest{
		Name: name,
		Type: instanceType,
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusCreated {
		return nil, errors.New(resp.Status())
	}

	instance := resp.JSON201

	log.Printf("Created instance %v", instance.Id)

	for instance.State != atlas.InstanceStateRunning {
		log.Printf("Waiting for instance %v to start. Current state: %v", instance.Id, strings.ToUpper(string(instance.State)))
		time.Sleep(10 * time.Second)

		resp, err := client.GetInstancesIdWithResponse(ctx, resp.JSON201.Id, &atlas.GetInstancesIdParams{
			XPROJECTID: projectID,
		})
		if err != nil {
			return nil, err
		}
		if resp.StatusCode() != http.StatusOK {
			return nil, errors.New(resp.Status())
		}

		instance = resp.JSON200
	}

	log.Printf("Started instance %v", instance.Id)
	return instance, nil
}

func stopInstance(ctx context.Context, client *atlas.ClientWithResponses, projectID, instanceID uuid.UUID) error {
	log.Printf("Stopping instance %v", instanceID)

	stopResp, err := client.PostInstancesIdActionsStopWithResponse(ctx, instanceID, &atlas.PostInstancesIdActionsStopParams{
		XPROJECTID: projectID,
	})
	if err != nil {
		return err
	}
	if stopResp.StatusCode() != http.StatusAccepted {
		return errors.New(stopResp.Status())
	}

	resp, err := client.GetInstancesIdWithResponse(ctx, instanceID, &atlas.GetInstancesIdParams{
		XPROJECTID: projectID,
	})
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.Status())
	}

	instance := resp.JSON200
	for instance.State != atlas.InstanceStateStopped {
		log.Printf("Waiting for instance %v to stop. Current state: %v", instance.Id, strings.ToUpper(string(instance.State)))
		time.Sleep(5 * time.Second)

		resp, err := client.GetInstancesIdWithResponse(ctx, instance.Id, &atlas.GetInstancesIdParams{
			XPROJECTID: projectID,
		})
		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			return errors.New(resp.Status())
		}

		instance = resp.JSON200
	}

	log.Printf("Stopped instance %v", instanceID)
	return nil
}

func startInstance(ctx context.Context, client *atlas.ClientWithResponses, projectID, instanceID uuid.UUID) error {
	log.Printf("Starting instance %v", instanceID)

	startResp, err := client.PostInstancesIdActionsStartWithResponse(ctx, instanceID, &atlas.PostInstancesIdActionsStartParams{
		XPROJECTID: projectID,
	})
	if err != nil {
		return err
	}
	if startResp.StatusCode() != http.StatusAccepted {
		log.Printf("here")
		return errors.New(startResp.Status())
	}

	resp, err := client.GetInstancesIdWithResponse(ctx, instanceID, &atlas.GetInstancesIdParams{
		XPROJECTID: projectID,
	})
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(resp.Status())
	}

	instance := resp.JSON200
	for instance.State != atlas.InstanceStateRunning {
		log.Printf("Waiting for instance %v to start. Current state: %v", instance.Id, strings.ToUpper(string(instance.State)))
		time.Sleep(5 * time.Second)

		resp, err := client.GetInstancesIdWithResponse(ctx, instanceID, &atlas.GetInstancesIdParams{
			XPROJECTID: projectID,
		})
		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			return errors.New(resp.Status())
		}

		instance = resp.JSON200
	}

	log.Printf("Started instance %v", instanceID)
	return nil
}

func deleteInstance(ctx context.Context, client *atlas.ClientWithResponses, projectID, instanceID uuid.UUID) error {
	log.Printf("Deleting instance %v", instanceID)

	respDelete, errDelete := client.DeleteInstancesIdWithResponse(ctx, instanceID, &atlas.DeleteInstancesIdParams{
		XPROJECTID: projectID,
	})
	if errDelete != nil {
		return errDelete
	}
	if respDelete.StatusCode() != http.StatusNoContent {
		return errors.New(respDelete.Status())
	}

	respGet, errGet := client.GetInstancesIdWithResponse(ctx, instanceID, &atlas.GetInstancesIdParams{
		XPROJECTID: projectID,
	})
	if errGet != nil {
		log.Fatal(errDelete)
	}

	for respGet.StatusCode() != http.StatusNotFound {
		log.Printf("Waiting for instance to be deleted. Current state: %v", strings.ToUpper(string(respGet.JSON200.State)))
		time.Sleep(5 * time.Second)

		respGet, errGet = client.GetInstancesIdWithResponse(ctx, instanceID, &atlas.GetInstancesIdParams{
			XPROJECTID: projectID,
		})
		if errGet != nil {
			return errGet
		}
	}

	log.Printf("Deleted instance %v", instanceID)
	return nil
}
