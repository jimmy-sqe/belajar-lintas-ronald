// Package adapter contains outbound adapters that implement domain ports.
// Each subtree is one axis (persistence, cache, auth, ...) with one folder
// per option (postgres, mysql, ...). Adapters MUST NOT import from sibling
// adapters; they may import from internal/domain and pkg/.
package adapter
