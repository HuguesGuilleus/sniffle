package service

import (
	"sniffle/front"
	"sniffle/service/about"
	"sniffle/service/eu_ec_ice"
	"sniffle/tool"
)

var List = []tool.Service{
	{Name: "about", Do: about.Do},
	{Name: "eu_ec_ice", Do: eu_ec_ice.Do},
	{Name: "front", Do: front.Do},
}
