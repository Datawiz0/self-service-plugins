package main

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo"
	"github.com/rightscale/godo"
)

// Setup routes for droplet actions
func SetupDropletsRoutes(e *echo.Echo) {
	e.Get("", listDroplets)
	e.Get("/:id", showDroplet)
	e.Post("", createDroplet)
	e.Delete("/:id", deleteDroplet)
	e.Get("/:id/kernels", listDropletKernels)
	e.Get("/:id/snapshots", listDropletSnapshots)
	e.Get("/:id/backups", listDropletBackups)
	e.Get("/:id/actions", listDropletActions)
	e.Get("/:id/neighbors", listDropletNeighbors)
}

// Helper function that builds a droplet resource href from its id
func dropletHref(id int) string {
	return fmt.Sprintf("/droplets/%d", id)
}

func listDroplets(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	list, err := paginateDroplets(client.Droplets.List)
	if err != nil {
		return err
	}
	return Respond(c, list)
}

func showDroplet(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	root, _, err := client.Droplets.Get(id)
	if err != nil {
		return err
	}
	droplet := DropletFromApi(root.Droplet)
	return Respond(c, droplet)
}

func createDroplet(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	req := godo.DropletCreateRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	root, _, err := client.Droplets.Create(&req)
	if err != nil {
		return err
	}
	c.Response.Header().Set("Location", dropletHref(root.Droplet.ID))
	return RespondNoContent(c)
}

func deleteDroplet(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	_, err = client.Droplets.Delete(id)
	if err != nil {
		return err
	}
	return RespondNoContent(c)
}

func listDropletKernels(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	list := []godo.Kernel{}
	opt := &godo.ListOptions{}
	for {
		kernels, resp, err := client.Droplets.ListKernels(id, opt)
		if err != nil {
			return err
		}
		list = append(list, kernels...)
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return err
		}
		opt.Page = page + 1
	}
	return Respond(c, list)
}

func listDropletSnapshots(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	return listDropletImages(c, client.Droplets.ListSnapshots)
}

func listDropletBackups(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	return listDropletImages(c, client.Droplets.ListBackups)
}

func listDropletImages(c *echo.Context, lister func(int, *godo.ListOptions) ([]godo.Image, *godo.Response, error)) error {
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	list, err := paginateImages(func(opt *godo.ListOptions) ([]godo.Image, *godo.Response, error) {
		return lister(id, opt)
	})
	if err != nil {
		return err
	}
	return Respond(c, list)
}

func listDropletActions(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	list, err := paginateActions(func(opt *godo.ListOptions) ([]godo.Action, *godo.Response, error) {
		return client.Droplets.ListActions(id, opt)
	})
	if err != nil {
		return err
	}
	return Respond(c, list)
}

func listDropletNeighbors(c *echo.Context) error {
	client, err := GetDOClient(c)
	if err != nil {
		return err
	}
	id, err := getIDParam(c)
	if err != nil {
		return err
	}
	droplets, _, err := client.Droplets.ListNeighbors(id)
	if err != nil {
		return err
	}
	var list []*Droplet
	for _, d := range droplets {
		list = append(list, DropletFromApi(&d))
	}
	return Respond(c, list)
}

// Helper function that retrieves the droplet id (number) from the request parameters
func getIDParam(c *echo.Context) (int, error) {
	sid := c.Param("id")
	if sid == "" {
		return 0, fmt.Errorf("missing droplet id")
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		return 0, fmt.Errorf("invalid droplet id '%s' - must be a number", sid)
	}
	return id, nil
}

// Paginate over droplet listing
func paginateDroplets(lister func(opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)) ([]godo.Droplet, error) {
	list := []godo.Droplet{}
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := lister(opt)
		if err != nil {
			return nil, err
		}
		list = append(list, droplets...)
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