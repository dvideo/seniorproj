
package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import "golang.org/x/crypto/bcrypt"

import "net/http"
import "fmt"
import("crypto/tls"
    "encoding/json"
    "io/ioutil"
    "net"
    "bytes"
    "log"
    "net/smtp"
    "strings"
    "github.com/rdegges/go-ipify"
    "github.com/mssola/user_agent"
    ) 

var db *sql.DB
var err error
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
func getMacAddr() (addr string) {
    interfaces, err := net.Interfaces()
    if err == nil {
        for _, i := range interfaces {
            if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
                // Don't use random as we have a real address
                addr = i.HardwareAddr.String()
                break
            }
        }
    }
    return
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
    var databaseLocation string
    var databaseDevice string
    var databaseEmail string
    var databaselocationkey string
    var databasedevicekey string
    var Temp int = 2 
    var issue string = "Location and Device" //default location if its a new device change
    err := db.QueryRow("SELECT Username, Location, Device,Email FROM allofusdbmysql2.UserTable WHERE Username=?", usr).Scan(&databaseUsername, &databaseLocation, &databaseDevice, &databaseEmail)
    if err != nil {
        fmt.Println("fill", err)
    }
    fmt.Println(databaseLocation)
    fmt.Println(IPfunction())
    databaselocationkey= (IPfunction()+usr)
    Device, OpSys, UserBrowser := UserAgentBot(req)
    fmt.Println(OpSys,UserBrowser)
    databasedevicekey= (Device+usr)
    //db.QueryRow("INSERT INTO allofusdbmysql2.userLocation values (?, ?,?)",usr,IPfunction(),databaselocationkey)// IF NOT EXISTS (SELECT * FROM 
    if (rowExists("SELECT UserInfoID FROM allofusdbmysql2.userLocation WHERE UserInfoID=?",databaselocationkey)) {
        Temp=Temp-1
        issue ="Device"
    }
        
    //db.QueryRow("INSERT INTO allofusdbmysql2.userdevice values (?, ?,?)",usr,Device,databasedevicekey)
    if rowExists("SELECT UserInfoID FROM allofusdbmysql2.userdevice WHERE UserInfoID=?",databasedevicekey) {
        issue = "Location"
        Temp=Temp-1
    }
   if (rowExists("SELECT UserInfoID FROM allofusdbmysql2.userdevice WHERE UserInfoID=?",databasedevicekey)&&(rowExists("SELECT UserInfoID FROM allofusdbmysql2.userLocation WHERE UserInfoID=?",databaselocationkey))){
        issue = "Location and Device "
        Temp=3
}
    if(Temp<=2){
    mail := Mail{}
    mail.senderId = "allofusnoreply@gmail.com" //defaul allofus email
    mail.toIds = []string{databaseEmail} //users we are sending alerts to email.
    mail.subject = "Security Alert"
    mail.body = "Dear "+usr+", \n\nYour AllOfUs account was just signed in from a new "+issue+". You are getting this email to make sure that this is you if this was you no action is needed. However, if it wasn't you please log in to your account and view your activity in the security section\n\nThank you, AllOfUs Team"

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
    username := req.FormValue("username")
    password := req.FormValue("password")
    question := req.FormValue("email2")
    var user string

    err := db.QueryRow("SELECT Username FROM allofusdbmysql2.UserTable WHERE Username=?", username).Scan(&user)

    switch {
    case err == sql.ErrNoRows:
        fmt.Println("user = " + username + " password = " + password)
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
        fmt.Println("hashedPassword = " , hashedPassword , " Err = " , err)
        fmt.Println("string(password) = " , string(hashedPassword))
        if err != nil {
            http.Error(res, "Server3 error, unable to create your account.", 500)
            return
        }

        _, err = db.Exec("INSERT INTO allofusdbmysql2.UserTable(Username, Password) VALUES(?, ?)", username, hashedPassword)

        var databaselocationkey string
        databaselocationkey = IPfunction()+ username
        db.QueryRow("INSERT INTO allofusdbmysql2.userLocation values (?, ?,?)",username,IPfunction(),databaselocationkey)
        var databasedevicekey string
        databasedevicekey = username+Device
        db.QueryRow("INSERT INTO allofusdbmysql2.userdevice values (?, ?,?)",username,Device,databasedevicekey)
            
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
    Device, OpSys, UserBrowser := UserAgentBot(req)
    fmt.Println(Device,OpSys,UserBrowser)
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
    var databaselocationkey string
    databaselocationkey = IPfunction()+ username
    db.QueryRow("INSERT INTO allofusdbmysql2.userLocation values (?, ?,?)",username,IPfunction(),databaselocationkey)
    var databasedevicekey string
    databasedevicekey = username+Device
    db.QueryRow("INSERT INTO allofusdbmysql2.userdevice values (?, ?,?)",username,Device,databasedevicekey)
    //fmt.Println("TESTING " + databaseUsername)
    //fmt.Println("database password = " + databasePassword + " password =  " + password)
    //fmt.Println("string(pw) = ", string(databasePassword))
    //fmt.Println("byte(pw) = ", []byte(databasePassword))
    err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))  //crypto/bcrypt: hashedSecret too short to be a bcrypted password
    fmt.Println(err)
    if err != nil {
        http.Redirect(res, req, "/login", 301)
        return
    }

    //fmt.Println("Hello " + databaseUsername)
    //res.Write([]byte("Hello " + databaseUsername))

    SendMessagemain(databaseUsername,req); //SendMessagemina(databaseUsername);
    http.ServeFile(res, req, "homepageAllofUs.html")

}

