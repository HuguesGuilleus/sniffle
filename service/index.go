package service

import (
	"sniffle/front"
	"sniffle/service/eu_ec_ice"
	"sniffle/tool"
)

var List = []tool.Service{
	{Name: "eu_ec_ice", Do: eu_ec_ice.Do},
	{Name: "front", Do: front.Do},
}
