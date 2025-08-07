# Performance Test for Sentiment API

This performance test is designed to benchmark the Sentiment Analysis API endpoints with the exact same calculations and methodology as the original Python version.

## Features

- **Comprehensive Testing**: Tests 30 different Indonesian text samples with various sentiments
- **Performance Metrics**: Measures response time, success rate, throughput, and statistical analysis
- **Accuracy Comparison**: Compares results between different API endpoints
- **JSON Output**: Saves detailed results to JSON file for further analysis
- **Real-time Progress**: Shows progress during testing with detailed output

## Performance Metrics Calculated

### Speed Metrics
- Total execution time
- Average response time
- Minimum/Maximum response time
- Median response time
- Requests per second (throughput)

### Reliability Metrics
- Success rate percentage
- Failed requests count
- Error analysis

### Comparison Analysis
- Agreement rate between different endpoints
- Detailed comparison of each test case
- Winner analysis across different criteria

## Usage

### Prerequisites
1. Make sure the Sentiment API server is running
2. The server should be accessible at `http://localhost:8000` (default) or specify custom URL

### Running the Test

#### Option 1: Run directly
```bash
go run performance_test.go
```

#### Option 2: Build and run executable
```bash
go build -o performance_test.exe performance_test.go
./performance_test.exe
```

#### Option 3: Custom URL
Modify the `baseURL` parameter in the `main()` function or update the code to accept command line arguments.

### Expected Output

The test will output:
1. **Real-time progress** showing each request being processed
2. **Summary statistics** including:
   - Success rates for each endpoint
   - Speed comparison metrics
   - Accuracy comparison results
   - Winner analysis
3. **JSON file** (`sentiment_performance_test.json`) with detailed results

### Sample Output
```
🧪 Sentiment Analysis Performance Test
Testing API v1 vs API v2 endpoints
Total requests: 30 per endpoint
============================================================
✅ Server is running, starting tests...

🚀 Starting Performance Test...
Testing 30 requests for each endpoint
============================================================
Testing API v1 endpoint (/api/v1/sentiment/analyze)...
API v1 Request 1/30: Produk ini sangat bagus dan berkualitas tinggi...
...

🏆 PERFORMANCE TEST RESULTS
================================================================================
📊 SUCCESS RATES:
   API v1: 100.0% (30/30)
   API v2: 100.0% (30/30)

⚡ SPEED COMPARISON:
   API v1 Total Time: 15.23s
   API v2 Total Time: 12.45s
   API v1 Avg Response: 0.508s
   API v2 Avg Response: 0.415s
   API v1 Throughput: 1.97 req/s
   API v2 Throughput: 2.41 req/s

🎯 ACCURACY COMPARISON:
   Agreement Rate: 93.3%
   Agreements: 28
   Disagreements: 2

🏅 WINNERS:
   Faster Average Response: API v2
   Higher Success Rate: Tie
   Higher Throughput: API v2
   Faster Total Time: API v2
```

## Test Data

The test uses 30 carefully selected Indonesian text samples covering:
- **Positive sentiments**: Product praise, satisfaction expressions
- **Negative sentiments**: Complaints, disappointments
- **Neutral sentiments**: Balanced or factual statements
- **Edge cases**: Extreme language, mixed sentiments

## API Endpoint Compatibility

The current implementation tests the existing API structure:
- **Endpoint**: `/api/v1/sentiment/analyze`
- **Method**: POST
- **Request Format**:
  ```json
  {
    "text_pertanyaan": "Bagaimana pendapat Anda?",
    "text_jawaban": "Text to analyze..."
  }
  ```
- **Response Format**:
  ```json
  {
    "success": true,
    "data": {
      "sentiment": "Positif"
    }
  }
  ```

## Customization

### Adding Different Endpoints
To test different endpoints (like the Python version's `/analyze-sentiment-llm` and `/analyze-sentiment-indobert`), modify the endpoint URLs in the `runPerformanceTest()` function:

```go
// Change these lines in runPerformanceTest()
result := spt.testAPIEndpoint(text, "/analyze-sentiment-llm")
result := spt.testAPIEndpoint(text, "/analyze-sentiment-indobert")
```

### Modifying Test Data
Update the `testTexts` array in `NewSentimentPerformanceTest()` to add your own test cases.

### Changing Base URL
Modify the default URL in the `main()` function:
```go
test := NewSentimentPerformanceTest("http://your-server:port")
```

## Error Handling

The test handles various error scenarios:
- Network connectivity issues
- HTTP errors (4xx, 5xx)
- JSON parsing errors
- Timeout errors (30-second timeout per request)
- Server unavailability

## Output Files

### JSON Results File
The test generates `sentiment_performance_test.json` containing:
- Test configuration
- Raw test data
- Detailed results for each endpoint
- Performance comparison metrics
- Accuracy comparison data
- Complete response data for debugging

## Performance Comparison with Python Version

This Go implementation provides identical functionality to the Python version:
- ✅ Same test data (30 Indonesian text samples)
- ✅ Same performance metrics calculations
- ✅ Same statistical analysis (mean, median, min, max)
- ✅ Same accuracy comparison methodology
- ✅ Same JSON output format
- ✅ Same real-time progress reporting
- ✅ Same error handling and timeout management

## Troubleshooting

### Server Not Running
```
❌ Cannot connect to server! Make sure server is running.
```
**Solution**: Start the Sentiment API server first

### All Requests Failing
- Check if the API endpoint URLs are correct
- Verify the request/response format matches your API
- Check server logs for errors
- Ensure proper authentication if required

### Low Success Rate
- Check server performance and capacity
- Monitor server logs for errors
- Consider increasing timeout values
- Verify network connectivity

## Dependencies

- Go 1.19+ (for JSON handling and HTTP client)
- No external dependencies required (uses only Go standard library)
