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
		`SELECT id, created_at, updated_at, email, group_id, group_type, is_accepted FROM invitations 
		 WHERE email = ? OR group_id = ?`,
		email, groupID,
	)
	if err != nil {
		return nil, fmt.Errorf("mysql.invitation_repository.GetContents: %w", err)
	}
	for rows.Next() {
		invite := &comusic.Invitation{Meta: &comusic.Meta{}}
		var groupType string
		err := rows.Scan(&invite.ID, &invite.CreatedAt, &invite.UpdatedAt, &invite.Email, &invite.GroupID, &groupType, &invite.IsAccepted)
		if err != nil {
			return nil, fmt.Errorf("mysql.invitation_repository.GetContents: %w", err)
		}
		invite.GroupType = comusic.NewGroupType(groupType)
		ret = append(ret, invite)
	}
	return ret, nil
}

func (r *InvitationRepository) GetByIDs(email, groupID string) (*comusic.Invitation, error) {
	invite := &comusic.Invitation{Meta: &comusic.Meta{}}
	var groupType string
	row := r.QueryRow(
		`SELECT id, created_at, updated_at, email, group_id, group_type, is_accepted FROM invitations 
		 WHERE email = ? AND group_id = ?`,
		email, groupID,
	)

	err := row.Scan(&invite.ID, &invite.CreatedAt, &invite.UpdatedAt, &invite.Email, &invite.GroupID, &groupType, &invite.IsAccepted)
	if err != nil {
		return nil, fmt.Errorf("mysql.invitation_repository.GetContents: %w", err)
	}
	invite.GroupType = comusic.NewGroupType(groupType)

	return invite, nil
}

func (r *InvitationRepository) Create(inv *comusic.Invitation) error {
	_, err := r.Exec(`
		INSERT INTO invitations (id, created_at, updated_at, group_id, email, group_type)
		VALUES (?, ?, ?, ?, ?, ?)`,
		inv.ID, inv.CreatedAt, inv.UpdatedAt, inv.GroupID, inv.Email, inv.GroupType,
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
