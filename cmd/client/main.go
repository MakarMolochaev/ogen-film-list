package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	filmsapi "ogen-film-list/gen/filmsapi"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()

	client, err := filmsapi.NewClient("http://localhost:8001")
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n=== Film Management System ===")
		fmt.Println("1. List all films")
		fmt.Println("2. Get film by ID")
		fmt.Println("3. Create new film")
		fmt.Println("4. Update film")
		fmt.Println("5. Delete film")
		fmt.Println("6. Exit")
		fmt.Print("Choose an option: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			listFilms(ctx, client)
		case "2":
			getFilm(ctx, client, scanner)
		case "3":
			createFilm(ctx, client, scanner)
		case "4":
			updateFilm(ctx, client, scanner)
		case "5":
			deleteFilm(ctx, client, scanner)
		case "6":
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func listFilms(ctx context.Context, client *filmsapi.Client) {
	fmt.Println("\n--- Listing all films ---")

	resp, err := client.ListFilms(ctx, filmsapi.ListFilmsParams{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch resp := resp.(type) {
	case *filmsapi.ListFilmsOKApplicationJSON:
		films := []filmsapi.Film(*resp)
		if len(films) == 0 {
			fmt.Println("No films found.")
			return
		}
		fmt.Printf("Found %d films:\n", len(films))
		for i, film := range films {
			fmt.Printf("%d. %s (%d) - %s\n", i+1, film.Title, film.Year, film.Director)
			fmt.Printf("   ID: %s, Rating: %.1f, Duration: %d min\n",
				film.ID.String(), film.Rating, film.Duration)
			fmt.Printf("   Country: %s, Age Rating: %s\n", film.Country, film.AgeRating)
			if len(film.Actors) > 0 {
				fmt.Printf("   Actors: %s\n", strings.Join(film.Actors, ", "))
			}
			fmt.Println()
		}
	default:
		fmt.Printf("Unexpected response: %T\n", resp)
	}
}

func getFilm(ctx context.Context, client *filmsapi.Client, scanner *bufio.Scanner) {
	fmt.Print("\nEnter film ID: ")
	scanner.Scan()
	idStr := strings.TrimSpace(scanner.Text())

	id, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Printf("Invalid ID format: %v\n", err)
		return
	}

	fmt.Printf("\n--- Getting film %s ---\n", idStr)

	resp, err := client.GetFilm(ctx, filmsapi.GetFilmParams{ID: id})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch resp := resp.(type) {
	case *filmsapi.Film:
		printFilmDetails(resp)
	case *filmsapi.GetFilmNotFound:
		fmt.Println("Film not found")
	default:
		fmt.Printf("Unexpected response: %T\n", resp)
	}
}

func createFilm(ctx context.Context, client *filmsapi.Client, scanner *bufio.Scanner) {
	fmt.Println("\n--- Creating new film ---")

	req := &filmsapi.CreateFilmRequest{}

	fmt.Print("Title: ")
	scanner.Scan()
	req.Title = strings.TrimSpace(scanner.Text())

	fmt.Print("Year: ")
	scanner.Scan()
	year, _ := strconv.Atoi(scanner.Text())
	req.Year = year

	fmt.Print("Country: ")
	scanner.Scan()
	req.Country = strings.TrimSpace(scanner.Text())

	fmt.Print("Director: ")
	scanner.Scan()
	req.Director = strings.TrimSpace(scanner.Text())

	fmt.Print("Duration (minutes): ")
	scanner.Scan()
	duration, _ := strconv.Atoi(scanner.Text())
	req.Duration = duration

	fmt.Println("Age Rating: 0+, 6+, 12+, 16+, 18+")
	fmt.Print("Choose age rating: ")
	scanner.Scan()
	ageRating := strings.TrimSpace(scanner.Text())
	req.AgeRating = filmsapi.CreateFilmRequestAgeRating(ageRating)

	resp, err := client.CreateFilm(ctx, req)
	if err != nil {
		fmt.Printf("Error creating film: %v\n", err)
		return
	}

	switch resp := resp.(type) {
	case *filmsapi.Film:
		fmt.Printf("\nFilm created successfully!\n")
		printFilmDetails(resp)
	default:
		fmt.Printf("Unexpected response: %T\n", resp)
	}
}

func updateFilm(ctx context.Context, client *filmsapi.Client, scanner *bufio.Scanner) {
	fmt.Print("\nEnter film ID to update: ")
	scanner.Scan()
	idStr := strings.TrimSpace(scanner.Text())

	id, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Printf("Invalid ID format: %v\n", err)
		return
	}

	fmt.Printf("\n--- Getting current film data ---\n")
	currentResp, err := client.GetFilm(ctx, filmsapi.GetFilmParams{ID: id})
	if err != nil {
		fmt.Printf("Error getting film: %v\n", err)
		return
	}

	var currentFilm *filmsapi.Film
	switch resp := currentResp.(type) {
	case *filmsapi.Film:
		currentFilm = resp
		fmt.Println("Current film data:")
		printFilmDetails(currentFilm)
	case *filmsapi.GetFilmNotFound:
		fmt.Println("Film not found")
		return
	default:
		fmt.Printf("Unexpected response: %T\n", resp)
		return
	}

	fmt.Println("\n--- Updating film (leave blank to keep current value) ---")

	req := &filmsapi.UpdateFilmRequest{}

	fmt.Printf("Title [%s]: ", currentFilm.Title)
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())
	if title != "" {
		req.Title = title
	} else {
		req.Title = currentFilm.Title
	}

	fmt.Printf("Year [%d]: ", currentFilm.Year)
	scanner.Scan()
	yearStr := strings.TrimSpace(scanner.Text())
	if yearStr != "" {
		year, _ := strconv.Atoi(yearStr)
		req.Year = year
	} else {
		req.Year = currentFilm.Year
	}

	fmt.Printf("Country [%s]: ", currentFilm.Country)
	scanner.Scan()
	country := strings.TrimSpace(scanner.Text())
	if country != "" {
		req.Country = country
	} else {
		req.Country = currentFilm.Country
	}

	fmt.Printf("Director [%s]: ", currentFilm.Director)
	scanner.Scan()
	director := strings.TrimSpace(scanner.Text())
	if director != "" {
		req.Director = director
	} else {
		req.Director = currentFilm.Director
	}

	fmt.Printf("Rating [%.1f]: ", currentFilm.Rating)
	scanner.Scan()
	ratingStr := strings.TrimSpace(scanner.Text())
	if ratingStr != "" {
		rating, _ := strconv.ParseFloat(ratingStr, 32)
		req.Rating = float32(rating)
	} else {
		req.Rating = currentFilm.Rating
	}

	fmt.Printf("Duration [%d]: ", currentFilm.Duration)
	scanner.Scan()
	durationStr := strings.TrimSpace(scanner.Text())
	if durationStr != "" {
		duration, _ := strconv.Atoi(durationStr)
		req.Duration = duration
	} else {
		req.Duration = currentFilm.Duration
	}

	fmt.Printf("Age Rating [%s]: ", currentFilm.AgeRating)
	scanner.Scan()
	ageRating := strings.TrimSpace(scanner.Text())
	if ageRating != "" {
		req.AgeRating = filmsapi.UpdateFilmRequestAgeRating(ageRating)
	} else {
		req.AgeRating = filmsapi.UpdateFilmRequestAgeRating(currentFilm.AgeRating)
	}

	fmt.Printf("Actors [%s]: ", strings.Join(currentFilm.Actors, ", "))
	scanner.Scan()
	actorsStr := strings.TrimSpace(scanner.Text())
	if actorsStr != "" {
		req.Actors = strings.Split(actorsStr, ",")
		for i := range req.Actors {
			req.Actors[i] = strings.TrimSpace(req.Actors[i])
		}
	} else {
		req.Actors = currentFilm.Actors
	}

	resp, err := client.UpdateFilm(ctx, req, filmsapi.UpdateFilmParams{ID: id})
	if err != nil {
		fmt.Printf("Error updating film: %v\n", err)
		return
	}

	switch resp := resp.(type) {
	case *filmsapi.Film:
		fmt.Printf("\nFilm updated successfully!\n")
		printFilmDetails(resp)
	default:
		fmt.Printf("Unexpected response: %T\n", resp)
	}
}

