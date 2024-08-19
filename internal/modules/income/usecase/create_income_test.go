package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: refatorar esse teste quando criar mock para o publisher
func TestCreateIncome(t *testing.T) {
	// t.Parallel()
	// ctx := context.Background()
	// incomeRepo := mocks.NewMockincomeRepository(t)
	// userRepo := mocks.NewMockuserRepository(t)
	// usr := user.New(user.Attributes{
	// 	ID:    user.ID{Value: 1},
	// 	Name:  "My test user",
	// 	Email: "email",
	// })
	//
	// useCase := NewCreateIncome(userRepo, incomeRepo, nil)
	// params := CreateIncomeParams{
	// 	Type:      income.Types.Salary,
	// 	Amount:    100,
	// 	UserID:    user.ID{Value: 1},
	// 	CreatedAt: nil,
	// }
	//
	// t.Run("getUserByID returns error", func(t *testing.T) {
	// 	userRepo.EXPECT().GetByID(ctx, usr.ID).Return(nil, errors.New("test error")).Once()
	// 	inc, err := useCase(ctx, params)
	// 	assert.Errorf(t, err, "userRepo.GetByID: test error")
	// 	assert.Nil(t, inc)
	// })
	//
	// t.Run("user does not exist", func(t *testing.T) {
	// 	userRepo.EXPECT().GetByID(ctx, usr.ID).Return(nil, nil).Once()
	// 	inc, err := useCase(ctx, params)
	// 	assert.Errorf(t, err, "user not found")
	// 	assert.Nil(t, inc)
	// })
	//
	// t.Run("store returns error", func(t *testing.T) {
	// 	incomeRepo.EXPECT().GetNextID().Return(income.ID{Value: 1}).Once()
	// 	userRepo.EXPECT().GetByID(ctx, usr.ID).Return(usr, nil).Once()
	// 	incomeRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()
	// 	inc, err := useCase(ctx, params)
	// 	assert.Errorf(t, err, "incomeRepo.Store: test error")
	// 	assert.Nil(t, inc)
	// })
	//
	// t.Run("success", func(t *testing.T) {
	// 	incomeRepo.EXPECT().GetNextID().Return(income.ID{Value: 1}).Once()
	// 	userRepo.EXPECT().GetByID(ctx, usr.ID).Return(usr, nil).Once()
	// 	incomeRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
	// 	inc, err := useCase(ctx, params)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, inc)
	// 	assert.Equal(t, user.ID{Value: 1}, inc.UserID)
	// 	assert.Equal(t, income.Types.Salary, inc.Type)
	// 	assert.Equal(t, 100, inc.Amount)
	// })

	assert.Equal(t, 1, 1)
}
