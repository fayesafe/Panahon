package server

import(
    "errors"
    "io/ioutil"
    "os"
    "net/http"
    "net/http/httptest"
    "testing"

    "Panahon/logger"
    "github.com/influxdata/influxdb/client/v2"
    "github.com/influxdata/influxdb/models"
    "github.com/gorilla/mux"
)

type DBClientMockOK struct {
    Server string
    Series string
}

func (influxClient DBClientMockOK) QueryAll(
    offset string) (*client.Response, error) {
    return createQueryReturn(), nil
}

func (influxClient DBClientMockOK) QueryAverage(
    col string,
    interval string,
    offset string,
    end string) (*client.Response, error) {
    return createQueryReturn(), nil
}

func (influx DBClientMockOK) QueryInterval(
    low string, high string) (*client.Response, error) {
    return createQueryReturn(), nil
}

func (influx DBClientMockOK) QueryMax(
    col string,
    interval string,
    offset string,
    end string) (*client.Response, error) {
    return createQueryReturn(), nil
}

func (influx DBClientMockOK) QueryMin(
    col string,
    interval string,
    offset string,
    end string) (*client.Response, error) {
    return createQueryReturn(), nil
}

type DBClientMockError struct {
    Server string
    Series string
}

func (influxClient DBClientMockError) QueryAll(
    offset string) (*client.Response, error) {
    return nil, errors.New("=== Test Error ===")
}

func (influxClient DBClientMockError) QueryAverage(
    col string,
    interval string,
    offset string,
    end string) (*client.Response, error) {
    return nil, errors.New("=== Test Error ===")
}

func (influx DBClientMockError) QueryInterval(
    low string, high string) (*client.Response, error) {
    return nil, errors.New("=== Test Error ===")
}

func (influx DBClientMockError) QueryMax(
    col string,
    interval string,
    offset string,
    end string) (*client.Response, error) {
    return nil, errors.New("=== Test Error ===")
}

func (influx DBClientMockError) QueryMin(
    col string,
    interval string,
    offset string,
    end string) (*client.Response, error) {
    return nil, errors.New("=== Test Error ===")
}

var expectedResponse string = `{"Series":[{"name":"Test","tags":{"test":"test`+
`"},"columns":["test","test"]}],"Messages":null,"error":"test"}`

func createQueryReturn() (*client.Response) {
    tags := map[string]string{
        "test":"test",
    }
    result := models.Row{
        Name: "Test",
        Tags: tags,
        Columns: []string{"test", "test"},
        Values: nil,
        Err: nil,
    }
    results := []client.Result{
        client.Result{
            Series: []models.Row{result},
            Messages: nil,
            Err: "test",
        },
    }
    returnVal := client.Response{
        Results: results,
        Err: "no error",
    }
    return &returnVal
}

//
// UNIT TESTS - GENERAL FUNCTIONS
//

func TestSendPayload(t *testing.T) {
    w := httptest.NewRecorder()
    sendPayload(createQueryReturn(), w)
    if w.Code != 200 {
        t.Errorf("Wrong status code for JSON response")
    }
    if w.Body.String() != expectedResponse {
        t.Errorf("Wrong body for JSON response")
    }
}

func TestIsStringNum(t *testing.T) {
    expectedTrue := "0"
    expectedFalse := "x"

    if !isStringNum(expectedTrue) {
        t.Errorf("0 is a num representation")
    }
    if isStringNum(expectedFalse) {
        t.Errorf("x is not a num representation")
    }
}

//
// UNIT TEST - HANDLER
//

func TestQueryLastOk(t *testing.T) {
    mockInflux := DBClientMockOK{
        Server : "",
        Series : "",
    }
    queryHandle := queryHandleLast(mockInflux)
    req, _ := http.NewRequest("GET", "", nil)
    w := httptest.NewRecorder()

    queryHandle.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("Home page didn't return %v", http.StatusOK)
    }
    if w.Body.String() != expectedResponse {
        t.Errorf("Incorrect Body.")
    }
}

