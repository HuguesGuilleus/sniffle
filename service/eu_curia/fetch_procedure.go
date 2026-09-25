package eu_curia

import (
	"fmt"

	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/fetch"
	"github.com/HuguesGuilleus/sniffle/tool/sch"
)

// Fetch a single affair by its affair ID.
func fetchProcedure(t *tool.Tool, affID string) {
	body := make([]byte, 0)
	body = append(body, `{"affId":"`...)
	body = append(body, affID...)
	body = append(body, `"}`...)
	request := fetch.R("POST", "https://infocuriaws.curia.europa.eu/elastic-connector/affairId/procedures", body,
		"Content-Type", "application/json; charset=utf-8",
		"Accept", "application/json",
	)

	dto := struct {
		// SearchHits []struct {
		// 	ID        string `json:"id"`
		// 	InnerHits struct {
		// 		Document struct {
		// 			MaxScore *json.RawMessage `json:"maxScore"`
		// 		} `json:"document"`
		// 	} `json:"innerHits"`
		// } `json:"searchHits"`
	}{}
	if tool.FetchJSON(t, procedureTypes, &dto, request) {
		fmt.Println("fail!")
		return
	}
}

var (
	affID       = sch.Regexp(`[TC]/\d{4}/\d\d/\w{10}/\d\d`)
	procedureID = sch.Regexp(`[TC]/\d{4}/\d\d/\w{10}/\d\d/[IPR]/\d\d`)
	otherID     = sch.Regexp(`[TC]/\d{4}/\d\d/\w{10}/\d\d/[IPR]/\d\d#[CT]-\d{1,4}/\d\d( DEP| RX)?`)

	matCode, matCodeHTML = sch.Def("matCode", sch.EnumString("ADHE", "AGRI.FA", "AGRI.OCMA", "AGRI.PADI", "AGRI.PHYT", "AGRI", "AXT.ACIN.ASSO", "AXT.ACIN", "AXT.COPT", "AXT.POCC.DUMP", "AXT.POCC.SAUV", "AXT.POCC", "AXT", "BEI", "COAD", "COES.FC", "COES", "COMP", "CONC.AIDE", "CONC.ENTE", "CONC.EPU", "CONC.MERG", "CONC.PODO", "CONC", "CONS", "CORE", "CULT", "DFON.CHDF", "DFON.CHDF", "DFON", "DGEN", "DRIN.ACTE", "DRIN.RESP", "DRIN", "EFJS", "ELSJ.COJC", "ELSJ.COJP", "ELSJ.COPO", "ELSJ.FAI.ASIL", "ELSJ.FAI.FRON", "ELSJ.FAI.IMMI", "ELSJ.FAI", "ELSJ", "EMPL", "ENER", "ENTR", "ENVI.DECH", "ENVI.POLL.EQEG", "ENVI.POLL", "ENVI", "FINA.BUDG", "FINA.LUFR", "FINA.RPRO", "FINA", "FISC.ACCI", "FISC.TVA", "FISC", "LCCA", "LCMA.CDOU.CODO", "LCMA.CDOU", "LCMA.RSTR.MEEQ", "LCMA.RSTR", "LCMA.UDOU.CLTA", "LCMA.UDOU.TEEQ", "LCMA.UDOU", "LCMA", "LCPE.ETAB", "LCPE.LCT.SESO", "LCPE.LCT", "LCPE.LCT", "LCPE", "LPSE", "MARC", "MARI", "NDCU.CITU.DDES", "NDCU.CITU", "NDCU.DISC.NAT", "NDCU.DISC", "NDCU", "PECH", "PEM.BCE", "PEM.EURO", "PEM.SEBC", "PEM.SURV", "PEM", "PESC.MRES", "PESC", "PIIC.BREV", "PIIC.DADV", "PIIC.DADV", "PIIC.DEMO", "PIIC.MARQ", "PIIC.OBVE", "PIIC", "PIND", "PRIN.DOCU", "PRIN.PDON", "PRIN", "PRIV", "PROJ.PROC", "PROJ.PROC", "PROJ", "PROJ", "PROP", "PSOC.EGA", "PSOC", "RAPL.ASSU", "RAPL.CHIM.REAC", "RAPL.CHIM", "RAPL.ICT", "RAPL.MAPU", "RAPL.MAPU", "RAPL.MAPU", "RAPL.NRT", "RAPL.PHAR", "RAPL.RFCE", "RAPL.RMSN.RMN", "RAPL.RMSN.RSN", "RAPL.RMSN", "RAPL.SEALI", "RAPL.SEFI", "RAPL", "RDTE", "RETR", "SANI", "SANT", "STAT", "TOUR", "TRAN.IAPA", "TRAN.SEAR", "TRAN"))

	formationTypes = sch.EnumString("_C", "_T", "10_T", "1CIN_T", "2_T", "3_T", "4_T", "5_T", "6_T", "7_T", "8_T", "9_T", "C_ADM_PV_C", "C_C", "I_C", "III_C", "IV_C", "IX_C", "IX_T", "P_C", "P_T", "V_C", "VI_C", "VIII_C", "VP_C", "VP_T", "X_C")

	resultTypes = sch.EnumString("ANNU", "ANNU=RI", "ASSI", "AVIS", "CARE", "COMP", "CONS", "DPEN", "FONC", "INTP", "INTV", "OMST", "OPPO", "PREJ", "PREJURG", "PVOI", "PVOIADM", "PVOIADM=OB", "RECT", "REEX", "REFE", "REFE=NL", "REFE=OB", "REFE=RF", "REFE=RI", "RENV", "RESP", "REVI", "SANC")
)

