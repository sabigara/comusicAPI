package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	comusic "github.com/sabigara/comusicAPI"
)

type InvitationRepository struct {
	*sqlx.DB
}

func NewInvitationRepository(db *sqlx.DB) *InvitationRepository {
	return &InvitationRepository{DB: db}
}

func (r *InvitationRepository) Filter(email, groupID string) ([]*comusic.Invitation, error) {
	ret := []*comusic.Invitation{}
	rows, err := r.Query(
		`SELECT email, group_id, group_type, is_accepted FROM invitations 
		 WHERE email = ? OR group_id = ?`,
		email, groupID,
	)
	if err != nil {
		return nil, fmt.Errorf("mysql.invitation_repository.GetContents: %w", err)
	}
	for rows.Next() {
		invite := &comusic.Invitation{}
		var groupType string
		err := rows.Scan(&invite.Email, &invite.GroupID, &groupType, &invite.IsAccepted)
		if err != nil {
			return nil, fmt.Errorf("mysql.invitation_repository.GetContents: %w", err)
		}
		invite.GroupType = comusic.NewGroupType(groupType)
		ret = append(ret, invite)
	}
	return ret, nil
}

func (r *InvitationRepository) GetByIDs(email, groupID string) (*comusic.Invitation, error) {
	invite := &comusic.Invitation{}
	var groupType string
	row := r.QueryRow(
		`SELECT email, group_id, group_type, is_accepted FROM invitations 
		 WHERE email = ? AND group_id = ?`,
		email, groupID,
	)

	err := row.Scan(&invite.Email, &invite.GroupID, &groupType, &invite.IsAccepted)
	if err != nil {
		return nil, fmt.Errorf("mysql.invitation_repository.GetContents: %w", err)
	}
	invite.GroupType = comusic.NewGroupType(groupType)

	return invite, nil
}

func (r *InvitationRepository) Create(email, groupID string, groupType comusic.GroupType) error {
	_, err := r.Exec(`
		INSERT INTO invitations (group_id, email, group_type)
		VALUES (?, ?, ?)`,
		groupID, email, groupType,
	)
	if err != nil {
		return fmt.Errorf("mysql.invitation_repository.Create: %w", err)
	}
	return nil
}

func (r *InvitationRepository) Accept(email, groupID string) error {
	res, err := r.Exec(`
		UPDATE invitations
		SET is_accepted = true
		WHERE email = ? AND group_id=?`,
		email, groupID,
	)
	if err != nil {
		return fmt.Errorf("mysql.invitation_repository.Update: %w", err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("mysql.invitation_repository.Update: %w", err)

	}
	return nil
}