func TestQueryLastError(t *testing.T) {
    mockInflux := DBClientMockError{
        Server : "",
        Series : "",
    }
    queryHandle := queryHandleLast(mockInflux)
    req, _ := http.NewRequest("GET", "", nil)
    w := httptest.NewRecorder()

    queryHandle.ServeHTTP(w, req)
    if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

func TestQueryIntervalOk(t *testing.T) {
    mockInflux := DBClientMockOK{
        Server : "",
        Series : "",
    }
    queryHandle := queryHandleInterval(mockInflux)
    req, _ := http.NewRequest("GET", "", nil)
    w := httptest.NewRecorder()

    queryHandle.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("Home page didn't return %v", http.StatusOK)
    }
    if w.Body.String() != expectedResponse {
        t.Errorf("Incorrect Body.")
    }
}

func TestQueryIntervalError(t *testing.T) {
    mockInflux := DBClientMockError{
        Server : "",
        Series : "",
    }
    queryHandle := queryHandleInterval(mockInflux)
    req, _ := http.NewRequest("GET", "", nil)
    w := httptest.NewRecorder()

    queryHandle.ServeHTTP(w, req)
    if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

func TestQueryAverageOk(t *testing.T) {
    mockInflux := DBClientMockOK{
        Server : "",
        Series : "",
    }
    queryHandle := queryHandleAverage(mockInflux)
    req, _ := http.NewRequest("GET", "", nil)
    w := httptest.NewRecorder()

    queryHandle.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("Home page didn't return %v", http.StatusOK)
    }
    if w.Body.String() != expectedResponse {
        t.Errorf("Incorrect Body.")
    }
}

