package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo-cli/analytics"
	"todo-cli/storage"
	"todo-cli/task"
)

func main() {
	store := storage.NewStorage("tasks.json")

	err := store.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки задач:", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n📜 TODo CLI Меню:")
		fmt.Println("1. Показать все задачи")
		fmt.Println("2. Добавить задачу")
		fmt.Println("3. Завершить задачу")
		fmt.Println("4. Удалить задачу")
		fmt.Println("5. Показать статистику")
		fmt.Println("6. Найти задачи по статусу")
		fmt.Println("7. Выход")
		fmt.Println("Выберите действие: ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			showTasks(store.Tasks)
		case "2":
			addTask(scanner, store)
		case "3":
			completeTask(scanner, store)
		case "4":
			deleteTask(scanner, store)
		case "5":
			showAnalytics(store.Tasks)
		case "6":
			filterByStatus(scanner, store)
		case "7":
			fmt.Println("👋 До свидания!")
			return
		default:
			fmt.Println("❌ Неверный выбор")
		}
	}
}

func showTasks(tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Println("📭 Нет задач")
		return
	}

	fmt.Printf("\n📌 Всего задач: %d\n", len(tasks))
	fmt.Println(strings.Repeat("-", 80))

	for _, t := range tasks {
		statusEmoji := map[task.Status]string{
			task.StatusPending:    "⏳",
			task.StatusInProgress: "🔄",
			task.StatusDone:       "✅",
		}

		fmt.Printf("%s [%s] %s (приоритет: %d)\n",
			statusEmoji[t.Status], t.ID[:8], t.Title, t.Priority)

		if len(t.Tags) > 0 {
			fmt.Printf("	🏷️ Теги: %s\n", strings.Join(t.Tags, ", "))
		}
	}
}

func addTask(scanner *bufio.Scanner, store *storage.Storage) {
	fmt.Print("Название задачи: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())

	fmt.Print("Описание: ")
	scanner.Scan()
	desc := strings.TrimSpace(scanner.Text())

	fmt.Print("Приоритет (1-5): ")
	scanner.Scan()
	priority, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || priority < 1 || priority > 5 {
		priority = 3
	}

	fmt.Print("Теги (через запятую): ")
	scanner.Scan()
	tagsStr := strings.TrimSpace(scanner.Text())

	tags := []string{}
	if tagsStr != "" {
		for _, tag := range strings.Split(tagsStr, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	t := task.NewTask(title, desc, priority, tags)

	store.Add(t)

	err = store.Save()
	if err != nil {
		fmt.Println("Ошибка сохранения:", err)
	}

	fmt.Println("✅ Задача добавлена!")
}

func completeTask(scanner *bufio.Scanner, store *storage.Storage) {
	fmt.Print("ID задачи для завершения: ")
	scanner.Scan()
	id := strings.TrimSpace(scanner.Text())

	t := store.FindByID(id)
	if t == nil {
		fmt.Println("❌ Задача не найдена")
		return
	}

	t.Complete()

	err := store.Save()
	if err != nil {
		fmt.Println("Ошибка сохранения:", err)
		return
	}

	fmt.Println("✅ Задача завершена!")
}

func deleteTask(scanner *bufio.Scanner, store *storage.Storage) {
	fmt.Print("ID задачи для удаления: ")
	scanner.Scan()
	id := strings.TrimSpace(scanner.Text())

	if store.Delete(id) {
		err := store.Save()
		if err != nil {
			fmt.Println("Ошибка сохранения:", err)
			return
		}
		fmt.Println("🗑️ Задача удалена!")
	} else {
		fmt.Println("❌ Задача не найдена")
	}
}

func showAnalytics(tasks []task.Task) {
	a := analytics.Calculate(tasks)

	fmt.Printf("\n📊 Статистика: \n")
	fmt.Printf("Всего задач: %d\n", a.TotalTasks)
	fmt.Printf("✅ Завершено: %d\n", a.CompletedTasks)
	fmt.Printf("⏳ В ожидании: %d\n", a.PendingTasks)
	fmt.Printf("🔄 В работе: %d\n", a.InProgress)
	fmt.Printf("⭐ Средний приоритет: %.2f\n", a.AvgPriority)
	fmt.Printf("📈 Прогресс: %.1f%%\n", a.CompletionRate)

	if len(a.TasksByTag) > 0 {
		fmt.Println("\n🏷️ Задачи по тегам:")
		for tag, count := range a.TasksByTag {
			fmt.Printf("  %s: %d\n", tag, count)
		}
	}
}

func filterByStatus(scanner *bufio.Scanner, store *storage.Storage) {
	fmt.Println("Выберите статус: pending, in_progress, done")
	scanner.Scan()
	statusStr := strings.TrimSpace(scanner.Text())

	status := task.Status(statusStr)

	filtered := store.GetByStatus(status)

	if len(filtered) == 0 {
		fmt.Printf("📭 Нет задач со статусом %s\n", status)
		return
	}

	fmt.Printf("\n📌 Задачи со статусом %s: %d\n", status, len(filtered))
	for _, t := range filtered {
		fmt.Printf("- %s (приоритет: %d)\n", t.Title, t.Priority)
	}
}
