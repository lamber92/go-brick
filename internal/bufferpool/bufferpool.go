package bufferpool

import "go.uber.org/zap/buffer"

// Referenced from: go.uber.org\zap@v1.24.0

var (
	_pool = buffer.NewPool()
	// Get retrieves a buffer from the pool, creating one if necessary.
	Get = _pool.Get
)
