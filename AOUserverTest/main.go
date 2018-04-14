package main

// import "time"
import("crypto/tls"
    "encoding/json"
    "io/ioutil"
    "log"
    "time"
    "net/smtp"
    "strings"
    "github.com/rdegges/go-ipify"
    "github.com/mssola/user_agent"
    "html/template"
    "fmt"
    "net/http"
    "golang.org/x/crypto/bcrypt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"

    ) 

var db *sql.DB
var err error
var sessionUser string
var templ *template.Template
type GeoIP struct {
        // The right side is the name of the JSON variable
    Ip          string  `json:"ip"`
    CountryCode string  `json:"country_code"`
    CountryName string  `json:"country_name""`
    RegionCode  string  `json:"region_code"`
    RegionName  string  `json:"region_name"`
    City        string  `json:"city"`
    Zipcode     string  `json:"zipcode"`
    Lat         float32 `json:"latitude"`
    Lon         float32 `json:"longitude"`
    MetroCode   int     `json:"metro_code"`
    AreaCode    int     `json:"area_code"`
}
type Cookie struct {
        Name       string
        Value      string
        Path       string
        Domain     string
        Expires    time.Time
        RawExpires string

    // MaxAge=0 means no 'Max-Age' attribute specified.
    // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
    // MaxAge>0 means Max-Age attribute present and given in seconds
        MaxAge   int
        Secure   bool
        HttpOnly bool
        Raw      string
        Unparsed []string // Raw text of unparsed attribute-value pairs
    }
    
type devLoc struct {
    Loc string
    Device string
    Date string
}

type Mail struct {
    senderId string
    toIds    []string
    subject  string
    body     string
}
var (
    address  string
    geo      GeoIP
    response *http.Response
    body     []byte
)

type SmtpServer struct {
    host string
    port string
}

func (s *SmtpServer) ServerName() string {
    return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
    message := ""
    message += fmt.Sprintf("From: %s\r\n", mail.senderId)
    if len(mail.toIds) > 0 {
        message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
    }

    message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
    message += "\r\n" + mail.body

    return message
}
func IPfunction() (addr string){
    ip, err := ipify.GetIp()
            if err != nil {
                fmt.Println("Couldn't get my IP address:", err)
            } else {
                //fmt.Println("My IP address is:", ip)
            }
        address = ip
        response, err = http.Get("https://freegeoip.net/json/" + address)
            if err != nil {
                fmt.Println(err)
            }
            defer response.Body.Close()
        // response.Body() is a reader type. We have
            // to use ioutil.ReadAll() to read the data
            // in to a byte slice(string)
            body, err = ioutil.ReadAll(response.Body)
            if err != nil {
                fmt.Println(err)
            }

            // Unmarshal the JSON byte slice to a GeoIP struct
            err = json.Unmarshal(body, &geo)
            if err != nil {
                fmt.Println(err)
            }
        fmt.Println("City:\t",geo.City)
        return geo.City
        
}
func SendMessagemain(usr string, req *http.Request) {
    //fmt.Println(IPfunction())
    //fmt.Println(getMacAddr()) //MAC ADDRESS
    
    var databaseUsername string
    var databaseEmail string
    var Temp int = 2 
    var issue string = "Location and Device" //default location if its a new device change
    err := db.QueryRow("SELECT Username, Email FROM allofusdbmysql2.UserTable WHERE Username=?", usr).Scan(&databaseUsername, &databaseEmail)
    if err != nil {
        fmt.Println("fill", err)
    }
    Device, OpSys, UserBrowser := UserAgentBot(req)
    fmt.Println(OpSys,UserBrowser)
    var key string
    key = IPfunction()+Device+usr
    //db.QueryRow("INSERT INTO allofusdbmysql2.userLocation values (?, ?,?)",usr,IPfunction(),databaselocationkey)// IF NOT EXISTS (SELECT * FROM 
    if(rowExists("SELECT UserInfoID From allofusdbmysql2.userlocationdevices where UserInfoID=?",key)){
        Temp=3;
    }
     
    fmt.Println(issue)
    if(Temp<=2){
    mail := Mail{}
    mail.senderId = "allofusnoreply@gmail.com" //defaul allofus email
    mail.toIds = []string{databaseEmail} //users we are sending alerts to email.
    mail.subject = "Security Alert"
    mail.body = "Dear "+usr+", \n\nYour AllOfUs account was just signed in from an unknown source. You are getting this email to make sure that this is you if this was you no action is needed. However, if it wasn't you please log in to your account and view your activity in the security section\n\nThank you, AllOfUs Team"

    messageBody := mail.BuildMessage()

    smtpServer := SmtpServer{host: "smtp.gmail.com", port: "465"}

    log.Println(smtpServer.host)
    //build an auth                            Password
    auth := smtp.PlainAuth("", mail.senderId, "AllOfUsNoRep", smtpServer.host)

    // Gmail will reject connection if it's not secure
    // TLS config
    tlsconfig := &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         smtpServer.host,
    }

    conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
    if err != nil {
        log.Panic(err)
    }

    client, err := smtp.NewClient(conn, smtpServer.host)
    if err != nil {
        log.Panic(err)
    }

    // step 1: Use Auth
    if err = client.Auth(auth); err != nil {
        log.Panic(err)
    }

    // step 2: add all from and to
    if err = client.Mail(mail.senderId); err != nil {
        log.Panic(err)
    }
    for _, k := range mail.toIds {
        if err = client.Rcpt(k); err != nil {
            log.Panic(err)
        }
    }

    // Data
    w, err := client.Data()
    if err != nil {
        log.Panic(err)
    }

    _, err = w.Write([]byte(messageBody))
    if err != nil {
        log.Panic(err)
    }

    err = w.Close()
    if err != nil {
        log.Panic(err)
    }

    client.Quit()

    log.Println("Mail sent successfully")
    }

}

