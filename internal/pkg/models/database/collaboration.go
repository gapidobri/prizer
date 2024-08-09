package database

type CreateCollaboration struct {
	CollaboratorId        string `db:"collaborator_id"`
	CollaborationMethodId string `db:"collaboration_method_id"`
}
