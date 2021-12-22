package googlecloud

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/spanner"
	"github.com/pottava/jaguer-cn-lottery/api/internal/lib"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Swag struct {
	ID    int64
	Name  string
	Stock int64
}

// ListSwags returns SWAG information
func ListSwags(ctx context.Context, projectID, spannerInstance, spannerDatabase string) ([]*Swag, error) {
	opts := []option.ClientOption{}
	if _, err := os.Stat(lib.Config.GcloudCreds); err == nil {
		opts = append(opts, option.WithCredentialsFile(lib.Config.GcloudCreds))
	}
	database := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		projectID, spannerInstance, spannerDatabase)
	client, err := spanner.NewClient(ctx, database, opts...)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	stmt := spanner.Statement{SQL: `SELECT id, name, stock FROM swags`}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	swags := []*Swag{}
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var id, stock int64
		var name string
		if err := row.Columns(&id, &name, &stock); err != nil {
			return nil, err
		}
		swags = append(swags, &Swag{id, name, stock})
	}
	return swags, nil
}
