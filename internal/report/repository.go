package report

import "sync"

type Repository struct {
	mutex *sync.Mutex

	last           *Report
	lastComplete   *Report
	lastSuccessful *Report
}

func NewRepository() Repository {
	return Repository{mutex: &sync.Mutex{}}
}

func (r *Repository) StartNewReport() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.last = NewReport()
}

func (r *Repository) Watch(messageChan <-chan Message) {
	r.last.Watch(messageChan)
}

func (r *Repository) GetLastJSON() ([]byte, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.last.ToJSON()
}

func (r *Repository) GetLastCompleteJSON() ([]byte, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.lastComplete.ToJSON()
}

func (r *Repository) GetLastSuccessfulJSON() ([]byte, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.lastSuccessful.ToJSON()
}

func (r *Repository) HasValidBuild() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.lastSuccessful != nil
}

func (r *Repository) UpdateStatus(status jobStatus) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.last.Status = status
}

func (r *Repository) UpdateStats(stats Stats) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.last.Stats = stats
}

func (r *Repository) MarkAsComplete() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.lastComplete = r.last
}

func (r *Repository) MarkAsSuccessful() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.lastSuccessful = r.last
}
