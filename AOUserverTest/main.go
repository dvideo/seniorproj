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
    //"os"

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

type SlideshowPhoto struct {
    Path string
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
    mail.body = "Dear "+usr+", \n\nYour AllOfUs account was just signed in from an unknown source in ("+ IPfunction()+") using a ("+Device+"). You are getting this email to make sure that this is you if this was you no action is needed. However, if it wasn't you please log in to your account and view your activity in the security section\n\nThank you, AllOfUs Team"

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
        //http.alert("Error")
        //log.Println("it exists already")
        //http.Error(res, "Username already exists. ", 500)
        fmt.Println("Error, username already exists.")
        http.ServeFile(res, req, "signup.html")
        return
    }

    if (rowExists("SELECT Username FROM allofusdbmysql2.UserTable WHERE Email=?",email)) {
        fmt.Println("Error, email already exists.")
        http.ServeFile(res, req, "signup.html")
        //http.Error(res, "Email already exists. ", 500)
        return
    }

    if (email == "") || (username=="")|| (fName=="")|| (lName=="")|| (password=="")|| (confPass=="")|| (bday=="") || (question==""){
        //http.Error(res, "Error, fields can't be blank. ", 500)
        //ok := dialog.Message("%s", "Do you want to continue?").Title("Are you sure?").YesNo()
        fmt.Println("Error, fields can't be blank.")
        http.ServeFile(res, req, "signup.html")
        return
    }

    if (!strings.ContainsAny(password, "123456789")) || (!strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")) || (len(password)<6){
        fmt.Println("Error, passwords must contain a number, a capital letter, and be at least 7 characters long.")
        http.ServeFile(res, req, "signup.html")
        //http.Error(res, "Error, passwords must contain a number, a capital letter, and be at least 7 characters long. ", 500)
        return
    }

    if password != confPass{
        fmt.Println("Error, passwords are not equal, please try again.")
        http.ServeFile(res, req, "signup.html")
        //http.Error(res, "Error, passwords are not equal, please try again. ", 500)
        return
    }

    if question != "Yes" {
        fmt.Println("Error, must by a person")
        http.ServeFile(res, req, "signup.html")
    }

    //if time.Now().Year()-18 >= bday{
    //}

    err := db.QueryRow("SELECT Username FROM allofusdbmysql2.UserTable WHERE Username=?", username).Scan(&user)

    switch {
        /*
    case err != nil || question!="Yes":
        http.Error(res, "Server error, unable to create your account.", 500)
        return
        */
    case err == sql.ErrNoRows:
        //fmt.Println("user = " + username + " password = " + password)
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
        //fmt.Println("hashedPassword = " , hashedPassword , " Err = " , err)
        //fmt.Println("string(password) = " , string(hashedPassword))
        if err != nil {
            http.Error(res, "Server error, unable to create your account 3.", 500)
            return
        }

        _, err = db.Exec("INSERT INTO allofusdbmysql2.UserTable(fName,lName,Username, Password, Email, DateOfBirth) VALUES(?,?, ?, ?, ?, ?)", fName, lName,username, hashedPassword,email,bday)
        db.Query("INSERT INTO allofusdbmysql2.userlocationdevices values(?,?,?,now(),?)",username,IPfunction(),Device,(IPfunction()+Device+username))
        
            
        if err != nil {
            http.Error(res, "Server error, unable to create your account 2.", 500)
            return
        }
        fmt.Println("User created!")
        res.Write([]byte("User created!"))
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
        fmt.Println("Username doesn't exist.")
        http.Redirect(res, req, "/login", 301)
        return
    }
    
    err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))  //crypto/bcrypt: hashedSecret too short to be a bcrypted password
    //fmt.Println(err)
    if err != nil {
        fmt.Println("Password incorrect.")
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
    // templ, err = templ.ParseGlob("templates/*.html")
     templ = template.Must(templ.ParseGlob("templates/*.html"))
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
    http.HandleFunc("/homepageAllofUs",homePage)
    http.ListenAndServe(":8080", nil)
}
/*
func slideshowPhotos(res http.ResponseWriter, req *http.Request){
    var un string
    var photo string
    cookie, _ := req.Cookie("username")  
    un = cookie.Value
    rows, err := db.Query("SELECT Photo FROM allofusdbmysql2.UserPost WHERE Username=?", un) //.Scan(&photo)   
    //fmt.Println(un)   
    if err != nil {
        //http.Redirect(res, req, "/login", 301)
        http.Error(res, "Error", 500)
        return
    }
    //var pic1 string - add 5 for the 5 photos?

    var ps []SlideshowPhoto 
    for rows.Next(){
        err = rows.Scan(&photo)
        if err != nil {
            log.Println(err)
            http.Error(res, " error", http.StatusInternalServerError)
            return
        }
        fmt.Println("testing this thing")
        ps = append(ps,SlideshowPhoto{Path: photo})
        //fmt.Println(ps)
    }
    fmt.Print(templ.ExecuteTemplate(res, "slideshow2.html", ps))

}


func hello(w http.ResponseWriter, r *http.Request){

    //Call to ParseForm makes form fields available.
    err := r.ParseForm()
    if err != nil {
        // Handle error here via logging and then return            
    }

    name := r.PostFormValue("name")
    fmt.Fprintf(w, "Hello, %s!", name)
}
*/


