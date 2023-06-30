package repository

import (
	"errors"
	"sync"

	"github.com/criteo/data-aggregation-api/internal/ingestor/cmdb"
	"github.com/criteo/data-aggregation-api/internal/ingestor/dcim"
	"github.com/criteo/data-aggregation-api/internal/report"
)

const ingestorNumber = 7

// FetchAssets get data from all ingestors.
func FetchAssets(reportCh chan report.Message) (*Assets, error) {
	wg := sync.WaitGroup{}
	var repo Assets

	var fetchFailure = make(chan report.Severity, ingestorNumber)

	// TODO: severity should be defined by the user via the configuration file
	// TODO: factorize

	// Devices
	wg.Add(1)
	go func() {
		defer wg.Done()
		if v, err := dcim.GetNetworkInventory(); err != nil {
			reportCh <- report.Message{
				Type:     report.IngestorMessage,
				Severity: report.Error,
				Text:     err.Error(),
			}
			fetchFailure <- report.Error
		} else {
			repo.DeviceInventory = v
		}
	}()

	// BGP Global configuration
	wg.Add(1)
	go func() {
		defer wg.Done()
		if v, err := cmdb.GetBGPGlobal(); err != nil {
			reportCh <- report.Message{
				Type:     report.IngestorMessage,
				Severity: report.Warning,
				Text:     err.Error(),
			}
			fetchFailure <- report.Warning
		} else {
			repo.CmdbBGPGlobal = v
		}
	}()

	// BGP sessions
	wg.Add(1)
	go func() {
		defer wg.Done()
		if v, err := cmdb.GetBGPSessions(); err != nil {
			reportCh <- report.Message{
				Type:     report.IngestorMessage,
				Severity: report.Error,
				Text:     err.Error(),
			}
			fetchFailure <- report.Error
		} else {
			repo.CmdbBGPSessions = v
		}
	}()

	// Peer groups
	wg.Add(1)
	go func() {
		defer wg.Done()
		if v, err := cmdb.GetPeerGroups(); err != nil { //nolint:staticcheck // to ignore deprecation notice
			reportCh <- report.Message{
				Type:     report.IngestorMessage,
				Severity: report.Warning,
				Text:     err.Error(),
			}
			fetchFailure <- report.Warning
		} else {
			repo.CmdbPeerGroups = v
		}
	}()

	// Route policies
	wg.Add(1)
	go func() {
		defer wg.Done()
		if v, err := cmdb.GetRoutePolicies(); err != nil {
			reportCh <- report.Message{
				Type:     report.IngestorMessage,
				Severity: report.Error,
				Text:     err.Error(),
			}
			fetchFailure <- report.Error
		} else {
			repo.CmdbRoutePolicies = v
		}
	}()

	// Prefix lists
	wg.Add(1)
	go func() {
		defer wg.Done()
		if v, err := cmdb.GetPrefixLists(); err != nil {
			reportCh <- report.Message{
				Type:     report.IngestorMessage,
				Severity: report.Error,
				Text:     err.Error(),
			}
			fetchFailure <- report.Error
		} else {
			repo.CmdbPrefixLists = v
		}
	}()

	// Community lists
	wg.Add(1)
	go func() {
		defer wg.Done()
		if v, err := cmdb.GetCommunityLists(); err != nil {
			reportCh <- report.Message{
				Type:     report.IngestorMessage,
				Severity: report.Error,
				Text:     err.Error(),
			}
			fetchFailure <- report.Error
		} else {
			repo.CmdbCommunityLists = v
		}
	}()

	// Wait for responses
	go func() {
		wg.Wait()
		close(fetchFailure)
	}()

	failed := false
	for e := range fetchFailure {
		if e == report.Error {
			failed = true
		}
	}
	if failed {
		return &repo, errors.New("ingestors failed to fetch the assets")
	}

	return &repo, nil
}
