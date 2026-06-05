package main // основной пакет

import ( // открываем блок импортов
	"database/sql" // пакет для работы с sql
	"errors"       // пакет для создания ошибок
	"time"         // пакет для времени

	"github.com/google/uuid"        // пакет для генерации уникальных айди
	_ "github.com/mattn/go-sqlite3" // анонимный импорт драйвера sqlite
) // закрываем импорты

func initDB() *sql.DB { // функция инициализации базы данных
	db, err := sql.Open("sqlite3", "library.db") // открываем файл базы данных напрямую
	if err != nil {                              // проверяем на ошибку
		panic(err) // критическая ошибка останавливает работу
	} // конец условия

	query := ` 
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		registration_date TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS books (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		isbn TEXT NOT NULL,
		year INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT 'Available'
	);
	CREATE TABLE IF NOT EXISTS issues (
		id TEXT PRIMARY KEY,
		book_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		issue_date TEXT NOT NULL,
		due_date TEXT NOT NULL,
		return_date TEXT,
		FOREIGN KEY(book_id) REFERENCES books(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);` // sql запрос на создание трех таблиц

	_, err = db.Exec(query) // выполняем запрос к базе
	if err != nil {         // если запрос упал
		panic(err) // останавливаем приложение
	} // конец проверки
	return db // возвращаем готовое подключение
} // конец функции

func createBook(db *sql.DB, b *Book) error { // функция добавления книги
	b.ID = uuid.New().String()                                                                                          // генерируем новый айди
	b.Status = bookAvailable                                                                                            // ставим статус по умолчанию
	_, err := db.Exec("INSERT INTO books VALUES (?, ?, ?, ?, ?, ?)", b.ID, b.Title, b.Author, b.ISBN, b.Year, b.Status) // вставляем в таблицу
	return err                                                                                                          // возвращаем ошибку если есть
} // конец функции

func getBooks(db *sql.DB) ([]Book, error) { // функция получения всех книг
	rows, err := db.Query("SELECT id, title, author, isbn, year, status FROM books") // делаем запрос
	if err != nil {
		return nil, err
	} // возвращаем ошибку
	defer rows.Close() // закрываем строки после чтения
	var books []Book   // создаем пустой срез книг
	for rows.Next() {  // бежим по строкам
		var b Book                                                         // создаем временную переменную
		rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Year, &b.Status) // читаем данные в структуру
		books = append(books, b)                                           // добавляем в срез
	} // конец цикла
	return books, nil // возвращаем результат
} // конец функции

func createUser(db *sql.DB, u *User) error { // функция создания читателя
	u.ID = uuid.New().String()                                                                                                 // генерируем айди
	u.RegistrationDate = time.Now()                                                                                            // берем текущее время
	_, err := db.Exec("INSERT INTO users VALUES (?, ?, ?, ?)", u.ID, u.Name, u.Email, u.RegistrationDate.Format(time.RFC3339)) // сохраняем
	return err                                                                                                                 // возвращаем статус
} // конец функции

func issueBook(db *sql.DB, userID, bookID string) error { // функция выдачи книги (транзакция)
	tx, _ := db.Begin()                                                        // начинаем транзакцию
	defer tx.Rollback()                                                        // откатываем если что-то пойдет не так
	var status string                                                          // переменная для статуса
	tx.QueryRow("SELECT status FROM books WHERE id = ?", bookID).Scan(&status) // ищем книгу
	if status != string(bookAvailable) {
		return errors.New("book not available")
	} // проверяем доступна ли

	now := time.Now()              // текущее время
	due := now.AddDate(0, 0, 14)   // прибавляем 14 дней
	issueID := uuid.New().String() // айди выдачи

	tx.Exec("INSERT INTO issues (id, book_id, user_id, issue_date, due_date) VALUES (?, ?, ?, ?, ?)", issueID, bookID, userID, now.Format(time.RFC3339), due.Format(time.RFC3339)) // пишем историю
	tx.Exec("UPDATE books SET status = ? WHERE id = ?", bookIssued, bookID)                                                                                                        // обновляем статус
	return tx.Commit()                                                                                                                                                             // сохраняем транзакцию
} // конец функции

func returnBook(db *sql.DB, bookID string) error { // функция возврата
	tx, _ := db.Begin()                                                                                   // открываем транзакцию
	defer tx.Rollback()                                                                                   // откатываем при ошибке
	var issueID string                                                                                    // айди записи выдачи
	tx.QueryRow("SELECT id FROM issues WHERE book_id = ? AND return_date IS NULL", bookID).Scan(&issueID) // ищем активную выдачу
	if issueID == "" {
		return errors.New("no active issue")
	} // если нет выдачи

	tx.Exec("UPDATE issues SET return_date = ? WHERE id = ?", time.Now().Format(time.RFC3339), issueID) // ставим дату возврата
	tx.Exec("UPDATE books SET status = ? WHERE id = ?", bookAvailable, bookID)                          // освобождаем книгу
	return tx.Commit()                                                                                  // подтверждаем изменения
} // конец функции
