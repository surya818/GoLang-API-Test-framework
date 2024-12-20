# Candidate Take-Home Exercise - SDET
# Kong Gateway API Tests 
Welcome to the Technical solution for the Candidate take home task for SDET</br>
This repository consists of Test Automation framework to testing APIs of the Service Catalog API developed in Go

**Tech Stack:**</br>
go</br>
Docker</br>
****</br>

**Prerequisites/Setup:**
1. Go and Docker installed installed</br>
2. gofumpt library instteled, which is a requirement for running the server in docker, via a make command 

**How to Run: 
Approach 1: Created a github action** </br>
<img src="https://github.com/user-attachments/assets/3f376401-2cd2-4084-a5d3-977510f97d21" alt="Image Description" width="1200" height="800">

</br>**Approach 2: (on Local Machine, Docker )** </br>
(If you want to just run the tests locally, using Docker) </br>

1. Clone this repo </br>
2. Navigate to root directory of the repo </br>
3. Run the command make docker-run to start the server
4. Verify the server started on docker container and exposed on port 18080. Verify running curl or equivalent http://localhost:18080 for a 404 error
6. For an easy to read HTML report , install the go library ==> ** go install github.com/vakenbolt/go-test-report@latest **
6. Run ** go test -v -json ./... | go-test-report** </br>
5. Verify Test results in test-report.html </br>

 
**What's happening behind the scenes of Github action:**
In a nutshell, in Github actions, we setup the Service catalog server from the application make file, setup all dependacies and once the server is up and running, run the tests against the localhost server
Below is the detailed sequence of steps
<ul type="square">
<li>Install go and dependancies like gofumpt</li>  
<li>Install and Set up docker</li> 
<li>Run the make docker-run command, which starts the server</li> 
<li>Poll for the server startup</li> 
<li>Install go-test-report for HTML reporting</li>  
<li>Run the tests</li> 
<li>Upload HTML report artifact</li> 


</br>**Main Packages used**:</br>
net/http --> For http client</br>
zap --> For logging</br>
testing --> For tests and for the assertions within</br>


</br>**Test Automation Architecture:**
The Test automation code consists of 3 parts
1. Framework layer: Which builds a re-usable http-client and has functionality for invocation of the HTTP Methods. This part of code is API agnostic and has no coupling or any relation with the Service Catalog API.
2. Service Layer: This layer hosts the code that is specific to the APIs of our Service Catalog. A seperate go file is created for each API (like Service API, Service Version API, Token API). The methods inside these files contain the actual APIs invocation within each API of Service catalog.
3. Test Layer: This layer has the actual tests, written utilizing the testing package. The Tests instantiates "apis" in the service layer and calls the methods which are the api calls within each area
   
![image](https://github.com/user-attachments/assets/5bc5de67-d519-41cd-b998-4d39a8d69f0c)


**Utils:**
Consists of utility code that could be used within tests or even service layer code that abstracts api supporting code. E.g, Parsing Http Response to strings and Generic structs, tokenizing JWT tokens and templating request payloads

**Configuration:**
Nothing is hard coded. Utilized existing configuration for some of the tests, by creating a Configuation object from the config.yml file. 

**Test Data:**
Again, no test data is hard coded. Everything is neatly randomized, using code in utils


<h3>**Test Details**</h3>
The tests focus on Kong Gateway functionality. In a nutshell tests deal with creation, modification, deletion of Control Planes, Services and routes. We have tests covering some of the CRUD operations that the rest api offers. And the tests are accpetance tests integrating flows accross different operations of Kong Gateway.
The Test scenarios covered via Test Automation:
  

**Successful Github Actions Run:** </br>
![image](https://github.com/user-attachments/assets/3f376401-2cd2-4084-a5d3-977510f97d21)

**Sample Test Reports:** </br>
![image](https://github.com/user-attachments/assets/a4260dbc-85a3-4dab-9d36-a92894c4ae45)



****Sample Logging**:**</br>
![image](https://github.com/user-attachments/assets/348b0d72-a29e-4952-928e-35b057c82941)








# kong-takehometask