func homePage(res http.ResponseWriter, req *http.Request) {
    http.ServeFile(res, req, "index.html")
}

func main() {
    db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/allofusdbmysql2") //3306 - johnny //8889 - josh //8889 - elijah
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()

    
    /*http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
        cookie, err := req.Cookie("cookie1")
        //cookie is not set 
        if err != nil{
            //id, _ := uuid.NewV4()
            cookie = &http.Cookie{
                Name: "ssession-ID",
               // Value: id.String(),
            }
        }
        if req.FormValue("username") != ""{
            cookie.Value = req.FormValue("username")
        }
        
        http.SetCookie(res, cookie)
    })*/
    

    err = db.Ping()
    if err != nil {
        panic(err.Error())
    }
    fmt.Println("TESTING")
    http.HandleFunc("/signup", signupPage)
    //fmt.Println("TESTING")
    http.HandleFunc("/login", loginPage)
    http.HandleFunc("/", homePage)
     http.HandleFunc("/settings", settings)
    http.HandleFunc("/locations", locations)
    http.HandleFunc("/profile", profile)
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
    
    //rating := req.FormValue("stat-info")
    
    var avgStat float32
    var fName string
    var lName string
    username := req.FormValue("userName")
    
    db.QueryRow("SELECT fName, lName FROM allofusdbmysql2.UserTable WHERE Username=?", username).Scan(&fName, &lName)
    if err != nil {
        http.Redirect(res, req, "/login", 301)
        return
    }
    
     err := db.QueryRow("SELECT StatAvg FROM allofusdbmysql2.statPost WHERE Username=?", username).Scan(&avgStat)
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

func locations(res http.ResponseWriter, req *http.Request) {
    
    rows, err := db.Query("SELECT * FROM allofusdbmysql2.userLocation")
    if err != nil {
        http.Redirect(res, req, "/locations", 301)
        return
    }
        
//    var Location string
    
    for rows.Next(){
        var Location string
        var userID int
        err = rows.Scan(&userID, &Location)
        if err != nil {
            http.Redirect(res, req, "/locations", 301)
            return
        }      
    }
    
    
    
    if req.Method != "POST" {
        http.ServeFile(res, req, "locations.html")
        return
    }
    
    if _, err := db.Exec("DROP TABLE allofusdbmysql2.userLocation"); err != nil{
    
    //if err != nil {
        //http.Redirect(res, req, "/login", 301)
        fmt.Println("Request failed.")
        return
    }else{
        fmt.Println("Susscessfully Deleted.")
        return
    }

}
