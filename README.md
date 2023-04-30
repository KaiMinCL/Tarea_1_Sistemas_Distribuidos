BD_aerolinea
============================

@autors:
 - ÁLVAREZ DAVID 
 - EVANS JULES
 - NAHUEL GUTIERREZ


The program is located at: `/tarea-1-sistemas-distribuidos/` to navigate to this folder:

` cd tarea-1-sistemas-distribuidos/`

If you type `ls -al` you should be able to see many files and sub-directories. The main ones are main.go and server.go

To resume the project architecture, there are two main programs that need to be executed each sepratly. The first to run is server.go with the command:

 - `go run server.go`
 
This program launches the gin API and sets the end points of the data base to interact with it. 

For the client side, the following program must be run:

 - `go run main.go`
 
This program sends HTTP request to the server to use the database. The interface is textual. To use the program now one must follow the instructions that is printed out by the main.go file.

For this project we made the following assumptions:
 - There are only two flights in one reservation. If more flights need to be booked they must be in a separate reservation.
 - The statistics print out the route of best and lowest earnings, as an example SCL - IQQ and not a plane number.
 - To generate a unique PNR, each time we create a reservation, we put in place a new end point that gets all the reservations and generates a PNR until it isn't in the reservation list.
 - The ancillaries remain the same if the flight date changes.
 - The cost of changing the flight date is only applied once.
 
 
