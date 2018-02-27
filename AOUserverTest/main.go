
package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import "golang.org/x/crypto/bcrypt"

import "net/http"
import "fmt"

var db *sql.DB
var err error

func signupPage(res http.ResponseWriter, req *http.Request) {
    fmt.Println("TESTING signup")
    if req.Method != "POST" {
        http.ServeFile(res, req, "signup.html")
        return
    }

    username := req.FormValue("username")
    password := req.FormValue("password")

    var user string

    err := db.QueryRow("SELECT Username FROM allofusdbmysql2.UserTable WHERE Username=?", username).Scan(&user)

    switch {
    case err == sql.ErrNoRows:
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(res, "Server error, unable to create your account.", 500)
            return
        }

        _, err = db.Exec("INSERT INTO allofusdbmysql2.UserTable(Username, Password) VALUES(?, ?)", username, hashedPassword)
        if err != nil {
            http.Error(res, "Server error, unable to create your account.", 500)
            return
        }
        fmt.Println("User created!")
        res.Write([]byte("User created!")) //not sure if this is working?
        return
    case err != nil:
        http.Error(res, "Server error, unable to create your account.", 500)
        return
    default:
        http.Redirect(res, req, "/", 301)
    }
}

func loginPage(res http.ResponseWriter, req *http.Request) {
    //fmt.Println("TESTING login")
    if req.Method != "POST" {
        http.ServeFile(res, req, "login.html")
        return
    }

    username := req.FormValue("username")
    password := req.FormValue("password")

    var databaseUsername string
    var databasePassword string
    
    err := db.QueryRow("SELECT Username, Password FROM allofusdbmysql2.UserTable WHERE Username=?", username).Scan(&databaseUsername, &databasePassword)
    if err != nil { // see below comment - remove the below if statement to get code to work
        http.Redirect(res, req, "/login", 301)
        return
    }
    //fmt.Println("TESTING"+username)
    fmt.Println("TESTING "+databaseUsername)
    //fmt.Println("TESTING")
    err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))  //for some reason wont work when both of the if err statements are in
    //if err != nil {
    //    http.Redirect(res, req, "/login", 301)
    //    return
    //}
    fmt.Println("Hello " + databaseUsername)
    res.Write([]byte("Hello " + databaseUsername))

}

func homePage(res http.ResponseWriter, req *http.Request) {
    http.ServeFile(res, req, "index.html")
}

func main() {
    db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/allofusdbmysql2") //3306 - johnny //8889 - josh
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("TESTING")
    http.HandleFunc("/signup", signupPage)
    //fmt.Println("TESTING")
    http.HandleFunc("/login", loginPage)
    http.HandleFunc("/", homePage)
    http.ListenAndServe(":8080", nil)
}

func settings(res http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        http.ServeFile(res, req, "settings.html")
        return
    }

    userName := req.FormValue("UserName")
    changeName := req.FormValue("changeName")
    changeUserName := req.FormValue("changeUserName")
    
    deactivate := req.FormValue("deactivate")
    reactiveate := req.FormValue("reactiveate")
    fName := req.FormValue("fName")
    lName := req.FormValue("lName")

    if reactiveate == "reactiveate"{
        err := db.QueryRow("INSET UserName FROM allofusdbmysql2.UserTable WHERE Username=?", userName)
        if err != nil {
            http.Redirect(res, req, "/login", 301)
            return
        }
    }
    if deactivate == "deactivate"{
        err := db.QueryRow("REMOVE UserName FROM allofusdbmysql2.UserTable WHERE Username=?", userName)//.Scan(&databaseUsername, &databasePassword)
            if err != nil {
                http.Redirect(res, req, "/login", 301)
                return
        }
    }
    if changeName == "changeName"{
        err := db.QueryRow("UPDATE allofusdbmysql2.UserTable SET fName, lName WHERE Username=?", lName, fName)//.Scan(&databaseUsername, &databasePassword)
            if err != nil {
                http.Redirect(res, req, "/login", 301)
                return
        }
    }
    if changeUserName == "changeUserName"{
        err := db.QueryRow("UPDATE allofusdbmysql2.UserTable SET Username WHERE Username=?", userName)//.Scan(&databaseUsername, &databasePassword)
            if err != nil {
                http.Redirect(res, req, "/login", 301)
                return
        }
    }
}

