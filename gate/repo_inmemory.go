package main

var _ ConfigRepo = (*InMemoryConfigRepo)(nil)

type InMemoryConfigRepo struct {
	data map[string]map[string][]byte
}

// NewInMemoryConfigRepo allocates PostgreSQL config reposistory.
func NewInMemoryConfigRepo() *InMemoryConfigRepo {
	return &InMemoryConfigRepo{
		data: make(map[string]map[string][]byte),
	}
}

func (r *InMemoryConfigRepo) SetParams(config, param string, payload []byte) {
	if _, found := r.data[config]; !found {
		r.data[config] = make(map[string][]byte)
	}
	r.data[config][param] = payload
}

// GetParams loads config params with given keys from DB.
func (r *InMemoryConfigRepo) GetParams(config, param string) ([]byte, error) {
	var data []byte

	params, found := r.data[config]
	if !found {
		return nil, ErrNotFound
	}
	data, found = params[param]
	if !found {
		return nil, ErrNotFound
	}

	return data, nil
}
