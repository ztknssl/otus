package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type envFile struct {
	envName string
	value   string
}

type envList []envFile

func main() {
	// Определяем комманду "go-envdir"
	sub := flag.NewFlagSet("go-envdir", flag.ContinueOnError)
	err := sub.Parse(os.Args[2:])
	if err != nil {
		fmt.Println(err)
	}
	// Аргументы комманды go-envdir - это путь к папке с файлами и имя комманды, которую нужно выполить
	dirPath := sub.Arg(0)
	cmdName := sub.Arg(1)
	// Читаем папку path
	var list envList
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		f, err := os.Open(file.Name())
		if err != nil {
			fmt.Printf("Can not open a file: %v\n", err)
		}
		defer f.Close()
		// Читаем строку в файле
		b, err := io.ReadAll(f)
		if err != nil {
			fmt.Printf("Can not read a file: %v\n", err)
		}
		list = append(list, envFile{
			envName: file.Name(),
			value:   string(b),
		})
	}
	// Устанавливаем переменные окружения
	for _, item := range list {
		os.Setenv(item.envName, item.value)
	}
	// Опеделяем комманду, которую нужно запустить с переменными окружения
	cmd := exec.Command(cmdName)
	cmd.Start()
	cmd.Wait()
	// Удаляем установленные ранее переменные окружения
	for _, item := range list {
		os.Unsetenv(item.envName)
	}
}
