package main // основной пакет программы

import ( // импорты
	"fmt"      // для вывода в консоль
	"net/http" // для запуска сервера
) // конец импортов

func main() { // главная функция
	db := initDB()   // инициализируем базу (функция сама подтянется из db.go)
	defer db.Close() // закрываем базу при выключении сервера

	mux := http.NewServeMux()                     // создаем стандартный роутер
	mux.HandleFunc("/books", handleBooks(db))     // привязываем обработчик (из handlers.go)
	mux.HandleFunc("/users", handleUsers(db))     // привязываем обработчик (из handlers.go)
	mux.HandleFunc("/issues", handleIssues(db))   // привязываем выдачу (из handlers.go)
	mux.HandleFunc("/returns", handleReturns(db)) // привязываем возврат (из handlers.go)

	fmt.Println("сервер запущен на порту 8080") // пишем в консоль
	http.ListenAndServe(":8080", mux)           // запускаем сервер навсегда
} // конец программы
