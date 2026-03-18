package cmd

import "context"

type mockBuilder struct {
	err error
}

func (m *mockBuilder) Build(ctx context.Context) error {
	return m.err
}