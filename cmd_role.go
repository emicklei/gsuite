package main

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli"
	"golang.org/x/net/context"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func cmdRoleList(c *cli.Context) error {
	client := sharedAuthClient(c)

	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	// TODO my_customer?
	roles, err := srv.Roles.List(myAccoutsCustomerID).MaxResults(int64(ifZero(c.Int("limit"), 100))).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve roles in domain: %v", err)
	}

	if optionJSON(c, roles.Items) {
		return nil
	}
	for _, u := range roles.Items {
		fmt.Println(u.RoleName)
	}
	return nil
}

func cmdRoleAssignment(c *cli.Context) error {
	client := sharedAuthClient(c)

	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	// Get all roles
	// TODO my_customer?
	roles, err := srv.Roles.List(myAccoutsCustomerID).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve roles in domain: %v", err)
	}

	// find role by name
	roleName := c.Args().Get(0)
	var roleID int64
	for _, each := range roles.Items {
		if each.RoleName == roleName {
			roleID = each.RoleId
			break
		}
	}

	// find all assigments per role
	ass, err := srv.RoleAssignments.List(myAccoutsCustomerID).RoleId(strconv.FormatInt(roleID, 10)).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve roles in domain: %v", err)
	}

	users := []*admin.User{}
	for _, each := range ass.Items {
		usr, err := srv.Users.Get(each.AssignedTo).Do()
		if err != nil {
			return fmt.Errorf("unable to retrieve user in domain: %v", err)
		}
		users = append(users, usr)
	}
	// find all the users

	if optionJSON(c, users) {
		return nil
	}
	for _, u := range users {
		fmt.Println(u.PrimaryEmail)
	}
	return nil
}
