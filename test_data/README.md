The test data here is based on the test data provided in the [Gtker guide](https://gtker.com/implementation-guide-for-the-world-of-warcraft-flavor-of-srp6/).

I opted to create new test inputs because I found a lot of the test data from the guide had mixed endianness. This was causing a lot of confusion and headaches while debugging.

The expected value included in the CSVs had to be added manually. I start by verifying that the test works using some known test inputs from the guide. Then, I run the test using my CSV data as input, and print the result. The results get copy/pasted into the CSV and I enable the assertion in the test.

For example, when generating the expected values for `calculate_x.csv`, I replaced the assertion with this:

```go
fmt.Printf("%x\n", calcX(username, password, salt).Bytes())
```
