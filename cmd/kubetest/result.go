package kubetest

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

//Result - represent Result struct
type Result struct {
	ID        string
	Kind      ClientType
	Messages  int
	Errors    []error
	Latencies []time.Duration
}

//NewResult
func NewResult(id string, n int) *Result {
	return &Result{
		ID:        id,
		Kind:      0,
		Messages:  n,
		Errors:    nil,
		Latencies: nil,
	}
}

//SetKind
func (r *Result) SetKind(kind ClientType) *Result {
	r.Kind = kind
	return r
}

//AddError
func (r *Result) AddError(err error) *Result {
	r.Errors = append(r.Errors, err)
	return r
}

//AddLatency
func (r *Result) AddLatency(t time.Duration) *Result {
	r.Latencies = append(r.Latencies, t)
	return r
}

//Latency
func (r *Result) Latency() time.Duration {
	var sum int64 = 0
	for i := 0; i < len(r.Latencies); i++ {
		sum += r.Latencies[i].Nanoseconds()
	}
	if len(r.Latencies) > 0 {
		avg := float64(sum) / float64(len(r.Latencies))
		return time.Duration(int64(avg)) * time.Nanosecond
	}
	return 0
}

//String
func (r *Result) String() string {
	return fmt.Sprintf("Results for %s: Messages: %d Errors: %d Latency: %s", ClientTypeNames[r.Kind], len(r.Errors), r.Messages, r.Latency().String())
}

//Verbose
func (r *Result) Verbose() string {

	list := []string{
		fmt.Sprintf("Results for %s: Messages: %d Errors: %d Latency: %s", ClientTypeNames[r.Kind], r.Messages, len(r.Errors), r.Latency().String()),
	}
	if r.Errors != nil {
		list = append(list, "Errors:")
		for i := 0; i < len(r.Errors); i++ {
			list = append(list, r.Errors[i].Error())
		}
	}
	if r.Latencies != nil {
		list = append(list, "Latencies:")
		for i := 0; i < len(r.Latencies); i++ {
			list = append(list, r.Latencies[i].String())
		}
	}
	return strings.Join(list, "\n\t")
}

//HasError
func (r *Result) HasError() bool {
	return len(r.Errors) > 0
}

//ResultsSet
type ResultsSet struct {
	sync.Mutex
	list []*Result
}

//NewResultsSet
func NewResultsSet() *ResultsSet {
	return &ResultsSet{
		Mutex: sync.Mutex{},
		list:  nil,
	}
}

//Add
func (rs *ResultsSet) Add(r *Result) {
	rs.Lock()
	defer rs.Unlock()
	rs.list = append(rs.list, r)
}

//Count
func (rs *ResultsSet) Count() int {
	rs.Lock()
	defer rs.Unlock()
	c := len(rs.list)
	return c
}

//Success
func (rs *ResultsSet) Success() int {
	rs.Lock()
	defer rs.Unlock()
	c := 0
	for _, result := range rs.list {
		if len(result.Errors) == 0 {
			c++
		}
	}
	return c
}

//Errors
func (rs *ResultsSet) Errors() int {
	rs.Lock()
	defer rs.Unlock()
	c := 0
	for _, result := range rs.list {
		c += len(result.Errors)
	}
	return c
}

//Latency
func (rs *ResultsSet) Latency() time.Duration {
	rs.Lock()
	defer rs.Unlock()
	var total time.Duration
	n := 0
	for _, result := range rs.list {
		if result.Latency() > 0 {
			total += result.Latency()
			n++
		}

	}
	if n == 0 {
		return 0
	}
	avg := float64(total) / float64(n)
	return time.Duration(avg) * time.Nanosecond
}

//HasErrors
func (rs *ResultsSet) HasErrors() bool {
	return rs.Errors() > 0
}

//Results
type Results struct {
	Producers *ResultsSet
	Consumers *ResultsSet
}

//NewResults
func NewResults(producers, consumers *ResultsSet) *Results {
	return &Results{
		Producers: producers,
		Consumers: consumers,
	}
}

//String
func (r *Results) String() string {
	var lines []string
	lines = append(lines,
		fmt.Sprintf("\tProducers Results - Success: %d, Errors: %d, Average Latency: %s", r.Producers.Success(), r.Producers.Errors(), r.Producers.Latency().String()),
		fmt.Sprintf("Consumers Results - Success: %d, Errors: %d, Average Latency: %s", r.Consumers.Success(), r.Consumers.Errors(), r.Consumers.Latency().String()))

	return strings.Join(lines, "\n\t")
}

//HasErrors
func (r *Results) HasErrors() bool {
	return r.Producers.HasErrors() || r.Consumers.HasErrors()
}

//Errors
func (r *Results) Errors() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("\tProducers Results - Success: %d, Errors: %d, Average Latency: %s", r.Producers.Success(), r.Producers.Errors(), r.Producers.Latency().String()))
	if r.Producers.HasErrors() {
		for _, res := range r.Producers.list {
			if res.HasError() {
				lines = append(lines, res.Verbose())
			}
		}
	}
	lines = append(lines, fmt.Sprintf("Consumers Results - Success: %d, Errors: %d, Average Latency: %s", r.Consumers.Success(), r.Consumers.Errors(), r.Consumers.Latency().String()))
	if r.Consumers.HasErrors() {
		for _, res := range r.Consumers.list {
			if res.HasError() {
				lines = append(lines, res.Verbose())
			}
		}
	}
	return strings.Join(lines, "\n\t")
}
