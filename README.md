# CloudChain Test Server

## About
The CloudChain Test Server was made to test my [CloudChain](https://github.com/peterxu30/cloudchain) project. 
The server runs on Google App Engine and CloudChain's backing store is Google's Cloud Firestore. Please read the CloudChain readme for additional details.

The server serves no purpose beyond testing the basic functionality of the CloudChain implementation and I see little value in setting up your own instance. But if you want to, you will need to change the testProjectId constant in cloudchain_utils.go to your GCP project ID.

## Tests
Currently, there are two tests that can be run. I run these tests through Postman.

1. Add Fifty Blocks Test
..* Usage: POST to [https://cloudchaintestserver.appspot.com/addblockstest](https://cloudchaintestserver.appspot.com/addblockstest)
No parameters or body required.

..* Adds 50 blocks to the CloudChain asynchronously to the cloudchain. The cloudchain is then iterated through to verify 50 blocks with the correct values were added.

2. Simultaneously Add And Read Fifty Blocks Test
..* Usage: POST to [https://cloudchaintestserver.appspot.com/addandreadblockstest](https://cloudchaintestserver.appspot.com/addandreadblockstest)
No parameters or body required.

..* Simultaneously add 50 blocks asynchronously to the cloudchain and read the blockchain from whatever the head currently is to the end. It then synchronously verifies that 50 blocks were added.

## Utility Endpoints
1. Reset
..* Usage: POST to [https://cloudchaintestserver.appspot.com/reset](https://cloudchaintestserver.appspot.com/reset)
No parameters or body required.

..* Deletes current blockchain and reinitializes it.