package controller

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"github.com/jeromelesaux/ethereum-training/client"
	"github.com/jeromelesaux/ethereum-training/config"
)

type Controller struct {
}

//
// curl -v -X POST http://localhost:8080/anchor   -F "file=@readme.txt"  -H "Content-Type: multipart/form-data"
//
func (ctr *Controller) Anchoring(c *gin.Context) {

	// get the file from multipart form
	f, err := c.FormFile("file")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while getting file multipart header.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	outfile := filepath.Join(config.MyConfig.GetFilepaths(), f.Filename)
	// save the file on system
	err = c.SaveUploadedFile(f, outfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not save file on system with error :%v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	sum256, hexa256, err := getSha256(outfile)
	if err != nil {
		sendJsonError(c, "Can not get the sha256 for file "+f.Filename, err)
		return
	}

	txHash, err := sendTransaction(sum256)
	if err != nil {
		sendJsonError(c, "Error in ethereum transaction", err)
		return
	}
	// return json ok result
	c.JSON(http.StatusOK, gin.H{
		"tx":      txHash,
		"message": "filename:" + f.Filename + " has hash256 " + hexa256,
	})
	return
}

//
// curl -v -X POST http://localhost:8080/verify   -F "file=@readme.txt" -f "hash=8753d45d70da590b0841392ac762161ac5230fa63a5b766759e6fd0d33a65631" -H "Content-Type: multipart/form-data"
//
func (ctr *Controller) Verify(c *gin.Context) {

	txHash := c.PostForm("txhash")

	// get the file from multipart form
	f, err := c.FormFile("file")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while getting file multipart header.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	outfile := filepath.Join(config.MyConfig.GetFilepaths(), f.Filename)
	// save the file on system
	err = c.SaveUploadedFile(f, outfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not save file on system with error :%v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	_, hexa256, err := getSha256(outfile)
	if err != nil {
		sendJsonError(c, "Can not get the sha256 for file "+f.Filename, err)
		return
	}

	// get the informations from the tx
	tx, isPending, err := client.EthClient.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		sendJsonError(c, "Can not get the transaction informations", err) // change to 404
		return
	}
	if isPending {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Transaction is still mining.",
		})
		return
	}
	data := tx.Data()
	hashInBlockChain := fmt.Sprintf("%x", data)

	if hexa256 != hashInBlockChain {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The hash of this document does not match from the transaction data.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "GREAT...OoOoOOOOoooo",
	})

	return
}

//
// curl -v -X POST http://localhost:8080/anchormultiple   -F "upload[]=@readme.txt" -F "upload[]=@hello.txt" -H "Content-Type: multipart/form-data"
//
func (ctr *Controller) AnchorMultiple(c *gin.Context) {
	form, _ := c.MultipartForm()
	var merkleContent []merkletree.Content
	var merkleHexas []merkleHexa
	files := form.File["upload[]"]
	for _, file := range files {
		now := time.Now()
		directoryName := now.Format(time.RFC3339Nano)
		directoryBase := filepath.Join(config.MyConfig.GetFilepaths(), directoryName)
		if err := os.MkdirAll(directoryBase, os.ModePerm); err != nil {
			sendJsonError(c, "Can not create local directory", err)
			return
		}
		outfile := filepath.Join(directoryBase, file.Filename)
		c.SaveUploadedFile(file, outfile)
		sum256, hexaHash, err := getSha256(outfile)
		if err != nil {
			sendJsonError(c, "error", err)
			return
		}
		merkleContent = append(merkleContent, NewMerkleContent(sum256, hexaHash))
		merkleHexas = append(merkleHexas, merkleHexa{Hexa: hexaHash})
	}

	t, err := merkletree.NewTree(merkleContent)
	if err != nil {
		sendJsonError(c, "Can not create merkletree", err)
		return
	}
	merkleRoot := t.MerkleRoot()
	//hexaMerkleRoot := fmt.Sprintf("%x", merkleRoot)
	txHash, err := sendTransaction(merkleRoot)
	if err != nil {
		sendJsonError(c, "Error in transaction", err)
		return
	}

	if err := saveToFile(merkleHexas, filepath.Join(config.MyConfig.GetHashpaths(), txHash)); err != nil {
		sendJsonError(c, "Can not save merkle content into file", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tx":      txHash,
		"message": "filenames uploaded on blockchain",
	})
}

