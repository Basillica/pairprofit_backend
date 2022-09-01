package responses

import (
	"time"
)

type ProfileDataUser struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	FullName  string `json:"full_name"`
	LastName  string `json:"last_name"`
}

type ProfileDataResponse struct {
	WorkspaceId       *int                   `json:"workspace_id"`
	User              ProfileDataUser        `json:"user"`
	Workspace         *string                `json:"workspace"`
	Company           *string                `json:"company"`
	Role              *string                `json:"role"`
	Url               *string                `json:"url"`
	ImageUri          *string                `json:"imageUri"`
	ProfileInboxID    *int                   `json:"profile_inbox_id"`
	PersonalInboxID   *int                   `json:"personal_inbox_id"`
	About             *string                `json:"about"`
	PreferredLanguage *string                `json:"preferred_language"`
	RecentGroups      []*GroupObjectResponse `json:"recent_groups"`
}

type GroupObjectResponse struct {
	ID                  int                  `json:"id"`
	IsPrivate           *bool                `json:"is_private"`
	PersonalGroupOwner  *int                 `json:"personal_group_owner"`
	GroupName           string               `json:"name"`
	WorkspaceID         *int                 `json:"workspace_id"`
	SingleUserObj       []*UserProfileV2User `json:"single_user_obj"`
	LastActivity        *time.Time           `json:"last_activity"`
	UnreadMessages      *int                 `json:"unread_messages"`
	ActivityCount       *int                 `json:"activity_count"`
	LatestTranscription *string              `json:"latest_transcription"`
	LatestUserID        *int                 `json:"latest_user_id"`
	LatestEntryS3ID     *string              `json:"latest_entry_s3_id"`
	PreferredLanguage   *string              `json:"preferred_language"`
	GroupUUID           *string              `json:"group_uuid"`
	LatestIsTopic       *bool                `json:"latest_is_topic"`
	LatestTopicName     *string              `json:"latest_topic_name"`
	GroupType           *string              `json:"group_type"`
	TalkingPoints       *string              `json:"talking_points"`
	ImageURI            *string              `json:"imageUri"`
}

type UserProfileV2User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	FullName  string `json:"full_name"`
	LastName  string `json:"last_name"`
}