func slideshow(res http.ResponseWriter, req *http.Request) {
    //slideshowPhotos(res,req)
    //req.ParseForm()
    if req.Method != "POST" {
        http.ServeFile(res, req, "slideshow2.html") //is this what i need to the submit button to work ?
        //fmt.Println("something happened")
        return
    }
    
    /*
    var photo string 

    err = db.QueryRow("SELECT Photo FROM allofusdbmysql2.UserPost WHERE Userid=? AND PostID=? ", 4, 83).Scan(&photo)      
    if err != nil {
        //http.Redirect(res, req, "/login", 301)
        http.Error(res, "Server error, unable to select photo.", 500)
        return
    } //send this photo back to html to display it
    fmt.Println(photo)

    type Test struct {
    Path     string
    }

    t, err := template.ParseFiles("slideshow2.gohtml")
    if err != nil {
        panic(err)
    }
    data := Test{Path: "stats.png"} //picture1 instead of photo
    fmt.Println(data)
    err = t.Execute(os.Stdout, data)
    if err != nil {
        panic(err)
    }
    */
    /*
    fmt.Println("was here")
    pic1 := req.FormValue("pic1")
    //picture1 = "static/img/" + picture1
    fmt.Println("here")
    fmt.Println(pic1)
    fmt.Println("there")
    //var picname string

    _, err = db.Exec("INSERT INTO allofusdbmysql2.UserPost (Username,Photo) VALUES (?,?)","dvideo",pic1)         
    if err != nil {
        http.Error(res, "Server error, unable to insert photo.", 500)
        return
    }
    fmt.Println("It works")
    */    
}

        
        
            

func loadSettings(res http.ResponseWriter, req *http.Request){
    cookie, _ := req.Cookie("username")
    cookieUserName := cookie.Value
    var ur []UserPerson
//    var fName string
//    
//    db.QueryRow("SELECT fName FROM allofusdbmysql2.UserTable WHERE Username=?", cookieUserName).Scan(&fName)
    
    ur = append(ur, UserPerson{UserN: cookieUserName})
    fmt.Print(templ.ExecuteTemplate(res, "settings.html", ur))
}

func settings(res http.ResponseWriter, req *http.Request) {
    cookie, _ := req.Cookie("username")
    cookieUserName := cookie.Value
    
    loadSettings(res, req)

    
    if req.Method != "POST" {
       // http.ServeFile(res, req, "templates/settings.html")
        return
    }

    fName := req.FormValue("fName")
    lName := req.FormValue("lName")
    usrN := req.FormValue("UserName")
    
    
    fmt.Println("change name")
    
    if (fName != "" && lName != ""){
        db.QueryRow("UPDATE allofusdbmysql2.UserTable SET fName = ?, lName = ? WHERE Username=?", fName, lName, cookieUserName)
        fmt.Println("change name")
    }
    if usrN != ""{
        db.QueryRow("UPDATE allofusdbmysql2.UserTable SET Username = ? WHERE Username=?", usrN, cookieUserName)
        fmt.Println("change UserName")
    }
    
}

