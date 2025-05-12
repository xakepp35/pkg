package xpgx

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestTxManager_Do_Success(t *testing.T) {
	ctx := context.Background()

	mockTx := NewMockTx(t)
	mockTx.EXPECT().Rollback(ctx).Return(nil)
	mockTx.EXPECT().Commit(ctx).Return(nil)

	mockPool := NewMockTransactional(t)
	mockPool.EXPECT().
		BeginTx(ctx, pgx.TxOptions{}).
		Return(mockTx, nil)

	manager := txManager{mockPool}

	err := manager.Do(ctx, func(tx pgx.Tx) error {
		assert.Equal(t, mockTx, tx)
		return nil
	})

	assert.NoError(t, err)
	mockTx.AssertExpectations(t)
	mockPool.AssertExpectations(t)
}

func TestTxManager_Do_ExecError(t *testing.T) {
	ctx := context.Background()
	execErr := errors.New("exec failed")

	mockTx := NewMockTx(t)
	mockTx.EXPECT().Rollback(ctx).Return(nil)
	mockTx.EXPECT().Commit(ctx).Maybe() // не будет вызван, но допустим

	mockPool := NewMockTransactional(t)
	mockPool.EXPECT().
		BeginTx(ctx, pgx.TxOptions{}).
		Return(mockTx, nil)

	manager := txManager{mockPool}

	err := manager.Do(ctx, func(tx pgx.Tx) error {
		return execErr
	})

	assert.ErrorIs(t, err, execErr)
	mockTx.AssertExpectations(t)
	mockPool.AssertExpectations(t)
}

func TestTxManager_Do_BeginTxError(t *testing.T) {
	ctx := context.Background()
	beginErr := errors.New("begin error")

	mockPool := NewMockTransactional(t)
	mockPool.EXPECT().
		BeginTx(ctx, pgx.TxOptions{}).
		Return(nil, beginErr)

	manager := txManager{mockPool}

	err := manager.Do(ctx, func(tx pgx.Tx) error {
		return nil
	})

	assert.ErrorIs(t, err, beginErr)
	mockPool.AssertExpectations(t)
}
