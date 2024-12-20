# Candidate Take-Home Exercise - SDET
# Kong Gateway API Tests 
Welcome to the Technical solution for the Candidate take home task for SDET</br>
This repository consists of Test Automation framework to testing APIs of the Service Catalog API developed in Go

**Tech Stack:**</br>
<ul type="square">
<li>go</li>  
<li>Docker</li> 
</ul>

**Prerequisites/Setup:**
1. Go and Docker installed installed</br>
2. gofumpt package (go install mvdan.cc/gofumpt@latest)
3. Make sure GOPATH/bin is accessible in PATH ==> **export PATH=$PATH:$(go env GOPATH)/bin**


**How to Run: 
Approach 1: Created a github action** </br>

![image](https://github.com/user-attachments/assets/ef69500e-4e4b-435a-8674-663a70d94c76)

</br>**What's happening behind the scenes of Github action:**
In a nutshell, in Github actions, we setup the Service catalog server from the application make file, setup all dependacies and once the server is up and running, run the tests against the localhost server
</br>Below is the detailed sequence of steps
<ul type="square">
<li>Install go and dependancies like gofumpt </li>  
<li>Install and Set up docker</li> 
<li>Run the make docker-run command, which starts the server</li> 
<li>Poll for the server startup</li> 
<li>Install go-test-report for HTML reporting</li>  
<li>Run the tests</li> 
<li>Upload HTML report artifact</li> 
</ul>

</br>**Approach 2: (on Local Machine, Docker )** </br>
(If you want to just run the tests locally, using Docker) </br>

1. Clone this repo </br>
2. Navigate to root directory of the repo </br>
3. Run the command make docker-run to start the server
4. Verify the server started on docker container and exposed on port 18080. Verify running curl or equivalent http://localhost:18080 for a 404 response
6. For an easy to read HTML report , install the go library ==> ** go install github.com/vakenbolt/go-test-report@latest **
6. Run ** go test -v -json ./... | go-test-report** </br>
5. Verify Test results in test-report.html </br>

 

</br>**Main Packages used**:</br>
net/http --> For http client</br>
zap --> For logging</br>
testing --> For tests and for the assertions within</br>


</br>**Test Automation Architecture:**
The Test automation code consists of 3 parts
</br>**1. Framework layer/Http Client layer:** Which builds a re-usable http-client and has functionality for invocation of the HTTP Methods. This part of code is API agnostic and has no coupling or any relation with the Service Catalog API.
</br>**2. Service Layer:** This layer hosts the code that is specific to the APIs of our Service Catalog. A seperate go file is created for each API (like Service API, Service Version API, Token API). The methods inside these files contain the actual APIs invocation within each API of Service catalog.
</br>**3. Test Layer:** This layer has the actual tests, written utilizing the testing package. The Tests instantiates "apis" in the service layer and calls the methods which are the api calls within each area
   
</br><img src="https://github.com/user-attachments/assets/5bc5de67-d519-41cd-b998-4d39a8d69f0c" alt="Image Description" width="400" height="450">

**Utils:**
Consists of utility code that could be used within tests or even service layer code that abstracts api supporting code. E.g, Parsing Http Response to strings and Generic structs, tokenizing JWT tokens and templating request payloads

**Configuration:**
Nothing is hard coded. Utilized existing configuration for some of the tests, by creating a Configuation object from the config.yml file. 

**Test Data:**
Again, no test data is hard coded. Everything is neatly randomized, using code in utils

**Overall Code structure:**

![image](https://github.com/user-attachments/assets/2d23de0f-64f6-423a-89e1-d6707e6658b8)



<h3>Test Details</h3>
The tests focus on aspects of functionality, business logic, error handling and finding unexpected behaviour. One of the core ideas is to not only test the happy path functionality of APIs but also provide a high quality APIs by finding corner cases.
A lot of test techniques like boundary value analysis, passing error input to find unexpected responses have been used in this effort. In addition to evaluating the API quality, documentation has also been considered while designing tests.
A common practice is to manually test the scenarios first and then plan the automation tests for the scenarios, so we get a clear idea of what should the auto-tests contain.</br>
</br> 

**Here's the exhaustive list of Test scenarios designed for this effort**
</br>

![image](https://github.com/user-attachments/assets/724c154f-a7bb-4eca-b001-01a6d21c67b0)




**Sample Test Reports:** </br>
</br> Here's a snip of test report generated. Also added the sample report in **test/sample-report** directory. The file name is **test-report.html** just for showing the reporting capability
</br>
![image](https://github.com/user-attachments/assets/a4260dbc-85a3-4dab-9d36-a92894c4ae45)


**Logging:**
Utmost care had been taken in logging everything concerning the api calls mode during the test. And majority of logging is written in the framework layer, which is reused across Service and Test layer, reducing the number of log statements in the code. One would establish a fair understanding of what's happening in the test by the sequence of log statements</br>

****Sample Logging**:**</br>
![image](https://github.com/user-attachments/assets/348b0d72-a29e-4952-928e-35b057c82941)


</br> **Bug-report** 
</br> A detailed bug report had been uploaded in the repository at repository root directory, the file name is **BugReport.xlsx**

# kong-takehometask
