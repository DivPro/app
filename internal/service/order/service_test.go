package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DivPro/app/internal/service/order/mocks"
	"github.com/DivPro/app/pkg/entity"
)

func TestService_Create(t *testing.T) {
	txMock := mocks.NewStorage(t)
	svc := NewService(txMock)

	txMock.On("Tx",
		context.Background(),
		mock.AnythingOfType("storage.TxFn"),
	).Return(nil)
	var order entity.Order
	orderCreated, err := svc.Create(context.Background(), order)
	assert.NoError(t, err)
	assert.Equal(t, order, *orderCreated)
}
