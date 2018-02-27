func settings(res http.ResponseWriter, req *http.Request) {
    fmt.Println("TESTING")
    if req.Method != "POST" {
        http.ServeFile(res, req, "settings.html")
        return
    }

    userName := req.FormValue("UserName")
    changeName := req.FormValue("changeName")
    changeUserName := req.FormValue("changeUserName")
    
    deactivate := req.FormValue("deactivate")
    reactiveate := req.FormValue("reactiveate")

   // var databaseUsername string
   // var databasePassword string
    if reactiveate != nil{
        err := db.QueryRow("INSET UserName FROM allofusdbmysql2.UserTable WHERE Username=?", userName)
        fmt.Println()
        if err != nil {
            http.Redirect(res, req, "/login", 301)
            return
        }
    }
    if deactivate != nil{
        err := db.QueryRow("REMOVE UserName FROM allofusdbmysql2.UserTable WHERE Username=?", userName)//.Scan(&databaseUsername, &databasePassword)
            fmt.Println()
            if err != nil {
                http.Redirect(res, req, "/login", 301)
                return
        }
    }
    if changeName != nil{
        err := db.QueryRow("UPDATE allofusdbmysql2.UserTable SET fName, lName WHERE Username=?", lName, fName)//.Scan(&databaseUsername, &databasePassword)
            fmt.Println()
            if err != nil {
                http.Redirect(res, req, "/login", 301)
                return
        }
    }
    if changeUserName != nil{
        err := db.QueryRow("UPDATE allofusdbmysql2.UserTable SET Username WHERE Username=?", userName)//.Scan(&databaseUsername, &databasePassword)
            fmt.Println()
            if err != nil {
                http.Redirect(res, req, "/login", 301)
                return
        }
    }

}

func profile(res http.ResponseWriter, req *http.Request) {
    fmt.Println("TESTING")
    if req.Method != "POST" {
        http.ServeFile(res, req, "profile.html")
        return
    }

    username := req.FormValue("username")

    var databaseUsername string
    
    err := db.QueryRow("SELECT Username, Password FROM allofusdbmysql2.UserTable WHERE Username=?", username).Scan(&databaseUsername, &databasePassword)
    fmt.Println()
    if err != nil {
        http.Redirect(res, req, "/login", 301)
        return
    }
    fmt.Println("TESTING"+username+"Sdfsv")
    fmt.Println("TESTING"+databaseUsername+"Sdfsv")
    fmt.Println("TESTING")
    err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
    if err != nil {
        http.Redirect(res, req, "/login", 301)
        
        return
    }
    fmt.Println("TESTING"+databaseUsername)
    res.Write([]byte("Hello" + databaseUsername))

}