func profile(res http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        http.ServeFile(res, req, "profile.html")
        return
    }
    
    username := req.FormValue("userName")
    //rating := req.FormValue("stat-info")

    
    err := db.QueryRow("SELECT Username, Password FROM allofusdbmysql2.UserTable WHERE Username=?", username)//.Scan(&databaseUsername)
    fmt.Println()
    if err != nil {
        http.Redirect(res, req, "/login", 301)
        return
    }
    
    /*err := db.QueryRow("INSERT INTO allofusdbmysql2.stats (PostID, statValue) VALUES (?, ?)", num rating)//.Scan(&databaseUsername)
    fmt.Println()
    if err != nil {
        http.Redirect(res, req, "/login", 301)
        return
    }
    num+1*/

}
/*

package main 

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import "golang.org/x/crypto/bcrypt"

import "net/http"
import "fmt"

// Global sql.DB to access the database by all handlers
var db *sql.DB 
var err error

func homePage(res http.ResponseWriter, req *http.Request) {
    http.ServeFile(res, req, "index.html")
}

func main() {
   // Create an sql.DB and check for errors
    db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/allofusdbmysql2")
    if err != nil {
        panic(err.Error())    
    }
    fmt.Println("HI 1")
    // sql.DB should be long lived "defer" closes it once this function ends
    defer db.Close()

    // Test the connection to the database
    err = db.Ping()
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("HI 2")

    http.HandleFunc("/", homePage)
    http.ListenAndServe(":3306", nil)   

    fmt.Println("HI 3") 
}

func login(res http.ResponseWriter, req *http.Request) {
    // If method is GET serve an html login page
    if req.Method != "POST" {
        http.ServeFile(res, req, "login.html")
        return
    }    

    // Grab the username/password from the submitted post form
    username := req.FormValue("username")
    password := req.FormValue("password")

    // Grab from the database 
    var databaseUsername  string
    var databasePassword  string

    // Search the database for the username provided
    // If it exists grab the password for validation
    err := db.QueryRow("SELECT Username, Password FROM allofusdbmysql2.UserTable WHERE Username=?", username).Scan(&databaseUsername, &databasePassword)
    // If not then redirect to the login page
    if err != nil {
        http.Redirect(res, req, "login.html", 301)
        return
    }
    
    // Validate the password
    err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
    // If wrong password redirect to the login
    if err != nil {
        http.Redirect(res, req, "login.html", 301)
        return
    }

    // If the login succeeded
    res.Write([]byte("Hello " + databaseUsername))
}

func singupPage(res http.ResponseWriter, req *http.Request) {

    // Serve signup.html to get requests to /signup
    if req.Method != "POST" {
        http.ServeFile(res, req, "signup.html")
        return
    }     
    
    username := req.FormValue("username")
    password := req.FormValue("password")

    var user string

    err := db.QueryRow("SELECT Username FROM allofusdbmysql2.UserTable WHERE Username=?", username).Scan(&user)

    switch {
    // Username is available
    case err == sql.ErrNoRows:
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(res, "Server error, unable to create your account.", 500)    
            return
        } 

        _, err = db.Exec("INSERT INTO allofusdbmysql2.UserTable(Username, Password) VALUES(?, ?)", username, hashedPassword)
        if err != nil {
            http.Error(res, "Server error, unable to create your account.", 500)    
            return
        }

        res.Write([]byte("User created!"))
        return
    case err != nil: 
        http.Error(res, "Server error, unable to create your account.", 500)    
        return
    default: 
        http.Redirect(res, req, "/", 301)
    }
}
*/