var procedureTypes = sch.Map(
	sch.FieldSR("totalHits", sch.StrictPositiveInt()),
	sch.FieldSR("maxScore", sch.ConstFloat(0.0)),
	sch.FieldSR("maxScoreOthersDocuments", sch.ConstFloat(0.0)),
	sch.FieldSR("totalHitsOthersDocs", sch.ConstInt(0)),
	sch.FieldSR("searchHits", sch.ArrayMin(1, sch.Map(
		sch.FieldSR("index", sch.Regexp(`index-affair@\d\d-\d\d-\d\d\d\d`)),
		sch.FieldSR("id", sch.Regexp(`[TC]/\d{4}/\d{2}/\w{10}/\d{2}`)),
		sch.FieldSR("score", sch.AnyFloat()),
		sch.FieldSR("sortValues", sch.EmptyArray()),
		sch.FieldSR("content", sch.Map(
			sch.FieldSR("id", procedureID),
			sch.FieldSR("refType", sch.String("affair")),
			sch.FieldSR("affId", affID),
			sch.FieldSR("publishedId", sch.Regexp(`(Avis |C-|T-)\d+/\d+`)),
			sch.FieldSR("publishedAffId", sch.Regexp(`(Avis |C-|T-)\d+/\d+`)),
			sch.FieldSR("procedureId", procedureID),
			sch.FieldSR("jurisdiction", sch.EnumString("Cour de justice", "Tribunal")),
			sch.FieldSR("jurisdictionCode", sch.EnumString("C", "T")),
			sch.FieldSR("order", sch.StrictPositiveInt()),
			sch.FieldSR("year", sch.StrictPositiveStringInt()),
			sch.FieldSO("nature", sch.EnumString(
				"Aide juridictionnelle",
				"Demandes d'avis",
				"Réexamen",
				"Opposition",
				"Interprétation",
				"Omission de statuer",
				"Pourvoi sur intervention",
				"Pourvoi sur référé",
				"Pourvoi",
				"Procédure préjudicielle",
				"Propriété intellectuelle",
				"Recours de fonctionnaires",
				"Recours direct",
				"Rectification",
				"Renvoi après annulation",
				"Révision",
				"Taxation des dépens",
			)),
			sch.FieldSO("natureCode", sch.EnumString(
				"AJ",
				"AVIS",
				"DEPE",
				"EVI",
				"ITRP",
				"OPPO",
				"OST",
				"PI",
				"PIV",
				"PV",
				"PVI",
				"PVR",
				"RAL",
				"RD",
				"RECT",
				"REVI",
				"RF",
				"RP",
				"RX",
			)),
			sch.FieldSR("suiteNumber", sch.StrictPositiveInt()),
			sch.FieldSR("procType", sch.EnumString("I", "P", "R")),
			sch.FieldSR("procNumber", sch.StrictPositiveInt()),
			sch.FieldSR("usualNameML", sch.ArrayMin(1, sch.Map(
				sch.Field(languesLowerType, sch.NotEmptyString(), true),
			))),
			sch.FieldSR("affObjectML", sch.EmptyArray()),
			sch.FieldSR("introductionDate", sch.Time("2006-01-02")),
			sch.FieldSR("introductionYear", sch.AnyInt()),
			sch.FieldSO("closeDate", sch.Time("2006-01-02")),
			sch.FieldSO("closeYear", sch.AnyInt()),
			sch.FieldSR("affairState", sch.EnumString("Affaire pendante", "Affaire clôturée")),
			sch.FieldSR("affairStateCode", sch.EnumString("ENC", "CLOTPUB")),
			sch.FieldSR("procLang", languesUpperTypes),
			sch.FieldSR("joinExist", sch.IntervalInt(0, 1)),
			sch.FieldSR("pilot", sch.IntervalInt(0, 1)),
			sch.FieldSR("joinAffairs", sch.Array(otherID)),
			sch.FieldSR("pourvoiAffIds", sch.Array(affID)),
			sch.FieldSR("originPourvoiIds", sch.Array(otherID)),
			sch.FieldSR("reExamenAffIds", sch.ArrayMax(1, otherID)),
			sch.FieldSR("originReExamanIds", sch.ArrayMax(1, otherID)),
			sch.FieldSR("transfertAffIds", sch.Array(otherID)),
			sch.FieldSR("originTransfertIds", sch.Array(otherID)),
			sch.FieldSO("numDocs", sch.Regexp(`[CT]\d{4}/\d{4}/[JOPR]`)),
			sch.FieldSR("matCode", sch.Array(matCode)),
			sch.FieldSO("matCodeML", sch.Array(sch.Map(
				sch.FieldSR("code", matCode),
				sch.FieldSR("label", sch.Array(sch.Map(
					sch.Field(languesLowerType, sch.NotEmptyString(), true),
				))),
			))),
			sch.FieldSO("reportingJudge", sch.NotEmptyString()),
			sch.FieldSR("formation", formationTypes),
			sch.FieldSO("avg", sch.NotEmptyString()),
			sch.FieldSO("reportingJudgeML", sch.Array(sch.Map(
				sch.FieldSR("code", sch.NotEmptyString()),
				sch.FieldSR("label", sch.ArrayMin(1, sch.Map(
					sch.Field(languesLowerType, sch.NotEmptyString(), true),
				))),
			))),
			sch.FieldSR("formationML", sch.Array(sch.Map(
				sch.FieldSR("code", formationTypes),
				sch.FieldSR("label", sch.ArrayMin(1, sch.Map(
					sch.Field(languesLowerType, sch.NotEmptyString(), true),
				))),
			))),
			sch.FieldSO("advocateML", sch.Array(sch.Map(
				sch.FieldSR("code", sch.NotEmptyString()),
				sch.FieldSR("label", sch.ArrayMin(1, sch.Map(
					sch.Field(languesLowerType, sch.NotEmptyString(), true),
				))),
			))),
			sch.FieldSR("resultTypes", sch.Array(resultTypes)),
			sch.FieldSO("procedureResultTypeML", sch.Array(sch.Map(
				sch.FieldSR("code", resultTypes),
				sch.FieldSR("label", sch.ArrayMin(1, sch.Map(
					sch.Field(languesLowerType, sch.NotEmptyString(), true),
				))),
			))),
			sch.FieldSR("parties", sch.AnyString()),
			sch.FieldSO("conclsDate", sch.ArrayRange(0, 1, sch.Time("2006-01-02"))),
			sch.FieldSR("conclLang", sch.Or(sch.String(""), languesUpperTypes)),
			sch.FieldSO("oqp", sch.NotEmptyString()),
			sch.FieldSO("oqpCountry", sch.NotEmptyString()), ///////////////
			sch.FieldSO("calendars", sch.ArrayRange(0, 3, sch.Map(
				sch.FieldSR("seanceId", sch.ConstInt(0)),
				sch.FieldSR("seanceType", sch.EnumString("PO", "ARRET")),
				sch.FieldSR("seanceDate", sch.Time("2006-01-02")),
			))),
			sch.FieldSR("oppObsFag", sch.ConstInt(0)),
			sch.FieldSO("procClosDate", sch.Time("2006-01-02")),
			sch.FieldSO("doctrineNotes", sch.ArrayMin(1, sch.NotEmptyString())),
			sch.FieldSR("lastModifiedDate", sch.Time("2006-01-02T15:04:05.999")),
			sch.FieldSR("contentLength", sch.ConstInt(0)),
			sch.FieldSR("docNoPart", sch.ConstInt(0)),

			sch.FieldSO("refDosper", sch.NotEmptyString()),
			sch.FieldSO("procIntroDate", sch.Time("2006-01-02")),

			sch.FieldSO("pourvoiLnsDec", sch.Any()),
			sch.FieldSO("procState", sch.Any()),
			sch.FieldSO("description", sch.Any()),
			sch.FieldSO("dispNatVisRef", sch.Any()),
			sch.FieldSO("dispNatVisRefWithDoscur", sch.Any()),
			sch.FieldSO("docTypeCode", sch.Any()),
			sch.FieldSO("idPilot", sch.Any()),
			sch.FieldSO("dispIntVisRefWithDoscur", sch.Any()),
			sch.FieldSO("dispIntVisRef", sch.Any()),
			sch.FieldSO("reExamenLnsDec", sch.Any()),
		)),
		sch.FieldSR("highlightFields", sch.Map(
			sch.FieldSO("doctrineNotes", sch.Array(sch.NotEmptyString())),
			sch.FieldSO("publishedId.keyword", sch.ArraySize(1, sch.NotEmptyString())),
		)),
		sch.FieldSR("innerHits", sch.Map(
			sch.FieldSR("document", sch.Map(
				sch.FieldSR("totalHits", sch.PositiveInt()),
				sch.FieldSR("totalHitsRelation", sch.String("EQUAL_TO")),
				sch.FieldSR("maxScore", sch.Or(
					sch.ConstFloat(1.0),
					sch.String("NaN"),
				)),
				sch.FieldSR("searchHits", sch.Array(sch.Map(
					sch.FieldSR("index", sch.Regexp(`index-affair@\d\d-\d\d-\d\d\d\d`)),
					sch.FieldSR("id", sch.Regexp(`[TC]/\d{4}/\d{2}/\w{10}/\d{2}/[IPR]/\d{2}-\d+`)),
					sch.FieldSR("score", sch.ConstFloat(1.0)),
					sch.FieldSO("sortValues", sch.EmptyArray()),
					sch.FieldSO("content", sch.Map(
						sch.FieldSR("id", sch.Regexp(`[TC]/\d{4}/\d{2}/\w{10}/\d{2}/[IPR]/\d{2}-\d+`)),
						sch.FieldSR("refType", sch.String("document")),
						sch.FieldSR("joinExist", sch.ConstInt(0)),
						sch.FieldSR("pilot", sch.ConstInt(0)),
						sch.FieldSR("oppObsFag", sch.ConstInt(0)),
						sch.FieldSR("docId", sch.StrictPositiveStringInt()),
						sch.FieldSR("logicDocId", sch.Regexp(`id_\d+`)),
						sch.FieldSO("jpLogicDocId", sch.Or(
							sch.String("id_null"),
							sch.Regexp(`id_\d+`),
						)),
						sch.FieldSO("idProcedure", sch.Regexp(`[TC]/\d{4}/\d{2}/\w{10}/\d{2}/[IPR]/\d{2}`)),
						sch.FieldSR("docType", sch.EnumString(
							"Arrêt (extraits)",
							"Arrêt (JO)",
							"Arrêt",
							"Avis (JO)",
							"Conclusions",
							"Décision de réexamen (JO)",
							"Décision de réexamen",
							"Demande (JO)",
							"Demande de décision préjudicielle - Addendum",
							"Demande de décision préjudicielle – Corrigendum",
							"Demande de décision préjudicielle",
							"Mémoire ou observations écrites",
							"Ordonnance (Information)",
							"Ordonnance (JO)",
							"Ordonnance (Sommaire)",
							"Ordonnance",
							"Radiation (JO)",
							"Requête (JO)",
							"Résumé",
						)),
						sch.FieldSR("docTypeCode", sch.EnumString(
							"ARR_COMM",
							"ARRET_EXT",
							"ARRET_NP",
							"ARRET",
							"AVI_COMM",
							"CONCL",
							"DDP_ADD",
							"DDP_COMM",
							"DDP_CORR",
							"DDP",
							"DECISION",
							"OBSRP_PUB",
							"ORD_COMM",
							"ORD_DR",
							"ORD_NP",
							"ORD_SOM",
							"ORD",
							"PV_COMM",
							"RAD_COMM",
							"REQ_COMM",
							"RES",
							"RX_COMM",
						)),
						sch.FieldSO("lnsDec", sch.Or(sch.String(""), sch.StrictPositiveStringInt())),
						sch.FieldSO("numDoc", sch.Or(sch.String(""), sch.Regexp(`^[CT]\d{4}/\d{4}/[ROPJ](/\d)?$`))),
						sch.FieldSR("docDate", sch.Time("2006-01-02")),
						sch.FieldSR("docFormat", sch.EnumString("HTML", "PDF")),
						sch.FieldSR("docFormats", sch.ArrayItems(sch.EnumString("HTML", "PDF"))),
						sch.FieldSR("docLang", languesUpperTypes),
						sch.FieldSR("docUrl", sch.String("")),
						sch.FieldSR("contentLength", sch.ConstInt(0)),
						sch.FieldSO("contentML", sch.ArrayMin(1, sch.Map(
							sch.Field(languesLowerType, sch.AnyString(), true),
						))),
						sch.FieldSO("description", sch.AnyString()),
						sch.FieldSO("docNoPart", sch.ConstInt(1)),
						sch.FieldSO("celex", sch.AnyString()),
						sch.FieldSO("docBelongToPilot", sch.String("NO")),
						sch.FieldSO("affairMatCode", sch.Array(matCode)),
						sch.FieldSO("groupByLogicalId", sch.Array(sch.MapExtra())),
						sch.FieldSR("published", sch.EnumString("yes", "no")),
						sch.FieldSR("bulleting", sch.AnyString()),
						sch.FieldSR("bulletingWithCorrespondence", sch.AnyString()),
						sch.FieldSR("descriptorML", sch.EmptyArray()),
						sch.FieldSR("dispNatVisRef", sch.AnyString()),
						sch.FieldSR("dispIntVisRef", sch.AnyString()),
						sch.FieldSR("dispNatVisRefWithDoscur", sch.AnyString()),
						sch.FieldSR("dispIntVisRefWithDoscur", sch.AnyString()),
						sch.FieldSR("refJoML", sch.EmptyArray()),
						sch.FieldSR("affairJurisdiction", sch.EnumString(
							"Cour de justice",
							"Tribunal",
						)),
						sch.FieldSR("affairJurisdictionCode", sch.EnumString("C", "T")),
						sch.FieldSR("decnat", sch.EmptyArray()),
						sch.FieldSR("routing", sch.Regexp(`[TC]/\d{4}/\d{2}/\w{10}/\d{2}/[IPR]/\d{2}`)),
						sch.FieldSO("lienPublicationRefJo", sch.URL("http://eur-lex.europa.eu/legal-content/{0}/TXT/PDF/")),
						sch.FieldSR("ecli", sch.Or(
							sch.String(""),
							sch.Regexp(`^ECLI:EU:[CT]:\d{4}:\d+$`),
						)),
						sch.FieldSO("idPublished", sch.Regexp(`^([TC]-|Avis )\d+/\d+( R| P(\([IR]\))?| | RX| RENV| RII| REC| DEP( II)?| \(PPU\)| PPU| P-R| R INTP)?$`)),
						sch.FieldSO("citationsMotif", sch.Array(sch.Map(
							sch.FieldSR("category", sch.AnyString()),
							sch.FieldSR("reference", sch.AnyString()),
							sch.FieldSR("referenceTransformed", sch.AnyString()),
							sch.FieldSR("relation", sch.AnyString()),
							sch.FieldSR("section", sch.AnyString()),
						))),
						sch.FieldSO("affairMatCodeML", sch.Any()),
						sch.FieldSO("matCodeML", sch.ArrayMin(1, sch.Map(
							sch.FieldSR("code", sch.NotEmptyString()),
							sch.FieldSR("label", sch.Array(sch.MapExtra(
							// todo
							))),
						))),
						sch.FieldSO("citationsDispositif", sch.Any()),
						sch.FieldSO("recueil", sch.Any()),
						sch.FieldSO("pubNum", sch.Any()),
						sch.FieldSO("repJurRefML", sch.Any()),
						sch.FieldSO("citationsConclusion", sch.Any()),
						sch.FieldSO("authorCode", sch.Any()),
						sch.FieldSO("author", sch.Any()),
					)),
					sch.FieldSO("highlightFields", sch.Any()),
					sch.FieldSO("innerHits", sch.Map()),
					sch.FieldSR("routing", sch.Regexp(`[TC]/\d{4}/\d{2}/\w{10}/\d{2}/[IPR]/\d{2}`)),
					sch.FieldSO("matchedQueries", sch.EmptyArray()),
				))),
				sch.FieldSR("empty", sch.AnyBool()),
			)),
		)),
		sch.FieldSR("matchedQueries", sch.EmptyArray()),
	))),
)
