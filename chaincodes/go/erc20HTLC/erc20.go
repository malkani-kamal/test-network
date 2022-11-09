package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Tokens struct {
	TokenId     string `json:"tokenId"`
	TokenName   string `json:"tokenName"`
	TokenOrg    string `json:"tokenOrg"`
	TokenSupply int    `json:"tokenAvailable"`
	TokenIssued int    `json:"tokenIssued"`
	TokenType   string `json:"tokenType"`
}

type Accounts struct {
	AccountId     string `json:"accountId"`
	AccountName   string `json:"accountName"`
	TokenId       string `json:"tokenId"`
	BalanceTokens int    `json:"balanceTokens"`
}

type Transaction struct {
	FromAccountId  string `json:"fromAccountId"`
	ToAccountId    string `json:"toAccountId"`
	TokenId        string `json:"tokenId"`
	ConversionRate int    `json:"conversionRate"`
	Amount         int    `json:"amount"`
}

//Create token definition. TokenSupply, TokenIssued, TokenType, TokenName.
//There can be multiple types of Token available
func (contract *SmartContract) CreateToken(ctx contractapi.TransactionContextInterface, tokenDefinitionsData string) (err error) {

	fmt.Printf("CreateToken start-->")
	stub := ctx.GetStub()
	var tokens Tokens

	err1 := json.Unmarshal([]byte(tokenDefinitionsData), &tokens)
	if err1 != nil {
		return fmt.Errorf("Failed to parse token argument. %s", err1.Error())
	}

	err = stub.PutState(tokens.TokenId, []byte(tokenDefinitionsData))
	if err != nil {
		return fmt.Errorf("Failed to insert token in ledger. %s", err.Error())
	}
	fmt.Println("Token created successfully")
	fmt.Errorf("Account details not found due to - %s", err1.Error())

	return nil
}

//MintToken will add the token to Minter's account if token's are available in Token Definition.
//It will also reduce the number of tokens from Token Definition.
func (contract *SmartContract) MintToken(ctx contractapi.TransactionContextInterface, transaction string) (string, error) {
	fmt.Printf("MintToken: %s", transaction)

	var transactionData Transaction
	errs := json.Unmarshal([]byte(transaction), &transactionData)
	if errs != nil {
		return "", fmt.Errorf("Failed to parse transaction data.", errs.Error())
	}

	var account Accounts
	account.AccountId = transactionData.ToAccountId
	account.TokenId = transactionData.TokenId
	key, _ := ctx.GetStub().CreateCompositeKey("account", []string{account.AccountId, account.TokenId})

	//checking the existing token balance and increasing if tokens are already available
	balance, _ := contract.GetBalance(ctx, key)
	fmt.Println("balance-", balance)
	if balance == -1 {
		account.BalanceTokens = transactionData.Amount
	} else {
		account.BalanceTokens = transactionData.Amount + balance
	}

	fmt.Println(transactionData.FromAccountId, balance)

	//Checking if Token definition have token supply available
	var token Tokens
	tokenAccount, err := ctx.GetStub().GetState(account.TokenId)
	err = json.Unmarshal(tokenAccount, &token)
	if err != nil {
		fmt.Errorf("Error in parsing the token data - %s", err.Error())
		return "fail", err
	}
	if (token.TokenSupply - token.TokenIssued) < transactionData.Amount {
		fmt.Errorf("Insufficient token!")
		return "fail", fmt.Errorf("Insufficient token!")
	} else {
		token.TokenIssued += transactionData.Amount
	}

	fmt.Println("Checked tokens has balance")

	//Adding token to Minters account balance
	accountTxn, err := json.Marshal(account)
	err = ctx.GetStub().PutState(key, accountTxn)
	if err != nil {
		return "", fmt.Errorf("Error while adding transaction to ledger - %s", err.Error())
	}
	fmt.Println("Account updated with token")

	//Reducing token from token supply.
	tokenTxn, err := json.Marshal(token)
	err = ctx.GetStub().PutState(account.TokenId, tokenTxn)
	if err != nil {
		return "", fmt.Errorf("Error while adding transaction to ledger - %s", err.Error())
	}

	fmt.Println("Tokens reduced from token balance.")

	return "success", nil
}

