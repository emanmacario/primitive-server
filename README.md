## Primitive Server


A web server that acts as an image transformation service. It allows a user to upload an image, and 
then guides them through a selection process using various options on the [primitive](https://github.com/fogleman/primitive) 
CLI to transform the image.

## Usage
1. Clone the repository
2. Install the ```primitive``` package
   ```
   go get -u github.com/fogleman/primitive
   ```
3. Run the server
  ```bash
  go run server.go
  ```
  
 ## Note:
 Currently a work in progress, will make a front-end client in React.js to view various image transformations
