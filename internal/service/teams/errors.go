package teams

import "errors"

var (
	ErrTeamExists = errors.New("Команда уже существует")
	ErrNotFound   = errors.New("Команда не найдена")
)
