package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/emicklei/go-restful"
    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type TrainResource struct {
    ID              int
    DriverName      string
    OperatingStatus bool
}

type StationResource struct {
    ID          int
    Name        string
    OpeningTime time.Time
    ClosingTime time.Time
}

type ScheduleResource struct {
    ID          int
    TrainID     int
    StationID   int
    ArrivalTime time.Time
}

func (t *TrainResource) Register(container *restful.Container) {
    ws := new(restful.WebService)
    ws.Path("/v1/trains").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
    ws.Route(ws.GET("/{train-id}").To(t.getTrain))
    ws.Route(ws.POST("").To(t.createTrain))
    ws.Route(ws.DELETE("/{train-id}").To(t.removeTrain))
    container.Add(ws)
}

func (t *TrainResource) getTrain(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("train-id")
    err := DB.QueryRow("SELECT id, driver_name, operating_status FROM train WHERE id = ?", id).
        Scan(&t.ID, &t.DriverName, &t.OperatingStatus)
    if err != nil {
        if err == sql.ErrNoRows {
            response.WriteErrorString(http.StatusNotFound, "Train could not be found")
        } else {
            response.WriteErrorString(http.StatusInternalServerError, err.Error())
        }
    } else {
        response.WriteEntity(t)
    }
}

func (t *TrainResource) createTrain(request *restful.Request, response *restful.Response) {
    var b TrainResource
    err := json.NewDecoder(request.Request.Body).Decode(&b)
    if err != nil {
        response.WriteErrorString(http.StatusBadRequest, "Invalid request payload")
        return
    }

    statement, err := DB.Prepare("INSERT INTO train(driver_name, operating_status) VALUES(?, ?)")
    if err != nil {
        response.WriteErrorString(http.StatusInternalServerError, err.Error())
        return
    }
    defer statement.Close()

    result, err := statement.Exec(b.DriverName, b.OperatingStatus)
    if err != nil {
        response.WriteErrorString(http.StatusInternalServerError, err.Error())
        return
    }

    newID, _ := result.LastInsertId()
    b.ID = int(newID)
    response.WriteHeaderAndEntity(http.StatusCreated, b)
}

func (t *TrainResource) removeTrain(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("train-id")
    statement, err := DB.Prepare("DELETE FROM train WHERE id = ?")
    if err != nil {
        response.WriteErrorString(http.StatusInternalServerError, err.Error())
        return
    }
    defer statement.Close()

    _, err = statement.Exec(id)
    if err != nil {
        response.WriteErrorString(http.StatusInternalServerError, err.Error())
        return
    }
    
    response.WriteHeader(http.StatusOK)
}

func main() {
    var err error
    DB, err = sql.Open("sqlite3", "./railapi.db")
    if err != nil {
        log.Fatal("Driver creation failed: ", err)
    }
    defer DB.Close()

    if err = DB.Ping(); err != nil {
        log.Fatal("Database connection failed: ", err)
    }
	
    wsContainer := restful.NewContainer()
    wsContainer.Router(restful.CurlyRouter{})
    
    t := &TrainResource{}
    t.Register(wsContainer)
    
    log.Printf("Starting server on localhost:8000")
    server := &http.Server{Addr: ":8000", Handler: wsContainer}
    log.Fatal(server.ListenAndServe())
}