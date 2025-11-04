package controller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/Beigelman/nossas-despesas/internal/pkg/pubsub"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
)

func TestRecalculateExpensesSplitRatioHandler(t *testing.T) {
	t.Parallel()

	// Definição dos casos de teste
	testCases := []struct {
		name          string
		mockSetup     func(subscriber *mocks.MockpubsubSubscriber, recalculateExpenses *mocks.MockusecaseRecalculateExpensesSplitRatio)
		expectedError bool
	}{
		{
			name: "should return error if subscriber fails",
			mockSetup: func(subscriber *mocks.MockpubsubSubscriber, recalculateExpenses *mocks.MockusecaseRecalculateExpensesSplitRatio) {
				subscriber.EXPECT().Subscribe(mock.Anything, pubsub.IncomesTopic).Return(
					nil, errors.New("subscription error"),
				).Once()
			},
			expectedError: true,
		},
		{
			name: "should start listening successfully when subscriber succeeds",
			mockSetup: func(subscriber *mocks.MockpubsubSubscriber, recalculateExpenses *mocks.MockusecaseRecalculateExpensesSplitRatio) {
				// Mock um canal vazio para simular sucesso na subscription
				msgChan := make(chan *message.Message)
				close(msgChan)

				subscriber.EXPECT().Subscribe(mock.Anything, pubsub.IncomesTopic).Return(
					msgChan, nil,
				).Once()
			},
			expectedError: false,
		},
	}

	// Execução dos casos de teste
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup dos mocks
			mockSubscriber := mocks.NewMockpubsubSubscriber(t)
			mockRecalculateExpenses := mocks.NewMockusecaseRecalculateExpensesSplitRatio(t)

			tc.mockSetup(mockSubscriber, mockRecalculateExpenses)

			// Create the handler
			handler := controller.NewRecalculateExpensesSplitRatio(mockSubscriber, mockRecalculateExpenses.Execute)

			// Execute the handler
			ctx := context.Background()
			err := handler(ctx)

			// Assertions
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
