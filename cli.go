package cli

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/banch0/core"
	"github.com/gookit/color"
	"github.com/howeyc/gopass"
)

// OperationsLoop ...
func OperationsLoop(db *sql.DB, loop func(db *sql.DB, cmd string) bool) {
	color.Warn.Println("\n  ===== Добро пожаловать! ===== \n")
	cyan := color.FgLightCyan.Render
	for {
		color.Bold.Println("Список доступных операций:\n")
		fmt.Printf("%s%s\n",
			cyan("1. Авторизация\n"),
			cyan("q. Выйти из приложения\n"))
		color.Bold.Print("Введите команду:")
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err)
		}
		if exit := loop(db, strings.TrimSpace(cmd)); exit {
			return
		}
	}
}

func loginedOperationsLoop(db *sql.DB, name string, loop func(db *sql.DB, cmd, userID string) bool) {
	color.Warn.Print("\n  Привет, ", name, "\n")
	cyan := color.FgLightCyan.Render
	for {
		color.Bold.Println("Список доступных операций:\n")
		fmt.Printf("%s%s%s%s\n",
			cyan("1. Посмотреть список счетов\n"),
			cyan("2. Перевести деньги другому пользователю\n"),
			cyan("3. Оплатить услугу\n"),
			cyan("4. Список банкоматов\n"))
		color.Bold.Print("Введите команду:")
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err)
		}
		if exit := loop(db, strings.TrimSpace(cmd), name); exit {
			return
		}
	}
}

// UnauthorizedOperationsLoop ...
func UnauthorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		ok, err := handleLogin(db)
		if err != nil {
			log.Printf("can't handle login: %v", err)
		}
		if !ok {
			log.Print("\n\n")
			color.Error.Print("Неправильно введён логин или пароль!")
			color.Bold.Println("\nПопробуйте ещё раз.\n")
			return false
		}
		loginedOperationsLoop(db, "John", authorizedOperationsLoop)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}

// authorized User operations
func authorizedOperationsLoop(db *sql.DB, cmd, userID string) (exit bool) {
	switch cmd {
	case "1":
		prod, err := core.AllUserAccounts(db, userID)
		if err != nil {
			return
		}
		showUsers(prod, false)
	case "2":
		fmt.Println("Перевести деньги:")
		fmt.Println("1. По номеру счета")
		fmt.Println("2. По номеру телефона")
	case "3":
		fmt.Println("Оплатить услугу")
		err := core.UseService("1", db)
		if err != nil {
			log.Println(err)
		}
	case "4":
		fmt.Println("Список банкоматов")
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}

// createUser ...
func createUser(db *sql.DB) (bool, error) {
	var (
		ok       bool
		Name     string
		Account  string
		Phone    string
		Login    string
		Password string
		Balance  string
	)
	color.Green.Println("Введите имя нового пользователя")
	color.Yellow.Print("Имя: ")
	_, err := fmt.Scan(&Name)
	scanHandleError(err)
	color.Green.Println("\nВведите счет нового пользователя")
	color.Yellow.Print("Счет: ")
	_, err = fmt.Scan(&Account)
	scanHandleError(err)
	color.Green.Println("\nВведите баланс нового пользователя")
	color.Yellow.Print("Баланс: ")
	_, err = fmt.Scan(&Balance)
	scanHandleError(err)
	color.Green.Println("\nВведите номер нового пользователя")
	color.Yellow.Print("Номер телефона: ")
	_, err = fmt.Scan(&Phone)
	scanHandleError(err)
	color.Green.Println("\nВведите логин нового пользователя")
	color.Yellow.Print("Логин: ")
	_, err = fmt.Scan(&Login)
	scanHandleError(err)
	color.Green.Println("\nВведите пароль нового пользователя")
	color.Yellow.Print("Пароль: ")
	_, err = fmt.Scan(&Password)
	scanHandleError(err)
	account, err := strconv.Atoi(Account)
	balance, err := strconv.Atoi(Balance)
	phone, err := strconv.Atoi(Phone)
	user := &core.UserType{
		Name:     Name,
		Account:  account,
		Password: Password,
		Phone:    phone,
		Login:    Login,
		Balance:  int64(balance),
	}
	err = core.CreateNewUser(db, user)
	if err != nil {
		log.Println(err)
		return ok, err
	}
	ok = true
	return ok, err
}

// working scanHandlerError ...
func scanHandleError(err error) (bool, error) {
	if err != nil {
		return false, err
	}
	return true, err
}

// ManagerOperationLoop ...
func ManagerOperationLoop(db *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		ok, err := createUser(db)
		if err != nil {
			log.Println(err)
		}
		return ok
	case "2":
		err := core.AddAccountToUser(db)
		if err != nil {
			log.Println(err)
		}
	case "3":
		err := core.CreateService(db)
		if err != nil {
			log.Println(err)
		}
	case "4":
		err := core.CreateATM(db)
		if err != nil {
			log.Println(err)
		}
	case "5":
		chooseFormat(db, exportFile)
	case "6":
		chooseFormat(db, importFile)
	}
	return exit
}

