package tweet

type User struct {
	// The integer representation of the unique identifier for this User. This number is greater than 53 bits and some programming languages may have difficulty/silent defects in interpreting it. Using a signed 64 bit integer for storing this identifier is safe. Use id_str for fetching the identifier to stay on the safe side.
	Id int64 `json:"id"`

	// The string representation of the unique identifier for this Tweet. Implementations should use this rather than the large integer in id.
	IdStr string `json:"id_str"`

	// The name of the user, as they've defined it. Not necessarily a person's name. Typically capped at 20 characters, but subject to change.
	Name string `json:"name"`

	// The screen name, handle, or alias that this user identifies themselves with. screen_names are unique but subject to change. Use id_str as a user identifier whenever possible. Typically a maximum of 15 characters long, but some historical accounts may exist with longer names.
	ScreenName string `json:"screen_name"`

	// The url of their profile image
	ProfileImgURL string `json:"profile_image_url"`
}

type RetweetedStatus struct {
	RetweetCount int64 `json:"retweeted_count"`
	User         User  `json:"user"`
}

type Url struct {
	Url         string `json:"url"`
	DisplayUrl  string `json:"display_url"`
	ExpandedUrl string `json:"expanded_url"`
	Indices     [2]int `json:"indices"`
}

type Entities struct {
	// Urls extracted from the tweet
	Urls []Url `json:"urls"`
}

type Tweet struct {
	// The integer representation of the unique identifier for this Tweet. This number is greater than 53 bits and some programming languages may have difficulty/silent defects in interpreting it. Using a signed 64 bit integer for storing this identifier is safe. Use id_str for fetching the identifier to stay on the safe side. See Twitter IDs, JSON and Snowflake.
	Id int64 `json:"id"`

	// The string representation of the unique identifier for this Tweet. Implementations should use this rather than the large integer in id.
	IdString string `json:"id_str"`

	RetweetedStatus RetweetedStatus `json:"retweeted_status"`

	// The actual UTF-8 text of the status update.
	Text string `json:"text"`

	// The user who posted this Tweet.
	User User `json:"user"`

	// Various entities we might be interested in
	Entities Entities `json:"entities"`
}