//
//  curl -v -X POST http://localhost:8080/verifymultiple   -F "file=@hello.txt" -F "txhash=0x3421c80698d64e2f227d1af487015032fa455707bfc4e5373236e78405a00dbc" -H "Content-Type: multipart/form-data"
//
func (ctr *Controller) VerifyMultiple(c *gin.Context) {
	txHash := c.PostForm("txhash")
	merklecontent, err := readFromFile(filepath.Join(config.MyConfig.GetHashpaths(), txHash))
	if err != nil {
		sendJsonError(c, "Can not read file containing the merklecontent", err)
		return
	}

	// get the file from multipart form
	f, err := c.FormFile("file")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while getting file multipart header.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	outfile := filepath.Join(config.MyConfig.GetFilepaths(), f.Filename)
	// save the file on system
	err = c.SaveUploadedFile(f, outfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not save file on system with error :%v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	_, hexa256, err := getSha256(outfile)
	if err != nil {
		sendJsonError(c, "Can not get the sha256 for file "+f.Filename, err)
		return
	}

	// get the informations from the tx
	tx, isPending, err := client.EthClient.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		sendJsonError(c, "Can not get the transaction informations", err) // change to 404
		return
	}
	if isPending {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Transaction is still mining.",
		})
		return
	}
	data := tx.Data()
	hashInBlockChain := fmt.Sprintf("%x", data)

	// check the hash in data of the transaction
	if hashInBlockChain == "" {
		sendJsonError(c, "Transaction does not containt merkltree.", errors.New("Empty data in transaction"))
		return
	}

	for _, v := range merklecontent {
		if v.Hexa == hexa256 {
			c.JSON(http.StatusOK, gin.H{
				"message": "GREAT...OoOoOOOOoooo",
			})
			return
		}
	}

	sendJsonNotFound(c, "file is not certified", errors.New(f.Filename+" not found."))
	return
}

func readFromFile(filePath string) (content []merkleHexa, err error) {
	fr, err := os.Open(filePath)
	if err != nil {
		return []merkleHexa{}, err
	}
	defer fr.Close()

	err = json.NewDecoder(fr).Decode(&content)
	return content, err
}

func saveToFile(content []merkleHexa, filePath string) error {
	fw, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fw.Close()

	return json.NewEncoder(fw).Encode(&content)
}

func sendJsonNotFound(c *gin.Context, msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s :%v\n", msg, err)
	c.JSON(http.StatusNotFound, gin.H{
		"error": err,
	})
	return
}

func sendJsonError(c *gin.Context, msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s :%v\n", msg, err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err,
	})
	return
}

func sendTransaction(data []byte) (string, error) {
	client.SafeNonceTx.NonceMutex.Lock() // lock nonce for operation

	nonce := client.SafeNonceTx.Nonce                                       // current nonce
	to := client.Auth.From                                                  // destination address
	value := big.NewInt(0)                                                  // value of the transaction                                                           // data to store in blockchain
	var gasLimit uint64 = 100000                                            // gas limit
	gasPrice, err := client.EthClient.SuggestGasPrice(context.Background()) // gas price
	if err != nil {
		log.Fatal(err)
	}

	// start ethereum transaction
	var tx = types.NewTransaction(
		nonce,
		to,
		value,
		gasLimit,
		gasPrice,
		data,
	)

	// find the id of the chain to use (instance test chainid)
	chainID, err := client.EthClient.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), client.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	// send the transaction to the network
	err = client.EthClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	// end tx
	client.SafeNonceTx.Nonce++
	client.SafeNonceTx.NonceMutex.Unlock()

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
	// return json ok result
	return signedTx.Hash().Hex(), nil
}

func getSha256(filePath string) ([]byte, string, error) {
	// re-open local file
	fh, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not open file from system with error :%v\n", err)
		return []byte{}, "", err
	}
	defer fh.Close()

	// compute sha256 sum
	h := sha256.New()

	if _, err := io.Copy(h, fh); err != nil {
		fmt.Fprintf(os.Stderr, "Can not read file from system with error :%v\n", err)
		return []byte{}, "", err
	}

	sum256 := h.Sum(nil)
	hexaSum := fmt.Sprintf("%x", sum256)
	// ok display result
	fmt.Fprintf(os.Stderr, "filename [%s] has sha256 [%s]\n",
		filePath,
		hexaSum)

	return sum256, hexaSum, nil
}

type MerkleContent struct {
	hexa string
	hash []byte
}

func (m *MerkleContent) Equals(other merkletree.Content) (bool, error) {
	return m.hexa == other.(*MerkleContent).hexa, nil
}

func NewMerkleContent(hash []byte, hexa string) *MerkleContent {
	return &MerkleContent{
		hexa: hexa,
		hash: hash,
	}
}

func (m *MerkleContent) CalculateHash() ([]byte, error) {
	if len(m.hash) > 0 {
		return m.hash, nil
	}
	h := sha256.New()
	if _, err := h.Write([]byte(m.hexa)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

type merkleHexa struct {
	Hexa string
}