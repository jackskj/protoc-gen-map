package mapper_test

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	ex "github.com/jackskj/protoc-gen-map/examples"
	td "github.com/jackskj/protoc-gen-map/testdata"
	"github.com/jackskj/protoc-gen-map/testdata/initdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	marsh     = jsonpb.Marshaler{}
	isVerbose = true
	initDB    = true
)

var (
	conn       *grpc.ClientConn
	ctx        context.Context
	db         *sql.DB
	grpcServer *grpc.Server
	lis        *bufconn.Listener
	requests   *td.Requests

	blogClient        ex.BlogQueryServiceClient
	reflectClient     td.TestReflectServiceClient
	testMappingClient td.TestMappingServiceClient
)

// Generate test data before running tests
// Start local server with bufconn
func setup() {
	requests = td.GenerateRequests()
	db = td.GetPG()
	ctx = context.Background()
	lis = bufconn.Listen(bufSize)
	grpcServer = grpc.NewServer()

	ex.RegisterBlogQueryServiceServer(grpcServer, &ex.BlogQueryServiceMapServer{DB: db})
	initdb.RegisterInitServiceServer(grpcServer, &initdb.InitServiceMapServer{DB: db})
	td.RegisterTestReflectServiceServer(grpcServer, &td.TestReflectServiceMapServer{DB: db})
	td.RegisterTestMappingServiceServer(grpcServer, &td.TestMappingServiceMapServer{DB: db})

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	if connection, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure()); err != nil {
		log.Fatalf("bufnet dial fail: %v", err)
	} else {
		conn = connection
	}
	if initDB {
		createDatabase()
	}
	blogClient = ex.NewBlogQueryServiceClient(conn)
	reflectClient = td.NewTestReflectServiceClient(conn)
	testMappingClient = td.NewTestMappingServiceClient(conn)
}

func TestMain(m *testing.M) {
	verbosePtr := flag.Bool("verbose", false, "run verboce, prints proto responses")
	initdbPtr := flag.Bool("initdb", true, "initialize and populate testing database")
	flag.Parse()
	isVerbose = *verbosePtr
	initDB = *initdbPtr
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func createDatabase() {
	initService := initdb.NewInitServiceClient(conn)
	initService.InitDB(ctx, &initdb.EmptyRequest{})
	for i := 0; i < len(requests.InsertAuthorRequests); i++ {
		initService.InsertAuthor(ctx, requests.InsertAuthorRequests[i])
	}
	for i := 0; i < len(requests.InsertBlogRequests); i++ {
		initService.InsertBlog(ctx, requests.InsertBlogRequests[i])
	}
	for i := 0; i < len(requests.InsertCommentRequests); i++ {
		initService.InsertComment(ctx, requests.InsertCommentRequests[i])
	}
	for i := 0; i < len(requests.InsertPostRequests); i++ {
		initService.InsertPost(ctx, requests.InsertPostRequests[i])
	}
	for i := 0; i < len(requests.InsertPostTagRequests); i++ {
		initService.InsertPostTag(ctx, requests.InsertPostTagRequests[i])
	}
	for i := 0; i < len(requests.InsertTagRequests); i++ {
		initService.InsertTag(ctx, requests.InsertTagRequests[i])
	}
}

func TestOneMessageStreamingResponse(t *testing.T) {
	req := ex.BlogIdsRequest{Ids: []uint32{1}, Titles: []string{"abc"}}
	resp, err := blogClient.SelectBlogs(ctx, &req)
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	blogStreamReader(resp)
}

func TestEmptyMessageStreamingResponse(t *testing.T) {
	req := ex.BlogIdsRequest{
		Ids:    []uint32{1},
		Titles: []string{"a"},
	}
	resp, err := blogClient.SelectBlogs(ctx, &req)
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	blogStreamReader(resp)
}

func TestStreamingResponse(t *testing.T) {
	req := ex.BlogIdsRequest{
		Ids:    []uint32{1, 2, 3, 4, 5},
		Titles: []string{"a"},
	}
	resp, err := blogClient.SelectBlogs(ctx, &req)
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	blogStreamReader(resp)
}

func TestComplexStreamingResponse(t *testing.T) {
	req := ex.BlogIdsRequest{
		Ids:    []uint32{1, 2, 3, 4, 5, 6, 7, 9},
		Titles: []string{"a"},
	}
	resp, err := blogClient.SelectDetailedBlogs(ctx, &req)
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	detailedBlogStreamReader(resp)
}

func TestMappingService(t *testing.T) {
	var (
		resp  proto.Message
		posts td.TestMappingService_NullResoultsForSubmapsClient
		err   error
	)
	resp, err = testMappingClient.RepeatedAssociations(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.EmptyQuery(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.InsertQueryAsExec(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.ExecAsQuery(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.UnclaimedColumns(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.MultipleRespForUnary(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.RepeatedPrimative(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Printf("stream error: %s\n", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.RepeatedEmpty(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.EmptyNestedField(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.NoMatchingColumns(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.AssociationInCollection(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.CollectionInAssociation(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.RepeatedTimestamp(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Printf("stream error: %s\n", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.SimpleEnum(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Printf("stream error: %s\n", err)
	}
	if isVerbose {
		printResp(resp)
	}
	resp, err = testMappingClient.NestedEnum(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Printf("stream error: %s\n", err)
	}
	if isVerbose {
		log.Printf("  asdasdasdsa: %s\n", err)
		printResp(resp)
	}
	posts, err = testMappingClient.NullResoultsForSubmaps(ctx, &td.EmptyRequest{})
	if err != nil {
		log.Fatalf("stream error: %s", err)
	}
	postReader(posts)
}

func blogStreamReader(stream ex.BlogQueryService_SelectBlogsClient) {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream error: %s", err)
		}
		if isVerbose {
			printResp(resp)
		}
	}
}

func detailedBlogStreamReader(stream ex.BlogQueryService_SelectDetailedBlogsClient) {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream error: %s", err)
		}
		if isVerbose {
			printResp(resp)
		}
	}
}

func postReader(stream td.TestMappingService_NullResoultsForSubmapsClient) {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream error: %s", err)
		}
		if isVerbose {
			printResp(resp)
		}
	}
}

func teardown() {
	defer conn.Close()
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}
func printResp(resp proto.Message) {
	fmt.Println(marsh.MarshalToString(resp))
	fmt.Println()
}