func signupPage(res http.ResponseWriter, req *http.Request) {
    fmt.Println("TESTING signup")
    if req.Method != "POST" {
        http.ServeFile(res, req, "signup.html")
        return
    }
    Device, OpSys, UserBrowser := UserAgentBot(req)
    fmt.Println(Device,OpSys,UserBrowser)

    email := req.FormValue("email")
    username := req.FormValue("username")
    fName := req.FormValue("firstname")
    lName := req.FormValue("lastname")
    password := req.FormValue("password")
    confPass := req.FormValue("confirmpassword")
    bday := req.FormValue("bday")
    question := req.FormValue("email2")

    var user string
     

    fmt.Println("email = " + email)
    fmt.Println("un = " + username)
    fmt.Println("fName = " + fName)
    fmt.Println("lName = " + lName)
    fmt.Println("pw = " + password)
    fmt.Println("confpw = " + confPass)
    fmt.Println("bday = " + bday)

    if (rowExists("SELECT Email FROM allofusdbmysql2.UserTable WHERE Username=?",username)) {
        http.Error(res, "Username already exists. ", 500)
        return
    }

    if (rowExists("SELECT Username FROM allofusdbmysql2.UserTable WHERE Email=?",email)) {
        http.Error(res, "Email already exists. ", 500)
        return
    }

    if (email == "") || (username=="")|| (fName=="")|| (lName=="")|| (password=="")|| (confPass=="")|| (bday=="") || (question==""){
        http.Error(res, "Error, fields can't be blank. ", 500)
        return

    }

    if (!strings.ContainsAny(password, "123456789")) || (!strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")) || (len(password)<6){
        http.Error(res, "Error, passwords must contain a number, a capital letter, and be at least 7 characters long. ", 500)
        return
    }

    if password != confPass{
        http.Error(res, "Error, passwords are not equal, please try again. ", 500)
        return
    }




    //if time.Now().Year()-18 >= bday{

    //}



    err := db.QueryRow("SELECT Username FROM allofusdbmysql2.UserTable WHERE Username=?", username).Scan(&user)

    switch {
    case err == sql.ErrNoRows:
        //fmt.Println("user = " + username + " password = " + password)
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
        //fmt.Println("hashedPassword = " , hashedPassword , " Err = " , err)
        //fmt.Println("string(password) = " , string(hashedPassword))
        if err != nil {
            http.Error(res, "Server3 error, unable to create your account.", 500)
            return
        }

        _, err = db.Exec("INSERT INTO allofusdbmysql2.UserTable(fName,lName,Username, Password, Email, DateOfBirth) VALUES(?,?, ?, ?, ?, ?)", fName, lName,username, hashedPassword,email,bday)
        db.Query("INSERT INTO allofusdbmysql2.userlocationdevices values(?,?,?,now(),?)",username,IPfunction(),Device,(IPfunction()+Device+username))
        
            
        if err != nil {
            http.Error(res, "S1erver error, unable to create your account.", 500)
            return
        }
        fmt.Println("User created!")
        res.Write([]byte("User created!"))
        return
    case err != nil || question!="Yes":
        http.Error(res, "Server error, unable to create your account.", 500)
        return
    default:
        http.Redirect(res, req, "/", 301)
    }
}
func UserAgentBot(req *http.Request)(string,string,string){
    ua := user_agent.New(req.UserAgent())    
    name, version := ua.Browser()         
    fmt.Printf("%v\n", version) //needs to print 
    return ua.Platform(), ua.OS(), name 
}
func rowExists(query string, args ...interface{}) bool {
    var exists bool
    query = fmt.Sprintf("SELECT exists (%s)", query)
    err := db.QueryRow(query, args...).Scan(&exists)
    if err != nil && err != sql.ErrNoRows {
            fmt.Printf("error checking if row exists '%s' %v", args, err)
    }
    return exists
}
func loginPage(res http.ResponseWriter, req *http.Request) {
    //fmt.Println("TESTING login")
    //seesionHandling()
    //session, _ := store.Get(req, "secretKey")
    
    Device, OpSys, UserBrowser := UserAgentBot(req)
    fmt.Println(Device,OpSys,UserBrowser)
    if req.Method != "POST" {
        http.ServeFile(res, req, "index.html")
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
    
    err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))  //crypto/bcrypt: hashedSecret too short to be a bcrypted password
    fmt.Println(err)
    if err != nil {
        http.Redirect(res, req, "/login", 301)
        return
    }

    //fmt.Println("Hello " + databaseUsername)
    //res.Write([]byte("Hello " + databaseUsername))

    SendMessagemain(databaseUsername,req); //SendMessagemina(databaseUsername);
     seesionHandling(res,req,username)
    db.QueryRow("INSERT INTO allofusdbmysql2.userlocationdevices values(?,?,?,now(),?)",username,IPfunction(),Device,(IPfunction()+Device+username))
        //sessionUser = username
    //session.Values["authenticated"] = true
    //session.Save(req, res)
    
    
    http.ServeFile(res, req, "homepageAllofUs.html")

}

