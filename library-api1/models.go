package main // объявляем основной пакет приложения

import "time" // импортируем пакет для работы с датой и временем

type BookStatus string // создаем свой тип для статуса книги

const ( // открываем блок констант
	bookAvailable BookStatus = "Available" // константа доступной книги
	bookIssued    BookStatus = "Issued"    // константа выданной книги
) // закрываем блок констант

type Book struct { // структура для описания книги
	ID     string     `json:"id"`     // уникальный идентификатор
	Title  string     `json:"title"`  // название книги
	Author string     `json:"author"` // автор книги
	ISBN   string     `json:"isbn"`   // международный номер
	Year   int        `json:"year"`   // год издания
	Status BookStatus `json:"status"` // текущий статус книги
} // конец структуры книги

type User struct { // структура для описания читателя
	ID               string    `json:"id"`                // уникальный идентификатор пользователя
	Name             string    `json:"name"`              // полное имя
	Email            string    `json:"email"`             // электронная почта
	RegistrationDate time.Time `json:"registration_date"` // дата регистрации в системе
} // конец структуры пользователя

type IssueReq struct { // структура для запроса на выдачу книги
	UserID string `json:"user_id"` // айди читателя
	BookID string `json:"book_id"` // айди книги
} // конец структуры запроса выдачи

type ReturnReq struct { // структура для запроса на возврат книги
	BookID string `json:"book_id"` // айди книги
} // конец структуры запроса возврата