type UserPerson struct {
    UserN string
    First string
    AvgS float32
    numV int
    Photo1 string
    Photo2 string
    Photo3 string
    Photo4 string
    Photo5 string
    Photo6 string
}

func loadUserInfo(res http.ResponseWriter, req *http.Request){
    cookie, _ := req.Cookie("username")
    cookieUserName := cookie.Value
    var ur []UserPerson
//    var avgStat float32
//    var numVotes int
//    var userID int
    
//    db.QueryRow("SELECT Userid FROM allofusdbmysql2.userTable WHERE Username=?",cookieUserName).Scan(&userID)
//    
//    db.QueryRow("SELECT StatAvg, NumVotes FROM allofusdbmysql2.statPost WHERE Userid=?", userID).Scan(&avgStat, &numVotes)
    
    ur = append(ur, UserPerson{UserN: cookieUserName})
    
    templ.ExecuteTemplate(res, "profile.html", ur)
    
    //fmt.Print(templ.ExecuteTemplate(res, "profile.html", ur))
    
}


func displayPost(res http.ResponseWriter, req *http.Request){
    cookie, _ := req.Cookie("username")
    cookieUserName := cookie.Value
    
    var postID int
    var userID int
    var photo string
    var statusUpdate string
    
    db.QueryRow("SELECT Userid FROM allofusdbmysql2.userTable WHERE Username=?",cookieUserName).Scan(&userID)
    fmt.Println(userID)
    
    //loop through userPost and print out each post
    rows, _ := db.Query("SELECT PostID, Photo, Status FROM allofusdbmysql2.userPost WHERE Userid=?", userID)
    for rows.Next() {
        err := rows.Scan(&postID, &photo, &statusUpdate)
        if err != nil {
            log.Fatal(err)
        }
        log.Println(postID, photo, statusUpdate)
        
    }
    
     if photo != "" && statusUpdate != ""{
        //loadData photo and statusUpdate
    } else if photo != ""{
        //loadData photo
    } else if statusUpdate != ""{
        //loadData status Update
    } else{
        fmt.Fprint(res, "<h1>No Posts at this time</h1>")
    }

                     
}