//Function to check the Toke balance of particular account.
func (contract *SmartContract) GetBalance(ctx contractapi.TransactionContextInterface, id string) (int, error) {

	fmt.Printf("GetBalance: %s", id)

	accountBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		fmt.Errorf("Error in fecthing product from world state - %s", err.Error())
		return -1, fmt.Errorf("Error in fecthing account from world state - %s", err.Error())
	}
	// fmt.Printf("accountBytes: %s", accountBytes)

	if accountBytes == nil {
		fmt.Errorf("Account does not exist with ID - %s", id)
		return -1, nil
	}

	var account Accounts
	errs := json.Unmarshal(accountBytes, &account)
	if errs != nil {
		fmt.Errorf("Error in parsing the account data - %s", errs.Error())
		return -1, errs
	}

	return account.BalanceTokens, nil
}

//Tansferring token from one account to other account
func (contract *SmartContract) transactToken(ctx contractapi.TransactionContextInterface, key string, amount int, action string) (string, error) {

	fmt.Printf("BurnToken-->", key, amount)

	var account Accounts
	accountData, _ := ctx.GetStub().GetState(key)
	err := json.Unmarshal(accountData, &account)
	if err != nil {
		fmt.Errorf("Error in parsing the token data - %s", err.Error())
		return "fail", err
	}

	//If operator is + then adding the balance else reducing the balance.
	if action == "+" {
		account.BalanceTokens += amount
	} else {
		//Checking the balance before reducing the number of tokens
		if account.BalanceTokens < amount {
			fmt.Errorf("Insufficient balance in account!")
			return "fail", fmt.Errorf("Insufficient balance in account!")
		}
		account.BalanceTokens -= amount
	}

	//Updating the token balance in the ledger.
	accountTxn, _ := json.Marshal(account)
	err = ctx.GetStub().PutState(key, accountTxn)
	if err != nil {
		return "", fmt.Errorf("Error while burn token to ledger - %s", err.Error())
	}

	return "success", nil
}

func (contract *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, transaction string) (string, error) {

	var transactionData Transaction

	errs := json.Unmarshal([]byte(transaction), &transactionData)
	if errs != nil {
		return "fail", nil
	}

	//creating composite key AccountId+TokenId
	fromKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{transactionData.FromAccountId, transactionData.TokenId})
	//First reducting the balance from From Account
	msg, err := contract.transactToken(ctx, fromKey, transactionData.Amount, "-")
	if err != nil {
		//Return incase insufficient balance error or any other error while reducing the balance
		return msg, fmt.Errorf("%s", err.Error())
	}

	toKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{transactionData.ToAccountId, transactionData.TokenId})
	//Adding the balance to To Account
	contract.transactToken(ctx, toKey, transactionData.Amount, "+")

	return "success", nil
}

// structure for the timelock
type HashTimeLock struct {
	LockID   string `json:"lockid"`
	FromID   string `json:"fromid"`
	ToID     string `json:"toid"`
	TokenId  string `json:"tokenid"`
	Amount   int    `json:"amount"`
	HashLock string `json:"hashlock"`
	TimeLock int64  `json:"timelock"`
}

