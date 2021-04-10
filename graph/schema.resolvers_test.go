package graph

import (
	"SpamhausResponseApi/db"
	"SpamhausResponseApi/model"
	"context"
	"reflect"
	"testing"
	"time"

	mocket "github.com/selvatico/go-mocket"
)

var curTime time.Time = time.Now()

type ActualExpected struct {
	source   []*model.IP
	expected []*model.IPDetails
}

func buildActualExpected() ActualExpected {
	return ActualExpected{
		source: []*model.IP{{Address: "46.102.177.99"}, {Address: "31.14.65.6"}},
		expected: []*model.IPDetails{
			{
				UUID:         "f1908fef-6d11-4b70-a22a-291274cdc4e3",
				ResponseCode: "127.0.0.2, 127.0.0.9",
				IPAddress:    "46.102.177.99",
				CreatedAt:    curTime,
				UpdatedAt:    curTime,
			},
			{
				UUID:         "26e3eca8-9395-4b74-84e4-cc3822fb292f",
				ResponseCode: "127.0.0.2, 127.0.0.9",
				IPAddress:    "31.14.65.6",
				CreatedAt:    curTime,
				UpdatedAt:    curTime,
			},
		},
	}
}

func TestInvalidIP(t *testing.T) {
	input1 := []*model.IP{{Address: "1111"}}
	input2 := []*model.IP{{Address: "185,77,248,1"}}
	input3 := []*model.IP{{Address: "192.168.1.999"}}
	input4 := []*model.IP{{Address: "192.168.1.aaa"}}

	var ctx context.Context = context.Background()
	mr := mutationResolver{
		&Resolver{},
	}

	if _, err := mr.Mutation().Enqueue(ctx, input1); err == nil {
		t.Logf("unable to validate invalid ip. %s", err)
		t.FailNow()
	}
	if _, err := mr.Mutation().Enqueue(ctx, input2); err == nil {
		t.Logf("unable to validate invalid ip. %s", err)
		t.FailNow()
	}
	if _, err := mr.Mutation().Enqueue(ctx, input3); err == nil {
		t.Logf("unable to validate invalid ip. %s", err)
		t.FailNow()
	}
	if _, err := mr.Mutation().Enqueue(ctx, input4); err == nil {
		t.Logf("unable to validate invalid ip. %s", err)
		t.FailNow()
	}
}

func TestEnqueue(t *testing.T) {
	ctx := context.Background()
	mr := mutationResolver{
		&Resolver{},
	}

	db.MockConnect()
	catcher := mocket.Catcher
	MockEnqueue(catcher)

	details := buildActualExpected()
	actual, err := mr.Enqueue(ctx, details.source)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(actual) != len(details.expected) {
		t.Log("actual and expected slice lengths are uneqaul.")
		t.FailNow()
	}
	for i := 0; i < len(actual); i++ {
		//since these are updates lets make sure that the updatedat fields are not equal
		if actual[i].UpdatedAt == details.expected[i].UpdatedAt {
			t.Log("updated values are equal")
			t.FailNow()
		}
		//update the records that were updated to make sure all else is equal
		details.expected[i].ResponseCode = actual[i].ResponseCode
		details.expected[i].UpdatedAt = actual[i].UpdatedAt

		if !reflect.DeepEqual(actual[i], details.expected[i]) {
			t.Logf("actual: %v, and expected: %v. are not equal", actual[i], details.expected[i])
			t.FailNow()
		}
	}
}

func MockEnqueue(catcher *mocket.MockCatcher) {
	detailsReply1 := []map[string]interface{}{
		{
			"uuid":          "f1908fef-6d11-4b70-a22a-291274cdc4e3",
			"response_code": "127.0.0.2, 127.0.0.9",
			"ip_address":    "46.102.177.99",
			"created_at":    curTime,
			"updated_at":    curTime,
		},
	}
	catcher.NewMock().WithQuery(`SELECT * FROM "ip_details"  WHERE (ip_address = 46.102.177.99)`).WithReply(detailsReply1)
	detailsReply2 := []map[string]interface{}{
		{
			"uuid":          "26e3eca8-9395-4b74-84e4-cc3822fb292f",
			"response_code": "127.0.0.2, 127.0.0.9",
			"ip_address":    "31.14.65.6",
			"created_at":    curTime,
			"updated_at":    curTime,
		},
	}
	catcher.NewMock().WithQuery(`SELECT * FROM "ip_details"  WHERE (ip_address = 31.14.65.6)`).WithReply(detailsReply2)

	catcher.NewMock().WithQuery(`UPDATE "ip_details" SET "created_at" = ?, "ip_address" = ?, "response_code" = ?, "updated_at" = ?, "uuid" = ?  WHERE (ip_address = ?)`)
}

func TestGetIPDetails(t *testing.T) {
	var ctx context.Context = context.Background()
	qr := queryResolver{
		&Resolver{},
	}

	db.MockConnect()

	catcher := mocket.Catcher
	MockGetIPQueries(catcher)

	details := buildActualExpected()
	actual, err := qr.GetIPDetails(ctx, *details.source[0])
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(actual, details.expected[0]) {
		t.Logf("actual: %v, and expected: %v. are not equal", actual, details.expected[0])
		t.FailNow()

	}
}

func MockGetIPQueries(catcher *mocket.MockCatcher) {
	detailsReply := []map[string]interface{}{
		{
			"uuid":          "f1908fef-6d11-4b70-a22a-291274cdc4e3",
			"response_code": "127.0.0.2, 127.0.0.9",
			"ip_address":    "46.102.177.99",
			"created_at":    curTime,
			"updated_at":    curTime,
		},
	}
	catcher.NewMock().WithQuery(`SELECT * FROM "ip_details"  WHERE (ip_address = 46.102.177.99)`).WithReply(detailsReply)
}
