package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type CatchAllStatus string

const (
	UnknownCatchAllStatus CatchAllStatus = "unknown"
	IsCatchAllStatus      CatchAllStatus = "catch-all"
	IsNotCatchAllStatus   CatchAllStatus = "not catch-all"
)

type Domain struct {
	Id         int
	DomainName string
	Delivered  int
	Bounced    int
}

// Strut which contains additional fields desired in the JSON output.
type domainJson struct {
	Domain
	IsCatchAll CatchAllStatus
}

// Returns whether the domain is a catch all domain or not
func (self *Domain) IsCatchAll() CatchAllStatus {
	if self.Id == 0 {
		return UnknownCatchAllStatus
	} else if self.Bounced > 0 {
		return IsNotCatchAllStatus
	} else if self.Delivered > 1000 {
		return IsCatchAllStatus
	} else {
		return UnknownCatchAllStatus
	}
}

// Increments the bounced count for the domain
func (self *Domain) IncrementBounced(pool *pgxpool.Pool) error {
	return createDomainOrIncrement(pool, self.DomainName, "bounced")
}

// Increments the delivered count for the domain
func (self *Domain) IncrementDelivered(pool *pgxpool.Pool) error {
	return createDomainOrIncrement(pool, self.DomainName, "delivered")
}

// Gets the domain details for the specified domain
func GetDomain(pool *pgxpool.Pool, logger *zap.SugaredLogger, domainName string) (*Domain, error) {
	result, err := pool.Query(context.Background(), "SELECT id, delivered, bounced FROM domains WHERE domain_name = $1", domainName)
	domain := &Domain{
		DomainName: domainName,
	}

	fmt.Print(domainName)

	if err != nil {
		logger.Errorf("Failed to get domain with error: %v", err, zap.String("domain", domainName))
		return nil, err
	}

	defer result.Close()

	if result.Next() {
		err := result.Scan(&domain.Id, &domain.Delivered, &domain.Bounced)
		if err != nil {
			return nil, err
		}
		logger.Infow("RESULT", zap.Int("delivered", domain.Delivered), zap.Int("bounced", domain.Bounced), zap.String("domain", domain.DomainName), zap.Int("domainId", domain.Id))
		return domain, nil
	}
	return domain, nil
}

// Returns a struct to be used as a Json serialization
func (self *Domain) Json() *domainJson {
	return &domainJson{
		Domain:     *self,
		IsCatchAll: self.IsCatchAll(),
	}
}

// Inserts a new domain record or increments the specified field if a record already exist for the domain
func createDomainOrIncrement(pool *pgxpool.Pool, domain string, incrementField string) error {
	query := fmt.Sprintf("INSERT INTO domains (domain_name, %[1]s) VALUES ($1, 1) ON CONFLICT (domain_name) DO UPDATE SET %[1]s = domains.%[1]s + 1", incrementField)
	_, err := pool.Exec(context.Background(), query, domain)
	return err
}
