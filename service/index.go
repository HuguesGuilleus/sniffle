package service

import (
	"sniffle/front"
	"sniffle/service/about"
	"sniffle/service/eu_ec_eci"
	"sniffle/tool"
)

var List = []tool.Service{
	{Name: "about", Do: about.Do},
	{Name: "eu_ec_eci", Do: eu_ec_eci.Do},
	{Name: "front", Do: front.Do},
}
