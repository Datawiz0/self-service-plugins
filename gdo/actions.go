package main

import (
	"github.com/labstack/echo"
	"github.com/rightscale/gdo/middleware"
	"github.com/rightscale/godo"
)

func SetupActionsRoutes(e *echo.Echo) {
	e.Get("", listActions)
	e.Get("/:id", showAction)
}

func listActions(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	list, er := paginateActions(client.Actions.List)
	return Respond(c, list, er)
}

func showAction(c *echo.Context) *echo.HTTPError {
	client, err := middleware.GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	action, _, er := client.Actions.Get(id)
	return Respond(c, action, er)
}

func paginateActions(lister func(opt *godo.ListOptions) ([]godo.Action, *godo.Response, error)) ([]godo.Action, error) {
	list := []godo.Action{}
	opt := &godo.ListOptions{}
	for {
		actions, resp, err := lister(opt)
		if err != nil {
			return nil, err
		}
		list = append(list, actions...)
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}
		opt.Page = page + 1
	}
	return list, nil
}
