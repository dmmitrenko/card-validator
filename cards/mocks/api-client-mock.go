package mocks

type MockApiClient struct {
	MockCheckINN func(iin string) error
}

func NewMockApiClient(mockCheckINN func(iin string) error) *MockApiClient {
	return &MockApiClient{
		MockCheckINN: mockCheckINN,
	}
}

func (m *MockApiClient) CheckINN(iin string) error {
	return m.MockCheckINN(iin)
}
