package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"math/rand"

	"github.com/moby/moby/pkg/namesgenerator"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/gorilla/handlers"

)

var (
	rdb           *redis.Client
	upgrader      = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	socketClients = make(map[*websocket.Conn]string)
	cfg, _        = config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	ecsClient     = ecs.NewFromConfig(cfg)
)

type ProjectRequest struct {
	GitURL string `json:"gitURL"`
	Slug   string `json:"slug"`
}

// Fill these values
const (
	clusterARN = ""
	taskDefARN = ""
	subnet1    = ""
	subnet2    = ""
	subnet3    = ""
	secGroup   = ""
)



func main() {

	// Fill TOKEN Value
	opt, err := redis.ParseURL("rediss://default:<TOKEN>@on-buffalo-26598.upstash.io:6379")
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	
	rdb = redis.NewClient(opt)

	go subscribeToRedis()

	router := mux.NewRouter()
	router.HandleFunc("/ws", handleWebSocket)
	router.HandleFunc("/project", handleProjectCreate).Methods("POST", "OPTIONS") // add OPTIONS here

	// Apply CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // or restrict to ["http://localhost:5173"]
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(router)

	log.Println("API Server running on :9000")
	log.Fatal(http.ListenAndServe(":9000", corsHandler))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		channel := string(msg)
		socketClients[conn] = channel
		conn.WriteMessage(websocket.TextMessage, []byte("Joined "+channel))
	}
}

func subscribeToRedis() {
	pubsub := rdb.PSubscribe(context.TODO(), "logs:*")
	ch := pubsub.Channel()

	log.Println("Subscribed to Redis logs:*")

	for msg := range ch {
		for conn, room := range socketClients {
			if strings.HasPrefix(msg.Channel, "logs:"+room) {
				conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			}
		}
	}
}

func generateSlug() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return namesgenerator.GetRandomName(int(r.Int63()))
}

func handleProjectCreate(w http.ResponseWriter, r *http.Request) {
	var req ProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	slug := req.Slug
	if slug == "" {
		slug = generateSlug()
	}

	_, err := ecsClient.RunTask(context.TODO(), &ecs.RunTaskInput{
		Cluster:        aws.String(clusterARN),
		LaunchType:     types.LaunchTypeFargate,
		Count:          aws.Int32(1),
		TaskDefinition: aws.String(taskDefARN),
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				AssignPublicIp: types.AssignPublicIpEnabled,
				Subnets:        []string{subnet1, subnet2, subnet3},
				SecurityGroups: []string{secGroup},
			},
		},
		Overrides: &types.TaskOverride{
			ContainerOverrides: []types.ContainerOverride{
				{
					// Fill image name
					Name: aws.String("<IMAGE_NAME>"),
					Environment: []types.KeyValuePair{
						{Name: aws.String("GIT_REPOSITORY__URL"), Value: aws.String(req.GitURL)},
						{Name: aws.String("PROJECT_ID"), Value: aws.String(slug)},
					},
				},
			},
		},
	})
	if err != nil {
		http.Error(w, "ECS Task failed to run: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"status": "queued",
		"data": map[string]string{
			"projectSlug": slug,
			"url":         "http://" + slug + ".localhost:8000",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}