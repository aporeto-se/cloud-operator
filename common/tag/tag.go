package tag

import (
	"fmt"
	"sync"

	"github.com/aporeto-se/cloud-operator/common/types"
)

// Matcher tag matcher. Matcher is NOT thread safe.
type Matcher struct {
	filter1 *types.Filter
	filter2 *types.Filter
	sync.Mutex
}

// NewMatcher returns new Matcher
func NewMatcher(filter1, filter2 *types.Filter) (*Matcher, error) {

	if filter1 == nil {
		return nil, fmt.Errorf("filter2 may be nil but filter1 may NOT")
	}

	return &Matcher{
		filter1: filter1,
		filter2: filter2,
	}, nil
}

func (t *Matcher) filter1MatchKubeCluster(clusterName string, clusterTags map[string]string) bool {

	if t.filter1 == nil {
		return false
	}

	if t.filter1.KubeMatchAny {
		return true
	}

	if t.filter1.HasKubeMatchName(clusterName) {
		return true
	}

	if clusterTags == nil {
		return false
	}

	for key, value := range clusterTags {
		if t.filter1.HasKubeMatchTag(key, value) {
			return true
		}
	}

	return false
}

func (t *Matcher) filter2MatchKubeCluster(clusterName string, clusterTags map[string]string) bool {

	if t.filter2 == nil {
		return true
	}

	if t.filter2.KubeMatchAny {
		return true
	}

	if t.filter2.HasKubeMatchName(clusterName) {
		return true
	}

	if clusterTags == nil {
		return false
	}

	for key, value := range clusterTags {
		if t.filter2.HasKubeMatchTag(key, value) {
			return true
		}
	}

	return false
}

// MatchKubeCluster returns true if match
func (t *Matcher) MatchKubeCluster(clusterName string, clusterTags map[string]string) bool {

	t.Lock()
	defer t.Unlock()

	if !t.filter1MatchKubeCluster(clusterName, clusterTags) {
		return false
	}

	if t.filter2MatchKubeCluster(clusterName, clusterTags) {
		return true
	}

	return false
}
