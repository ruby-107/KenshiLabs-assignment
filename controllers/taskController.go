package controllers

import (
	"context"
	"kenshilabs/database"
	"kenshilabs/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// POST /tasks: Create a new task (requires authentication)
func CreateTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	taskCollection := database.GetCollection(database.DB, "tasks")

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	task.UserID = userID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := taskCollection.InsertOne(ctx, task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
	}

	task.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return c.Status(fiber.StatusCreated).JSON(task)
}

// GET /tasks: Retrieve all tasks for the authenticated user
func GetTasks(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	taskCollection := database.GetCollection(database.DB, "tasks")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := taskCollection.Find(ctx, bson.M{"userid": userID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tasks"})
	}
	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tasks"})
	}

	return c.JSON(tasks)
}

// GET /tasks/:id: Retrieve a task by its ID (requires authentication)
func GetTaskByID(c *fiber.Ctx) error {
	// Get the task ID from the request params
	taskID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	userID := c.Locals("user_id").(string)

	taskCollection := database.GetCollection(database.DB, "tasks")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var task models.Task

	filter := bson.M{"_id": objID, "userid": userID}
	err = taskCollection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	return c.JSON(task)
}

// PUT /tasks/:id: Update a task by its ID (requires authentication)
func UpdateTask(c *fiber.Ctx) error {
	taskID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	userID := c.Locals("user_id").(string)

	taskCollection := database.GetCollection(database.DB, "tasks")

	var updatedTask models.Task
	if err := c.BodyParser(&updatedTask); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	filter := bson.M{"_id": objID, "userid": userID}

	update := bson.M{"$set": updatedTask}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = taskCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
	}

	return c.JSON(fiber.Map{"message": "Task updated"})
}

// DELETE /tasks/:id: Delete a task by its ID (requires authentication)
func DeleteTask(c *fiber.Ctx) error {
	taskID := c.Params("id")

	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task ID"})
	}

	userID := c.Locals("user_id").(string)

	taskCollection := database.GetCollection(database.DB, "tasks")

	filter := bson.M{"_id": objID, "userid": userID}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := taskCollection.DeleteOne(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete task"})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	return c.JSON(fiber.Map{"message": "Task deleted"})
}