func profile(res http.ResponseWriter, req *http.Request) {
    cookie, _ := req.Cookie("username")
    cookieUserName := cookie.Value
    loadUserInfo(res, req)
    //displayPost(res, req)
    fmt.Println("\nThis is the userName: ",cookieUserName)

    if req.Method != "POST" {
        //http.ServeFile(res, req, "/profile.html")
       // return
   }
    
    //rating := req.FormValue("stat-info")
    

    var userID int
    var postID int

    db.QueryRow("SELECT Userid FROM allofusdbmysql2.userTable WHERE Username=?",cookieUserName).Scan(&userID)
    //db.QueryRow("SELECT StatAvg, NumVotes FROM allofusdbmysql2.statPost WHERE Userid=?", userID).Scan(&avgStat, &numVotes)
    
    fmt.Println("This is the user id: ",userID)
    statusUpdate := req.FormValue("statusUpdate")
    statusPhoto := req.FormValue("picture")
    
//    req.ParseMultipartForm(32 << 20)
//    file, statusPhoto,err := req.FormFile("picture")
//    defer file.Close()

    
    fmt.Println("hi, This is your status: ", statusUpdate)
    
    var jNoPost = `<script>confirm(\'No post was made!\')</script>`
    
    //send error that no data was taken
    if (statusPhoto == "") && (statusUpdate == ""){
        //http.Error(res, "No post was made", 500)
        fmt.Fprint(res, jNoPost)
    }
    
    //Post with a both a photo and a picture
    if (statusPhoto != "") && (statusUpdate != ""){
        if(rowExists("SELECT PostID FROM allofusdbmysql2.userPost WHERE Userid=?",userID)){
            //Get last postID
            //Add 1 and use the postID result
            err := db.QueryRow("SELECT PostID FROM allofusdbmysql2.userPost WHERE Userid=?", userID).Scan(&postID)
            postID =+ 1
            db.Exec("INSET INTO allofusdbmysql2.userPost(PostID, Picture, Status) Values(?,?,?)", postID, statusPhoto, statusUpdate)
            if err != nil{
               // http.Redirect(res, req, "templates/profile.html", 301)
                http.ServeFile(res, req, "templates/profile.html")
            }
            fmt.Println("This is the user postID: ",postID)
        }else{
            //NO post yet, postID equals 1
            postID = 1
            db.Exec("INSET INTO allofusdbmysql2.userPost(PostID, Picture, Status) Values(?,?,?)", postID, statusPhoto, statusUpdate)
            if err != nil{
                //http.Redirect(res, req, "templates/profile.html", 301)
                http.ServeFile(res, req, "templates/profile.html")
            }
        }
        
        fmt.Println("Post with a status & picture: ",postID)
    }

    //Post with a status update done
    if statusUpdate != ""{
            if(rowExists("SELECT PostID FROM allofusdbmysql2.userPost WHERE Userid=?",userID)){
                //Get last postID
                //Add 1 and use the postID result
                err := db.QueryRow("SELECT PostID FROM allofusdbmysql2.userPost WHERE Userid=?", userID).Scan(&postID)
                postID =+ 1
                db.Exec("INSET INTO allofusdbmysql2.userPost(PostID, Status) Values(?,?)", postID, statusUpdate)
                if err != nil{
                    //http.Redirect(res, req, "templates/profile.html", 301)
                    http.ServeFile(res, req, "templates/profile.html")
                    fmt.Println("Post with post ID 1 status: ",postID)
                }
                fmt.Println("This is the user postID 4: ",postID)
            }else{
                //NO post yet, postID equals 1
                postID = 1
                db.Exec("INSET INTO allofusdbmysql2.userPost(PostID, Status) Values(?,?)", postID, statusUpdate)
                if err != nil{
                    //http.Redirect(res, req, "templates/profile.html", 301)
                    http.ServeFile(res, req, "templates/profile.html")
                    fmt.Println("Post with a postID 2: ",postID)
                }
                //http.Redirect(res, req, "templates/profile.html", 301)
                http.ServeFile(res, req, "templates/profile.html")
                fmt.Println("Post with a postID 3: ",postID)
                
            }
            fmt.Println("Post with a status: ",postID)
    }
    
    //Post with a photo
    if statusPhoto != ""{
            if(rowExists("SELECT PostID FROM allofusdbmysql2.userPost WHERE Userid=?",userID)){
                //Get last postID
                //Add 1 and use the postID result
                err := db.QueryRow("SELECT PostID FROM allofusdbmysql2.userPost WHERE Userid=?", userID).Scan(&postID)
                postID =+ 1
                db.Exec("INSET INTO allofusdbmysql2.userPost(PostID, Photo) Values(?,?)", postID, statusPhoto)
                if err != nil{
                    //http.Redirect(res, req, "templates/profile.html", 301)
                    http.ServeFile(res, req, "templates/profile.html")
                }
                fmt.Println("This is the user postID: ",postID)
            }else{
                //NO post yet, postID equals 1
                postID = 1
                    db.Exec("INSET INTO allofusdbmysql2.userPost(PostID, Photo) Values(?,?)", postID, statusPhoto)
                if err != nil{
                   // http.Redirect(res, req, "templates/profile.html", 301)
                    http.ServeFile(res, req, "templates/profile.html")
                }
            }
        fmt.Println("Post with a picture: ",postID)
    }

    

    
//    db.Exec("INSET INTO allofusdbmysql2.statPost(StatAvg, NumVotes) Values(?,?) FROMWHERE Userid=?", userID).Scan(&avgStat, &numVotes)


    
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
        //http.ServeFile(res, req, "/locations.html")
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
