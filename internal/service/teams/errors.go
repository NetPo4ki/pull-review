package teams

import "errors"

var (
	ErrTeamExists = errors.New("команда уже существует")
	ErrNotFound   = errors.New("команда не найдена")
)
