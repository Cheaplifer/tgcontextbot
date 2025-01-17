package storage

import (
	"context"
	db "github.com/jackc/pgx/v4"
	"log"
	"strings"
)

//TODO - rewrite all functions that work with databases into (error, string) format
func CheckIfPresentInChats(idd int64) bool {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")
	if err != nil {
		log.Println(err)
		return false
	}

	defer conn.Close(context.Background())

	var chat_id int64

	err = conn.QueryRow(context.Background(), "select autoinc_id from added_to_chats where id = $1", idd).Scan(&chat_id)

	if err != nil {
		log.Println(err)
		return false
	}

	return true

}

func AddChatIDToDatabase(idd int64) error {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")
	if err != nil {
		return err
	}

	defer conn.Close(context.Background())

	_, Err := conn.Exec(context.Background(), "insert into added_to_chats (id, badword) values ($1, 'placeholder')", idd)
	if Err != nil {
		return Err
	}

	return nil
}

func AddWordToBlacklist(idd int64, badWord string) error {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
		return err
	}

	defer conn.Close(context.Background())

	var allBadWords []string = strings.Split(badWord, " ")

	var badWordByChat []string = nil

	newRows, Err := conn.Query(context.Background(), "select badword from added_to_chats where id = $1", idd)

	if Err != nil {
		return Err
	}

	if newRows == nil {
		return nil
	}

	for newRows.Next() {
		var txt string
		errWithParse := newRows.Scan(&txt)
		if errWithParse != nil {
			return errWithParse
		}
		badWordByChat = append(badWordByChat, txt)
	}

	var flag bool = false
	var i int = 0
	var j int = 0
	for i = 0; i < len(allBadWords); i++ {
		flag = false
		for j = 0; j < len(badWordByChat); j++ {
			if allBadWords[i] == badWordByChat[j] {
				flag = true
			}
		}

		if !flag {
			_, err = conn.Exec(context.Background(), "insert into added_to_chats (id, badword) values ($1, $2)", idd, allBadWords[i])

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetAllBadWordsByChat(idd int64) (*[]string, error) {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
		return nil, err
	}

	defer conn.Close(context.Background())

	newRows, Err := conn.Query(context.Background(), "select badword from added_to_chats where id = $1", idd)

	if Err != nil {
		return nil, Err
	}

	if newRows == nil {
		return nil, nil
	}

	var allWords []string = nil

	for newRows.Next() {
		var txt string

		errWithParse := newRows.Scan(&txt)

		if errWithParse != nil {
			return nil, errWithParse
		}

		allWords = append(allWords, txt)

	}

	return &allWords, nil

}

func DeleteWordFromBlacklist(idd int64, badWord string) error {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
		return err
	}

	defer conn.Close(context.Background())

	var allBadWords []string = strings.Split(badWord, " ")

	var badWordByChat []string = nil

	newRows, Err := conn.Query(context.Background(), "select badword from added_to_chats where id = $1", idd)

	if Err != nil {
		return Err
	}

	if newRows == nil {
		return nil
	}

	for newRows.Next() {
		var txt string
		errWithParse := newRows.Scan(&txt)
		if errWithParse != nil {
			return errWithParse
		}
		badWordByChat = append(badWordByChat, txt)
	}

	var flag = false
	var i = 0
	var j = 0
	for i = 0; i < len(allBadWords); i++ {
		flag = true
		for j = 0; j < len(badWordByChat); j++ {
			if allBadWords[i] == badWordByChat[j] {
				flag = false
			}
		}

		if !flag {
			_, err = conn.Exec(context.Background(), "delete from added_to_chats where id = $1 and badword = $2", idd, allBadWords[i])

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func AddWordToID(keystring string, idd int64) (error, string) {
	var s string
	s = strings.Trim(keystring, "/setsubstitutewith ")
	s = strings.TrimSpace(s)
	var allString []string

	allString = strings.Split(s, "||")

	if len(allString) != 3 {
		return nil, "В команде не два слова, или они разделены неправильными символами(не ||). Пример - /setsubstitutewith aboba||amongus||"
	}
	var first string
	var second string
	first = allString[0]
	second = allString[1]
	var checker string
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
		return err, "Ошибка при соединении с базой данных"
	}
	defer conn.Close(context.Background())

	newRows, Err := conn.Query(context.Background(), "select replace_phrase from chat_phrases where chat_id = $1 and find_phrase = $2", idd, first)

	if Err != nil {
		return Err, "Мы не сумели подключиться к базе данных. Повторите запрос через какое-то время."
	}
	var cnt int64 = 0
	for newRows.Next() {
		cnt++
		ErrWithParse := newRows.Scan(&checker)
		if ErrWithParse != nil {
			return ErrWithParse, "Ошибка при работе с базой данных."
		}
	}
	if cnt != 0 {
		return nil, "Уже есть заменитель на эту фразу"
	}
	_, Err = conn.Exec(context.Background(), "insert into chat_phrases (chat_id, find_phrase, replace_phrase) values ($1, $2, $3)", idd, first, second)

	if Err != nil {
		return Err, "Мы не сумели подключиться к базе данных. Повторите запрос через какое-то время."
	}

	return nil, "Набор фраз добавлен"
}

func GetAllPairsFromChat(idd int64) (*[]string, *[]string, error, string) {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
		return nil, nil, err, "Мы не сумели установить соединение с базой данных"
	}
	defer conn.Close(context.Background())
	newRows, err := conn.Query(context.Background(), "select find_phrase, replace_phrase from chat_phrases where chat_id = $1", idd)
	var firstPair []string = nil
	var secondPair []string = nil
	var firstWordInPair string
	var secondWordInPair string

	if newRows == nil {
		return nil, nil, nil, "В данном чате нету замен"
	}
	for newRows.Next() {
		ErrWithParse := newRows.Scan(&firstWordInPair, &secondWordInPair)

		if ErrWithParse != nil {
			return nil, nil, ErrWithParse, "Ошибка при работе с базой данных"
		}

		firstPair = append(firstPair, firstWordInPair)

		secondPair = append(secondPair, secondWordInPair)
	}

	return &firstPair, &secondPair, nil, "Все успешно."
}

func DeleteWordFromChat(idd int64, key string) (error, string) {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")
	var s string = strings.TrimLeft(key, "/deletesubstitute")
	s = strings.TrimSpace(s)
	if err != nil {
		return err, "Мы не сумели установить соединение с базой данных"
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "delete from chat_phrases where chat_id = $1 and find_phrase = $2", idd, s)

	if err != nil {
		return err, "Что-то пошло не так. Попробуйте позже."
	} else {
		return nil, "Все успешно."
	}
}