func deleteFilm(ctx context.Context, client *filmsapi.Client, scanner *bufio.Scanner) {
	fmt.Print("\nEnter film ID to delete: ")
	scanner.Scan()
	idStr := strings.TrimSpace(scanner.Text())

	id, err := uuid.Parse(idStr)
	if err != nil {
		fmt.Printf("Invalid ID format: %v\n", err)
		return
	}

	fmt.Printf("Are you sure you want to delete film %s? (y/N): ", idStr)
	scanner.Scan()
	confirm := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}

	fmt.Printf("\n--- Deleting film %s ---\n", idStr)

	resp, err := client.DeleteFilm(ctx, filmsapi.DeleteFilmParams{ID: id})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch resp.(type) {
	case *filmsapi.DeleteFilmNoContent:
		fmt.Println("Film deleted successfully")
	case *filmsapi.DeleteFilmNotFound:
		fmt.Println("Film not found")
	default:
		fmt.Printf("Unexpected response: %T\n", resp)
	}
}

func printFilmDetails(film *filmsapi.Film) {
	fmt.Printf("ID: %s\n", film.ID.String())
	fmt.Printf("Title: %s\n", film.Title)
	fmt.Printf("Year: %d\n", film.Year)
	fmt.Printf("Country: %s\n", film.Country)
	fmt.Printf("Director: %s\n", film.Director)
	fmt.Printf("Rating: %.1f/10\n", film.Rating)
	fmt.Printf("Duration: %d minutes\n", film.Duration)
	fmt.Printf("Age Rating: %s\n", film.AgeRating)
	if len(film.Actors) > 0 {
		fmt.Printf("Actors: %s\n", strings.Join(film.Actors, ", "))
	} else {
		fmt.Println("Actors: None")
	}
}
