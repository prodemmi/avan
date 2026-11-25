package notifiers

import (
	"avan/internal/db/sqlc"
	"errors"
	"reflect"
	"testing"
)

func TestWeightedFailoverAndSend(t *testing.T) {
	primary := sqlc.Provider{ID: 1, WorkspaceID: 7, AppID: 9, ConfigType: 0}
	candidates := []sqlc.Provider{
		{ID: 2, WorkspaceID: 7, AppID: 9, ConfigType: 0, Priority: 10, Weight: 1},
		{ID: 3, WorkspaceID: 7, AppID: 9, ConfigType: 0, Priority: 1, Weight: 1},
		{ID: 4, WorkspaceID: 7, AppID: 9, ConfigType: 0, Priority: 1, Weight: 3},
		{ID: 5, WorkspaceID: 7, AppID: 99, ConfigType: 0, Priority: 0, Weight: 100},
		{ID: 6, WorkspaceID: 7, AppID: 9, ConfigType: 1, Priority: 0, Weight: 100},
	}
	lastWeight := func(n int64) int64 { return n - 1 }
	ordered := weightedFailoverOrder(primary, candidates, lastWeight)
	ids := make([]int64, len(ordered))
	for i, provider := range ordered {
		ids[i] = provider.ID
	}
	if !reflect.DeepEqual(ids, []int64{4, 3, 2}) {
		t.Fatalf("failover order = %v, want [4 3 2]", ids)
	}

	var calls []int64
	err := sendWithFailover(primary, candidates, sqlc.Message{ProviderID: primary.ID}, false, lastWeight, func(provider sqlc.Provider, message sqlc.Message) error {
		calls = append(calls, provider.ID)
		if message.ProviderID != provider.ID {
			t.Fatalf("message provider ID = %d, want %d", message.ProviderID, provider.ID)
		}
		if provider.ID == 4 {
			return nil
		}
		return errors.New("unavailable")
	})
	if err != nil {
		t.Fatalf("sendWithFailover() error = %v", err)
	}
	if !reflect.DeepEqual(calls, []int64{1, 4}) {
		t.Fatalf("send attempts = %v, want [1 4]", calls)
	}
}
