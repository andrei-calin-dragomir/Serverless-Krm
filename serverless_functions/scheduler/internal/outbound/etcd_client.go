package outbound

import (
	"context"
	"fmt"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ETCDClient struct {
	client *clientv3.Client
	config clientv3.Config
	mu     sync.RWMutex // Protect client access
	// retryTimeout time.Duration
}

func NewETCDClient(endpoints []string, timeout time.Duration /*retryTimeout time.Duration*/) (*ETCDClient, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &ETCDClient{
		client: client,
		config: cfg,
		// retryTimeout: retryTimeout,
	}, nil
}

func (e *ETCDClient) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}

// Check if the client is connected
func (e *ETCDClient) IsConnected() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Use an etcd endpoint to verify connectivity
	_, err := e.client.Status(ctx, e.config.Endpoints[0])
	return err == nil
}

// Retry connection to etcd cluster
func (e *ETCDClient) RetryConnection() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.client != nil {
		e.client.Close()
	}

	client, err := clientv3.New(e.config)
	if err != nil {
		return err
	}
	e.client = client
	return nil
}

// // Provide a reference to the client (internal use)
// func (e *ETCDClient) getClient() *clientv3.Client {
// 	e.mu.RLock()
// 	defer e.mu.RUnlock()
// 	return e.client
// }

// SaveKey saves a key-value pair to etcd.
func (ec *ETCDClient) SaveKey(ctx context.Context, key, value string) error {
	_, err := ec.client.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to save key %s: %w", key, err)
	}
	return nil
}

// DeleteKey deletes a key from etcd.
func (ec *ETCDClient) DeleteKey(ctx context.Context, key string) error {
	_, err := ec.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Extend and implement caching on get calls
// GetKeysWithPrefix retrieves all keys with a given prefix from etcd.
func (ec *ETCDClient) GetKeysWithPrefix(ctx context.Context, prefix string) (map[string]string, error) {
	resp, err := ec.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get keys with prefix %s: %w", prefix, err)
	}

	result := make(map[string]string)
	for _, kv := range resp.Kvs {
		result[string(kv.Key)] = string(kv.Value)
	}
	return result, nil
}
