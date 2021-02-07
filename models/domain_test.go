package models

import (
	"testing"

	"github.com/ersmith/mailgun-coding-challenge/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestIncrementDomainDelivered(t *testing.T) {
	dbPool := test.CreateTestPgxPool(t)
	domainName := test.RandomDomainName(20)
	domain := Domain{
		DomainName: domainName,
	}
	domain.IncrementDelivered(dbPool)
	fetchedDomain, err := GetDomain(dbPool, zap.NewNop().Sugar(), domainName)

	assert.Nil(t, err)
	assert.Equal(t, 1, fetchedDomain.Delivered)
}

func TestIncrementDomainBounced(t *testing.T) {
	dbPool := test.CreateTestPgxPool(t)
	domainName := test.RandomDomainName(20)
	domain := Domain{
		DomainName: domainName,
	}
	domain.IncrementBounced(dbPool)
	fetchedDomain, err := GetDomain(dbPool, zap.NewNop().Sugar(), domainName)

	assert.Nil(t, err)
	assert.Equal(t, 1, fetchedDomain.Bounced)
}

func TestBouncedResultsInNoCatchall(t *testing.T) {
	dbPool := test.CreateTestPgxPool(t)
	domainName := test.RandomDomainName(20)
	domain := Domain{
		DomainName: domainName,
	}
	domain.IncrementBounced(dbPool)
	fetchedDomain, err := GetDomain(dbPool, zap.NewNop().Sugar(), domainName)

	assert.Nil(t, err)
	assert.Equal(t, IsNotCatchAllStatus, fetchedDomain.IsCatchAll())
}

func TestDomainResultsInUnknownCatchallNoData(t *testing.T) {
	dbPool := test.CreateTestPgxPool(t)
	domainName := test.RandomDomainName(20)

	fetchedDomain, err := GetDomain(dbPool, zap.NewNop().Sugar(), domainName)

	assert.Nil(t, err)
	assert.Equal(t, UnknownCatchAllStatus, fetchedDomain.IsCatchAll())
}

func TestDomainResultsInUnknownCatchallWithDelivered(t *testing.T) {
	dbPool := test.CreateTestPgxPool(t)
	domainName := test.RandomDomainName(20)
	domain := Domain{
		DomainName: domainName,
	}

	for i := 0; i < 900; i++ {
		domain.IncrementDelivered(dbPool)
	}

	fetchedDomain, err := GetDomain(dbPool, zap.NewNop().Sugar(), domainName)

	assert.Nil(t, err)
	assert.Equal(t, UnknownCatchAllStatus, fetchedDomain.IsCatchAll())
}

func TestDomainIsCatchAllTrueWithManyDelivered(t *testing.T) {
	dbPool := test.CreateTestPgxPool(t)
	domainName := test.RandomDomainName(20)
	domain := Domain{
		DomainName: domainName,
	}

	for i := 0; i < 1001; i++ {
		domain.IncrementDelivered(dbPool)
	}

	fetchedDomain, err := GetDomain(dbPool, zap.NewNop().Sugar(), domainName)

	assert.Nil(t, err)
	assert.Equal(t, IsCatchAllStatus, fetchedDomain.IsCatchAll())
}

func TestDomainIsCatchAllFalseWithManyDeliveredAndBounced(t *testing.T) {
	dbPool := test.CreateTestPgxPool(t)
	domainName := test.RandomDomainName(20)
	domain := Domain{
		DomainName: domainName,
	}

	for i := 0; i < 1001; i++ {
		domain.IncrementDelivered(dbPool)
	}
	domain.IncrementBounced(dbPool)
	fetchedDomain, err := GetDomain(dbPool, zap.NewNop().Sugar(), domainName)

	assert.Nil(t, err)
	assert.Equal(t, IsNotCatchAllStatus, fetchedDomain.IsCatchAll())
}
