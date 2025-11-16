# what-to-watch
Finds a series, TV show, or film to watch

This program provides a menu to view either TV shows you are currently watching or films you want to watch.

**Shows**: Prints out TV shows that a user is currently watching with their provider, current series, and current episode. The program prompts the user to enter the show they watched and updates the data with the new series and episode number.

**Films**: Displays a list of films from your collection with genre and provider information.

The plans directory contains AI generated plans for implementations.

See GitHub Issues for future plans. New Issues and Pull Requests are welcome.

## Usage

### Run

Run `go run .` from the repository to run the program.

You will see a menu:
```
What would you like to view?
1. Currently watching shows
2. Films
Enter your choice (1 or 2):
```

Select option 1 to view and update currently watching shows, or option 2 to view your films collection.

### Build & Run

1. Run `go build`
2. Run `go install`
3. `what-to-watch` should be available to run from any directory
