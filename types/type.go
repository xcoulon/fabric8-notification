package types

type NotificationType string

const (
	WorkitemCreate         NotificationType = "workitem.create"
	WorkitemUpdate         NotificationType = "workitem.update"
	CommentCreate          NotificationType = "comment.create"
	CommentUpdate          NotificationType = "comment.update"
	UserEmailUpdate        NotificationType = "user.email.update"
	UserDeactivation       NotificationType = "user.deactivation"
	InvitationTeamNoorg    NotificationType = "invitation.team.noorg"
	InvitationSpaceNoorg   NotificationType = "invitation.space.noorg"
	AnalyticsNotifyCVE     NotificationType = "analytics.notify.cve"
	AnalyticsNotifyVersion NotificationType = "analytics.notify.version"
)

var notifiers = map[NotificationType][]string{
	AnalyticsNotifyCVE:     {"fabric8-gemini-server"},
	UserEmailUpdate:        {"fabric8-auth"},
	AnalyticsNotifyVersion: {"fabric8-gemini-server"},
}

func (nType NotificationType) Notifiers() []string {
	return notifiers[nType]
}