func (contract *SmartContract) TransferConditional(ctx contractapi.TransactionContextInterface, lockId string, hashlock string, timelock string, transaction string) (string, error) {

	var transactionData Transaction

	errs := json.Unmarshal([]byte(transaction), &transactionData)
	if errs != nil {
		return "fail", nil
	}

	//Reducing the balance from From Account
	fromKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{transactionData.FromAccountId, transactionData.TokenId})
	msg, err := contract.transactToken(ctx, fromKey, transactionData.Amount, "-")
	if err != nil {
		return msg, fmt.Errorf("%s", err.Error())
	}

	// timeLockconv, err := strconv.ParseInt(timelock, 10, 64)
	// fmt.Println("timeLockconv:", timeLockconv)

	timeInt, err := strconv.Atoi(timelock)
	if err != nil {
		return "Error converting timeLock.", fmt.Errorf("%s", err.Error())
	}

	expiryTime := time.Now().Add(time.Minute * time.Duration(timeInt)).Unix()

	var hashTimeLock HashTimeLock

	hashTimeLock.FromID = transactionData.FromAccountId
	hashTimeLock.ToID = transactionData.ToAccountId
	hashTimeLock.Amount = transactionData.Amount
	hashTimeLock.TokenId = transactionData.TokenId
	hashTimeLock.LockID = lockId
	hashTimeLock.HashLock = strings.ToLower(hashlock)
	hashTimeLock.TimeLock = expiryTime

	hashTimeLockAsBytes, _ := json.Marshal(hashTimeLock)

	//Parking the transaction in hashTimeLock
	ctx.GetStub().PutState(lockId, hashTimeLockAsBytes)

	return "Conditional transfer successful! Hash Lock created.", nil

}

func (contract *SmartContract) Claim(ctx contractapi.TransactionContextInterface, lock_id string, pwd string) string {

	hashTimeLockAsBytes, _ := ctx.GetStub().GetState(lock_id)

	hashTimeLock := new(HashTimeLock)
	_ = json.Unmarshal(hashTimeLockAsBytes, hashTimeLock)

	hash := sha256.Sum256([]byte(pwd))

	hashString := fmt.Sprintf("%x", hash)

	fmt.Println("Hash String:", hashString)
	fmt.Println("hashTimeLock---", hashTimeLock)

	if hashTimeLock.HashLock != hashString {

		fmt.Println("Invalid password:", hashString, hashTimeLock.HashLock)
		fmt.Println("Transaction to be reverted:")

		return "invalid password"
	}

	currTime := time.Now().Unix()

	fmt.Println("currTime-", currTime)

	//Checking if hashTimeLock expired or no
	if hashTimeLock.TimeLock > currTime {
		//Adding balance to ToAccount from hashLockTime
		toKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{hashTimeLock.ToID, hashTimeLock.TokenId})
		contract.transactToken(ctx, toKey, hashTimeLock.Amount, "+")

		// Deleting hashTimeLock after adding the balance to claimer account
		ctx.GetStub().DelState(lock_id)

		fmt.Println("Timelock still active. Actual transaction timestamp:", hashTimeLock.TimeLock, "Actual timelock:", currTime)
		return "Tokens claimed successfully!"
	} else {
		return "Tokens claimed unsuccessful!  Timelock expired."
	}

}

func (contract *SmartContract) Revert(ctx contractapi.TransactionContextInterface, lock_id string) string {

	hashTimeLockAsBytes, _ := ctx.GetStub().GetState(lock_id)

	hashTimeLock := new(HashTimeLock)
	_ = json.Unmarshal(hashTimeLockAsBytes, hashTimeLock)

	currTime := time.Now().Unix()
	fmt.Println("currTime-", currTime)
	if hashTimeLock.TimeLock > currTime {
		return "Timelock still not expired and token are yet open for claim."
	} else {

		//Adding balance to FromAccount from hashLockTime
		fromKey, _ := ctx.GetStub().CreateCompositeKey("account", []string{hashTimeLock.FromID, hashTimeLock.TokenId})
		contract.transactToken(ctx, fromKey, hashTimeLock.Amount, "+")

		// Deleting hashTimeLock after adding the balance to Minters account
		ctx.GetStub().DelState(lock_id)
		fmt.Println("Revert of Tokens successful.")
		return "Revert of Tokens successful to Minter."
	}
}
