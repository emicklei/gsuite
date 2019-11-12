package main

import (
	"errors"
	"fmt"

	"github.com/urfave/cli"
	admin "google.golang.org/api/admin/directory/v1"
)

func cmdDomainList(c *cli.Context) error {

	client := sharedAuthClient()

	srv, err := admin.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	r, err := srv.Domains.List("my_customer").Do() // ??
	if err != nil {
		return fmt.Errorf("unable to retrieve domains: %v", err)
	}
	if optionJSON(c, r.Domains) {
		return nil
	}
	for _, each := range r.Domains {
		fmt.Println(each.DomainName)
	}
	return nil
}

func primaryDomain() (string, error) {
	client := sharedAuthClient()

	srv, err := admin.New(client)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	r, err := srv.Domains.List("my_customer").Do() // ??
	if err != nil {
		return "", fmt.Errorf("unable to retrieve domains: %v", err)
	}

	for _, each := range r.Domains {
		if each.IsPrimary {
			return each.DomainName, nil
		}
	}
	return "", errors.New("no primary domain found")
}
