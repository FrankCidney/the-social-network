package follow

type Service interface {
	SendFollowRequest(senderID, targetID string) error // checks is_public: auto-follows or creates pending request
	AcceptRequest(recipientID, senderID string) error
	DeclineRequest(recipientID, senderID string) error
	Unfollow(followerID, followedID string) error
	IsFollowing(followerID, followedID string) (bool, error)
}
