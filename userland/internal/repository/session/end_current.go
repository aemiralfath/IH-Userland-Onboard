package session

import "context"

func (r *Repository) EndCurrentSession(ctx context.Context, id string) error {
	
	if err := r.SessionStore.EndSession(ctx, id); err != nil {
		return err
	}

	return nil
}
