package types

import "time"

type Contacts struct {
	ID      string
	Name    string
	Email   string
	Details string

	Status int
}

type ContactActivity struct {
	ContactsID   string
	CampaignID   int
	ActivityType int
	ActivityDate time.Time
}

type ResultRow struct {
	CampaignId int
	Count      int
}
