package store

import (
	"context"

	"github.com/pkg/errors"
	lk "github.com/ulule/loukoum"
	"github.com/ulule/makroud"

	"git.iiens.net/edouardparis/town/failures"
	"git.iiens.net/edouardparis/town/models"
)

type Addresses struct {
	Store
}

func NewAddresses(s Store) *Addresses {
	return &Addresses{s}
}

func (a Addresses) GetByValue(ctx context.Context, v string) (*models.Address, error) {
	address := &models.Address{}

	q := lk.Select(columns(address)).
		From(address.TableName()).
		Where(lk.Condition("value").Equal(v))

	err := a.Get(ctx, q, address)
	if err != nil {
		if !makroud.IsErrNoRows(err) {
			return nil, errors.Wrapf(err, "cannot retrieve address with value: %s", v)
		}
		return nil, failures.ErrNotFound
	}

	return address, nil
}

func (a *Addresses) Create(ctx context.Context, address *models.Address) error {
	query := lk.Insert(address.TableName()).
		Set(
			lk.Pair("value", address.Value),
			lk.Pair("amount_collected", address.AmountCollected),
			lk.Pair("amount_received", address.AmountReceived),
		).Returning("id", "created_at", "updated_at")
	err := a.Save(ctx, query, address)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}