func TestQueryAverageOkError(t *testing.T) {
    mockInflux := DBClientMockError{
        Server : "",
        Series : "",
    }
    queryHandle := queryHandleAverage(mockInflux)
    req, _ := http.NewRequest("GET", "", nil)
    w := httptest.NewRecorder()

    queryHandle.ServeHTTP(w, req)
    if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

//
// INTEGRATION TESTS
//

type MockInterface interface {
    QueryAll(offset string) (*client.Response, error)
    QueryAverage(
        col string,
        interval string,
        offset string,
        end string) (*client.Response, error)
    QueryInterval(low string, high string) (*client.Response, error)
    QueryMax(col string,
        interval string,
        offset string,
        end string) (*client.Response, error)
    QueryMin(col string,
        interval string,
        offset string,
        end string) (*client.Response, error)
}

type DBClientMockIntegrationOK struct {
    Server string
    Series string
}

func (influxClient DBClientMockIntegrationOK) QueryAll(
    offset string) (*client.Response, error) {
    return createQueryReturnWithParams(offset, "x", "y", "z"), nil
}

func (influxClient DBClientMockIntegrationOK) QueryAverage(
    col string, interval string, offset string, end string) (*client.Response, error) {
    return createQueryReturnWithParams(col, interval, offset, end), nil
}

func (influx DBClientMockIntegrationOK) QueryInterval(
    low string, high string) (*client.Response, error) {
    return createQueryReturnWithParams(low, high, "x", "y"), nil
}

func (influx DBClientMockIntegrationOK) QueryMax(
    col string, interval string, offset string, end string) (*client.Response, error) {
    return createQueryReturnWithParams(col, interval, offset, end), nil
}

func (influx DBClientMockIntegrationOK) QueryMin(
    col string, interval string, offset string, end string) (*client.Response, error) {
    return createQueryReturnWithParams(col, interval, offset, end), nil
}

type DBClientMockIntegrationError struct {
    Server string
    Series string
}

func (influxClient DBClientMockIntegrationError) QueryAll(
    offset string) (*client.Response, error) {
    return nil, errors.New("=== Test Error ===")
}

func (influxClient DBClientMockIntegrationError) QueryAverage(
    col string, interval string, offset string, end string) (*client.Response, error) {
    return nil, errors.New("=== Test Error ===")
}

func (influx DBClientMockIntegrationError) QueryInterval(
    low string, high string) (*client.Response, error) {
    return nil, errors.New("=== Test Error ===")
}

func (influx DBClientMockIntegrationError) QueryMax(
    col string, interval string, offset string, end string) (*client.Response, error) {
    return nil, errors.New("=== Test Error ===")
}

func (influx DBClientMockIntegrationError) QueryMin(
    col string, interval string, offset string, end string) (*client.Response, error) {
    return nil, errors.New("=== Test Error ===")
}

var expectedResponseOneParamEmpty string = `{"Series":[{"name":"Test","tags":`+
`{"":"","test":"test","x":"x","y":"y","z":"z"},"columns":["test","test"]}],"M`+
`essages":null,"error":"test"}`
var expectedResponseOneParam string = `{"Series":[{"name":"Test","tags":{"1":`+
`"1","test":"test","x":"x","y":"y","z":"z"},"columns":["test","test"]}],"Mess`+
`ages":null,"error":"test"}`
var expectedResponseTwoParamEmpty string = `{"Series":[{"name":"Test","tags":`+
`{"0":"0","2147483647000":"2147483647000","test":"test","x":"x","y":"y"},"columns":`+
`["test","test"]}],"Messages":null,"error":"test"}`
var expectedResponseTwoParamRangeOne string = `{"Series":[{"name":"Test","tag`+
`s":{"10":"10","2147483647000":"2147483647000","test":"test","x":"x","y":"y"},"colu`+
`mns":["test","test"]}],"Messages":null,"error":"test"}`
var expectedResponseTwoParamRangeFull string = `{"Series":[{"name":"Test","ta`+
`gs":{"10":"10","100":"100","test":"test","x":"x","y":"y"},"columns":["test",`+
`"test"]}],"Messages":null,"error":"test"}`
var expectedResponseAvWoOffset string = `{"Series":[{"name":"Test","tags":{"0`+
`":"0","10w":"10w","2147483647000":"2147483647000","temp":"temp","test":"test"},"co`+
`lumns":["test","test"]}],"Messages":null,"error":"test"}`
var expectedResponseAvWOffset string = `{"Series":[{"name":"Test","tags":{"10`+
`0":"100","10w":"10w","2147483647000":"2147483647000","temp":"temp","test":"test"},`+
`"columns":["test","test"]}],"Messages":null,"error":"test"}`
var expectedResponseAvFull string = `{"Series":[{"name":"Test","tags":{"100":`+
`"100","10w":"10w","123":"123","temp":"temp","test":"test"},"columns":["test"`+
`,"test"]}],"Messages":null,"error":"test"}`

type TestRoute struct {
    Prefix string
    Route string
    Handler func(dbClient) http.Handler
    Mock MockInterface
    expectedCode int
    expectedBody string
}

type routesToTest []TestRoute

func createQueryReturnWithParams(
    param1 string, param2 string, param3 string, param4 string) (*client.Response) {
    tags := map[string]string{
        "test":"test",
        param1:param1,
        param2:param2,
        param3:param3,
        param4:param4,
    }
    result := models.Row{
        Name: "Test",
        Tags: tags,
        Columns: []string{"test", "test"},
        Values: nil,
        Err: nil,
    }
    results := []client.Result{
        client.Result{
            Series: []models.Row{result},
            Messages: nil,
            Err: "test",
        },
    }
    returnVal := client.Response{
        Results: results,
        Err: "no error",
    }
    return &returnVal
}

func testHelper(t *testing.T, route TestRoute) {
    req, _ := http.NewRequest("GET", route.Route, nil)
    res := httptest.NewRecorder()
    m := SetupMux(route.Prefix, route.Handler(route.Mock))
    m.ServeHTTP(res, req)
    if res.Body.String() != route.expectedBody {
        t.Errorf("%s does not match %s", res.Body.String(), route.expectedBody)
        t.Log(route.Prefix, route.Route)
    }
    if res.Code != route.expectedCode {
        t.Errorf("%d does not match %d", res.Code, route.expectedCode)
    }
}

func TestIntegrateRoutes(t *testing.T) {
    mockInfluxHappy := DBClientMockIntegrationOK{
        Server : "",
        Series : "",
    }
    mockInfluxError := DBClientMockIntegrationError{
        Server : "",
        Series : "",
    }
    testingRoutes := routesToTest{
        TestRoute{
            Prefix : "/last",
            Route : "/last",
            Handler : queryHandleLast,
            Mock : mockInfluxHappy,
            expectedCode : http.StatusOK,
            expectedBody : expectedResponseOneParamEmpty,
        },
        TestRoute{
            Prefix : "/last",
            Route : "/last",
            Handler : queryHandleLast,
            Mock : mockInfluxError,
            expectedCode : http.StatusInternalServerError,
            expectedBody : "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/last/{last:[0-9]+}",
            Route : "/last/1",
            Handler : queryHandleLast,
            Mock : mockInfluxHappy,
            expectedCode : http.StatusOK,
            expectedBody : expectedResponseOneParam,
        },
        TestRoute{
            Prefix : "/last/{last:[0-9]+}",
            Route : "/last/1",
            Handler : queryHandleLast,
            Mock : mockInfluxError,
            expectedCode : http.StatusInternalServerError,
            expectedBody : "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/range",
            Route : "/range",
            Handler : queryHandleInterval,
            Mock : mockInfluxHappy,
            expectedCode : http.StatusOK,
            expectedBody : expectedResponseTwoParamEmpty,
        },
        TestRoute{
            Prefix : "/range",
            Route : "/range",
            Handler : queryHandleInterval,
            Mock : mockInfluxError,
            expectedCode : http.StatusInternalServerError,
            expectedBody : "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/range/{low:[0-9]+}",
            Route : "/range/10",
            Handler : queryHandleInterval,
            Mock : mockInfluxHappy,
            expectedCode : http.StatusOK,
            expectedBody : expectedResponseTwoParamRangeOne,
        },
        TestRoute{
            Prefix : "/range/{low:[0-9]+}",
            Route : "/range/10",
            Handler : queryHandleInterval,
            Mock : mockInfluxError,
            expectedCode : http.StatusInternalServerError,
            expectedBody : "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/range/{low:[0-9]+}/{high:[0-9]+}",
            Route : "/range/10/100",
            Handler : queryHandleInterval,
            Mock : mockInfluxHappy,
            expectedCode : http.StatusOK,
            expectedBody : expectedResponseTwoParamRangeFull,
        },
        TestRoute{
            Prefix : "/range/{low:[0-9]+}/{high:[0-9]+}",
            Route : "/range/10/100",
            Handler : queryHandleInterval,
            Mock : mockInfluxError,
            expectedCode : http.StatusInternalServerError,
            expectedBody : "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/av/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}",
            Route : "/av/temp/10w",
            Handler : queryHandleAverage,
            Mock : mockInfluxHappy,
            expectedCode : http.StatusOK,
            expectedBody : expectedResponseAvWoOffset,
        },
        TestRoute{
            Prefix : "/av/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}",
            Route : "/av/temp/10w",
            Handler : queryHandleAverage,
            Mock : mockInfluxError,
            expectedCode : http.StatusInternalServerError,
            expectedBody : "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/av/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}",
            Route : "/av/temp/10w/100",
            Handler : queryHandleAverage,
            Mock : mockInfluxHappy,
            expectedCode : http.StatusOK,
            expectedBody : expectedResponseAvWOffset,
        },
        TestRoute{
            Prefix : "/av/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}",
            Route : "/av/temp/10w/100",
            Handler : queryHandleAverage,
            Mock : mockInfluxError,
            expectedCode : http.StatusInternalServerError,
            expectedBody : "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/av/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}/{high:[0-9]+}",
            Route : "/av/temp/10w/100/123",
            Handler : queryHandleAverage,
            Mock : mockInfluxHappy,
            expectedCode : http.StatusOK,
            expectedBody : expectedResponseAvFull,
        },
        TestRoute{
            Prefix : "/av/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}/{high:[0-9]+}",
            Route : "/av/temp/10w/100/123",
            Handler : queryHandleAverage,
            Mock : mockInfluxError,
            expectedCode : http.StatusInternalServerError,
            expectedBody : "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}",
            Route : "/max/temp/10w",
            Handler : queryHandleMax,
            Mock : mockInfluxHappy,
            expectedCode: http.StatusOK,
            expectedBody: expectedResponseAvWoOffset,
        },
        TestRoute{
            Prefix : "/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}",
            Route : "/max/temp/10w",
            Handler : queryHandleMax,
            Mock : mockInfluxError,
            expectedCode: http.StatusInternalServerError,
            expectedBody: "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}",
            Route : "/max/temp/10w/100",
            Handler : queryHandleMax,
            Mock : mockInfluxHappy,
            expectedCode: http.StatusOK,
            expectedBody: expectedResponseAvWOffset,
        },
        TestRoute{
            Prefix : "/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}",
            Route : "/max/temp/10w/100",
            Handler : queryHandleMax,
            Mock : mockInfluxError,
            expectedCode: http.StatusInternalServerError,
            expectedBody: "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}/{high:[0-9]+}",
            Route : "/max/temp/10w/100/123",
            Handler : queryHandleMax,
            Mock : mockInfluxHappy,
            expectedCode: http.StatusOK,
            expectedBody: expectedResponseAvFull,
        },
        TestRoute{
            Prefix : "/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}/{high:[0-9]+}",
            Route : "/max/temp/10w/100/123",
            Handler : queryHandleMax,
            Mock : mockInfluxError,
            expectedCode: http.StatusInternalServerError,
            expectedBody: "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/min/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}",
            Route : "/min/temp/10w",
            Handler : queryHandleMin,
            Mock : mockInfluxHappy,
            expectedCode: http.StatusOK,
            expectedBody: expectedResponseAvWoOffset,
        },
        TestRoute{
            Prefix : "/min/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}",
            Route : "/min/temp/10w",
            Handler : queryHandleMin,
            Mock : mockInfluxError,
            expectedCode: http.StatusInternalServerError,
            expectedBody: "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/min/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}",
            Route : "/min/temp/10w/100",
            Handler : queryHandleMin,
            Mock : mockInfluxHappy,
            expectedCode: http.StatusOK,
            expectedBody: expectedResponseAvWOffset,
        },
        TestRoute{
            Prefix : "/min/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}",
            Route : "/min/temp/10w/100",
            Handler : queryHandleMin,
            Mock : mockInfluxError,
            expectedCode: http.StatusInternalServerError,
            expectedBody: "Internal Server Error\n",
        },
        TestRoute{
            Prefix : "/min/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}/{high:[0-9]+}",
            Route : "/min/temp/10w/100/123",
            Handler : queryHandleMin,
            Mock : mockInfluxHappy,
            expectedCode: http.StatusOK,
            expectedBody: expectedResponseAvFull,
        },
        TestRoute{
            Prefix : "/min/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}/{high:[0-9]+}",
            Route : "/min/temp/10w/100/123",
            Handler : queryHandleMin,
            Mock : mockInfluxError,
            expectedCode: http.StatusInternalServerError,
            expectedBody: "Internal Server Error\n",
        },
    }
    for _, TestRoute := range testingRoutes {
        testHelper(t, TestRoute)
    }
}


func TestMain(m *testing.M) {
	SetupLogger()
	exitStatus := m.Run()
	os.Exit(exitStatus)
}

func SetupMux(prefix string, handler http.Handler) *mux.Router {
    m := mux.NewRouter()
    m.PathPrefix(prefix).Handler(handler).Methods("GET")
    return m
}

func SetupLogger() {
	logger.Init(ioutil.Discard, ioutil.Discard, ioutil.Discard)
}
