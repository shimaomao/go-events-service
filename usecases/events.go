// Package usecases implements application-specific business logic
// for the events service.
package usecases

import (
	"github.com/declantraynor/go-events-service/domain"
)

type EventInteractor struct {
	Store domain.EventStore
}

func (interactor *EventInteractor) AddEvent(name string, timestamp string) error {

	parsedTimestamp, err := ParseTimestamp(timestamp)
	if err != nil {
		return err
	}

	event := domain.Event{Name: name, Timestamp: parsedTimestamp.Unix()}
	if err := interactor.Store.Put(event); err != nil {
		return err
	}

	return nil
}

func (interactor *EventInteractor) CountEventsInTimeRange(from, to string) (map[string]int, error) {
	parsedFrom, fromerr := ParseTimestamp(from)
	if fromerr != nil {
		return map[string]int{}, fromerr
	}

	parsedTo, toerr := ParseTimestamp(to)
	if toerr != nil {
		return map[string]int{}, toerr
	}

	if !parsedFrom.Before(parsedTo) {
		return map[string]int{}, InvalidTimeRangeError{from: from, to: to}
	}

	eventNames, err := interactor.Store.Names()
	if err != nil {
		return map[string]int{}, err
	}

	counts := map[string]int{}
	for _, name := range eventNames {
		count, err := interactor.Store.CountInTimeRange(name, parsedFrom.Unix(), parsedTo.Unix())
		if err != nil {
			return map[string]int{}, err
		}
		counts[name] = count
	}

	return counts, nil
}
