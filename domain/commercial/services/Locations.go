package services

import "strings"

// UgandaLocations is a curated, self-hosted table of major Ugandan towns with
// their APPROXIMATE road distance (km) from KAMPALA CITY CENTRE. It exists so
// transport can be resolved locally without any external map/geocoding service
// (required for offline / on-premises operation). Distances are approximate and
// are NOT a verified road distance; the transport fee itself is a standard
// amount (see InstallationPackage.TransportStandardUGX), not distance-derived.
type UgandaLocation struct {
	Name       string `json:"name"`
	DistanceKm int    `json:"distance_km"`
}

// KampalaCentreName is the reference origin.
const KampalaCentreName = "Kampala City Centre"

var ugandaLocations = []UgandaLocation{
	{KampalaCentreName, 0},
	{"Nansana", 12},
	{"Kira", 15},
	{"Bweyogerere", 16},
	{"Wakiso", 18},
	{"Bulenga", 18},
	{"Gayaza", 20},
	{"Mukono", 21},
	{"Seeta", 22},
	{"Nsangi", 22},
	{"Matugga", 25},
	{"Entebbe", 40},
	{"Mpigi", 45},
	{"Lugazi", 55},
	{"Buwama", 60},
	{"Mityana", 75},
	{"Kayabwe", 75},
	{"Njeru", 78},
	{"Jinja", 80},
	{"Kiboga", 120},
	{"Iganga", 120},
	{"Masaka", 130},
	{"Mubende", 150},
	{"Nakasongola", 150},
	{"Busia", 195},
	{"Hoima", 200},
	{"Lyantonde", 200},
	{"Masindi", 215},
	{"Tororo", 230},
	{"Mbale", 245},
	{"Mbarara", 270},
	{"Soroti", 300},
	{"Bushenyi", 320},
	{"Ishaka", 310},
	{"Fort Portal", 320},
	{"Apac", 320},
	{"Gulu", 330},
	{"Lira", 340},
	{"Ntungamo", 340},
	{"Kasese", 375},
	{"Rukungiri", 380},
	{"Nebbi", 400},
	{"Kabale", 410},
	{"Kanungu", 420},
	{"Kitgum", 430},
	{"Adjumani", 460},
	{"Moroto", 470},
	{"Arua", 475},
	{"Yumbe", 480},
	{"Kisoro", 480},
	{"Koboko", 510},
	{"Kotido", 520},
}

// UgandaLocations returns the location catalogue (ascending by distance).
func UgandaLocations() []UgandaLocation {
	out := make([]UgandaLocation, len(ugandaLocations))
	copy(out, ugandaLocations)
	return out
}

// LocationDistanceKm resolves a location name to its approximate distance from
// Kampala city centre. The second result is false when the location is unknown.
func LocationDistanceKm(name string) (int, bool) {
	target := strings.ToLower(strings.TrimSpace(name))
	if target == "" {
		return 0, false
	}
	for _, loc := range ugandaLocations {
		if strings.ToLower(loc.Name) == target {
			return loc.DistanceKm, true
		}
	}
	return 0, false
}
