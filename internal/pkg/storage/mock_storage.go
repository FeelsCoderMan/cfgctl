package storage

type MockFileStore struct {
	*FileStore
	loaded bool
}

func NewMockFileStore() *MockFileStore {
	fs := NewFileStore()
	return &MockFileStore{
		FileStore: fs,
	}
}

func (m *MockFileStore) GetPath() string {
	m.mu.RLock()
	path := m.FileStore.path
	m.mu.RUnlock()
	return path
}

func (m *MockFileStore) IsLoaded() bool {
	m.mu.RLock()
	loaded := m.loaded
	m.mu.RUnlock()
	return loaded
}

func (m *MockFileStore) IsDataValid(data any) bool {
	m.mu.RLock()
	storage := m.FileStore.storage
	m.mu.RUnlock()

	if len(data.(map[string]any)) != len(storage) {
		return false
	}

	for key, value := range data.(map[string]any) {
		if storage[key] != value {
			return false
		}
	}

	return true
}

func (m *MockFileStore) LoadData(data map[string]any) error {
	m.mu.Lock()
	if err := m.FileStore.LoadData(data); err != nil {
		m.mu.Unlock()
		return err
	}

	m.loaded = true
	m.mu.Unlock()
	return nil
}
