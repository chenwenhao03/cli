package route

import (
	"cf"
	"cf/api"
	"cf/requirements"
	"cf/terminal"
	"errors"
	"github.com/codegangsta/cli"
)

type ParkRoute struct {
	ui        terminal.UI
	routeRepo api.RouteRepository
	spaceReq  requirements.SpaceRequirement
	domainReq requirements.DomainRequirement
}

func NewParkRoute(ui terminal.UI, routeRepo api.RouteRepository) (cmd *ParkRoute) {
	cmd = new(ParkRoute)
	cmd.ui = ui
	cmd.routeRepo = routeRepo
	return
}

func (cmd *ParkRoute) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {

	if len(c.Args()) != 2 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "park-route")
		return
	}

	cmd.spaceReq = reqFactory.NewSpaceRequirement(c.Args()[0])
	cmd.domainReq = reqFactory.NewDomainRequirement(c.Args()[1])
	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
		cmd.spaceReq,
		cmd.domainReq,
	}
	return
}

func (cmd *ParkRoute) Run(c *cli.Context) {
	space := cmd.spaceReq.GetSpace()
	domain := cmd.domainReq.GetDomain()
	route := cf.Route{Host: c.String("n"), Domain: domain}

	cmd.ui.Say("Adding url route %s to space %s...",
		terminal.EntityNameColor(route.URL()), terminal.EntityNameColor(space.Name))

	_, apiStatus := cmd.routeRepo.CreateInSpace(route, domain, space)
	if apiStatus.IsError() {
		cmd.ui.Failed(apiStatus.Message)
		return
	}

	cmd.ui.Ok()
}