// AllOperations ...
func AllOperations(db *sql.DB, loop func(db *sql.DB, cmd string) bool) {
	color.Warn.Println("\n ====== Добро пожаловать! ======\n")
	magenta := color.FgLightBlue.Render
	for {
		color.Bold.Println("Список доступных операций:\n")
		fmt.Printf("%s%s%s%s%s%s%s\n",
			magenta("1. Добавить ползователя\n"),
			magenta("2. Добавить счет пользователю\n"),
			magenta("3. Добавить услуги\n"),
			magenta("4. Добавить банкомат\n"),
			magenta("5. Экспорт\n"),
			magenta("6. Импорт\n"),
			magenta("q. Выйти из приложения\n"))
		color.Bold.Print("Введите команду:")
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err)
		}
		if exit := loop(db, strings.TrimSpace(cmd)); exit {
			return
		}
	}
}

func chooseFormat(db *sql.DB, loop func(db *sql.DB, cmd string) bool) {
	color.Warn.Println("\nВыберите формат файла: \n")
	cyan := color.FgLightCyan.Render
	for {
		color.Bold.Println("Список доступных форматов:\n")
		fmt.Printf("%s%s%s\n",
			cyan("1. XML \n"),
			cyan("2. JSON \n"),
			cyan("q. Вернутся на главную\n"))
		color.Bold.Print("Введите команду:")
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err)
		}
		if exit := loop(db, strings.TrimSpace(cmd)); exit {
			return
		}
	}
}

func chooseType(db *sql.DB, types, message string, loop func(db *sql.DB, cmd, message string) bool) {
	color.Warn.Println("\nВыберите данные для", types, ": \n")
	cyan := color.FgLightCyan.Render
	for {
		color.Bold.Println("Список доступных данных для", types, ":\n")
		fmt.Printf("%s%s%s%s\n",
			cyan("1. список пользователей \n"),
			cyan("2. список счетов (с пользователями) \n"),
			cyan("3. список банкоматов \n"),
			cyan("q. Вернутся на главную\n"))
		color.Bold.Print("Введите команду:")
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err)
		}

		if exit := loop(db, strings.TrimSpace(cmd), message); exit {
			return
		}
	}
}

// Import & Export
func exportFile(db *sql.DB, command string) bool {
	switch command {
	case "1":
		chooseType(db, "экспорта", "\nГотово! файл экспортирован в JSON\n", chooseData)
	case "2":
		chooseType(db, "экспорта", "\nГотово! файл экспортирован в XML\n", chooseData)
	}
	return true
}

func importFile(db *sql.DB, command string) bool {
	switch command {
	case "1":
		chooseType(db, "импорта", "\nГотово! файл импортирован в JSON\n", chooseData)
	case "2":
		chooseType(db, "импорта", "\nГотово! файл импортирован в XML\n", chooseData)
	}
	return true
}

// prepare
func chooseData(db *sql.DB, cmd, message string) bool {
	query := `SELECT name, account FROM users;`
	switch cmd {
	case "1":
		user, err := core.ExportData(db, query)
		if err != nil {
			log.Println(err)
		}
		showUsers(user, true)
		log.Println("query to database: 1")
		color.Success.Println(message)
		if false {
			color.Error.Println("Извините возникла ошибка, попробуйте позже")
		}
	case "2":
		user, err := core.ExportData(db, query)
		if err != nil {
			log.Println(err)
		}
		showUsers(user, false)
		log.Println("query to database: 2")
		color.Success.Println(message)
		if false {
			color.Error.Println("Извините возникла ошибка, попробуйте позже")
		}
	case "3":
		// exportData(db, query)
		log.Println("query to database: 3")
		color.Success.Println(message)
		if false {
			color.Error.Println("Извините возникла ошибка, попробуйте позже")
		}
	case "q":
		return true
	default:
		fmt.Printf("Вы ввели неверную команду: %s\n", cmd)
	}
	return false
}

func showUsers(datas []core.UserType, flag bool) {
	for _, data := range datas {
		if !flag {
			fmt.Printf("Имя пользователя: %s счёт: %v", data.Name, data.Account)
			continue
		}
		fmt.Printf("Имя пользователя: %s ", data.Name)
	}
}

func printATMs(atms []core.ATM) {
	for _, atm := range atms {
		fmt.Println("Список банкоматов: ", atm)
	}
}

// Working login func
func handleLogin(db *sql.DB) (ok bool, err error) {
	color.Green.Println("Введите Ваш логин и пароль")
	var login string
	color.Yellow.Print("Логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return false, err
	}
	color.Yellow.Print("Пароль: ")
	password, _ := gopass.GetPasswdMasked()

	ok, err = core.Login(login, string(password), db)
	if err != nil {
		return false, err
	}

	return ok, err
}

// NEED FINISHED ...
func handleSale(db *sql.DB) (err error) {
	fmt.Println("Введите ваш логин и пароль")
	var id int64
	fmt.Print("Id of products")
	_, err = fmt.Scan(&id)
	if err != nil {
		return err
	}
	var qty int64
	fmt.Print("quantity:")
	_, err = fmt.Scan(&qty)
	if err != nil {
		return err
	}
	return nil
}
