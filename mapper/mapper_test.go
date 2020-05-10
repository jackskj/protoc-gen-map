package mapper_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	marsh  = jsonpb.Marshaler{}
	update = false
	initDB = true
)

var (
	conn        *grpc.ClientConn
	ctx         context.Context
	db          *sql.DB
	grpcServer  *grpc.Server
	lis         *bufconn.Listener
	requests    *td.Requests
	testResults map[string]interface{}

	blogClient        ex.BlogQueryServiceClient
	reflectClient     td.TestReflectServiceClient
	testMappingClient td.TestMappingServiceClient
	tdSrv             td.TestMappingServiceMapServer
	blogServer        *ex.BlogQueryServiceMapServer
)

// Generate test data before running tests
// Start local server with bufconn
func setup() {
	requests = td.GenerateRequests()
	testResults = make(map[string]interface{})
	db = td.GetPG()
	ctx = context.Background()
	lis = bufconn.Listen(bufSize)
	grpcServer = grpc.NewServer()

	blogServer = &ex.BlogQueryServiceMapServer{DB: db, Dialect: "postgres"}
	ex.RegisterBlogQueryServiceServer(grpcServer, blogServer)
	initdb.RegisterInitServiceServer(grpcServer, &initdb.InitServiceMapServer{DB: db, Dialect: "postgres"})
	td.RegisterTestReflectServiceServer(grpcServer, &td.TestReflectServiceMapServer{DB: db, Dialect: "postgres"})
	tdSrv = td.TestMappingServiceMapServer{DB: db, Dialect: "postgres"}
	td.RegisterTestMappingServiceServer(grpcServer, &tdSrv)

	registerCallbacks()
	registerFailedCallbacks()

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
	updatePtr := flag.Bool("update", false, "update the golden file, results are always considered correct")
	initdbPtr := flag.Bool("initdb", true, "initialize and populate testing database")
	flag.Parse()
	update = *updatePtr
	initDB = *initdbPtr
	setup()
	code := m.Run()
	goldenFile := "../testdata/mapper.golden"
	if update {
		// update golden file
		updateGoldenFile(goldenFile)
	} else {
		// compare existing results
		compareResults(goldenFile)
	}
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

func TestParamsResponse(t *testing.T) {
	req := ex.BlogRequest{
		Id:       uint32(1),
		AuthorId: uint32(1),
	}

	resp, err := blogClient.SelectBlog(ctx, &req)
	protoResult("blogClient.SelectParam", resp, err, nil, false)

	blogServer.Dialect = "mysql"
	resp, err = blogClient.SelectBlog(ctx, &req)
	protoResult("blogClient.IncorrectParam", resp, err, nil, true)

	blogServer.Dialect = ""
	resp, err = blogClient.SelectBlog(ctx, &req)
	protoResult("blogClient.MissingDialectParam", resp, nil, err, true)

	blogServer.Dialect = "postgres"
}

func TestOneMessageStreamingResponse(t *testing.T) {
	req := ex.BlogIdsRequest{Ids: []uint32{1}, Titles: []string{"abc"}}
	resp, err := blogClient.SelectBlogs(ctx, &req)
	sResp, sErr := blogStreamReader(resp)
	protoResult("blogClient.SelectBlogs_1", sResp, err, sErr, false)
}

func TestEmptyMessageStreamingResponse(t *testing.T) {
	req := ex.BlogIdsRequest{
		Ids:    []uint32{1},
		Titles: []string{"a"},
	}
	resp, err := blogClient.SelectBlogs(ctx, &req)
	sResp, sErr := blogStreamReader(resp)
	protoResult("blogClient.SelectBlogs_2", sResp, err, sErr, false)
}

func TestStreamingResponse(t *testing.T) {
	req := ex.BlogIdsRequest{
		Ids:    []uint32{1, 2, 3, 4, 5},
		Titles: []string{"a"},
	}
	resp, err := blogClient.SelectBlogs(ctx, &req)
	sResp, sErr := blogStreamReader(resp)
	protoResult("blogClient.SelectBlogs_3", sResp, err, sErr, false)
}

func TestComplexStreamingResponse(t *testing.T) {
	req := ex.BlogIdsRequest{
		Ids:    []uint32{1, 2, 3, 4, 5, 6, 7, 9},
		Titles: []string{"a"},
	}
	resp, err := blogClient.SelectDetailedBlogs(ctx, &req)
	sResp, sErr := detailedBlogStreamReader(resp)
	protoResult("blogClient.SelectDetailedBlogs", sResp, err, sErr, false)
}

func TestMappingService(t *testing.T) {
	var (
		resp  proto.Message
		posts td.TestMappingService_NullResoultsForSubmapsClient
		err   error
	)
	resp, err = testMappingClient.RepeatedAssociations(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.RepeatedAssociations", resp, err, nil, false)
	resp, err = testMappingClient.EmptyQuery(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.EmptyQuery", resp, err, nil, false)
	resp, err = testMappingClient.InsertQueryAsExec(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.InsertQueryAsExec", resp, err, nil, false)
	resp, err = testMappingClient.ExecAsQuery(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.ExecAsQuery", resp, err, nil, false)
	resp, err = testMappingClient.UnclaimedColumns(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.UnclaimedColumns", resp, err, nil, false)
	resp, err = testMappingClient.MultipleRespForUnary(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.MultipleRespForUnary", resp, err, nil, false)
	resp, err = testMappingClient.NoRespForUnary(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.NoRespForUnary", resp, err, nil, false)
	resp, err = testMappingClient.RepeatedEmpty(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.RepeatedEmpty", resp, err, nil, false)
	resp, err = testMappingClient.EmptyNestedField(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.EmptyNestedField", resp, err, nil, false)
	resp, err = testMappingClient.NoMatchingColumns(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.NoMatchingColumns", resp, err, nil, false)
	resp, err = testMappingClient.AssociationInCollection(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.AssociationInCollection", resp, err, nil, false)
	resp, err = testMappingClient.CollectionInAssociation(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.CollectionInAssociation", resp, err, nil, false)
	resp, err = testMappingClient.SimpleEnum(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.SimpleEnum", resp, err, nil, false)
	resp, err = testMappingClient.NestedEnum(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.NestedEnum", resp, err, nil, false)

	resp, err = testMappingClient.RepeatedTimestamp(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.RepeatedTimestamp", resp, err, nil, true)
	resp, err = testMappingClient.RepeatedPrimative(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.RepeatedPrimative", resp, err, nil, true)

	posts, err = testMappingClient.NullResoultsForSubmaps(ctx, &td.EmptyRequest{})
	sResp, sErr := postReader(posts)
	protoResult("testMappingClient.NullResoultsForSubmaps", sResp, err, sErr, false)
}

func TestUnaryCallbacks(t *testing.T) {
	resp, err := testMappingClient.BlogB(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.BlogB", resp, err, nil, false)
	resp, err = testMappingClient.BlogA(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.BlogA", resp, err, nil, false)
	resp, err = testMappingClient.BlogC(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.BlogC", resp, err, nil, false)
}

func TestStreamingCallbacks(t *testing.T) {
	resp, err := testMappingClient.BlogsB(ctx, &td.EmptyRequest{})
	sResp, sErr := blogStreamReader(resp)
	protoResult("testMappingClient.BlogsB", sResp, err, sErr, false)
	sResp, sErr = blogStreamReader(resp)
	resp, err = testMappingClient.BlogsA(ctx, &td.EmptyRequest{})
	sResp, sErr = blogStreamReader(resp)
	protoResult("testMappingClient.BlogsA", sResp, err, sErr, false)
	sResp, sErr = blogStreamReader(resp)
	resp, err = testMappingClient.BlogsC(ctx, &td.EmptyRequest{})
	sResp, sErr = blogStreamReader(resp)
	protoResult("testMappingClient.BlogsC", sResp, err, sErr, false)
}

func TestFailedUnaryCallbacks(t *testing.T) {
	resp, err := testMappingClient.BlogBF(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.BlogBF", resp, err, nil, true)
	resp, err = testMappingClient.BlogAF(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.BlogAF", resp, err, nil, true)
	resp, err = testMappingClient.BlogCF(ctx, &td.EmptyRequest{})
	protoResult("testMappingClient.BlogCF", resp, err, nil, true)
}

func TestFailedStreamingCallbacks(t *testing.T) {
	resp, err := testMappingClient.BlogsBF(ctx, &td.EmptyRequest{})
	sResp, sErr := blogStreamReader(resp)
	protoResult("testMappingClient.BlogsBF", sResp, err, sErr, true)
	sResp, sErr = blogStreamReader(resp)
	resp, err = testMappingClient.BlogsAF(ctx, &td.EmptyRequest{})
	sResp, sErr = blogStreamReader(resp)
	protoResult("testMappingClient.BlogsAF", sResp, err, sErr, true)
	sResp, sErr = blogStreamReader(resp)
	resp, err = testMappingClient.BlogsCF(ctx, &td.EmptyRequest{})
	sResp, sErr = blogStreamReader(resp)
	protoResult("testMappingClient.BlogsCF", sResp, err, sErr, true)
}

func blogStreamReader(stream ex.BlogQueryService_SelectBlogsClient) ([]proto.Message, error) {
	var responses []proto.Message
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		} else {
			responses = append(responses, resp)
		}
		if err != nil {
			return responses, err
		}
	}
	return responses, nil
}

func detailedBlogStreamReader(stream ex.BlogQueryService_SelectDetailedBlogsClient) ([]proto.Message, error) {
	var responses []proto.Message
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		} else {
			responses = append(responses, resp)
		}
		if err != nil {
			return responses, err
		}
	}
	return responses, nil
}

func postReader(stream td.TestMappingService_NullResoultsForSubmapsClient) ([]proto.Message, error) {
	var responses []proto.Message
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		} else {
			responses = append(responses, resp)
		}
		if err != nil {
			return responses, err
		}
	}
	return responses, nil
}

func registerCallbacks() {
	tdSrv.RegisterBlogBBeforeQueryCallback(td.BlogBeforeQueryCallback)
	tdSrv.RegisterBlogsBBeforeQueryCallback(td.BlogsBeforeQueryCallback)
	tdSrv.RegisterBlogAAfterQueryCallback(td.BlogAfterQueryCallback)
	tdSrv.RegisterBlogsAAfterQueryCallback(td.BlogsAfterQueryCallback)
	tdSrv.RegisterBlogCCache(td.BlogCache)
	tdSrv.RegisterBlogsCCache(td.BlogsCache)
}

func registerFailedCallbacks() {
	tdSrv.RegisterBlogBFBeforeQueryCallback(td.FailedBlogBeforeQueryCallback)
	tdSrv.RegisterBlogsBFBeforeQueryCallback(td.FailedBlogsBeforeQueryCallback)
	tdSrv.RegisterBlogAFAfterQueryCallback(td.FailedBlogAfterQueryCallback)
	tdSrv.RegisterBlogsAFAfterQueryCallback(td.FailedBlogsAfterQueryCallback)
	tdSrv.RegisterBlogCFCache(td.FailedBlogCache)
	tdSrv.RegisterBlogsCFCache(td.FailedBlogsCache)
}

func protoResult(testName string, resp interface{}, err error, sErr error, expectsErr bool) {
	if expectsErr == false {
		if err != nil {
			log.Fatalln("protoc-gen-map error with "+testName+":%v", err)
		}
		if sErr != nil {
			log.Fatalln("protoc-gen-map error with "+testName+":%v", sErr)
		}
		testResults[testName] = resp
	} else if expectsErr == true {
		testResults[testName] = fmt.Sprintf("%v", []error{err, sErr})
	}
}

func updateGoldenFile(goldenFile string) {
	jsonResult := generateResultBytes()
	if err := ioutil.WriteFile(goldenFile, jsonResult, 0644); err != nil {
		log.Fatalln(err)
	}
}

func compareResults(goldenFile string) {
	goldenFileJson, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		log.Fatalln(err)
	}

	jsonResult := generateResultBytes()

	resultDiff := diff.New()
	d, err := resultDiff.Compare(goldenFileJson, jsonResult)
	if err != nil {
		log.Fatalln(err)
	}
	formatter := formatter.NewDeltaFormatter()
	diffString, err := formatter.Format(d)
	if diffString != "{}\n" {
		log.Println("Results Do Not Match Golden File, " +
			"if this is expecred result with go test with --update")
		log.Fatalln(diffString)
	}
}

func generateResultBytes() []byte {
	var jsonResult []byte
	testResults["callbacks"] = td.CallbackPool
	if r, err := json.MarshalIndent(testResults, "", "    "); err != nil {
		log.Fatalln(err)
	} else {
		jsonResult = r
	}
	return jsonResult
}

func teardown() {
	defer conn.Close()
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}
