package notifiers

import (
	"avan/internal/db/sqlc"
	"fmt"
	"math/rand"
	"sort"

	"go.uber.org/zap"
)

func (n Notifier) SendWithFailover(primary sqlc.Provider, message sqlc.Message, retry bool) error {
	providers, err := n.ProviderService.GetProviders(primary.WorkspaceID, primary.OwnerID)
	if err != nil {
		n.Logger.Warn("failed to load provider failover candidates", zap.Error(err))
	}
	return sendWithFailover(primary, providers, message, retry, rand.Int63n, func(provider sqlc.Provider, message sqlc.Message) error {
		return n.SendMessage(provider, message)
	})
}

func sendWithFailover(primary sqlc.Provider, candidates []sqlc.Provider, message sqlc.Message, retry bool, draw func(int64) int64, send func(sqlc.Provider, sqlc.Message) error) error {
	// ponytail: a timeout can hide a successful delivery; duplicate suppression needs Phase 3 idempotency.
	if err := SendWithRetry(func() error {
		message.ProviderID = primary.ID
		return send(primary, message)
	}, retry); err == nil {
		return nil
	} else {
		lastErr := err
		ordered := weightedFailoverOrder(primary, candidates, draw)
		for _, provider := range ordered {
			err := SendWithRetry(func() error {
				message.ProviderID = provider.ID
				return send(provider, message)
			}, retry)
			if err == nil {
				return nil
			}
			lastErr = err
		}
		if len(ordered) == 0 {
			return fmt.Errorf("provider %q failed with no eligible failover provider: %w", primary.Name, err)
		}
		return fmt.Errorf("provider %q and %d failover provider(s) failed: %w", primary.Name, len(ordered), lastErr)
	}
}

func weightedFailoverOrder(primary sqlc.Provider, candidates []sqlc.Provider, draw func(int64) int64) []sqlc.Provider {
	groups := make(map[int64][]sqlc.Provider)
	for _, provider := range candidates {
		if provider.ID == primary.ID || provider.WorkspaceID != primary.WorkspaceID || provider.AppID != primary.AppID || provider.ConfigType != primary.ConfigType {
			continue
		}
		groups[provider.Priority] = append(groups[provider.Priority], provider)
	}

	priorities := make([]int64, 0, len(groups))
	for priority := range groups {
		priorities = append(priorities, priority)
	}
	sort.Slice(priorities, func(i, j int) bool { return priorities[i] < priorities[j] })

	var ordered []sqlc.Provider
	for _, priority := range priorities {
		group := groups[priority]
		for len(group) > 0 {
			totalWeight := int64(0)
			for _, provider := range group {
				weight := provider.Weight
				if weight < 1 {
					weight = 1
				}
				totalWeight += weight
			}
			choice := draw(totalWeight)
			selected := 0
			for index, provider := range group {
				weight := provider.Weight
				if weight < 1 {
					weight = 1
				}
				if choice < weight {
					selected = index
					break
				}
				choice -= weight
			}
			ordered = append(ordered, group[selected])
			group = append(group[:selected], group[selected+1:]...)
		}
	}
	return ordered
}
