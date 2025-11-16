package prs

import "errors"

var (
	ErrPRExists    = errors.New("pr уже существует")
	ErrNotFound    = errors.New("pr не найден")
	ErrPRMerged    = errors.New("нельзя менять после MERGED")
	ErrNotAssigned = errors.New("пользователь не был назначен ревьювером")
	ErrNoCandidate = errors.New("нет доступных кандидатов")
)
