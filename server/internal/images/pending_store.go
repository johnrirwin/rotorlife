package images

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// InMemoryPendingStore keeps approved uploads in memory until user confirms Save.
type InMemoryPendingStore struct {
	mu      sync.RWMutex
	uploads map[string]PendingUpload
	ttl     time.Duration
}

// NewInMemoryPendingStore creates a pending store with the provided TTL.
func NewInMemoryPendingStore(ttl time.Duration) *InMemoryPendingStore {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}

	return &InMemoryPendingStore{
		uploads: make(map[string]PendingUpload),
		ttl:     ttl,
	}
}

// Put stores an approved upload and returns its token id.
func (s *InMemoryPendingStore) Put(upload PendingUpload) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.cleanupLocked(now)

	id := uuid.NewString()
	upload.ID = id
	upload.ExpiresAt = now.Add(s.ttl)
	s.uploads[id] = upload

	return id
}

// Get fetches an upload token scoped to a user.
func (s *InMemoryPendingStore) Get(ownerUserID, uploadID string) (*PendingUpload, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.cleanupLocked(now)

	upload, ok := s.uploads[uploadID]
	if !ok {
		return nil, false
	}
	if upload.OwnerUserID != ownerUserID {
		return nil, false
	}

	copyUpload := upload
	return &copyUpload, true
}

// Delete removes an upload token.
func (s *InMemoryPendingStore) Delete(uploadID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.uploads, uploadID)
}

func (s *InMemoryPendingStore) cleanupLocked(now time.Time) {
	for id, upload := range s.uploads {
		if now.After(upload.ExpiresAt) {
			delete(s.uploads, id)
		}
	}
}
