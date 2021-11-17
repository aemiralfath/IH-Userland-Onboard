package session

import "context"

func (r *Repository) DeleteOtherSession(ctx context.Context, id string) error {
	
	if err := r.SessionStore.DeleteOtherSession(ctx, id); err != nil {
		return err
	}

	return nil
}
