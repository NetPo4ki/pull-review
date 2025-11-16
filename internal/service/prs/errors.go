package prs

import "errors"

var (
	ErrPRExists    = errors.New("PR уже существует")
	ErrNotFound    = errors.New("PR не найден")
	ErrPRMerged    = errors.New("Нельзя менять после MERGED")
	ErrNotAssigned = errors.New("Пользователь не был назначен ревьювером")
	ErrNoCandidate = errors.New("Нет доступных кандидатов")
)
