package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"SpamhausResponseApi/db"
	"SpamhausResponseApi/graph/generated"
	"SpamhausResponseApi/helpers"
	"SpamhausResponseApi/model"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (r *mutationResolver) Enqueue(ctx context.Context, ip []*model.IP) ([]*model.IPDetails, error) {
	var success []*model.IPDetails
	for _, address := range ip {
		if !helpers.IpReg.MatchString(address.Address) {
			return nil, fmt.Errorf("invalid entry. %s", address.Address)
		}

		reverseIP, err := helpers.ReverseUsingSeperator(address.Address, ".")
		if err != nil {
			log.Println(err)
			return nil, err
		}
		fqdn := reverseIP + helpers.Zen

		out, err := exec.Command("dig", "+short", fqdn).Output()
		if err != nil {
			return nil, err
		}

		details := model.IPDetails{}
		db.Conn.Where("ip_address = ?", address.Address).Find(&details)

		codes := ""
		for _, code := range strings.Split(string(out), "\n") {
			codes += ", " + code
		}
		codes = strings.Trim(codes, ", ")

		if details.UUID == "" {
			details.UUID = uuid.NewString()
			details.IPAddress = address.Address
			details.ResponseCode = strings.Trim(codes, ", ")
			db.Conn.Save(&details)
		} else {
			if details.ResponseCode != strings.Trim(codes, ", ") {
				details.ResponseCode = strings.Trim(codes, ", ")
			}
			details.UpdatedAt = time.Now()
			db.Conn.Model(details).Where("ip_address = ?", details.IPAddress).Update(&details)
		}

		success = append(success, &details)
	}
	return success, nil
}

func (r *queryResolver) GetIPDetails(ctx context.Context, ip model.IP) (*model.IPDetails, error) {
	details := &model.IPDetails{}

	db.Conn.Where("ip_address = ?", ip.Address).Find(&details)

	if details.UUID == "" {
		return nil, fmt.Errorf("no data has been returned for ip_address %s", ip.Address)
	}
	return details, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
