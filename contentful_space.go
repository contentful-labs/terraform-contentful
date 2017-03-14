package main

type collectionProperties struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
	Sys   struct {
		Type string `json:"type"`
	} `json:"sys"`
}

type link struct {
	Sys struct {
		Type     string `json:"type"`
		LinkType string `json:"linkType"`
		ID       string `json:"id"`
	} `json:"sys"`
}

type spaceSys struct {
	Type      string `json:"type"`
	ID        string `json:"id"`
	Version   int    `json:"version"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	CreatedBy link   `json:"createdBy"`
	UpdatedBy link   `json:"updatedBy"`
}

type spaceData struct {
	Sys  spaceSys `json:"sys"`
	Name string   `json:"name"`
}

type spaceCollection struct {
	collectionProperties
	Items []spaceData `json:"items"`
}
