package managers

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"git.iiens.net/edouardparis/town/app"
	"git.iiens.net/edouardparis/town/opennode"
	"git.iiens.net/edouardparis/town/resources"
)

var _ = uuid.Must(uuid.NewV4())

func CreateCharge(c *app.Config) (*resources.Charge, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "Error during uuid creation")
	}

	payload := &opennode.ChargePayload{
		Amount:   int64(6),
		Currency: "EUR",
		OrderID:  uid.String(),
	}

	payload.CallbackURL = fmt.Sprintf("%s/api/webhooks/checkout", c.InfoConfig.URLs.Website)

	charge, err := opennode.NewClient(&c.PaymentConfig).CreateCharge(payload)
	if err != nil {
		return nil, errors.Wrap(err, "Error during charge creation")
	}

	return resources.NewCharge(charge), nil
}