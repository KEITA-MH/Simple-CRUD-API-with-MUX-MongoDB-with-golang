package main

import (

    "fmt"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
    "context"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
   /* jwt "github.com/dgrijalva/jwt-go"
    "time"*/
)

var collection *mongo.Collection
// var mySecretKey = []byte("supersecretkeytokensigned")

func init(){
    
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    
    if err != nil {
        log.Fatal(err)
    }
    
    err = client.Ping(context.TODO(), nil)
    
    if err != nil {
        log.Fatal(err)
    }
    
    collection = client.Database("sugu").Collection("users")
}

type User struct{
    ID          primitive.ObjectID  `json:"_id" bson:"_id"`
    Fullname   string              `json:"fullname" bson:"fullname"`
    Email       string              `json:"email" bson:"email"`
    Password    string              `json:"password" bson:"password"`
    Image       string              `json:"image" bson:"image"`
    Sexe        string              `json:"sexe" bson:"sexe"`
}


/*func generateJWT(user User) (string, error){
    
    token := jwt.New(jwt.SinningMethodHS256)
    
    claims := jwt.Claims.(jwt.MapClaims)
    
    claims["authorized"] = true
    claims["user"] = user.ID
    claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
    
    tokenString, err := jwt.SignedString(mySecretKey)
    
    if err != nil {
        log.Fatal(err)
        return "", nil
    }
    
    return tokenString, nil
}*/

func getAll(res http.ResponseWriter, req *http.Request){
    
    res.Header().Set("content-type", "application/json")
    
    var person []User
    cur, err := collection.Find(context.TODO(), bson.D{{}})
    
    if err != nil {
        res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
    }
    defer cur.Close(context.TODO())
    
    for cur.Next(context.TODO()){
        var user User
        err := cur.Decode(&user)
        if err != nil{
            log.Fatal(err)
        }
        
        person = append(person, user)
    }
    
    if err := cur.Err(); err != nil {
        res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
    }
    
    json.NewEncoder(res).Encode(person)
}

func insertUser(res http.ResponseWriter, req *http.Request) {
    res.Header().Set("content-type", "application/json")
    var user User
    _= json.NewDecoder(req.Body).Decode(&user)
    insertResult, err := collection.InsertOne(context.TODO(), bson.M{
        "fullname": user.Fullname,
        "email": user.Email,
        "password": user.Password,
        "image": user.Image,
        "sexe": user.Sexe,
    })
    if err != nil {
        res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
    }
    json.NewEncoder(res).Encode(insertResult)
}
func deleteUsers(res http.ResponseWriter, req *http.Request){
    deleteResult,err := collection.DeleteMany(context.TODO(), bson.D{{}})
    if err != nil {
        res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
    }
    json.NewEncoder(res).Encode(deleteResult)
}

func getUser(res http.ResponseWriter, req *http.Request){
    var user User
    params := mux.Vars(req)
    id, _ :=    primitive.ObjectIDFromHex(params["id"])
    
    
    res.Header().Set("content-type", "application/json")
    err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
    if err != nil {
        res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
    }
    json.NewEncoder(res).Encode(user)
}

func updateUser(res http.ResponseWriter, req *http.Request){
    params := mux.Vars(req)
    id, _ := primitive.ObjectIDFromHex(params["id"])
    
     var user User
    _= json.NewDecoder(req.Body).Decode(&user)
    
    res.Header().Set("content-type","application/json")
    log.Print(user)
    updateResult := collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": id},bson.M{"$set":req.Body})
    
    json.NewEncoder(res).Encode(updateResult)
}

func deleteUser(res http.ResponseWriter, req *http.Request) {
    
    params := mux.Vars(req)
    id, _ := primitive.ObjectIDFromHex(params["id"])
    
    res.Header().Set("content-type", "application/json")
    deleteResult := collection.FindOneAndDelete(context.TODO(), bson.M{"_id": id})
   
    json.NewEncoder(res).Encode(deleteResult)
}

func main(){
    fmt.Println("Starting the application...")
    
    router := mux.NewRouter()
    
    /* Here PUT operation don't have any importances */
    router.HandleFunc("/user",getAll).Methods("GET")
    router.HandleFunc("/user/add",insertUser).Methods("POST")
    router.HandleFunc("/user/delete",deleteUsers).Methods("DELETE")
   
    /* Here POST operation don't have any importances */
    router.HandleFunc("/user/{id}",getUser).Methods("GET")
    router.HandleFunc("/user/update/{id}",updateUser).Methods("PUT")
    router.HandleFunc("/user/delete/{id}",deleteUser).Methods("DELETE")
    
    http.ListenAndServe(":3030", router)
    
}
