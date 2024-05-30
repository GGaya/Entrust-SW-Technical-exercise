# Entrust-SW-Technical-exercise


This is a .txt pagination application made with golang (Go) where: 

  -	Each line consists of a maximum of 80 characters.

  - Each page has 25 lines, at the end of the page theres separation mark that includes the page number.

  - The paginated document can be in .txt or .pdf


# Dependencies

`github.com/jung-kurt/gofpdf` (already installed in the repository)

# Environment configuration

1. You need to have Go installed, you can install it [here](https://go.dev/doc/install).
2. Clone the repository in a local directory.


# How to run it 

1. go build pagination.go
2. ./pagination document_name [txt | pdf]

`document_name:` the name of the .txt file you want to paginate with the file extension.

`[txt | pdf]:` it's the file extension of the output of the paginated document, the default mode it's txt.

