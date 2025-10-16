package main

import (
	"context"
	"slices"

	"github.com/hay-kot/hookfeed/backend/cmd/tools/seed/client"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/pkgs/utils"
	"github.com/rs/zerolog"
)

func seed(ctx context.Context, l zerolog.Logger, hookfeed client.Client, data Data) error {
	// Auto-register, then try to login. Ignore errors because we'll
	// catch them on login
	_, _ = hookfeed.Auth.Register(ctx, data.User)

	sess, err := hookfeed.Auth.Login(ctx, dtos.UserAuthenticate{
		Email:    data.User.Email,
		Password: data.User.Password,
	})
	if err != nil {
		return err
	}

	utils.Dump(sess)

	user, err := hookfeed.User.Self(ctx)
	if err != nil {
		return err
	}

	utils.Dump(user)

	l.Info().Str("email", data.User.Email).Msg("user authenticated")
	l.Info().Msg("processing packing lists")

	all, err := hookfeed.Lists.GetAll(ctx, dtos.PackingListQuery{
		Pagination: dtos.Pagination{
			Skip:  0,
			Limit: 300,
		},
		OrderBy: "",
	})
	if err != nil {
		return err
	}

	// Create any that don't already exist
	for _, list := range data.Lists {
		exists := slices.ContainsFunc(all.Items, func(v dtos.PackingList) bool {
			return v.Name == list.Name
		})
		if exists {
			l.Info().Str("name", list.Name).Msg("packing list already exists")
			continue
		}
		l.Info().Str("name", list.Name).Msg("creating packing list")
		created, err := hookfeed.Lists.Create(ctx, dtos.PackingListCreate{
			Name:        list.Name,
			Description: list.Description,
			DueDate:     &list.DueDate,
			Days:        list.Days,
			Tags:        list.Tags,
		})
		if err != nil {
			return err
		}
		l.Info().Str("id", created.ID.String()).Msg("packing list created")

		if list.Status != "not-started" {
			updated, err := hookfeed.Lists.Update(ctx, created.ID, dtos.PackingListUpdate{
				Status: &list.Status,
			})
			if err != nil {
				return err
			}

			created = updated
		}

		l.Info().Str("id", created.ID.String()).Msg("creating items")
		for i, item := range list.Items {
			createdItem, err := hookfeed.ListItems.Create(ctx, created.ID, dtos.PackingListItemCreate{
				Name:     item.Name,
				Category: item.Category,
				Quantity: item.Quantity,
				Notes:    item.Notes,
			})
			if err != nil {
				return err
			}

			// Mark odd as packed
			if i%2 == 1 {
				_, err := hookfeed.ListItems.Update(ctx, created.ID, createdItem.ID, dtos.PackingListItemUpdate{
					IsPacked: utils.Ptr(true),
				})
				if err != nil {
					return err
				}
			}
		}
	}

	l.Info().Msg("packing lists processed")
	return nil
}