func homePage(res http.ResponseWriter, req *http.Request) {
    http.ServeFile(res, req, "homepageAllofUs.html")
}

func logout(res http.ResponseWriter, req *http.Request) {
    http.ServeFile(res, req, "index.html")
    //http.ServeFile(res, req, "logout.html")
   // session, _ := store.Get(req, "secretKey")
    
   // session.Values["authenticated"] = false
   // session.Save(req, res)
}

func seesionHandling(w http.ResponseWriter, r *http.Request,username string){
    expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie := http.Cookie{Name: "username", Value: username, Expires: expiration}
    http.SetCookie(w, &cookie)
    r.Cookie("username")
    fmt.Println(cookie) //should print our the username
}


func main() {
    templ, err = templ.ParseGlob("templates/*.html")
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
    db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/allofusdbmysql2") //3306 - johnny //8889 - josh //8889 - elijah
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
    http.HandleFunc("/", loginPage)
    http.HandleFunc("/settings", settings)
    http.HandleFunc("/slideshow", slideshow)
    http.HandleFunc("/locations", locations)
    http.HandleFunc("/profile", profile)
    http.ListenAndServe(":8080", nil)
}

func slideshow(res http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {


        //http.FileServer(http.Dir("./SlideshowStuff")))


        //http.Handle("/SlideshowStuff/", http.StripPrefix("/SlideshowStuff/", http.FileServer(http.Dir("./SlideshowStuff"))))
        
        http.ServeFile(res, req, "slideshow2.html")
        return
    }
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

type UserPerson struct {
    UserN string
    First string
}

func loadUserInfo(res http.ResponseWriter, req *http.Request){
    cookie, _ := req.Cookie("username")
    cookieUserName := cookie.Value
    var ur []UserPerson
    
    ur = append(ur, UserPerson{UserN: cookieUserName, First: "AllOfUS"})
    fmt.Print(templ.ExecuteTemplate(res, "profile.html", ur))
}

func profile(res http.ResponseWriter, req *http.Request) {
    cookie, _ := req.Cookie("username")
    cookieUserName := cookie.Value
    
    fmt.Println("\nThis is the userName: ",cookieUserName)
    
    loadUserInfo(res, req)

    
    //rating := req.FormValue("stat-info")
    
    var avgStat float32
    var numVotes int
    var userID int
    var postID int
    

    if req.Method != "POST" {
//        http.ServeFile(res, req, "profile.html")
        return
    }
    
    db.QueryRow("SELECT Userid FROM allofusdbmysql2.userTable WHERE Username=?",cookieUserName).Scan(&userID)
    
    
    statusUpdate := req.FormValue("statusUpdate")
    if(rowExists("SELECT PostID FROM allofusdbmysql2.userPost WHERE Userid=?",userID)){
        //Get last postID
        //Add 1 and use the postID result
        err := db.QueryRow("SELECT PostID FROM allofusdbmysql2.userPost WHERE Userid=?", userID).Scan(&postID)
        postID =+ 1
        db.Exec("INSET INTO allofusdbmysql2.userPost(PostID, Status) Values(?,?)", postID, statusUpdate)
        if err != nil{
            http.Redirect(res, req, "/profile", 301)
            fmt.Println("Error")
        }
        fmt.Println("Error")
    }else{
        //NO post yet, postID equals 1
        postID = 1
        db.Exec("INSET INTO allofusdbmysql2.userPost(PostID, Status) Values(?,?)", postID, statusUpdate)
        if err != nil{
            http.Redirect(res, req, "/profile", 301)
        }
        fmt.Println("Error")
    }
    
    
     db.QueryRow("SELECT StatAvg, NumVotes FROM allofusdbmysql2.statPost WHERE Userid=?", userID).Scan(&avgStat, &numVotes)
//    db.Exec("INSET  StatAvg, NumVotes FROM allofusdbmysql2.statPost WHERE Userid=?", userID).Scan(&avgStat, &numVotes)
//    
//    if err != nil {
//        http.Redirect(res, req, "/profile", 301)
//        return
//    }

    
    /*err := db.QueryRow("INSERT INTO allofusdbmysql2.stats (PostID, statValue) VALUES (?, ?)", num rating)//.Scan(&databaseUsername)
    fmt.Println()
    if err != nil {
        http.Redirect(res, req, "/login", 301)
        return
    }
    num+1*/

}

func loadlocationstable(res http.ResponseWriter, req *http.Request){
    var un string
    cookie, _ := req.Cookie("username")  
    un = "%"+cookie.Value
    //rows, err := db.Query("SELECT username,location FROM allofusdbmysql2.userLocation Where UserInfoID LIKE ?",un)//was supposed to query specific columns tho
    rows, err := db.Query("Select location,device,CreatedDate FROM allofusdbmysql2.userlocationdevices Where UserInfoID LIKE ?",un)
    if err != nil {
    log.Println(err)
    http.Error(res, "there was an error", http.StatusInternalServerError)
    return
    }
    var device string
    var loc string
    var date string
    /*if req.Method != "POST" {
        }*/
    var ps []devLoc
    //loop through the db
    for rows.Next() {
    err = rows.Scan( &loc, &device, &date)
    if err != nil {
        log.Println(err)
        http.Error(res, "there was an error", http.StatusInternalServerError)
        return
    }
    //fmt.Print(append(ps, devLoc{ Username: username, Loc: loc, Dev: dev}))
    ps = append(ps, devLoc{ Loc: loc, Device: device, Date: date})
    }
    fmt.Print(templ.ExecuteTemplate(res, "locations.html", ps))
    
}
func locations(res http.ResponseWriter, req *http.Request) {
    var temploc string
    var tempdev string
    loadlocationstable(res,req)
    if req.Method != "POST" {
        //http.ServeFile(res, req, "locations.html")
        return
    }
    var cookie,err = req.Cookie("location")
        if err == nil {
            var cookievalue = cookie.Value
            fmt.Println(cookievalue)
        }
    temploc = cookie.Value
    cookie, err = req.Cookie("device")
        if err == nil {
            var cookievalue = cookie.Value
            fmt.Println(cookievalue)
        }
   tempdev = cookie.Value
    cookie, err = req.Cookie("username")
    if err == nil {
        var cookievalue = cookie.Value
        fmt.Println(cookievalue)
    }
    fmt.Println(temploc+tempdev+cookie.Value)
    fmt.Println(tempdev+temploc+cookie.Value)
    if _, err := db.Exec("DELETE FROM allofusdbmysql2.userlocationdevices WHERE UserInfoID =?",(tempdev+temploc+cookie.Value))
         err != nil{
        fmt.Println("Request failed.")
        fmt.Println(err)
        return
    }else{
        fmt.Println("Susscessfully Deleted.")
        return
    }

}
